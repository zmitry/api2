package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/starius/api2"
	"github.com/starius/api2/example"

	"go.uber.org/fx"
)

func NewLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)
	logger.Print("Executing NewLogger.")
	return logger
}

func NewMux(lc fx.Lifecycle, logger *log.Logger) *http.ServeMux {
	logger.Print("Executing NewMux.")
	// First, we construct the mux and server. We don't want to start the server
	// until all handlers are registered.
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// If NewMux is called, we know that another function is using the mux. In
	// that case, we'll use the Lifecycle type to register a Hook that starts
	// and stops our HTTP server.
	//
	// Hooks are executed in dependency order. At startup, NewLogger's hooks run
	// before NewMux's. On shutdown, the order is reversed.
	//
	// Returning an error from OnStart hooks interrupts application startup. Fx
	// immediately runs the OnStop portions of any successfully-executed OnStart
	// hooks (so that types which started cleanly can also shut down cleanly),
	// then exits.
	//
	// Returning an error from OnStop hooks logs a warning, but Fx continues to
	// run the remaining hooks.
	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 15 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			logger.Print("Starting HTTP server.")
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Print("Stopping HTTP server.")
			return server.Shutdown(ctx)
		},
	})

	return mux
}

func Register(mux *http.ServeMux, routes []api2.Route) {
	api2.BindRoutes(mux, routes)
}

var Module = fx.Options(
	fx.Provide(
		example.NewEchoRepository,
		example.NewEchoService,
		example.GetRoutes,
		func(s *example.EchoService) example.IEchoService { return s },
		NewLogger,
		NewMux,
	))

func main() {
	fx.New(
		Module,
		fx.Invoke(Register),
	).Run()
}
