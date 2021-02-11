package app

import (
	"context"
	"github.com/bal3000/BalStreamer.API/infrastructure"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	RabbitMQ infrastructure.RabbitMQ
	Echo     *echo.Echo
	Config   infrastructure.Configuration
}

func NewServer(rabbit infrastructure.RabbitMQ, e *echo.Echo, config infrastructure.Configuration) *Server {
	return &Server{RabbitMQ: rabbit, Echo: e, Config: config}
}

func (s *Server) Run() error {
	// Middleware
	s.Echo.Use(middleware.Logger())
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	// Routes
	s.SetRoutes()

	// Start server
	go func() {
		if err := s.Echo.Start(":8080"); err != nil {
			s.Echo.Logger.Info("Shutting down server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Echo.Shutdown(ctx); err != nil {
		s.Echo.Logger.Fatal(err)
		return err
	}

	return nil
}
