package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// versionString is overwritten using compiler ldflags
var versionString = "0.0.0-local"
var healthy = true

const envPrefix = "TESTDUMMY"
const terminationLogPath = "/dev/termination-log"

type RuntimeConfig struct {
	Healthy              bool   `default:"true"`
	HealthyAfterSeconds  *int   `split_words:"true"`
	BindAddress          string `split_words:"true" default:"localhost:8000"`
	PanicSeconds         *int   `split_words:"true"`
	EnableRequestLogging bool   `split_words:"true" default:"false"`
	EnableEnv            bool   `split_words:"true" default:"false"`
	RootPath             string `split_words:"true" default:"/"`
	StressCpuDuration    string `split_words:"true" default:"0s"`
}

func main() {
	var rc RuntimeConfig
	err := envconfig.Process(envPrefix, &rc)
	ExitIfErr(err, "Error loading environment variables")

	healthy = rc.Healthy

	logger := log.New(os.Stdout, "", log.LstdFlags)

	if rc.HealthyAfterSeconds != nil {
		healthy = false
		time.AfterFunc(time.Duration(*rc.HealthyAfterSeconds)*time.Second, func() {
			healthy = true
		})
	}

	if rc.PanicSeconds != nil {
		time.AfterFunc(time.Duration(*rc.PanicSeconds)*time.Second, func() {
			panic("Panicing due to TESTDUMMY_PANIC_SECONDS being set")
		})
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    rc.BindAddress,
		Handler: logging(rc, logger)(mux),
	}

	addRoute := func(path string, handler func(http.ResponseWriter, *http.Request)) {
		mux.HandleFunc(filepath.Join(rc.RootPath, path), handler)
	}

	stressDuration, err := time.ParseDuration(rc.StressCpuDuration)
	ExitIfErr(err, "Invalid StressCpuDuration")

	addRoute("", func(w http.ResponseWriter, r *http.Request) { pingHandler(w, r, stressDuration) })
	addRoute("ping", func(w http.ResponseWriter, r *http.Request) { pingHandler(w, r, stressDuration) })
	addRoute("echo", echoHandler)
	addRoute("health", healthHandler)
	addRoute("version", versionHandler)
	addRoute("exit", exitHandler)
	addRoute("status", statusHandler)
	if rc.EnableEnv {
		addRoute("env", envHandler)
	}

	logger.Printf("TestDummy v%s", versionString)
	logger.Printf("Listening on %s\n", rc.BindAddress)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		ExitIfErr(err, "Unable to start server")
	}
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	for _, env := range os.Environ() {
		_, err := w.Write([]byte(fmt.Sprintf("%s\n", env)))
		LogIfErr(err, "Error writing response to /env")
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request, stressDuration time.Duration) {
	if stressDuration > 0 {
		log.Printf("Stressing CPU for %s", stressDuration)
		stressCpu(stressDuration)
	}
	w.WriteHeader(200)
	_, err := w.Write([]byte("pong"))
	LogIfErr(err, "Error writing response to /ping")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	healthyParam := r.URL.Query().Get("healthy")
	if healthyParam != "" {
		healthyParamParsed, err := strconv.ParseBool(healthyParam)
		if err == nil {
			log.Printf("Setting healthy to: %v", healthyParamParsed)
			healthy = healthyParamParsed
		}
	}

	if healthy {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		bodyBytes = []byte(fmt.Sprintf("Unable to read body: %s", err))
	}
	_, err = w.Write(bodyBytes)
	LogIfErr(err, "Error writing response to /echo")
}

func stressCpu(duration time.Duration) {
	done := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
				}
			}
		}()
	}

	time.Sleep(duration)
	close(done)
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	exitCode := 1
	codeParam := r.URL.Query().Get("code")
	if codeParam != "" {
		codeParamParsed, err := strconv.ParseInt(codeParam, 10, 32)
		if err == nil {
			exitCode = int(codeParamParsed)
		}
	}

	terminationError := errors.New("Fatal error")

	err := ioutil.WriteFile(terminationLogPath, []byte(fmt.Sprintf("%+v", terminationError)), 0666)
	LogIfErr(err, "Error writing to termination log")

	os.Exit(exitCode)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte(versionString))
	LogIfErr(err, "Error writing response to /version")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	status := 400
	statusParam := r.URL.Query().Get("status")
	if statusParam != "" {
		statusParamParsed, err := strconv.ParseInt(statusParam, 10, 32)
		if err == nil {
			status = int(statusParamParsed)
		}
	}

	w.WriteHeader(status)
}

func LogIfErr(err error, message string) {
	if err != nil {
		log.Printf("%s: %s", message, err)
	}
}

func ExitIfErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func logging(rc RuntimeConfig, logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rc.EnableRequestLogging {
					logger.Println(r.Method, r.URL.Path, r.RemoteAddr)
				}
				next.ServeHTTP(w, r)
			}()
		})
	}
}
