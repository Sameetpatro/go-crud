package main

import (
	"fmt"
	"github.com/Sameetpatro/go-crud/internal/config"
	"net/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"log/slog"
	"context"
	"time"
	"github.com/Sameetpatro/go-crud/internal/http/handlers/students"
)

func main(){
	fmt.Println("HELLO WORLD")
	//4 steps to create an API
	//load config
	//database load
	//setup router
	//set server

	cfg := config.MustLoad() // eithi config load heigala

	//aeithi server ra endpoints banauchi
	router := http.NewServeMux() //aeita multiplexer banauchi endpoints ra
	router.HandleFunc("POST /api/students", students.New()) // aeita gote endpoint

	//server bana hela 
	server :=http.Server{
		Addr: cfg.HTTPServer.Addr, //aeita config re thiba storage path ku server address re set karuchi
		Handler: router,
	}
	slog.Info("server started", slog.String("address", cfg.HTTPServer.Addr)) //aeita server start hela ki hela na ta ku print karuchi

	//graceful shutdown:- jetebele ame control + C press kariki shutdown karuchu gote server ku, we need to check if kichi process run karuchi ki nahi and if run karuchi then taku complete kariki end karibara achi so
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func(){
		err := server.ListenAndServe() //server ku read kara hauchi
		if err != nil{
			log.Fatal("Failed to start server")
		}
	} ()
	<- done 

	slog.Info("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)	
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Shutdown failed", slog.String("error", err.Error()))
 	}
	slog.Info("Server exited properly")
}