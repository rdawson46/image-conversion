package server

import (
    "github.com/rdawson46/pic-conversion/internal/conversion"
	"context"
	"fmt"
	"image"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type ServerType int

const (
    Test ServerType = iota
    Prod
)

type ServerConfig struct {
    Port int
    MaxUploadSize int64
    RateLimitReq float64
    RateLimitBurst int
    ShutdownTimeout time.Duration
    sType ServerType
}

func NewConfig(port int, maxUpload int64, rateLimitReq float64, rateLimitBurst int, shutdown time.Duration, sType ServerType) ServerConfig {
    return ServerConfig{
        Port: port,
        MaxUploadSize: maxUpload << 20,
        RateLimitReq: rateLimitReq,
        RateLimitBurst: rateLimitBurst,
        ShutdownTimeout: shutdown,
        sType: sType,
    }
}

type Server struct {
    config ServerConfig
    httpServer *http.Server
    logger *zap.Logger
    rateLimiter *rate.Limiter
}

func NewServer(config ServerConfig) *Server {
    logger, _ := zap.NewProduction()

    // TODO: connect to db by checking config

    return &Server{
        config: config, 
        logger: logger,
        rateLimiter: rate.NewLimiter(
            rate.Limit(config.RateLimitReq),
            config.RateLimitBurst,
        ),
    }
}

func (s *Server) rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !s.rateLimiter.Allow() {
            http.Error(w, "Rate limit Exceeded", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    }
}

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowd", http.StatusMethodNotAllowed)
        return
    }

    r.Body = http.MaxBytesReader(w, r.Body, s.config.MaxUploadSize)
    if err := r.ParseMultipartForm(s.config.MaxUploadSize); err != nil {
        http.Error(w, "File too large", http.StatusBadRequest)
        return
    }

    file, _, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Error retrieving file", http.StatusBadRequest)
        return
    }

    defer file.Close()

    // TODO: fix this
    img, _, err := image.Decode(file)
    if err != nil {
        http.Error(w, "Invalid image format", http.StatusBadRequest)
        return
    }

    // HACK: hard coding 100 temp
    ansiArt := conversion.ConvertImage(img, 100)

    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprint(w, ansiArt)
}

func (s *Server) Start() error {
    mux := http.NewServeMux()


    // TODO: need to set up routing for UI at /
    // mux.HandleFunc("/", s.index())

    mux.HandleFunc("/upload", s.rateLimitMiddleware(s.uploadHandler))

    s.httpServer = &http.Server{
        Addr: fmt.Sprintf(":%d", s.config.Port),
        Handler: mux,
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    go func() {
        s.logger.Info("Starting Server", zap.Int("port", s.config.Port))

        if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
            s.logger.Error("Server Failed", zap.Error(err))
        }
    }()

    return nil
}


func (s *Server) Shutdown() error {
    ctx, cancel := context.WithTimeout(
        context.Background(),
        s.config.ShutdownTimeout,
    )

    defer cancel()

    if err := s.httpServer.Shutdown(ctx); err != nil {
        s.logger.Error("Server shutdown", zap.Error(err))
        return err
    }

    s.logger.Info("Server Shutdown gracefully")
    return nil
}
