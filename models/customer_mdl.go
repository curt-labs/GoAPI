package models

import (
	"../helpers/database"
	"net/url"
)

var (
	customerPriceStmt = `select distinct cp.price from ApiKey as ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join CustomerPricing cp on c.customerID = cp.cust_id
					where api_key = '%s'
					and cp.partID = %d`

	customerPartStmt = `select distinct ci.custPartID from ApiKey as ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join CartIntegration ci on c.customerID = ci.custID
					where ak.api_key = '%s'
					and ci.partID = %d`

	customerStmt = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
				c.latitude, c.longitude, c.searchURL, c.logo, c.website,
				c.postal_code, s.state, s.abbr as state_abbr, cty.name as country_name, cty.abbr as country_abbr,
				d_types.type as dealer_type, d_tier.tier as dealer_tier, mpx.code as mapix_code, mpx.description as mapic_desc,
				sr.name as rep_name, sr.code as rep_code, c.parentID
				from Customer as c
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as d_types on c.dealer_type = d_types.dealer_type
				left join DealerTiers d_tier on c.tier = d_tier.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where c.customerID = %d`

	customerLocationsStmt = `select cl.locationID, cl.name, cl.email, cl.address, cl.city,
					cl.postalCode, cl.phone, cl.fax, cl.latitude, cl.longitude,
					cl.cust_id, cl.contact_person, cl.isprimary, cl.ShippingDefault,
					s.state, s.abbr as state_abbr, cty.name as cty_name, cty.abbr as cty_abbr
					from CustomerLocations as cl
					left join States as s on cl.stateID = s.stateID
					left join Country as cty on s.countryID = cty.countryID
					where cl.cust_id = %d`
)

type Customer struct {
	Id                                   int
	Name, Email, Address, Address2, City string
	State, StateAbbreviation             string
	Country, CountryAbbreviation         string
	PostalCode                           string
	Phone, Fax                           string
	ContactPerson                        string
	Latitude, Longitude                  float64
	Website                              *url.URL
	Parent                               *Customer
	SearchUrl, Logo                      *url.URL
	DealerType, DealerTier               string
	SalesRepresentative                  string
	SalesRepresentativeCode              int
	MapixCode, MapixDescription          string
	Locations                            *[]CustomerLocation
	Users                                []CustomerUser
}

type CustomerLocation struct {
	Id                                     int
	Name, Email, Address, City, PostalCode string
	State, StateAbbreviation               string
	Country, CountryAbbreviation           string
	Phone, Fax                             string
	Latitude, Longitude                    float64
	CustomerId                             int
	ContactPerson                          string
	IsPrimary, ShippingDefault             bool
}

func (c *Customer) GetCustomer() (err error) {

	locationChan := make(chan int)
	go func() {
		if locErr := c.GetLocations(); locErr != nil {
			err = locErr
		}
		locationChan <- 1
	}()

	err = c.Basics()

	<-locationChan

	return err
}

func (c *Customer) Basics() error {

	row, res, err := database.Db.QueryFirst(customerStmt, c.Id)
	if database.MysqlError(err) {
		return err
	}

	customerID := res.Map("customerID")
	name := res.Map("name")
	email := res.Map("email")
	address := res.Map("address")
	address2 := res.Map("address2")
	city := res.Map("city")
	phone := res.Map("phone")
	fax := res.Map("fax")
	contact := res.Map("contact_person")
	lat := res.Map("latitude")
	lon := res.Map("longitude")
	search := res.Map("searchURL")
	site := res.Map("website")
	logo := res.Map("logo")
	zip := res.Map("postal_code")
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	country := res.Map("country_name")
	country_abbr := res.Map("country_abbr")
	dealer_type := res.Map("dealer_type")
	dealer_tier := res.Map("dealer_tier")
	mpx_code := res.Map("mapix_code")
	mpx_desc := res.Map("mapic_desc")
	rep_name := res.Map("rep_name")
	rep_code := res.Map("rep_code")
	parentID := res.Map("parentID")

	sURL, _ := url.Parse(row.Str(search))
	websiteURL, _ := url.Parse(row.Str(site))
	logoURL, _ := url.Parse(row.Str(logo))

	c.Id = row.Int(customerID)
	c.Name = row.Str(name)
	c.Email = row.Str(email)
	c.Address = row.Str(address)
	c.Address2 = row.Str(address2)
	c.City = row.Str(city)
	c.State = row.Str(state)
	c.StateAbbreviation = row.Str(state_abbr)
	c.Country = row.Str(country)
	c.CountryAbbreviation = row.Str(country_abbr)
	c.PostalCode = row.Str(zip)
	c.Phone = row.Str(phone)
	c.Fax = row.Str(fax)
	c.ContactPerson = row.Str(contact)
	c.Latitude = row.ForceFloat(lat)
	c.Longitude = row.ForceFloat(lon)
	c.Website = websiteURL
	c.SearchUrl = sURL
	c.Logo = logoURL
	c.DealerType = row.Str(dealer_type)
	c.DealerTier = row.Str(dealer_tier)
	c.SalesRepresentative = row.Str(rep_name)
	c.SalesRepresentativeCode = row.Int(rep_code)
	c.MapixCode = row.Str(mpx_code)
	c.MapixDescription = row.Str(mpx_desc)

	if row.Int(parentID) != 0 {
		parent := Customer{
			Id: row.Int(parentID),
		}
		if err = parent.GetCustomer(); err == nil {
			c.Parent = &parent
		}
	}

	return nil
}

func (c *Customer) GetLocations() error {
	rows, res, err := database.Db.Query(customerLocationsStmt, c.Id)
	if database.MysqlError(err) {
		return err
	}

	locationID := res.Map("locationID")
	name := res.Map("name")
	email := res.Map("email")
	address := res.Map("address")
	city := res.Map("city")
	phone := res.Map("phone")
	fax := res.Map("fax")
	contact := res.Map("contact_person")
	lat := res.Map("latitude")
	lon := res.Map("longitude")
	zip := res.Map("postalCode")
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	country := res.Map("cty_name")
	country_abbr := res.Map("cty_abbr")
	customerID := res.Map("cust_id")
	isPrimary := res.Map("isprimary")
	shipDefault := res.Map("ShippingDefault")

	var locs []CustomerLocation
	for _, row := range rows {
		l := CustomerLocation{
			Id:                  row.Int(locationID),
			Name:                row.Str(name),
			Email:               row.Str(email),
			Address:             row.Str(address),
			City:                row.Str(city),
			State:               row.Str(state),
			StateAbbreviation:   row.Str(state_abbr),
			Country:             row.Str(country),
			CountryAbbreviation: row.Str(country_abbr),
			PostalCode:          row.Str(zip),
			Phone:               row.Str(phone),
			Fax:                 row.Str(fax),
			ContactPerson:       row.Str(contact),
			CustomerId:          row.Int(customerID),
			Latitude:            row.ForceFloat(lat),
			Longitude:           row.ForceFloat(lon),
			IsPrimary:           row.ForceBool(isPrimary),
			ShippingDefault:     row.ForceBool(shipDefault),
		}
		locs = append(locs, l)
	}
	c.Locations = &locs
	return nil
}

func GetCustomerPrice(api_key string, part_id int) (price float64, err error) {
	db := database.Db

	row, _, err := db.QueryFirst(customerPriceStmt, api_key, part_id)
	if database.MysqlError(err) {
		return
	}
	if len(row) == 1 {
		price = row.Float(0)
	}

	return
}

func GetCustomerCartReference(api_key string, part_id int) (ref int, err error) {
	db := database.Db

	row, _, err := db.QueryFirst(customerPartStmt, api_key, part_id)
	if database.MysqlError(err) {
		return
	}

	if len(row) == 1 {
		ref = row.Int(0)
	}

	return
}