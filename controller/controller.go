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
	RetrieveAllTaxes(context.Context) (TaxBill, error)
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

	tax := m.Tax{Name: taxData.Name, TaxCode: taxData.TaxCode, Price: taxData.Price}
	tax, err := c.database.CreateTax(tax)

	if err != nil {
		return m.Tax{}, err
	}

	return tax, nil
}

// TaxDetail holds the struture of a Tax detail
type TaxDetail struct {
	Name       string
	TaxCode    uint
	TaxType    string
	Refundable bool
	Price      int
	TaxFee     float64
	TotalPrice float64
}

// TaxBill holds the structure of a Bill object
type TaxBill struct {
	TaxDetails    []TaxDetail
	PriceSubTotal int
	TaxSubTotal   float64
	GrandTotal    float64
}

func (c *controller) createTaxDetail(tax m.Tax) TaxDetail {
	name := tax.Name
	taxCode, _ := c.database.RetrieveTaxCodeByCode(tax.TaxCode)
	taxCodeName := taxCode.Name
	code := taxCode.Code
	price := tax.Price

	refundable := true
	if code != 1 {
		refundable = false
	}

	taxFee := 0.0
	if code == 1 {
		taxFee = float64(price) * 0.1
	} else if code == 2 {
		taxFee = 10 + (float64(price) * 0.02)
	} else {
		if price >= 100 {
			taxFee = float64(price-100) * 0.01
		}
	}

	totalPrice := float64(price) + taxFee

	return TaxDetail{
		Name:       name,
		TaxCode:    code,
		TaxType:    taxCodeName,
		Price:      price,
		Refundable: refundable,
		TaxFee:     taxFee,
		TotalPrice: totalPrice,
	}
}

// RetrieveAllTaxes is a function for getting taxes
func (c *controller) RetrieveAllTaxes(ctx context.Context) (TaxBill, error) {
	select {
	case <-ctx.Done():
		return TaxBill{}, ctx.Err()
	default:
	}

	taxes, err := c.database.RetrieveAllTaxes()
	if err != nil {
		return TaxBill{}, err
	}

	priceSubTotal, taxSubTotal, grandTotal := 0, 0.0, 0.0
	numTaxes := len(taxes)
	var bills = make([]TaxDetail, numTaxes)
	for i, tax := range taxes {
		bills[i] = c.createTaxDetail(tax)
		priceSubTotal += bills[i].Price
		taxSubTotal += bills[i].TaxFee
		grandTotal += bills[i].TotalPrice
	}

	taxBill := TaxBill{TaxDetails: bills, PriceSubTotal: priceSubTotal, TaxSubTotal: taxSubTotal, GrandTotal: grandTotal}
	return taxBill, nil
}
