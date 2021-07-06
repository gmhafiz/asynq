package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"

	"tasks/configs"
	"tasks/internal/middleware"
	db "tasks/third_party/database"
	"tasks/third_party/validate"
)

type Server struct {
	cfg        *configs.Configs
	db         *sqlx.DB
	redis      *asynq.Client
	router     *chi.Mux
	httpServer *http.Server
	validator  *validator.Validate
}

func New(version string) *Server {
	log.Printf("Starting API version: %s\n", version)
	return &Server{}
}

func (s *Server) Init() {
	s.newConfig()
	s.newDatabase()
	s.newRedis()
	s.newValidator()
	s.newRouter()
	s.setGlobalMiddleware()
	s.initDomains()
}

func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr:           s.cfg.Api.Host.String() + ":" + s.cfg.Api.Port,
		Handler:        s.router,
		ReadTimeout:    s.cfg.Api.ReadTimeout,
		WriteTimeout:   s.cfg.Api.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	go func() {
		log.Printf("Serving at %s:%s\n", s.cfg.Api.Host, s.cfg.Api.Port)
		printAllRegisteredRoutes(s.router)
		err := s.httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-quit

	timer := time.NewTimer(8 * time.Second)
	done := make(chan bool)
	go func() {
		log.Println("shutting down routine commenced...")
		<-timer.C
		done <- true
		log.Println("...server is shut down.")
	}()
	<-done

	ctx, shutdown := context.WithTimeout(context.Background(), s.cfg.Api.IdleTimeout*time.Second)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) newConfig() {
	s.cfg = configs.New()
}

func (s *Server) newDatabase() {
	if s.cfg.Database.Driver == "" {
		log.Fatal("please fill in database credentials in .env file")
	}
	s.db = db.NewSqlx(s.cfg)
	s.db.SetMaxOpenConns(s.cfg.Database.MaxConnectionPool)
}

func (s *Server) newValidator() {
	s.validator = validate.New()
}

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}

func (s *Server) setGlobalMiddleware() {
	s.router.Use(middleware.Json)
	s.router.Use(middleware.Cors)
	if s.cfg.Api.RequestLog {
		s.router.Use(chiMiddleware.Logger)
	}
	s.router.Use(middleware.Recovery)
}

func (s *Server) initDomains() {
	s.initHealth()
	s.initEmail()
}

func (s *Server) newRedis() {
	// for single redis instance
	s.redis = asynq.NewClient(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", s.cfg.Redis.Host, s.cfg.Redis.Port)})

	// for connecting to the Redis cluster
	//redis := asynq.NewClient(asynq.RedisClusterClientOpt{Addrs: []string{
	//	"0.0.0.0:7000", "0.0.0.0:7001", "0.0.0.0:7002",
	//}})
}

func printAllRegisteredRoutes(router *chi.Mux) {
	walkFunc := func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%-6s %s ", method, path)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Print(err)
	}
}

func (s *Server) Config() *configs.Configs {
	return s.cfg
}

func (s *Server) DB() *sqlx.DB {
	return s.db
}
