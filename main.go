package main

import (
	"flag"
	"os"

	"fmt"
	"github.com/MICSTI/imsazon/auth"
	"github.com/MICSTI/imsazon/hello"
	"github.com/MICSTI/imsazon/inmemory"
	"github.com/go-kit/kit/log"
	"net/http"
	"os/signal"
	"syscall"
	"github.com/MICSTI/imsazon/mail"
	"github.com/creamdog/gonfig"
	log2 "log"
	"github.com/MICSTI/imsazon/stock"
	"github.com/MICSTI/imsazon/payment"
	"github.com/MICSTI/imsazon/cart"
	"github.com/MICSTI/imsazon/order"
)

const (
	defaultPort = "8605"
)

func main() {
	var (
		// read environment variables or use the default values from above
		addr = envString("PORT", defaultPort)
		configFilePath = envString("CONFIG_FILE_PATH", "")

		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")

		// TODO do we need this?
		//ctx = context.Background()
	)

	flag.Parse()

	f, err := os.Open(configFilePath)
	if err != nil {
		log2.Fatal("Could not read config file")
	}
	defer f.Close();
	config, err := gonfig.FromJson(f)
	if err != nil {
		log2.Fatal("Could not get config from config file")
	}

	// read JWT secret from config file
	jwtSecretString, err := config.GetString("jwt/secret", "")
	if err != nil {
		log2.Fatal("Could not get JWT secret config value")
	}

	// parse string to byte array
	jwtSecret := []byte(jwtSecretString)

	// Mail configuration
	mailHost, err := config.GetString("mail/host", "")
	if err != nil {
		log2.Fatal("Could not get mail host config value")
	}

	mailPort, err := config.GetInt("mail/port", 0)
	if err != nil {
		log2.Fatal("Could not get mail port config value")
	}

	mailUsername, err := config.GetString("mail/username", "")
	if err != nil {
		log2.Fatal("Could not get mail username config value")
	}

	mailPassword, err := config.GetString("mail/password", "")
	if err != nil {
		log2.Fatal("Could not get mail password config value")
	}

	mailServerCredentials := mail.MailServerCredentials{
		Host: 		mailHost,
		Port:		mailPort,
		Username:	mailUsername,
		Password:	mailPassword,
	}

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// init in-memory repository stores here
	var (
		users = inmemory.NewUserRepository()
		products = inmemory.NewProductRepository()
		carts = inmemory.NewCartRepository()
		orders = inmemory.NewOrderRepository()
	)

	// all services are initialized here
	var hs hello.Service
	hs = hello.NewService()
	hs = hello.NewLoggingService(log.With(logger, "component", "hello"), hs)

	var as auth.Service
	as = auth.NewService(jwtSecret, users)
	as = auth.NewLoggingService(log.With(logger, "component", "auth"), as)

	var ms mail.Service
	ms = mail.NewService(mailServerCredentials)
	ms = mail.NewLoggingService(log.With(logger, "component", "mail"), ms)

	var sts stock.Service
	sts = stock.NewService(products)
	sts = stock.NewLoggingService(log.With(logger, "component", "stock"), sts)

	var ps payment.Service
	ps = payment.NewService()
	ps = payment.NewLoggingService(log.With(logger, "component", "payment"), ps)

	var cs cart.Service
	cs = cart.NewService(carts)
	cs = cart.NewLoggingService(log.With(logger, "component", "cart"), cs)

	var ors order.Service
	ors = order.NewService(orders)
	ors = order.NewLoggingService(log.With(logger, "component", "order"), ors)

	// now comes the HTTP REST API stuff
	httpLogger := log.With(logger, "component", "http")

	// init router
	mux := http.NewServeMux()

	mux.Handle("/hello/", hello.MakeHandler(hs, httpLogger))
	mux.Handle("/auth/", auth.MakeHandler(as, httpLogger))
	mux.Handle("/mail/", mail.MakeHandler(ms, httpLogger))
	mux.Handle("/stock/", stock.MakeHandler(sts, httpLogger))
	mux.Handle("/payment/", payment.MakeHandler(ps, httpLogger))
	mux.Handle("/cart/", cart.MakeHandler(cs, httpLogger))
	mux.Handle("/order/", order.MakeHandler(ors, httpLogger))

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
		// TODO auth code would go here

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
