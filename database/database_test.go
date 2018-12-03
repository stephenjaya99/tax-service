package database_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	d "gitlab.com/stephenjaya99/tax-service/database"
	m "gitlab.com/stephenjaya99/tax-service/model"
)

type DatabaseSuite struct {
	suite.Suite
	database d.Database
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, &DatabaseSuite{})
}

func (ds *DatabaseSuite) SetupSuite() {
	client, _ := gorm.Open("sqlite3", "test.sqlite")

	ds.database = d.New(client)

	// migrate mock database
	ds.database.Migrate(&m.Tax{})
	ds.database.Migrate(&m.TaxCode{})

	// insert TaxCode
	code1 := m.TaxCode{
		Code: 1,
		Name: "Food",
	}
	code2 := m.TaxCode{
		Code: 2,
		Name: "Tobacco",
	}
	code3 := m.TaxCode{
		Code: 3,
		Name: "Entertaiment",
	}

	client.Create(code1)
	client.Create(code2)
	client.Create(code3)
}

func (ds *DatabaseSuite) TearDownSuite() {
	os.Remove("test.sqlite")
}

func (ds *DatabaseSuite) TestCreateTaxOnSuccess() {
	dummyTax := m.Tax{
		Name:    "dummy_setup2345",
		TaxCode: m.TaxCode{Code: 1, Name: "TestCode"},
		Price:   5000,
	}

	_, err := ds.database.CreateTax(dummyTax)

	assert.Nil(ds.T(), err, "Error should be nil!")
}

func (ds *DatabaseSuite) TestCreateTaxInvalidTaxCode() {
	dummyTax := m.Tax{
		Name:      "dummy_setup2",
		TaxCode:   m.TaxCode{Code: 5, Name: "TestCode"},
		TaxCodeID: 5,
		Price:     5000,
	}

	_, err := ds.database.CreateTax(dummyTax)

	assert.NotNil(ds.T(), err, "Error should not be nil!")
}

func (ds *DatabaseSuite) TestRetrieveAllTaxesOnSuccess() {

	_, err := ds.database.RetrieveAllTaxes()

	assert.Nil(ds.T(), err, "Error should be nil!")
}

func (ds *DatabaseSuite) TestCreateTaxCodeOnSuccess() {

	_, err := ds.database.CreateTaxCode(10, "dummy_code")

	assert.Nil(ds.T(), err, "Error should be nil!")
}
