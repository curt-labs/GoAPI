package faq_model

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetFaqs(t *testing.T) {
	var f Faq
	f.Question = "testQuestion"
	f.Answer = "testAnswer"
	var err error
	Convey("Testing Create", t, func() {
		err = f.Create()
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		So(f.Question, ShouldEqual, "testQuestion")
		So(f.Answer, ShouldEqual, "testAnswer")
	})
	Convey("Testing Update", t, func() {
		f.Question = "testQuestion222"
		f.Answer = "testAnswer222"
		err = f.Update()
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		So(f.Question, ShouldEqual, "testQuestion222")
		So(f.Answer, ShouldEqual, "testAnswer222")
	})
	Convey("Testing Get", t, func() {
		err = f.Get()
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		So(f.Question, ShouldHaveSameTypeAs, "str")
		So(f.Answer, ShouldHaveSameTypeAs, "str")

		page := "1"
		result := "1"

		qs, err := GetQuestions(page, result)
		So(qs, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(qs.Pagination.Page, ShouldEqual, 1)
		So(qs.Pagination.ReturnedCount, ShouldNotBeNil)
		So(qs.Pagination.PerPage, ShouldEqual, len(qs.Objects))

		as, err := GetAnswers(page, result)
		So(as, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(as.Pagination.Page, ShouldEqual, 1)
		So(as.Pagination.ReturnedCount, ShouldNotBeNil)
		So(as.Pagination.PerPage, ShouldNotBeNil)
		So(as.Pagination.PerPage, ShouldEqual, len(as.Objects))

		var fs Faqs
		fs, err = GetAll()
		So(fs, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(len(fs), ShouldNotBeNil)

		as, err = Search(f.Question, "", "1", "0")
		So(as, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(as.Pagination.Page, ShouldEqual, 1)
		So(as.Pagination.ReturnedCount, ShouldNotBeNil)
		So(as.Pagination.PerPage, ShouldNotBeNil)
		So(as.Pagination.PerPage, ShouldEqual, len(as.Objects))
	})
	Convey("Testing Delete", t, func() {
		err = f.Delete()
		So(err, ShouldBeNil)

	})

}
