package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	c "gitlab.com/stephenjaya99/tax-service/controller"
	m "gitlab.com/stephenjaya99/tax-service/model"
)

type ControllerSuite struct {
	suite.Suite
	controller c.Controller
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, &ControllerSuite{})
}

func (cs *ControllerSuite) SetupSuite() {
	db := &DBMock{}
	cs.controller = c.New(c.ControllerOpt{Database: db})
}

func (cs *ControllerSuite) TestInstantiation() {
	db := &DBMock{}
	controller := c.New(c.ControllerOpt{Database: db})

	assert.NotNil(cs.T(), controller, "Instance should not be nil!")
}

func (cs *ControllerSuite) TestCreateTaxOnContextTimeoutError() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	_, err := cs.controller.CreateTax(ctx, c.TaxRequest{Name: "abc", TaxCode: 1, Price: 1000})

	assert.NotNil(cs.T(), err, "Error should not be nil!")
}

func (cs *ControllerSuite) TestCreateTaxOnDatabaseError() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, errs := cs.controller.CreateTax(ctx, c.TaxRequest{Name: "abc", TaxCode: 1, Price: 1000})

	assert.NotNil(cs.T(), errs, "Error should not be nil!")
}

func (cs *ControllerSuite) TestCreateTaxOnSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, errs := cs.controller.CreateTax(ctx, c.TaxRequest{Name: "test", TaxCode: 1, Price: 1000})

	assert.Nil(cs.T(), errs, "Error should be nil!")
}

func (cs *ControllerSuite) TestRetrieveAllTaxesOnContextTimeoutError() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	_, errs := cs.controller.RetrieveAllTaxes(ctx)

	assert.NotNil(cs.T(), errs, "Error should not be nil!")
}

func (cs *ControllerSuite) TestRetrieveAllTaxesOnSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	bills, errs := cs.controller.RetrieveAllTaxes(ctx)
	assert.Equal(cs.T(), "Burger", bills.TaxDetails[0].Name)
	assert.Equal(cs.T(), "Food", bills.TaxDetails[0].TaxType)
	assert.Equal(cs.T(), 100, int(bills.TaxDetails[0].TaxFee))
	assert.Equal(cs.T(), true, bills.TaxDetails[0].Refundable)
	assert.Nil(cs.T(), errs, "Error should be nil!")

}

// Struct for Mock DB
type DBMock struct{}

func (dm *DBMock) Migrate(i interface{}) {
	return
}

func (dm *DBMock) CreateTax(tax m.Tax) (m.Tax, error) {
	if tax.Name == "abc" {
		return m.Tax{}, errors.New("Tax Name doesn't exist!")
	}
	return m.Tax{}, nil
}

func (dm *DBMock) CreateTaxCode(code uint, name string) (m.TaxCode, error) {
	return m.TaxCode{}, nil
}

func (dm *DBMock) RetrieveAllTaxes() ([]m.Tax, error) {
	var taxes = make([]m.Tax, 1)
	taxes[0] = m.Tax{Name: "Burger", TaxCode: 1, Price: 1000}
	return taxes, nil
}

func (dm *DBMock) RetrieveTaxCodeByCode(code uint) (m.TaxCode, error) {
	return m.TaxCode{Code: 1, Name: "Food"}, nil
}
