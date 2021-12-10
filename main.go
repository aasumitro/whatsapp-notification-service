package main

import (
	"fmt"
	"github.com/aasumitro/gowa/docs"
	httpHandlers "github.com/aasumitro/gowa/internal/delivery/http/handlers"
	"github.com/aasumitro/gowa/internal/delivery/http/middlewares"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"runtime"
)

func init() {
	if os.Getenv("SERVER_URL") == "" {
		exitF("SERVER_URL env is required")
	}
	if os.Getenv("SERVER_ENV") == "" {
		exitF("SERVER_ENV env is required")
	}
	if os.Getenv("SERVER_READ_TIMEOUT") == "" {
		exitF("SERVER_READ_TIMEOUT env is required")
	}
	if os.Getenv("JWT_SECRET_KEY") == "" {
		exitF("JWT_SECRET_KEY env is required")
	}
	if os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES") == "" {
		exitF("JWT_SECRET_KEY_EXPIRE_MINUTES env is required")
	}
	if os.Getenv("WHATSAPP_CLIENT_VERSION_MAJOR") == "" {
		exitF("WHATSAPP_CLIENT_VERSION_MAJOR env is required")
	}
	if os.Getenv("WHATSAPP_CLIENT_VERSION_MINOR") == "" {
		exitF("WHATSAPP_CLIENT_VERSION_MINOR env is required")
	}
	if os.Getenv("WHATSAPP_CLIENT_VERSION_BUILD") == "" {
		exitF("WHATSAPP_CLIENT_VERSION_BUILD env is required")
	}
	if os.Getenv("WHATSAPP_CLIENT_SESSION_PATH") == "" {
		exitF("WHATSAPP_CLIENT_SESSION_PATH env is required")
	}
}

// @title WhatsApp Web API with Golang
// @version 1.0
// @description Golang, Gin, Whatsapp Web API and Swagger.
// @termsOfService http://swagger.io/terms/
// @contact.name @developer.gowa
// @contact.email hello@aasumitro.id
// @BasePath /
func main() {
	// sets the maximum number of CPUs that can be executing
	runtime.GOMAXPROCS(runtime.NumCPU())
	// set server mode
	gin.SetMode(os.Getenv("SERVER_ENV"))
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	appEngine := gin.Default()
	// register custom middleware
	httpMiddleware := middlewares.InitHttpMiddleware()
	// use custom middleware
	appEngine.Use(httpMiddleware.CORS())
	// swagger info base path
	docs.SwaggerInfo.BasePath = appEngine.BasePath()
	// initialize http handler
	httpHandlers.NewHomeHttpHandler(appEngine)
	appEngine.Use(httpMiddleware.Auth())
	httpHandlers.NewWhatsappMessageHttpHandler(appEngine)
	// initialize ws handler
	// wsHandlers.NewWhatsappAuthenticationWsHandler()
	// Running the server
	log.Fatal(appEngine.Run(os.Getenv("SERVER_URL")))
}

func exitF(s string, args ...interface{}) {
	errorF(s, args...)
	os.Exit(1)
}

func errorF(s string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, s+"\n", args...)
	if err != nil {
		return
	}
}
