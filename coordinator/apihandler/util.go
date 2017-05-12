package apihandler

import "net/url"

func queryToMap(vals url.Values) map[string]string {
	rMap := make(map[string]string)
	for k, v := range vals {
		rMap[k] = v[0]
	}
	return rMap
}
