package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/zeiss/pkg/server"

	"github.com/gofiber/fiber/v2"
)

type srv struct{}

func (s *srv) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		app := fiber.New()
		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		})

		run(func() error {
			ticker := time.NewTimer(100 * time.Second)
			defer ticker.Stop()

			<-ticker.C

			return server.NewServerError(errors.New("timeout"))
		})

		ready()

		err := app.Listen(":3000")
		if err != nil {
			return err
		}

		return nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, _ := server.WithContext(ctx)
	s.SetLimit(3)

	s.Listen(&srv{}, true)
	s.Listen(server.NewDebug(server.WithPprof()), true)

	log.Printf("starting %s", server.Service.Name())
	serverErr := &server.ServerError{}
	if err := s.Wait(); errors.As(err, &serverErr) {
		log.Print(err)
		os.Exit(1)
	}
}
