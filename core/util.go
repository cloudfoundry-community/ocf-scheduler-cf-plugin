package core

import (
	"net/url"
	"strings"
)

//func GetBearer() string {

//}

func GetScheduler(api string) string {
	u, err := url.Parse(api)
	if err != nil {
		panic("couldn't parse the api endpoint")
	}

	parts := strings.Split(u.Host, ".")
	parts[0] = "scheduler"
	u.Host = strings.Join(parts, ".")

	return u.String()
}
