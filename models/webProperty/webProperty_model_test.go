package webProperty_model

import (
	"database/sql"
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var once sync.Once
var MockedDTX *apicontext.DataContext

func initDtx() {
	MockedDTX, _ = apicontextmock.Mock()
}

func TestWebPropertiesBetter(t *testing.T) {
	var w WebProperty
	var wr WebPropertyRequirement
	var wn WebPropertyNote
	var wt WebPropertyType
	var err error
	once.Do(initDtx)
	defer apicontextmock.DeMock(MockedDTX)
	Convey("Testing WebProperties", t, func() {
		//New WebProperty
		w.Name = "test prop"
		w.Url = "www.hotdavid.com"
		w.CustID = MockedDTX.CustomerID

		//make up badge
		seed := int64(time.Now().Second() + time.Now().Minute() + time.Now().Hour() + time.Now().Year())
		rand.Seed(seed)
		w.BadgeID = strconv.Itoa(rand.Int()) //random badge
		w.CustID = MockedDTX.CustomerID
		//Test Requirement
		wr.ReqType = "Req Type"
		err = wr.Create(MockedDTX)
		So(err, ShouldBeNil)
		//Test Note
		wn.Text = "Note text"
		err = wn.Create(MockedDTX)
		So(err, ShouldBeNil)

		//Test Type
		wt.Type = "A type"
		wt.TypeID = 1
		err = wt.Create(MockedDTX)
		So(err, ShouldBeNil)

		//Create Web Property
		w.WebPropertyRequirements = append(w.WebPropertyRequirements, wr)
		w.WebPropertyNotes = append(w.WebPropertyNotes, wn)
		w.WebPropertyType = wt
		w.CustID = MockedDTX.CustomerID
		err = w.Create(MockedDTX)
		So(err, ShouldBeNil)
		So(w, ShouldNotBeNil)

		wr.Compliance = true
		err = wr.Update(MockedDTX)
		So(err, ShouldBeNil)
		wn.Text = "New Text"
		err = wn.Update(MockedDTX)
		So(err, ShouldBeNil)
		wt.Type = "B type"
		err = wt.Update(MockedDTX)
		So(err, ShouldBeNil)
		//Update Property
		w.Name = "New Name"
		err = w.Update(MockedDTX)
		So(err, ShouldBeNil)

		err = wr.Get()
		So(err, ShouldBeNil)

		err = wn.Get()
		So(err, ShouldBeNil)

		err = wt.Get()
		So(err, ShouldBeNil)

		w.WebPropertyRequirements = append(w.WebPropertyRequirements, wr)
		w.WebPropertyNotes = append(w.WebPropertyNotes, wn)
		w.WebPropertyType = wt

		//Search
		obj, err := Search(w.Name, "", "", "", "", "", "", "", "", "", "", "", "1", "1")
		So(err, ShouldBeNil)
		So(len(obj.Objects), ShouldEqual, 0)

		//Get Property
		err = w.Get(MockedDTX)
		So(err, ShouldBeNil)

		ws, err := GetAll(MockedDTX)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ws), ShouldBeGreaterThan, 0)
		}
		ns, err := GetAllWebPropertyNotes(MockedDTX)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ns), ShouldBeGreaterThan, 0)
		}

		rs, err := GetAllWebPropertyRequirements(MockedDTX)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(rs), ShouldBeGreaterThan, 0)
		}
		ts, err := GetAllWebPropertyTypes(MockedDTX)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ts), ShouldBeGreaterThan, 0)
		}

		//Deletes
		// err = w.Delete()
		// So(err, ShouldBeNil)
		// err = wn.Delete()
		// So(err, ShouldBeNil)
		// err = wt.Delete()
		// So(err, ShouldBeNil)

		// err = wr.Delete()
		// So(err, ShouldBeNil)

	})
	// _ = apicontextmock.DeMock(MockedDTX)

}
func BenchmarkCreateDeleteWebProperty(b *testing.B) {
	Convey("Testing WebProperties", b, func() {
		b.ResetTimer()
		var w WebProperty
		w.Name = "test"
		w.CustID = 1
		w.BadgeID = "666"
		w.Url = "www.test.com"
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Create(MockedDTX)
			w.Delete(MockedDTX)
		}
		b.StopTimer()

	})
}

func BenchmarkCreateDeleteWebPropertyRequirement(b *testing.B) {
	Convey("Testing Requirements", b, func() {
		b.ResetTimer()
		var w WebPropertyRequirement
		w.Requirement = "test req"
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Create(MockedDTX)
			w.Delete(MockedDTX)
		}
		b.StopTimer()

	})
}
func BenchmarkCreateDeleteWebPropertyNote(b *testing.B) {
	Convey("Testing Note", b, func() {
		b.ResetTimer()
		var w WebPropertyNote
		w.Text = "test note"
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Create(MockedDTX)
			w.Delete(MockedDTX)
		}
		b.StopTimer()

	})
}
func BenchmarkCreateDeleteWebPropertyType(b *testing.B) {
	Convey("Testing Type", b, func() {
		b.ResetTimer()
		var w WebPropertyType
		w.Type = "test type"
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Create(MockedDTX)
			w.Delete(MockedDTX)
		}
		b.StopTimer()

	})
}

func BenchmarkGetWebProperty(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	Convey("Testing WebProperties", b, func() {
		b.ResetTimer()
		var w WebProperty
		w.Name = "test"
		w.CustID = 1
		w.BadgeID = "666"
		w.Url = "www.test.com"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Get(MockedDTX)
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
	//_ = apicontextmock.DeMock(MockedDTX)
}
func BenchmarkGetWebPropertyRequirement(b *testing.B) {

	Convey("Testing Requirements", b, func() {
		b.ResetTimer()
		var w WebPropertyRequirement
		w.Requirement = "test req"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Get()
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})

}
func BenchmarkGetWebPropertyNote(b *testing.B) {
	Convey("Testing Note", b, func() {
		b.ResetTimer()
		var w WebPropertyNote
		w.Text = "test note"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Get()
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
}
func BenchmarkGetWebPropertyType(b *testing.B) {
	Convey("Testing Type", b, func() {
		b.ResetTimer()
		var w WebPropertyType
		w.Type = "test type"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Get()
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
}

func BenchmarkGetAllWebProperty(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	Convey("Testing WebProperties", b, func() {
		b.ResetTimer()
		var w WebProperty
		w.Name = "test"
		w.CustID = 1
		w.BadgeID = "666"
		w.Url = "www.test.com"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetAll(MockedDTX)
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
	//_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetAllWebPropertyRequirement(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	Convey("Testing Requirements", b, func() {
		b.ResetTimer()
		var w WebPropertyRequirement
		w.Requirement = "test req"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetAllWebPropertyRequirements(MockedDTX)
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
	//_ = apicontextmock.DeMock(MockedDTX)

}
func BenchmarkGetAllWebPropertyNote(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	Convey("Testing Note", b, func() {
		b.ResetTimer()
		var w WebPropertyNote
		w.Text = "test note"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetAllWebPropertyNotes(MockedDTX)
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
	//_ = apicontextmock.DeMock(MockedDTX)

}
func BenchmarkGetAllWebPropertyType(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	Convey("Testing Type", b, func() {
		b.ResetTimer()
		var w WebPropertyType
		w.Type = "test type"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetAllWebPropertyTypes(MockedDTX)
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
	//_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkUpdateWebProperty(b *testing.B) {
	Convey("Testing WebProperties", b, func() {
		b.ResetTimer()
		var w WebProperty
		w.Name = "test"
		w.CustID = 1
		w.BadgeID = "666"
		w.Url = "www.test.com"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Update(MockedDTX)
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
}
func BenchmarkUpdateWebPropertyRequirement(b *testing.B) {
	Convey("Testing Requirements", b, func() {
		b.ResetTimer()
		var w WebPropertyRequirement
		w.Requirement = "test req"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Update(MockedDTX)
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
}
func BenchmarkUpdateWebPropertyNote(b *testing.B) {
	Convey("Testing Note", b, func() {
		b.ResetTimer()
		var w WebPropertyNote
		w.Text = "test note"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Update(MockedDTX)
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
}
func BenchmarkUpdateWebPropertyType(b *testing.B) {
	Convey("Testing Type", b, func() {
		b.ResetTimer()
		var w WebPropertyType
		w.Type = "test type"
		w.Create(MockedDTX)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Update(MockedDTX)
		}
		b.StopTimer()
		w.Delete(MockedDTX)
	})
}
