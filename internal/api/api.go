package api

import (
	"autoJoosik-market-data-fetcher/internal/scheduler"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Api struct {
	Runner *scheduler.Runner
	Port   int
}

func (api *Api) Init() {
	// 운영이면 ReleaseMode 추천
	// gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger())

	// 헬스체크
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	// 헬스체크
	r.GET("/running", func(c *gin.Context) {
		isRunning := api.Runner.IsRunning()
		c.JSON(http.StatusOK, gin.H{"running": isRunning})
	})

	srv := &http.Server{
		Addr:              ":" + fmt.Sprint(api.Port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 서버 시작
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		_ = start // 필요하면 latency 찍기
	}
}

func parseID(s string) (int64, error) {
	// strconv.ParseInt 래퍼
	// (여기서 에러 메시지/범위 처리 등 커스텀 가능)
	return strconv.ParseInt(s, 10, 64)
}
