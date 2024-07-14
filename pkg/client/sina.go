package client

import (
	"encoding/json"
	"io"
	"net/http"
)

type Point struct {
	Day    string `json:"day"` // e.g. '2001-12-10', '2024-02-01 14:20:00'
	Open   string `json:"open"`
	Close  string `json:"close"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Volume string `json:"volume"`

	// ma_pricesN
	// ma_volumeN
}

type Points []Point

type SinaSeriesGetter interface {
	QueryCodeSeriesByDay(code string) (Points, error)
}

type SinaClient struct {
}

func (c *SinaClient) QueryCodeSeriesByDay(code string) (Points, error) {
	resp, err := http.Get(`https://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?symbol=` + code + `&scale=240&ma=no&datalen=10000`)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	points := Points{}
	if err = json.Unmarshal(body, &points); err != nil {
		return nil, err
	}
	return points, err
}
