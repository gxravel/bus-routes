package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/pkg/errors"
)

func ParseQueryInt64Slice(r *http.Request, field string) ([]int64, error) {
	q := r.URL.Query()
	params := q[field]

	if len(params) == 0 {
		return nil, nil
	}

	var vals []int64

	for _, p := range params {
		slice := strings.Split(p, ",")
		if vals == nil {
			vals = make([]int64, 0, len(slice))
		}

		for _, s := range slice {
			if s == "" {
				continue
			}
			val, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, errors.Errorf("can't parse %v to int", s)
			}
			vals = append(vals, val)
		}
	}

	return vals, nil
}

// ParseBusFilter returns filter to buses store depending on client's type.
func ParseBusFilter(r *http.Request) (*dataprovider.BusFilter, error) {
	userIDs, err := ParseQueryInt64Slice(r, "ids")
	if err != nil {
		return nil, err
	}

	return dataprovider.NewBusFilter().
		ByIDs(userIDs...), nil

}
