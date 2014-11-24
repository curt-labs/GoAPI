package faq_controller

//reference
//https://groups.google.com/forum/#!topic/golang-nuts/DARY7HY-pbY
//http://blog.wercker.com/2014/02/06/RethinkDB-Gingko-Martini-Golang.html
//http://golang.org/pkg/net/http/httptest/#example_ResponseRecorder
//https://github.com/mies/martini-rethink

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/pagination"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/curt-labs/GoAPI/models/faq"
	"github.com/go-martini/martini"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestFaqs(t *testing.T) {
	var f faq_model.Faq
	var err error
	Convey("Test Faqs", t, func() {
		//test create
		form := url.Values{"question": {"test"}, "answer": {"testAnswer"}}
		v := form.Encode()
		body := strings.NewReader(v)
		testThatHttp.Request("post", "/faqs", "", "", Create, body, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &f)
		So(f, ShouldHaveSameTypeAs, faq_model.Faq{})

		//test update
		form = url.Values{"question": {"test new"}, "answer": {"testAnswer new"}}
		v = form.Encode()
		body = strings.NewReader(v)
		testThatHttp.Request("put", "/faqs/", ":id", strconv.Itoa(f.ID), Update, body, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &f)
		So(f, ShouldHaveSameTypeAs, faq_model.Faq{})

		//test get
		testThatHttp.Request("get", "/faqs/", ":id", strconv.Itoa(f.ID), Get, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &f)
		So(f, ShouldHaveSameTypeAs, faq_model.Faq{})
		So(f.Question, ShouldEqual, "test new")

		//test getall
		testThatHttp.Request("get", "/faqs", "", "", GetAll, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var fs faq_model.Faqs
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &fs)
		So(len(fs), ShouldBeGreaterThan, 0)

		//test search - responds w/ horrid pagination object
		testThatHttp.Request("get", "/faqs/search", "", "?question=test new", Search, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		var l pagination.Objects
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &l)
		So(len(l.Objects), ShouldBeGreaterThan, 0)

		//test delete
		testThatHttp.Request("delete", "/faqs/", ":id", strconv.Itoa(f.ID), Delete, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &f)
		So(f, ShouldHaveSameTypeAs, faq_model.Faq{})
	})
}

func BenchmarkGetFaqs(b *testing.B) {
	//no faq object creation (fails on empty db)
	var cu customer_new.CustomerUser
	cu.Name = "test cust user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}

	RequestBenchmark(b.N, "GET", "/faqs/1?key="+apiKey, nil, Get)
	RequestBenchmark(b.N, "GET", "/faqs?key="+apiKey, nil, GetAll)

	cu.Delete()
}

func RequestBenchmark(runs int, method, route string, body *url.Values, handler martini.Handler) {

	opts := httprunner.ReqOpts{
		Body:    body,
		Handler: handler,
		URL:     route,
		Method:  method,
	}

	(&httprunner.Runner{
		Req: &opts,
		N:   runs,
		C:   1,
	}).Run()

}
