package middleware

import (
	"github.com/curt-labs/GoAPI/helpers/slack"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
	"github.com/segmentio/analytics-go"
	"net/http"
	"time"
)

func Meddler() martini.Handler {
	return func(res http.ResponseWriter, r *http.Request, c martini.Context) {
		start := time.Now()

		authed := checkAuth(r)
		if !authed {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		c.Next()
		go logRequest(r, time.Since(start))
	}
}

func checkAuth(r *http.Request) bool {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}
	if key == "" {
		key = r.Header.Get("key")
	}
	if key == "" {
		return false
	}

	user, err := customer.GetCustomerUserFromKey(key)
	if err != nil || user.Id == "" {
		return false
	}

	go user.LogApiRequest(r)

	return true
}

func logRequest(r *http.Request, reqTime time.Duration) {
	client := analytics.New("oactr73lbg")

	key := r.Header.Get("key")
	if key == "" {
		vals := r.URL.Query()
		key = vals.Get("key")
	}
	if key == "" {
		key = r.FormValue("key")
	}

	err := client.Track(map[string]interface{}{
		"event":       r.URL.String(),
		"userId":      key,
		"method":      r.Method,
		"header":      r.Header,
		"query":       r.URL.Query().Encode(),
		"referer":     r.Referer(),
		"userAgent":   r.UserAgent(),
		"form":        r.Form,
		"requestTime": int64((reqTime.Nanoseconds() * 1000) * 1000),
	})
	if err != nil {
		m := slack.Message{
			Channel:  "debugging",
			Username: "GoAPI",
			Text:     err.Error(),
		}
		m.Send()
	}
}
