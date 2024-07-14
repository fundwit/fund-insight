package servehttp

import (
	"context"
	"fundinsight/pkg/series"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type RouteHandler interface {
	RegisterToRoute(r *gin.Engine, middleWares ...gin.HandlerFunc)
}

type seriesHandler struct {
	validator    *validator.Validate
	seriesGetter series.SeriesGetter
}

func NewSeriesHandler(seriesGetter series.SeriesGetter) *seriesHandler {
	return &seriesHandler{
		validator:    validator.New(),
		seriesGetter: seriesGetter,
	}
}

func (h *seriesHandler) RegisterToRoute(r *gin.Engine, middleWares ...gin.HandlerFunc) {
	g := r.Group("/v1/series", middleWares...)
	g.GET("", h.QuerySeries)
}

func (h *seriesHandler) QuerySeries(c *gin.Context) {
	query := series.SeriesQuery{}
	_ = c.MustBindWith(&query, binding.Query)

	result, err := h.seriesGetter.QuerySeries(&query, context.Background())
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, result)
}
