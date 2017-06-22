package leveldb

import (
	"errors"
	"strconv"
	"time"
)

const (
	QueryDateFormat = "2006-01-02 15:04"
)

func parseDate(date string) (int, error) {

	//attempt to parseDate
	t, err := time.Parse(QueryDateFormat, date)
	if err == nil {
		return int(t.Unix()), nil
	}

	di, err := strconv.Atoi(date)
	if err == nil {
		return di, nil
	}

	return -1, errors.New("Invalid time provided. Must be unix timestamp or in format yyyy-mm-dd hh:mm")
}

func stringInSlice(in []string, find string) bool {
	for _, s := range in {
		if s == find {
			return true
		}
	}
	return false
}

func intersectSlice(one, two []string) []string {
	intersec := make([]string, 0)
	for _, v := range one {
		if stringInSlice(two, v) {
			intersec = append(intersec, v)
		}
	}
	return intersec
}
