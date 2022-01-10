package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/classic-massok/classic-massok-be/api/core"
	"github.com/classic-massok/classic-massok-be/config"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	cfg, err := config.RenderConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "Classic Massok BE",
		RunE: func(cmd *cobra.Command, args []string) error {
			errChan := make(chan error)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("%s:%d", cfg.Database.URI, cfg.Database.Port)))
			if err != nil {
				return fmt.Errorf("error creating mongo client: %w", err)
			}

			if err = client.Ping(ctx, readpref.Primary()); err != nil {
				return fmt.Errorf("error connecting to mongo: %w", err)
			}

			db := client.Database(cfg.Database.Name)

			go func() {
				err := serveHTTP(getEchoRouter(db, cfg), cfg)
				errChan <- err
			}()

			for err := range errChan {
				return err
			}

			return nil
		},
	}

	if err = rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serveHTTP(echoRouter http.Handler, cfg *config.Config) error {
	handler := createHandler(
		echoRouter.ServeHTTP, cfg,
		panicsReturn500,
	)

	port := fmt.Sprintf(":%d", cfg.Server.Port)

	server := &http.Server{
		Addr:           port,
		Handler:        handler,
		ReadTimeout:    180 * time.Minute,
		WriteTimeout:   180 * time.Minute,
		MaxHeaderBytes: 16384,
	}

	fmt.Println("server listening on ", port)
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}

		return errors.Wrap(err, "failed to listen")
	}

	return nil
}

func panicsReturn500(next http.HandlerFunc, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				// TODO: create logger for errors and log stack traces
				var trace interface{}
				if cfg.Logging.HTTPVerbose {
					trace = fmt.Sprintf("%s\n\n%s", r, string(debug.Stack()))
				}

				data, err := json.Marshal(core.JSON(req.Context().Value(lib.EchoContextKey).(echo.Context), 500, trace, lib.ErrServerError))
				if err != nil {
					// TODO: log error returning panic 500 here
					return
				}

				if _, err = w.Write(data); err != nil {
					// TODO: log error writing http error response
					return
				}
			}

			debug.PrintStack()
		}()

		next(w, req)
	}
}

func createHandler(handler http.HandlerFunc, cfg *config.Config, middleware ...func(fn http.HandlerFunc, cfg *config.Config) http.HandlerFunc) http.HandlerFunc {
	for index := range middleware {
		handler = middleware[len(middleware)-index-1](handler, cfg)
	}
	return handler
}
