package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/pkg/errors"
)

func parseQueryUint64(r *http.Request, field string) (uint64, error) {
	value, err := ParseQueryParam(r, field)
	if err != nil {
		return 0, err
	}

	if value == "" {
		return 0, nil
	}

	iv, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, errors.Errorf("%v is not a uint", value)
	}

	return iv, nil
}

func ParseQueryInt(r *http.Request, field string) (int, error) {
	value, err := ParseQueryParam(r, field)
	if err != nil {
		return 0, err
	}

	if value == "" {
		return 0, nil
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.Errorf("%v is not a int", value)
	}

	return int(i), nil
}

// ParseQueryParam parses query params for specific field.
func ParseQueryParam(r *http.Request, field string) (string, error) {
	q := r.URL.Query()

	param := q.Get(field)
	if param == "" {
		return "", nil
	}

	return param, nil
}

// ParseQueryParams parses query params for specific field.
func ParseQueryParams(r *http.Request, field string) ([]string, error) {
	q := r.URL.Query()
	params := q[field]

	if len(params) == 0 {
		return nil, nil
	}

	return params, nil
}

// ParsePaginator parses paginator from request.
func ParsePaginator(r *http.Request) (*dataprovider.Paginator, error) {
	limit, err := parseQueryUint64(r, "limit")
	if err != nil {
		return nil, err
	}

	offset, err := parseQueryUint64(r, "offset")
	if err != nil {
		return nil, err
	}

	return dataprovider.NewPaginator(offset, limit), nil
}

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

func ParseQueryIntSlice(r *http.Request, field string) ([]int, error) {
	i64, err := ParseQueryInt64Slice(r, field)
	if err != nil {
		return nil, err
	}
	var result = make([]int, 0, len(i64))
	for _, i := range i64 {
		result = append(result, int(i))
	}

	return result, nil
}

// ParseBusFilter returns filter to buses store depending on client's type.
func ParseBusFilter(r *http.Request) (*dataprovider.BusFilter, error) {
	ids, err := ParseQueryInt64Slice(r, "ids")
	if err != nil {
		return nil, err
	}

	cities, err := ParseQueryParams(r, "cities")
	if err != nil {
		return nil, err
	}

	nums, err := ParseQueryParams(r, "nums")
	if err != nil {
		return nil, err
	}

	return dataprovider.NewBusFilter().
		ByIDs(ids...).ByCities(cities...).ByNums(nums...), nil

}

func ParseCityFilter(r *http.Request) (*dataprovider.CityFilter, error) {
	ids, err := ParseQueryIntSlice(r, "ids")
	if err != nil {
		return nil, err
	}

	names, err := ParseQueryParams(r, "names")
	if err != nil {
		return nil, err
	}

	return dataprovider.NewCityFilter().
		ByIDs(ids...).ByNames(names...), nil

}

func ParseDeleteCityFilter(r *http.Request) (*dataprovider.CityFilter, error) {
	id, err := ParseQueryInt(r, "id")
	if err != nil {
		return nil, err
	}

	name, err := ParseQueryParam(r, "name")
	if err != nil {
		return nil, err
	}

	return dataprovider.NewCityFilter().
		ByIDs(id).ByNames(name), nil
}
