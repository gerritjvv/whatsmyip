package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gerritjvv/whatsmyip/internal"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type RespErr struct {
	Message string `json:"message"`
}

func main() {
	defer os.Exit(0)

	var wait time.Duration
	var port int
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.IntVar(&port, "port", 8007, "port=8007")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/myip", func(writer http.ResponseWriter, request *http.Request) {
		resp := internal.ReturnMyIpHandler(request)
		respBytes, err := json.Marshal(resp)
		if err != nil {
			writer.WriteHeader(500)
			writeError(writer, err)
			return
		}

		writer.WriteHeader(200)
		_, err = writer.Write(respBytes)

		// we already wrote the header, and cannot write again
		// if any error, write the error object
		if err != nil {
			writeError(writer, err)
		}
	})

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 10,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		fmt.Printf("Listening on %d\n", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
}

func writeError(writer http.ResponseWriter, err error) {
	errObj := RespErr{
		Message: err.Error(),
	}
	errObjBytes, err2 := json.Marshal(errObj)
	if err2 != nil {
		fmt.Println(err2)
	}

	_, err2 = writer.Write(errObjBytes)
	if err2 != nil {
		fmt.Println(err2)
	}
}
