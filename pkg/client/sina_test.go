package client

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestQueryCodeSeriesByDay(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should be able to catch errors", func(t *testing.T) {

		sinaClient := SinaClient{}
		points, err := sinaClient.QueryCodeSeriesByDay("sh600519")
		Expect(err).To(BeNil())
		Expect(points).ToNot(BeNil())
	})
}
