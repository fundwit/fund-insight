package main

import (
	"fundinsight/pkg/bizerror"
	"fundinsight/pkg/series"
	"fundinsight/pkg/servehttp"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("service start")

	engine := gin.Default()
	engine.Use(bizerror.ErrorHandling())

	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "fund-insight")
	})

	servehttp.NewSeriesHandler(series.NewSeriesService()).RegisterToRoute(engine)

	err := engine.Run(":80")
	if err != nil {
		panic(err)
	}
}
