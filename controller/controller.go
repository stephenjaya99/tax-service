package controller

import (
	"context"

	d "gitlab.com/stephenjaya99/tax-service/database"
	m "gitlab.com/stephenjaya99/tax-service/model"
)

// controller holds the structure of Contoller
type controller struct {
	database d.Database
}

// Controller holds the contract of Contoller Interface
type Controller interface {
	CreateTax(context.Context, TaxRequest) (m.Tax, error)
	RetrieveAllTaxes(context.Context) ([]TaxBill, error)
}

type ControllerOpt struct {
	Database d.Database
}

func New(opt ControllerOpt) Controller {
	return &controller{
		database: opt.Database,
	}
}

// TaxRequest holds the structure of a TaxRequest
type TaxRequest struct {
	Name    string `json:"name" binding:"required"`
	TaxCode uint   `json:"tax_code" binding:"required"`
	Price   int    `json:"price" binding:"required"`
}

// CreateTax is a function for creating tax and store it in db
func (c *controller) CreateTax(ctx context.Context, taxData TaxRequest) (m.Tax, error) {
	select {
	case <-ctx.Done():
		return m.Tax{}, ctx.Err()
	default:
	}

	tax := m.Tax{Name: taxData.Name, TaxCodeID: taxData.TaxCode, Price: taxData.Price}
	tax, err := c.database.CreateTax(tax)

	if err != nil {
		return m.Tax{}, err
	}

	return tax, nil
}

// TaxBill holds the struture of a Bill Object
type TaxBill struct {
	Name       string
	TaxCode    uint
	TaxType    string
	Refundable bool
	Price      int
	TaxFee     float64
	TotalPrice float64
}

func (c *controller) createTaxBill(tax m.Tax) TaxBill {
	name := tax.Name
	// fmt.Println(tax.TaxCode.Code, tax.TaxCode.Name)
	taxCode := tax.TaxCode.Code
	taxType := tax.TaxCode.Name
	price := tax.Price

	refundable := true
	if taxCode != 1 {
		refundable = false
	}

	taxFee := 0.0
	if taxCode == 1 {
		taxFee = float64(price) * 0.1
	} else if taxCode == 2 {
		taxFee = 10 + (float64(price) * 0.02)
	} else {
		if price >= 100 {
			taxFee = float64(price-100) * 0.01
		}
	}

	totalPrice := float64(price) + taxFee

	return TaxBill{
		Name:       name,
		TaxCode:    taxCode,
		TaxType:    taxType,
		Price:      price,
		Refundable: refundable,
		TaxFee:     taxFee,
		TotalPrice: totalPrice,
	}
}

// RetrieveAllTaxes is a function for getting taxes
func (c *controller) RetrieveAllTaxes(ctx context.Context) ([]TaxBill, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	taxes, err := c.database.RetrieveAllTaxes()
	if err != nil {
		return nil, err
	}

	numTaxes := len(taxes)
	var bills = make([]TaxBill, numTaxes)
	for i, tax := range taxes {
		bills[i] = c.createTaxBill(tax)
	}
	return bills, nil
}
