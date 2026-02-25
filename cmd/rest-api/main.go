package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mars-60/mini-project1.git/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	//setup router
	router:= http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Welcome to rest-api"))
	})

	//setup server
	server:=http.Server{
		Addr: cfg.Address,
		Handler: router,
	}
    
	
	fmt.Println("server started ",slog.String("address",cfg.Address))
	fmt.Printf("server started %s",cfg.Address)

	done:=make(chan os.Signal,1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func(){
    err:=server.ListenAndServe()
	if err!=nil{
		log.Fatal("Failed to start server")
	}
	}()
	
    <-done

	slog.Info("Shutting down the server")

	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	if err:=server.Shutdown(ctx);err!=nil{
		slog.Error("Failed to shutdown server",slog.String("error",err.Error()))
	}

	slog.Info("server shutdown successfully")
}