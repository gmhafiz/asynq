package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	redisLib "github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"

	"tasks/config"
	"tasks/internal/middleware"
	"tasks/third_party/database"
	"tasks/third_party/validate"
)

// Server struct holds all dependency reference required in this microservice.
// From the necessary httpServer, router, and asynq library to optional
// dependencies like a database.
type Server struct {
	cfg            *config.Config
	asynq          *asynq.Client
	srv            *asynq.Server
	redisClientOpt asynq.RedisClientOpt
	router         *chi.Mux
	httpServer     *http.Server

	// These are optional dependencies
	db        *sqlx.DB
	validator *validator.Validate
	redis     *redisLib.Client
}

// New is a constructor that returns a Server struct.
func New(version string) *Server {
	log.Printf("Starting API version: %s\n", version)
	return &Server{}
}

// Init initializes all dependencies for Producer APIs. Order of initialization
// is important. Configuration parsing is usually the first thing that happens
// so that all other dependencies can be configured correctly.
// GlobalMiddlewares happens before route registration and after router
// initialization.
//
// To register new tasks, figure out which domain it belongs to, create if it
// does not exist yet. initDomain() is usually done last.
func (s *Server) Init() {
	s.newConfig()
	s.newAsynq()
	s.newDatabase()
	s.newValidator()
	s.newRouter()
	s.setGlobalMiddleware()
	s.initDomains()
}

// InitConsumer initialize settings specific to Consumer APIs
func (s *Server) InitConsumer() {
	s.newConfig()
	s.newDatabase()
	s.newAsynq()
}

// Run runs the server. There is a graceful shutdown mechanism that listens
// to operating system signals.
func (s *Server) Run() {
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.cfg.Api.Host, s.cfg.Api.Port),
		Handler:        s.router,
		MaxHeaderBytes: 1 << 20,
	}

	// errs is an unbuffered channel that holds all errors of our go-routines.
	errs := make(chan error)

	// Launch go routine that starts the Producer API
	go func() {
		log.Printf("Serving at %s:%d\n", s.cfg.Api.Host, s.cfg.Api.Port)
		printAllRegisteredRoutes(s.router)
		errs <- s.httpServer.ListenAndServe()
		err := <-errs
		log.Fatal(err)
	}()

	// Launch another go routine that listens to operating signal for
	// termination.
	go func() {
		errs <- gracefulShutdown(context.Background(), s)
	}()

	err := <-errs
	log.Fatal(err)
}

func gracefulShutdown(ctx context.Context, s *Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

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

	ctx, shutdown := context.WithTimeout(ctx, 8*time.Second)
	defer shutdown()

	// Close any other opened resources here
	_ = s.DB().Close()

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) newConfig() {
	s.cfg = config.New()
}

func (s *Server) newDatabase() {
	if s.cfg.Database.Driver == "" {
		log.Fatal("please fill in database credentials in .env file")
	}
	s.db = database.NewSqlx(s.cfg)
	s.db.SetMaxOpenConns(s.cfg.Database.MaxConnectionPool)
}

func (s *Server) newValidator() {
	s.validator = validate.New()
}

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}

// setGlobalMiddleware is where all Handlers (controllers) will go through.
func (s *Server) setGlobalMiddleware() {
	s.router.Use(middleware.Json)
	s.router.Use(middleware.Cors)
	if s.cfg.Api.RequestLog {
		s.router.Use(chiMiddleware.Logger)
	}
	s.router.Use(middleware.Recovery)
}

func (s *Server) newAsynq() {
	// Defaults to connecting to Redis cluster if setting is set.
	// srv is used through its public getter AsyncServer().
	srv := &asynq.Server{} //nolint:staticcheck

	cfg := asynq.Config{
		// Specify how many concurrent workers to use. Ideally the number is
		// number of threads + spindle cont
		Concurrency: 12,
		// Optionally specify multiple queues with different priority.
		Queues: map[string]int{
			"critical": 7,
			"default":  4,
			"low":      1,
		},
	}

	if s.cfg.Redis.Addresses != "" {
		srv = asynq.NewServer(
			// (Option 1) for connecting to the Redis cluster (Default)
			asynq.RedisClusterClientOpt{Addrs: strings.Split(s.cfg.Redis.Addresses, ",")},
			cfg,
		)
	} else if s.cfg.Redis.Host != "" && s.cfg.Redis.Port != "" {
		srv = asynq.NewServer(
			// (Option 2) for single redis instance
			asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", s.cfg.Redis.Host, s.cfg.Redis.Port)},
			cfg,
		)
	} else {
		log.Fatal("must set redis credentials")
	}

	s.asynq = asynq.NewClient(s.redisClientOpt)
	s.srv = srv
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

func (s *Server) Config() *config.Config {
	return s.cfg
}

func (s *Server) DB() *sqlx.DB {
	return s.db
}

func (s *Server) Redis() *redisLib.Client {
	return s.redis
}

func (s *Server) AsynqServer() *asynq.Server {
	return s.srv
}
