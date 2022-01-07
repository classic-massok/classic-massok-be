package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "Classic Massok BE",
		RunE: func(cmd *cobra.Command, args []string) error {
			errChan := make(chan error)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:5555"))
			if err != nil {
				return fmt.Errorf("error creating mongo client: %w", err)
			}

			if err = client.Ping(ctx, readpref.Primary()); err != nil {
				return fmt.Errorf("error connecting to mongo: %w", err)
			}

			db := client.Database("classic-massok")

			go func() {
				err := serveHTTP(getEchoRouter(db))
				errChan <- err
			}()

			for err := range errChan {
				return err
			}

			return nil
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serveHTTP(echoRouter http.Handler) error { // TODO: make this configurable
	handler := createHandler(echoRouter.ServeHTTP)
	port := ":8080"

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

func createHandler(handler http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for index := range middleware {
		handler = middleware[len(middleware)-index-1](handler)
	}
	return handler
}
