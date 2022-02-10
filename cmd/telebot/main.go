package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/rasulov-emirlan/pukbot/config"
	"github.com/rasulov-emirlan/pukbot/internal/puk"
	"github.com/rasulov-emirlan/pukbot/internal/server"
	"github.com/rasulov-emirlan/pukbot/internal/telegrambot"
	"github.com/rasulov-emirlan/pukbot/pkg/db"
	"github.com/rasulov-emirlan/pukbot/pkg/fs"
	"github.com/rasulov-emirlan/pukbot/pkg/logger"
)

func main() {
	l := logger.NewLogger()
	cfg, err := config.NewConfig(false)
	if err != nil {
		l.Infof("error occured: %v", err)
		return
	}
	fileSystem, err := fs.NewFileSystem(cfg.GoogleFS)
	if err != nil {
		l.Infof("error occured: %v;", err)
		return
	}

	if err = godotenv.Load(".env"); err != nil {
		l.Infof("error occured: %v;", err)
		return
	}
	pgxconn, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		l.Infof("error occured: %v;", err)
		return
	}
	defer pgxconn.Close(context.Background())
	pukRepository := puk.NewRepository(fileSystem, pgxconn)
	pukService := puk.NewService(l, pukRepository)

	bot, err := telegrambot.NewBot(cfg.BotToken, pukService)
	if err != nil {
		l.Infof("error occured: %v", err)
		return
	}

	srvr, err := server.NewServer(pukService, ":"+cfg.ServerPort, time.Second*15)
	if err != nil {
		l.Infof("error occured: %v", err)
		return
	}
	go func() {
		bot.Start()
		defer bot.Close()
	}()
	go func() {
		if err := srvr.Start(); err != nil && err != http.ErrServerClosed {
			l.Infof("server fell due to: %v", err)
		}
		defer srvr.Close()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Gracefully closed")
}
