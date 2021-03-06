package main

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"github.com/aasumitro/gowa/docs"
	_ "github.com/aasumitro/gowa/internal/delivery"
	httpHandlers "github.com/aasumitro/gowa/internal/delivery/http/handlers"
	"github.com/aasumitro/gowa/internal/delivery/http/middlewares"
	"github.com/aasumitro/gowa/internal/domain/contracts"
	"github.com/aasumitro/gowa/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

// ae stand for App Engine
var ginEngine *gin.Engine

// wac stand for Whatsapp Client
var whatsappService contracts.WhatsappService

func init() {
	// sets the maximum number of CPUs that can be executing
	runtime.GOMAXPROCS(runtime.NumCPU())

	// validate environment variables are set
	validateEnvironment()

	// set server mode
	gin.SetMode(os.Getenv("SERVER_ENV"))

	// Create a new WhatsApp connection
	whatsappService = newWhatsappClient()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	ginEngine = gin.Default()

	// swagger info base path
	docs.SwaggerInfo.BasePath = ginEngine.BasePath()
}

// @title WhatsApp Web API with Golang
// @version 1.0
// @description Golang, Gin, Whatsapp Web API and Swagger.
// @termsOfService http://swagger.io/terms/

// @contact.name @developer.gowa
// @contact.url https://aasumitro.id/
// @contact.email hello@aasumitro.id

// @BasePath /api/v1

// @license.name  MIT
// @license.url   https://github.com/aasumitro/gowa/blob/master/LICENSE
func main() {
	// initialize home http handler
	httpHandlers.NewHomeHttpHandler(ginEngine)

	// create a new router group for the handler to register routes to and apply the middleware to it.
	// The middleware will be applied to all the routes registered in this group.
	v1Group := ginEngine.Group("/api/v1")
	{
		// initialize whatsapp http handler
		httpMiddleware := middlewares.InitHttpMiddleware()
		waGroup := v1Group.Group("/whatsapp").Use(httpMiddleware.CORS())
		{
			// whatsapp account handler
			httpHandlers.NewWhatsappAccountHttpHandler(waGroup, whatsappService)

			// whatsapp message handler
			httpHandlers.NewWhatsappMessageHttpHandler(waGroup, whatsappService)
		}
	}

	// Running the server
	log.Fatal(ginEngine.Run(os.Getenv("SERVER_URL")))
}

func validateEnvironment() {
	if os.Getenv("SERVER_NAME") == "" {
		exitF("SERVER_NAME env is required")
	}
	if os.Getenv("SERVER_DESCRIPTION") == "" {
		exitF("SERVER_DESCRIPTION env is required")
	}
	if os.Getenv("SERVER_URL") == "" {
		exitF("SERVER_URL env is required")
	}
	if os.Getenv("SERVER_ENV") == "" {
		exitF("SERVER_ENV env is required")
	}
	if os.Getenv("SERVER_READ_TIMEOUT") == "" {
		exitF("SERVER_READ_TIMEOUT env is required")
	}
	if os.Getenv("SERVER_UPLOAD_LIMIT") == "" {
		exitF("SERVER_UPLOAD_LIMIT env is required")
	}
	if os.Getenv("WAC_MAJOR_VERSION") == "" {
		exitF("WAC_MAJOR_VERSION env is required")
	}
	if os.Getenv("WAC_MINOR_VERSION") == "" {
		exitF("WAC_MINOR_VERSION env is required")
	}
	if os.Getenv("WAC_BUILD_VERSION") == "" {
		exitF("WAC_BUILD_VERSION env is required")
	}
	if os.Getenv("WAC_SESSION_PATH") == "" {
		exitF("WAC_SESSION_PATH env is required")
	}
	if os.Getenv("WAC_UPLOAD_PATH") == "" {
		exitF("WAC_UPLOAD_PATH env is required")
	}
}

func newWhatsappClient() contracts.WhatsappService {
	wac, err := whatsapp.NewConnWithOptions(&whatsapp.Options{
		Timeout:         20 * time.Second,
		ShortClientName: os.Getenv("SERVER_NAME"),
		LongClientName:  os.Getenv("SERVER_DESCRIPTION"),
	})
	if err != nil {
		exitF("WhatsApp connection error: ", err)
	}

	waClientVerMajInt, err := strconv.Atoi(
		os.Getenv("WAC_MAJOR_VERSION"))
	if err != nil {
		exitF("Error conversion", err)
	}

	waClientVerMinInt, err := strconv.Atoi(
		os.Getenv("WAC_MINOR_VERSION"))
	if err != nil {
		exitF("Error conversion", err)
	}

	waClientVerBuildInt, err := strconv.Atoi(
		os.Getenv("WAC_BUILD_VERSION"))
	if err != nil {
		exitF("Error conversion", err)
	}

	wac.SetClientVersion(
		waClientVerMajInt,
		waClientVerMinInt,
		waClientVerBuildInt,
	)

	whatsappService := services.WhatsappService{Conn: wac}

	//Restore session if exists
	err = whatsappService.RestoreSession()
	if err != nil {
		exitF("Error restoring whatsapp session. ", err)
	}

	return &whatsappService
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
