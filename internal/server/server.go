package server

import (
	"errors"
	"strings"
	"time"

	_ "github.com/rasulov-emirlan/pukbot/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rasulov-emirlan/pukbot/internal/puk"
	"github.com/rasulov-emirlan/pukbot/pkg/logger"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Puk Bot Server
// @version 1.0
// @description This is THE web api of THE PUK bot.

// @contact.name Rasulov Emirlan
// @contact.email rasulov-emirlan@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @Accept json
type Server struct {
	r          *echo.Echo
	l          logger.Logger
	pukService puk.Service
	port       string
	timeout    time.Duration
}

func NewServer(pukService puk.Service, port string, timeout time.Duration) (*Server, error) {
	router := echo.New()
	if timeout > time.Second*60 {
		return nil, errors.New("timout is too big")
	}
	router.Server.ReadTimeout = timeout
	router.Server.WriteTimeout = timeout
	return &Server{
		r:          router,
		pukService: pukService,
		port:       port,
		timeout:    timeout,
	}, nil
}

func (s *Server) Start() error {
	s.plugRoutes()
	return s.r.Start(s.port)
}

func (s *Server) Close() error {
	return s.r.Close()
}

func (s *Server) plugRoutes() {
	s.r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderAccessControlAllowOrigin},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	s.r.Use(middleware.RemoveTrailingSlash())
	s.r.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))
	s.r.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	s.r.GET("/api/puks", pukList(s.pukService, s.l))
	s.r.GET("/swagger/*", echoSwagger.WrapHandler)
}
