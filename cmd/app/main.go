package main

import (
	"context"
	"dogker/lintang/monitor-service/cmd/di"
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
	"go.uber.org/zap"

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
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// init app
	wireApp := di.InitApp(cfg, handler)

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"postgres": func(ctx context.Context) error {
			return postgres.ClosePostgres(wireApp.PG.Pool)
		},
		"rmq": func(ctx context.Context) error {
			return wireApp.RMQ.Close()
		},
		"grpc": func(ctx context.Context) error {
			wireApp.GRPCServer.GracefulStop()
			return nil
		},
	})

	select {
	case _ = <-wait:
		fmt.Println("")
	case err = <-httpServer.Notify():

		gracefulShutdownHttpNotify(context.Background(), 2*time.Second, map[string]operation{
			"postgres": func(ctx context.Context) error {
				return postgres.ClosePostgres(wireApp.PG.Pool)
			},

			"rmq": func(ctx context.Context) error {
				return wireApp.RMQ.Close()
			},
			"grpc": func(ctx context.Context) error {
				wireApp.GRPCServer.GracefulStop()
				return nil
			},
		})
		time.Sleep(1 * time.Second)
		zap.L().Info(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())

	}

	httpServer.Shutdown()

}

type operation func(ctx context.Context) error

func gracefulShutdownHttpNotify(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {

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

// gracefulShutdown waits for termination syscalls and doing clean up operations after received it
func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
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
