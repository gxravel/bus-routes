package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/pkg/errors"
)

type intSize uint8

const (
	i intSize = iota
	i64
	ui
	ui64
)

func parseQueryInt(r *http.Request, field string, is intSize) (interface{}, error) {
	value, err := ParseQueryParam(r, field)
	if err != nil {
		return 0, err
	}

	if value == "" {
		return 0, nil
	}

	iv, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.Errorf("%v is not an int", value)
	}

	switch is {
	case i64:
		return int64(iv), nil
	case ui:
		return uint(iv), nil
	case ui64:
		return uint64(iv), nil
	default:
		return int(iv), nil
	}
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
	limit, err := parseQueryInt(r, "limit", ui64)
	if err != nil {
		return nil, err
	}

	offset, err := parseQueryInt(r, "offset", ui64)
	if err != nil {
		return nil, err
	}

	return dataprovider.NewPaginator(offset.(uint64), limit.(uint64)), nil
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
	id, err := parseQueryInt(r, "id", i)
	if err != nil {
		return nil, err
	}

	name, err := ParseQueryParam(r, "name")
	if err != nil {
		return nil, err
	}

	return dataprovider.NewCityFilter().
		ByIDs(id.(int)).ByNames(name), nil
}

func ParseStopFilter(r *http.Request) (*dataprovider.StopFilter, error) {
	ids, err := ParseQueryInt64Slice(r, "ids")
	if err != nil {
		return nil, err
	}

	cities, err := ParseQueryParams(r, "cities")
	if err != nil {
		return nil, err
	}

	addresses, err := ParseQueryParams(r, "addresses")
	if err != nil {
		return nil, err
	}

	return dataprovider.NewStopFilter().
		ByIDs(ids...).ByCities(cities...).ByAddresses(addresses...), nil

}

func ParseDeleteStopFilter(r *http.Request) (*dataprovider.StopFilter, error) {
	id, err := parseQueryInt(r, "id", i64)
	if err != nil {
		return nil, err
	}

	address, err := ParseQueryParam(r, "address")
	if err != nil {
		return nil, err
	}

	return dataprovider.NewStopFilter().
		ByIDs(id.(int64)).ByAddresses(address), nil
}
