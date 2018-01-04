package main

import (
	"context"
	"os"
	"flag"

	"github.com/go-kit/kit/log"
	"net/http"
	"os/signal"
	"syscall"
	"fmt"
)

const (
	defaultPort = "8605"
	defaultRoutingServiceURL = "http://localhost;8606"
)

func main() {
	var (
		// read environment variables or use the default values from above
		addr = envString("PORT", defaultPort)
		rsurl = envString("ROUTINGSERVICE_URL", defaultRoutingServiceURL)

		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
		routingServiceURL = flag.String("service.routing", rsurl, "routing service URL")

		ctx = context.Background()
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// TODO init some stuff here?

	// TODO init all service here

	// now comes the HTTP REST API stuff
	httpLogger := log.With(logger, "component", "http")

	// init mux router
	mux := http.NewServeMux()

	// all services are declared here
	var hs 


	// error handling
	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", "httpAddr", "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}