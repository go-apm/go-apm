package main

import (
	"github.com/facebookgo/grace/gracehttp"
	_ "github.com/go-apm/go-apm/conf"
	"github.com/go-apm/go-apm/endpoint"
	_ "github.com/go-apm/go-apm/port/pprof/http"
	_ "github.com/go-apm/go-apm/port/store"
	"github.com/go-apm/go-apm/util"
	"github.com/go-apm/go-apm/util/xhttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"github.com/uber-go/zap"
	"net/http"
	_ "net/http/pprof"
)

var boostLogger = zap.NewJSON()

var a = []string{}

func main() {

	go util.WithRecover(listenProf)

	e := echo.New()
	e.Use(
		xhttp.RequestID,
		xhttp.RequestLogger,
		middleware.Recover(),
		middleware.StaticWithConfig(middleware.StaticConfig{Root: "static", Browse: false}),
	)
	go func() {
		for {
			for i := 0; i < 1000; i++ {
				a = append(a, "1212121212")
			}
			a = []string{}
		}
	}()
	httpRouter(e)
	std := standard.New(viper.GetString("httpPort"))
	std.SetHandler(e)
	boostLogger.Info("Start listen http", zap.String("port", viper.GetString("httpPort")))
	err := gracehttp.Serve(std.Server)
	if err != nil {
		boostLogger.Fatal("Server Errror Occured", zap.Error(err))
	}
}

func listenProf() {
	boostLogger.Info("Start listen pprof", zap.String("port", viper.GetString("pprofPort")))
	http.ListenAndServe(viper.GetString("pprofPort"), nil)
}

func httpRouter(e *echo.Echo) {
	api := e.Group("/api")
	{
		heap := api.Group("/heaps")
		{
			heap.GET("/", endpoint.ListHeaps)
			heap.PUT("/", endpoint.NewHeap)
			heap.GET("/:id", endpoint.ViewHeap)
		}
	}
}
