package api

import (
	"autoJoosik-market-data-fetcher/internal/autoSellerService"
	"autoJoosik-market-data-fetcher/internal/scheduler"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Api struct {
	Runner *scheduler.Runner
	Port   int
}
type BuyRequest struct {
	StkCd string  `json:"stkCd" binding:"required"`
	Qty   float64 `json:"qty" binding:"required"`
}

func (api *Api) Init() {
	// 운영이면 ReleaseMode 추천
	// gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "ajax"},
		AllowCredentials: false, // 쿠키/인증 필요 없으면 false로 둬도 됨
		MaxAge:           12 * time.Hour,
	}))

	// 헬스체크
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	// 헬스체크
	r.GET("/running", func(c *gin.Context) {
		isRunning := api.Runner.IsRunning()
		c.JSON(http.StatusOK, gin.H{"running": isRunning})
	})
	// 스케줄러 상태
	r.GET("/scheduler/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"running": api.Runner.IsRunning(),
			"enabled": api.Runner.IsEnabled(),
		})
	})
	// 스케줄러 중지
	r.POST("/scheduler/stop", func(c *gin.Context) {
		if err := api.Runner.Disable(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"running": api.Runner.IsRunning(),
			"enabled": api.Runner.IsEnabled(),
		})
	})
	// 스케줄러 시작
	r.POST("/scheduler/start", func(c *gin.Context) {
		if err := api.Runner.Enable(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"running": api.Runner.IsRunning(),
			"enabled": api.Runner.IsEnabled(),
		})
	})
	// 주식 구매
	r.POST("/market/buy", func(c *gin.Context) {
		var req BuyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := autoSellerService.Buy(req.StkCd, req.Qty); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"stkCd":  req.StkCd,
			"qty":    req.Qty,
		})
	})

	// 주식 판매
	r.POST("/market/sell", func(c *gin.Context) {
		var req BuyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := autoSellerService.Sell(req.StkCd, req.Qty); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"stkCd":  req.StkCd,
			"qty":    req.Qty,
		})
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
