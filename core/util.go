package core

import (
	"net/url"
	"strings"
)

//func GetBearer() string {

//}

func GetScheduler(api string) string {
	u := url.Parse(api)

	parts := strings.Split(u.Host, ".")
	parts[0] = "scheduler"
	u.Host = strings.Join(parts, ".")

	return u.String()
}
