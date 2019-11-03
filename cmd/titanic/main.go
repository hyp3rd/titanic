package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	titanic "gitlab.com/hyperd/titanic"
	"gitlab.com/hyperd/titanic/cockroachdb"
	titanicsvc "gitlab.com/hyperd/titanic/implementation"
	"gitlab.com/hyperd/titanic/inmemory"
	"gitlab.com/hyperd/titanic/middleware"
	httptransport "gitlab.com/hyperd/titanic/transport/http"
)

func main() {
	var (
		httpAddr     = flag.String("http.addr", ":3000", "HTTP listen address")
		httpsAddr    = flag.String("https.addr", ":8443", "HTTPS listen address")
		databaseType = flag.String("database.type", "cockroachdb", "Database type")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)

		logger = log.With(logger,
			"svc", "titanic",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")

	defer level.Info(logger).Log("msg", "service ended")

	selectedBackend := *databaseType // safe to dereference the *string
	isInMemory := (selectedBackend == "inmemory")

	var db *gorm.DB
	{
		if !isInMemory {
			var err error

			const addr = "postgresql://d4gh0s7@roach1:26257/titanic?sslmode=disable"

			db, err = gorm.Open("postgres", addr)

			if err != nil {
				level.Error(logger).Log("exit", err)
				os.Exit(-1)
			}
			defer db.Close()

			// Set to `true` and GORM will print out all DB queries.
			db.LogMode(true)

			// Disable table name's pluralization globally
			db.SingularTable(true)
			db.AutoMigrate(&titanic.People{})
		}
	}

	var svc titanic.Service
	{
		if isInMemory {
			level.Info(logger).Log("backend", "database", "type", "inmemory")

			repository, err := inmemory.NewInmemService(logger)
			if err != nil {
				level.Error(logger).Log("exit", err)
				os.Exit(-1)
			}
			svc = titanicsvc.NewService(repository, logger)
		} else {
			level.Info(logger).Log("backend", "database", "type", "cockroachdb")

			repository, err := cockroachdb.New(db, logger)
			if err != nil {
				level.Error(logger).Log("exit", err)
				os.Exit(-1)
			}
			svc = titanicsvc.NewService(repository, logger)
		}
		// Service middleware: Logging
		svc = middleware.LoggingMiddleware(logger)(svc)

	}

	var h http.Handler
	{
		h = httptransport.MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	go func() {
		logger.Log("transport", "HTTPS", "addr", *httpsAddr)
		errs <- http.ListenAndServeTLS(*httpsAddr, "/etc/tls/certs/tls.crt", "/etc/tls/certs/tls.key", h)
		// errs <- http.ListenAndServeTLS(*httpsAddr, "./tls/tls.crt", "./tls/tls.key", h)

	}()

	logger.Log("exit", <-errs)
}
