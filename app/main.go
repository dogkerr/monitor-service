package main

import (
	"context"
	"dogker/lintang/monitor-service/app/di"
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/internal/rest/middleware"
	"dogker/lintang/monitor-service/pkg/httpserver"
	"dogker/lintang/monitor-service/pkg/postgres"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gin-gonic/gin"
)

func init() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	//
	if err := middleware.InitLogger(cfg.LogConfig); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

}

func main() {
	cfg, err := config.NewConfig()

	// init logger

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	handler := gin.New()
	httpServer := httpserver.New(handler, httpserver.Port("5000"))

	// init app
	wireApp := di.InitApp(cfg, handler)

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"postgres": func(ctx context.Context) error {
			return postgres.ClosePostgres(wireApp.PG.Pool)
		},
		"http-server": func(ctx context.Context) error {
			return httpServer.Shutdown()
		},
		"rmq": func(ctx context.Context) error {
			return wireApp.RMQ.Close()
		},
	})

	<-wait

	// // Waiting signal
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// select {
	// case s := <-interrupt:
	// 	zap.L().Fatal("app - Run - signal: " + s.String())
	// case err := <-httpServer.Notify():
	// 	zap.L().Fatal(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	// }

	// // Shutdown
	// err = httpServer.Shutdown()
	// if err != nil {
	// 	zap.L().Fatal(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	// }
}

type operation func(ctx context.Context) error

// gracefulShutdown waits for termination syscalls and doing clean up operations after received it
func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Println("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
