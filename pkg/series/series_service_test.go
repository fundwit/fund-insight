package series

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
)

func TestQuerySeries(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should be able to catch errors", func(t *testing.T) {
		seriesSvc := NewSeriesService()
		points, err := seriesSvc.QuerySeries(&SeriesQuery{Code: "sh600519"}, context.Background())
		Expect(err).To(BeNil())
		Expect(points).ToNot(BeNil())
	})
}
