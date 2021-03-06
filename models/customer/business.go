package customer

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/sortutil"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getBusinessClassesStmt = `select b.BusinessClassID, b.name, b.sort, b.showOnWebsite from BusinessClass as b
		where (b.brandID = ? or 0 = ?) && b.showOnWebsite = 1
		group by b.name
		order by b.sort`
	createBusinessClass = `insert into BusinessClass (name, sort, showOnWebsite, brandID) values (?,?,?,?)`
	deleteBusinessClass = `delete from BusinessClass where BusinessClassID = ?`
)

type BusinessClasses []BusinessClass
type BusinessClass struct {
	ID            int    `json:"id" xml:"id"`
	Name          string `json:"name" xml:"name"`
	Sort          int    `json:"sort" xml:"sort"`
	ShowOnWebsite bool   `json:"show" xml:"show"`
	BrandID       int    `json:"brandID,omitempty" xml:"brandID,omitempty"`
}

func GetAllBusinessClasses(dtx *apicontext.DataContext) (classes BusinessClasses, err error) {
	err = database.Init()
	if err != nil {
		return
	}

	stmt, err := database.DB.Prepare(getBusinessClassesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(dtx.BrandID, dtx.BrandID)
	if err != nil {
		return
	}
	var bc BusinessClass
	for rows.Next() {
		bc = BusinessClass{}
		err = rows.Scan(
			&bc.ID,
			&bc.Name,
			&bc.Sort,
			&bc.ShowOnWebsite,
		)
		if err != nil {
			return
		}
		classes = append(classes, bc)
	}
	defer rows.Close()

	sortutil.AscByField(classes, "Sort")
	return
}

func (b *BusinessClass) Create() error {
	var err error
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(createBusinessClass)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(b.Name, b.Sort, b.ShowOnWebsite, b.BrandID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	b.ID = int(id)
	return err
}

func (b *BusinessClass) Delete() error {
	var err error
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteBusinessClass)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(b.ID)
	if err != nil {
		return err
	}
	return err
}
