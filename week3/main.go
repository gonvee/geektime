package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	g, ctx        = errgroup.WithContext(context.Background())
	server        = &http.Server{Addr: ":8080"}
	debug         = &http.Server{Addr: "127.0.0.1:8081", Handler: http.DefaultServeMux}
	done          = make(chan os.Signal, 1)
	graceShutdown = errors.New("shutdown graceful")
)

func main() {
	// 信号量处理
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go signalHandle()

	// 应用服务
	g.Go(func() error {
		router := http.NewServeMux()
		router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprint(rw, "Hi")
		})

		server.Handler = router

		return server.ListenAndServe()
	})

	// debug服务
	g.Go(func() error {
		return debug.ListenAndServe()
	})

	if err := g.Wait(); err != nil && err == graceShutdown {
		log.Println(err)
	}
}

func signalHandle() error {
	<-done

	c, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := server.Shutdown(c); err != nil {
		return err
	}

	if err := debug.Shutdown(c); err != nil {
		return err
	}

	return graceShutdown
}
