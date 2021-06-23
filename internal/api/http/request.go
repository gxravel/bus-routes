package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/pkg/errors"
)

func parseQueryInt64(r *http.Request, field string) (int64, error) {
	value, err := ParseQueryParam(r, field)
	if err != nil {
		return 0, err
	}

	if value == "" {
		return 0, nil
	}

	iv, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.Errorf("%v it not an int", value)
	}
	return iv, nil
}

func parseQueryInt(r *http.Request, field string) (int, error) {
	result, err := parseQueryInt64(r, field)
	if err != nil {
		return 0, err
	}
	return int(result), nil
}
func parseQueryInt8(r *http.Request, field string) (int8, error) {
	result, err := parseQueryInt64(r, field)
	if err != nil {
		return 0, err
	}
	return int8(result), nil
}
func parseQueryUint64(r *http.Request, field string) (uint64, error) {
	result, err := parseQueryInt64(r, field)
	if err != nil {
		return 0, err
	}
	return uint64(result), nil
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
	vals, err := ParseQueryInt64Slice(r, field)
	if err != nil {
		return nil, err
	}
	var result = make([]int, 0, len(vals))
	for _, val := range vals {
		result = append(result, int(val))
	}
	return result, nil
}

func ParseQueryUint8Slice(r *http.Request, field string) ([]int8, error) {
	vals, err := ParseQueryInt64Slice(r, field)
	if err != nil {
		return nil, err
	}
	var result = make([]int8, 0, len(vals))
	for _, val := range vals {
		result = append(result, int8(val))
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
	id, err := parseQueryInt(r, "id")
	if err != nil {
		return nil, err
	}

	name, err := ParseQueryParam(r, "name")
	if err != nil {
		return nil, err
	}

	filter := dataprovider.NewCityFilter()
	if id != 0 {
		filter = filter.ByIDs(id)
	}
	if name != "" {
		filter = filter.ByNames(name)
	}

	return filter, nil
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
	id, err := parseQueryInt64(r, "id")
	if err != nil {
		return nil, err
	}

	address, err := ParseQueryParam(r, "address")
	if err != nil {
		return nil, err
	}

	filter := dataprovider.NewStopFilter()
	if id != 0 {
		filter = filter.ByIDs(id)
	}
	if address != "" {
		filter = filter.ByAddresses(address)
	}
	return filter, nil
}

func ParseRouteFilter(r *http.Request) (*dataprovider.RouteFilter, error) {
	busIDs, err := ParseQueryInt64Slice(r, "bus_ids")
	if err != nil {
		return nil, err
	}

	stopIDs, err := ParseQueryInt64Slice(r, "stop_ids")
	if err != nil {
		return nil, err
	}

	steps, err := ParseQueryUint8Slice(r, "steps")
	if err != nil {
		return nil, err
	}

	return dataprovider.NewRouteFilter().
		ByBusIDs(busIDs...).ByStopIDs(stopIDs...).BySteps(steps...), nil

}

func ParseDeleteRouteFilter(r *http.Request) (*dataprovider.RouteFilter, error) {
	busID, err := parseQueryInt64(r, "bus_id")
	if err != nil {
		return nil, err
	}

	stopID, err := parseQueryInt64(r, "stop_id")
	if err != nil {
		return nil, err
	}

	step, err := parseQueryInt8(r, "step")
	if err != nil {
		return nil, err
	}

	fmt.Printf("step: %v\n", step)

	filter := dataprovider.NewRouteFilter()
	if busID != 0 {
		filter = filter.ByBusIDs(busID)
	}
	if stopID != 0 {
		filter = filter.ByStopIDs(stopID)
	}
	if step != 0 {
		filter = filter.BySteps(step)
	}
	fmt.Printf("step filter: %v\n", filter.Steps)
	return filter, nil
}
