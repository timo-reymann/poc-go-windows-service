package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type HttpServer struct {
	port int
	srv  *http.Server
}

func NewHttpServer(port int) *HttpServer {
	return &HttpServer{
		port: port,
	}
}

func (hs *HttpServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! This is a Windows service running an HTTP server on port 2025")
}

func (hs *HttpServer) Start() error {
	hs.srv = &http.Server{
		Addr:    ":" + strconv.Itoa(hs.port),
		Handler: http.HandlerFunc(hs.handleRequest),
	}

	if err := hs.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (hs *HttpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := hs.srv.Shutdown(ctx); err != nil {
		return err
	}
	hs.srv = nil
	return nil
}
