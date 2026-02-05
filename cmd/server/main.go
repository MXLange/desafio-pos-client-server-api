package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MXLange/desafio-pos-client-server-api/cmd/server/handlers"
	"github.com/MXLange/desafio-pos-client-server-api/cmd/server/infra/db"
	"github.com/MXLange/desafio-pos-client-server-api/cmd/server/repository"
)



const (
	dbConnString = "file:app.db?mode=rwc"
	serverAddr   = ":8080"
	priceAPIURL  = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	callTimeout  = time.Duration(200) * time.Millisecond
	dbTimeout    = time.Duration(10) * time.Millisecond
)

func main() {

	db, err := db.NewDBConnection(dbConnString)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	repository, err := repository.New(db)
	if err != nil {
		panic(err)
	}

	err = repository.CreateTables()
	if err != nil {
		panic(err)
	}

	priceHandler, err := handlers.NewPriceHandler(repository, callTimeout, dbTimeout, priceAPIURL)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /cotacao", priceHandler.GetPrice)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Println("[SERVER] Server is running on", serverAddr)
		errCh <- server.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Println(err.Error())
			os.Exit(1)
		}
	case <-sigCh:
		log.Println("[SERVER] shutdown requested")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("[SERVER] "+err.Error())
		os.Exit(1)
	}

	log.Println("[SERVER] shutdown complete")
}