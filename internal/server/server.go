package server

import (
	"context"
	"errors"
	"fmt"
	"image"
	"net/http"
	"os"
	"time"

	"github.com/rdawson46/pic-conversion/internal/conversion"
	"github.com/rdawson46/pic-conversion/internal/storage"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var width int = 150

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
    db storage.Client
}

func NewServer(config ServerConfig) *Server {
    logger, _ := zap.NewDevelopment()

    var client storage.Client

    // TODO: command line arg 
    switch config.sType {
    case Test:
        client = storage.NewSampleDB()
    case Prod:
        // TODO: implement
        client = storage.NewMongo()
    default:
        fmt.Println("Not valid db type")
        os.Exit(1)
    }

    return &Server{
        config: config, 
        logger: logger,
        rateLimiter: rate.NewLimiter(
            rate.Limit(config.RateLimitReq),
            config.RateLimitBurst,
        ),
        db: client,
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

    // r.Body = http.MaxBytesReader(w, r.Body, s.config.MaxUploadSize)
    /*
    if err := r.ParseMultipartForm(s.config.MaxUploadSize); err != nil {
        s.logger.Info("File too large")
        http.Error(w, "File too large", http.StatusBadRequest)
        return
    }
    */

    file, _, err := r.FormFile("image")
    if err != nil {
        s.logger.Info("Error retrieving file")
        http.Error(w, "Error retrieving file", http.StatusBadRequest)
        return
    }

    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        http.Error(w, "Invalid image format", http.StatusBadRequest)
        return
    }

    hashId := storage.CalculateImageHash(img)


    var compressedAnsi string
    var cached bool
    compressedAnsi, err = s.db.CheckForImage(hashId)


    if err != nil {
        // image not stored
        if errors.Is(err, storage.NotStoredError){
            // produce
            cached = false
            //ansiArt := conversion.ConvertColorImage(img, width)
            ansiArt := conversion.OldConvertImage(img, width)
            s.logger.Info("Image made:", zap.Int("len(ansiArt)", len(ansiArt)))

            // compress
            // compressedAnsi, err = conversion.Compress(ansiArt)
            compressedAnsi = ansiArt

            // compression error check
            /*
            if err != nil {
                http.Error(w, "Error compressing image", http.StatusInternalServerError)
                return
            }
            */

        } else {
            // other error, write out and return
            http.Error(w, "Error converting image to ansi", http.StatusInternalServerError)
            return
        }
    } else {
        cached = true
    }

    w.Header().Set("Content-Type", "text/plain")

    // store image
    err = s.db.StoreImage(hashId, compressedAnsi)

    if err != nil {
        s.logger.Warn("Image could not be stored")
    }

    // change
    s.logger.Info("Image sent:", zap.Int("len(compressedAnsi)", len(compressedAnsi)), zap.Bool("cached", cached))
    fmt.Fprint(w, compressedAnsi)
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
        s.logger.Info("Starting Server", zap.Int("port", s.config.Port), zap.String("address", s.httpServer.Addr))

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
