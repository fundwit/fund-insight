package servehttp_test

import (
	"context"
	"fmt"
	"fundinsight/pkg/bizerror"
	"fundinsight/pkg/domain"
	"fundinsight/pkg/series"
	"fundinsight/pkg/servehttp"
	"fundinsight/pkg/testinfra"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

func TestQueryWorkflowsRestAPI(t *testing.T) {
	RegisterTestingT(t)

	router := gin.Default()
	router.Use(bizerror.ErrorHandling())
	servehttp.NewSeriesHandler(&MockSeriesService{}).RegisterToRoute(router)

	t.Run("should return all series", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/series?code=sh600519", nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(MatchJSON(`[{"code":"sh600519","open_time":null,"series":[["2022-01-12 13:45:56",80.11,90.22,70.33,85.44,1234567]],"create_time":null,"last_update_time":null}]`))
	})

	t.Run("should be able to handle error when query workflows", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/series?code=err01235", nil)
		status, body, _ := testinfra.ExecuteRequest(req, router)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"common.internal_server_error","message":"a mocked error","data":null}`))
	})
}

type MockSeriesService struct {
}

func (s *MockSeriesService) QuerySeries(q *series.SeriesQuery, c context.Context) ([]domain.Series, error) {
	if strings.HasPrefix(q.Code, "err") {
		return nil, fmt.Errorf("a mocked error")
	}
	return []domain.Series{
		{
			Code: q.Code,
			Series: []interface{}{
				[]interface{}{"2022-01-12 13:45:56", 80.11, 90.22, 70.33, 85.44, 1234567},
			},
		},
	}, nil
}
