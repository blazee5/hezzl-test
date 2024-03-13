package handler

import (
	"context"
	"github.com/blazee5/hezzl-test/internal/config"
	"github.com/blazee5/hezzl-test/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	log        *zap.SugaredLogger
	cfg        *config.Config
	httpServer *http.Server
	db         *pgxpool.Pool
	rdb        *redis.Client
	nc         *nats.Conn
	services   *service.Service
}

func NewServer(log *zap.SugaredLogger, cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client, nc *nats.Conn, services *service.Service) *Server {
	return &Server{log: log, cfg: cfg, db: db, rdb: rdb, nc: nc, services: services}
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         s.cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) InitRoutes() *gin.Engine {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}))

	api := r.Group("/api/v1")

	good := api.Group("/good")
	{
		good.POST("/create", s.CreateGood)
		good.PATCH("/update", s.UpdateGood)
		good.DELETE("/remove", s.DeleteGood)
		good.PATCH("/reprioritize", s.ReprioritizeGood)
	}

	api.GET("/goods/list", s.GetGoods)

	return r
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
