package main

import (
	"os"
	"flag"

	"github.com/go-kit/kit/log"
	"net/http"
	"os/signal"
	"syscall"
	"fmt"
	"github.com/MICSTI/imsazon/hello"
	"github.com/MICSTI/imsazon/inmemory"
	"github.com/MICSTI/imsazon/auth"
)

const (
	defaultPort = "8605"
)

func main() {
	var (
		// read environment variables or use the default values from above
		addr = envString("PORT", defaultPort)

		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")

		// TODO do we need this?
		//ctx = context.Background()
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// init in-memory repository stores here
	var (
		users			= inmemory.NewUserRepository()
	)

	// all services are initialized here
	var hs hello.Service
	hs = hello.NewService()
	hs = hello.NewLoggingService(log.With(logger, "component", "hello"), hs)

	var as auth.Service
	as = auth.NewService(users)
	as = auth.NewLoggingService(log.With(logger, "component", "auth"), as)

	// now comes the HTTP REST API stuff
	httpLogger := log.With(logger, "component", "http")

	// init router
	mux := http.NewServeMux()

	mux.Handle("/hello/", hello.MakeHandler(hs, httpLogger))
	mux.Handle("/auth/", auth.MakeHandler(as, httpLogger))

	http.Handle("/", accessControl(mux))

	// error handling
	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
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