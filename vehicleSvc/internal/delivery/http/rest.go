package rest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	//_ "app/docs"

	"app/config"
	"app/internal/delivery/http/handler"
	"app/internal/repository"
	"app/internal/usecase"
	"github.com/gin-gonic/gin"
)

// NewRest
// @title Vehicle Service API
// @version 1.0
// @description Vehicle Service API
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
func NewRest(cfg *config.MainConfig) *gin.Engine {
	if cfg.Environment != "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	app.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	})

	return app
}

func Serve(app *gin.Engine, cfg *config.MainConfig) (err error) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServicePort),
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error: %s\n", err)
		}
	}()

	log.Println("-------------------------------------------")
	log.Println("server started")
	log.Printf("running on port %d\n", cfg.ServicePort)
	log.Println("-------------------------------------------")

	return gracefulShutdown(server)
}

func gracefulShutdown(srv *http.Server) error {
	done := make(chan os.Signal)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Println("Shutting down server...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Error while shutting down Server. Initiating force shutdown...")
		return err
	}

	log.Println("Server exiting...")

	return nil
}

type Route struct {
	Config     *config.MainConfig
	App        *gin.Engine
	UseCase    usecase.Usecase
	Repository repository.RepositoryImpl
}

func (r *Route) RegisterRoutes() {
	r.App.Use(gin.Recovery())

	//auth := middleware.NewAuthMiddleware(r.Config, r.Repository)

	handler := handler.NewHandler(r.UseCase)

	r.App.GET("/vehicles/:vehicle_id/location", handler.GetVehicleLocationLatest)
	r.App.GET("vehicles/:vehicle_id/history", handler.GetVehicleHistory)
}
