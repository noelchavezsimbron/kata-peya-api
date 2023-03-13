package main

import (
	"fmt"
	stdHttp "net/http"
	"net/http/pprof"
	"os"
	"strconv"

	"kata-peya/config"
	"kata-peya/docs"
	"kata-peya/internal/http"
	"kata-peya/internal/http/handler"
	petRepository "kata-peya/internal/pet/repository"
	petUseCase "kata-peya/internal/pet/usecase"
	"kata-peya/internal/storage/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	conf = config.Get()
)

// @title Kata Peya API
// @version         1.0
// @description     This is a kata api server.
// @host localhost:8080
// @BasePath /api/kata-peya/v1
// @schemes http
func main() {
	tp, err := jaegerTracerProvider(conf.InstanceName)
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)

	var (
		server = newServer()
		router = server.Group(conf.Server.BasePath)
		db     = mysql.New(mysql.Config{
			Host:     conf.DB.Host,
			Port:     conf.DB.Port,
			Database: conf.DB.Database,
			User:     conf.DB.User,
			Password: conf.DB.Password,
		})

		petsRepo    = petRepository.NewPetMysqlRepository(db)
		petsUC      = petUseCase.NewUseCase(petsRepo)
		petsHandler = handler.NewPetsHandler(petsUC)
		api         = http.NewApi(petsHandler)
	)

	serveDocs(router)
	serveProfRoutes(router)
	api.Routes(router)

	if err := server.Start(fmt.Sprintf(":%d", conf.Server.Port)); err != nil {
		log.Fatal("server started error: ", err)
	}

}

func newServer() *echo.Echo {
	server := echo.New()
	server.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `"method":"${method}","uri":"${uri}","status":"${status}",` +
				`"latency_human":"${latency_human}", "error":"{${error}}"` + "\n",
		}),
		otelecho.Middleware(conf.InstanceName),
		middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize: 1 << 10, // 1 KB
			LogLevel:  4,       // Error
		}),
	)
	return server
}

func serveDocs(router *echo.Group) {
	docs.SwaggerInfo.InfoInstanceName = conf.InstanceName + "-openapi"
	swag.Register(docs.SwaggerInfo.InstanceName(), docs.SwaggerInfo)
	router.GET("/docs/*", echoSwagger.EchoWrapHandler(
		func(c *echoSwagger.Config) {
			c.InstanceName = conf.InstanceName + "-openapi"
		}),
	)
}

func serveProfRoutes(router *echo.Group) {
	pprofEnabled, _ := strconv.ParseBool(os.Getenv("PPROF_ENABLED"))
	if !pprofEnabled {
		return
	}

	debug := router.Group("/debug/pprof")
	debug.GET("/", echo.WrapHandler(stdHttp.HandlerFunc(pprof.Index)))
	debug.GET("/cmdline", echo.WrapHandler(stdHttp.HandlerFunc(pprof.Cmdline)))
	debug.GET("/profile", echo.WrapHandler(stdHttp.HandlerFunc(pprof.Profile)))
	debug.GET("/symbol", echo.WrapHandler(stdHttp.HandlerFunc(pprof.Symbol)))
	debug.POST("/symbol", echo.WrapHandler(stdHttp.HandlerFunc(pprof.Symbol)))
	debug.GET("/trace", echo.WrapHandler(stdHttp.HandlerFunc(pprof.Trace)))
	debug.GET("/allocs", echo.WrapHandler(pprof.Handler("allocs")))
	debug.GET("/block", echo.WrapHandler(pprof.Handler("block")))
	debug.GET("/goroutine", echo.WrapHandler(pprof.Handler("goroutine")))
	debug.GET("/heap", echo.WrapHandler(pprof.Handler("heap")))
	debug.GET("/mutex", echo.WrapHandler(pprof.Handler("mutex")))
	debug.GET("/threadcreate", echo.WrapHandler(pprof.Handler("threadcreate")))
}

func jaegerTracerProvider(service string) (*tracesdk.TracerProvider, error) {

	URL := conf.Otel.ExporterEndpoint

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(URL)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
			attribute.String("environment", conf.Environment),
			attribute.Int64("ID", 1),
		)),
	)
	return tp, nil
}
