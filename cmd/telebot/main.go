package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rasulov-emirlan/pukbot/config"
	"github.com/rasulov-emirlan/pukbot/internal/delivery/graphql"
	"github.com/rasulov-emirlan/pukbot/internal/delivery/http/server"
	"github.com/rasulov-emirlan/pukbot/internal/puk"
	"github.com/rasulov-emirlan/pukbot/pkg/db"
	"github.com/rasulov-emirlan/pukbot/pkg/fs"
	"github.com/rasulov-emirlan/pukbot/pkg/logger"
)

var isConfigFromFile = false

func main() {
	l := logger.NewLogger()
	if len(os.Args) > 1 {
		isConfigFromFile = true
	}
	cfg, err := config.NewConfig(isConfigFromFile, ".env")
	if err != nil {
		l.Errorf("error occured: %v", err)
		return
	}
	fileSystem, err := fs.NewFileSystem(cfg.GoogleFS)
	if err != nil {
		l.Errorf("error occured: %v;", err)
		return
	}

	pgxconn, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		l.Errorf("error occured: %v;", err)
		return
	}
	defer pgxconn.Close(context.Background())
	pukRepository := puk.NewRepository(fileSystem, pgxconn)
	pukService := puk.NewService(l, pukRepository)

	// bot, err := telegrambot.NewBot(cfg.BotToken, pukService)
	// if err != nil {
	// 	l.Errorf("error occured: %v", err)
	// 	return
	// }

	gqlhandler, gqlplayground := graphql.NewHandler(pukService)
	srvr, err := server.NewServer(pukService, ":"+cfg.ServerPort, time.Second*15, gqlhandler.ServeHTTP, gqlplayground)
	if err != nil {
		l.Errorf("error occured: %v", err)
		return
	}
	// go func() {
	// 	bot.Start()
	// 	defer bot.Close()
	// }()
	go func() {
		if err := srvr.Start(); err != nil && err != http.ErrServerClosed {
			l.Errorf("server fell due to: %v", err)
		}
		defer srvr.Close()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Gracefully closed")
}
