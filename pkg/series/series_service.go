package series

import (
	"context"
	"fundinsight/pkg/client"
	"fundinsight/pkg/domain"
	"strconv"
)

type SeriesQuery struct {
	Code string `json:"code" form:"code"`
}

type SeriesGetter interface {
	QuerySeries(q *SeriesQuery, c context.Context) ([]domain.Series, error)
}

type SeriesService struct {
	seriesGetter client.SinaSeriesGetter
}

func NewSeriesService() *SeriesService {
	return &SeriesService{seriesGetter: &client.SinaClient{}}
}

func (s *SeriesService) QuerySeries(q *SeriesQuery, c context.Context) ([]domain.Series, error) {
	rawPoints, err := s.seriesGetter.QueryCodeSeriesByDay(q.Code)
	if err != nil {
		return nil, err
	}

	seriesData := []interface{}{}
	for _, p := range rawPoints {
		// 1 open, 2 high, 3 low, 4 close, 5 volume
		data := make([]interface{}, 6)
		data[0] = p.Day

		d, err := strconv.ParseFloat(p.Open, 64)
		if err != nil {
			return nil, err
		}
		data[1] = d

		d, err = strconv.ParseFloat(p.High, 64)
		if err != nil {
			return nil, err
		}
		data[2] = d

		d, err = strconv.ParseFloat(p.Low, 64)
		if err != nil {
			return nil, err
		}
		data[3] = d

		d, err = strconv.ParseFloat(p.Close, 64)
		if err != nil {
			return nil, err
		}
		data[4] = d

		d, err = strconv.ParseFloat(p.Volume, 64)
		if err != nil {
			return nil, err
		}
		data[5] = d

		seriesData = append(seriesData, data)
	}

	series := domain.Series{Code: q.Code, Series: seriesData}
	return []domain.Series{series}, nil
}
