package domain

import "github.com/fundwit/go-commons/types"

type Series struct {
	Code     string          `json:"code"    binding:"required" gorm:"primary_key"`
	OpenTime types.Timestamp `json:"open_time" binding:"required"`

	// arrays of array: e.g. [[timestamp, open, close, high, low, volume],...]
	SeriesBlob     string          `json:"-"  sql:"type:LONGBLOB NOT NULL"`
	Series         []interface{}   `json:"series" sql:"-"`
	CreateTime     types.Timestamp `json:"create_time" sql:"type:DATETIME(6) NOT NULL"`
	LastUpdateTime types.Timestamp `json:"last_update_time" sql:"type:DATETIME(6) NOT NULL"`
}
