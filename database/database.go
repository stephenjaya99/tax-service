package database

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	m "gitlab.com/stephenjaya99/tax-service/model"
)

// database holds the structure for Database
// its encapsulating GORM ORM
type database struct {
	client *gorm.DB
}

// Database holds the contract for DB interface
type Database interface {
	// Migrate should holds implementation for migrating any struct to database table
	Migrate(interface{})

	// CreateTax should holds implementation for storing tax in database
	CreateTax(m.Tax) (m.Tax, error)

	// RetrieveAlltaxes should holds implementation retrieving tax from database.
	RetrieveAllTaxes() ([]m.Tax, error)

	// CreateTaxCode should hols implementation for storing tax code in database
	CreateTaxCode(uint, string) (m.TaxCode, error)
}

// New is used to initiate interface for the DB client
// It returns the interface for the client
func New(client *gorm.DB) Database {
	return &database{client: client}
}

// Migrate is an utility to migrate any golang struct into database table
// Encapsulating GORM Automigrate
func (d *database) Migrate(i interface{}) {
	d.client.AutoMigrate(i)
}

// CreateTax is a function for storing tax in database
func (d *database) CreateTax(tax m.Tax) (m.Tax, error) {

	//set created date
	tax.CreatedAt = time.Now()
	var taxCode m.TaxCode
	if err := d.client.Where(&m.TaxCode{Code: tax.TaxCodeID}).First(&taxCode).Error; err != nil {
		return m.Tax{}, err
	}
	tax.TaxCode = m.TaxCode{}
	tax.TaxCodeID = taxCode.Code

	// put tax to database
	if err := d.client.Where(&m.Tax{Name: tax.Name}).First(&m.Tax{}).Error; err != nil {
		d.client.Create(&m.Tax{Name: tax.Name, Price: tax.Price, TaxCodeID: tax.TaxCodeID, TaxCode: taxCode})

		var taxes []m.Tax
		d.client.Find(&taxes)
		for _, elem := range taxes {
			fmt.Println(elem.Name, elem.TaxCode.Code, elem.TaxCode.Name)
		}
	}

	// update tax to latest
	d.client.Where(&m.Tax{Name: tax.Name}).First(&tax)
	tax.TaxCode = taxCode
	fmt.Println("HASIL DB ANJING", tax.Name, tax.TaxCode.Code, tax.TaxCode.Name)

	return tax, nil
}

// CreateTaxCode is a funtion for creating initial tax codes
func (d *database) CreateTaxCode(code uint, name string) (m.TaxCode, error) {

	taxCode := m.TaxCode{Code: code, Name: name}
	if err := d.client.Where(m.TaxCode{Code: code}).First(&m.TaxCode{}).Error; err != nil {
		d.client.Create(&taxCode)
	}

	return taxCode, nil
}

// RetrieveAllTax is a function for retrieving tax from database
func (d *database) RetrieveAllTaxes() ([]m.Tax, error) {

	var taxes []m.Tax
	d.client.Find(&taxes)
	for _, elem := range taxes {
		fmt.Println(elem.Name, elem.TaxCode.Code, elem.TaxCode.Name)
	}
	return taxes, nil
}
