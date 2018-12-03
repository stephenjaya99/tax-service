## Tax Service API Documentation
Simple Tax Service API. Stacks used :
- Golang (Gin Web Framework, Gorm for ORM, Govendor for handling dependencies)
- Docker for Containerization (Web API & Database)
- Postgres for Database
- Using Model-Controller-Handler pattern for handling API calls
- Utilize Golang's Interfaces for Abstraction
- Unit Tests

### Usages
- `docker-compose up` will start pulling and creating container for database, web api, and migration 
- `make test` will run unit test and calculate the code coverage
- Web API will available at `localhost:8000`

### API Calls
### `GET /tax` Response
```json
{
  "Meta": {
      "Code": 200,
      "Message": "Success",
      "Error": ""
  },
  "Body": [
      {
        "Name": "Burger",
        "TaxCode": 1,
        "TaxType": "Food",
        "Refundable": true,
        "Price": 1000,
        "TaxFee": 100,
        "TotalPrice": 1100
      }
  ]
}
```

### `POST /tax`
#### Request
```JSON
{
	"name": "Burger",
	"tax_code": 1,
	"price": 1000
}
```
- `name` : `String`
- `tax_code` : `Integer`
- `price` : `Integer`

#### Response
```JSON
{
  "Meta": {
      "Code": 200,
      "Message": "Success",
      "Error": ""
  },
  "Body": {
      "id": 2,
      "created_at": "2018-12-03T22:07:25.174399Z",
      "name": "KFC",
      "price": 200,
      "tax_code": {
          "code": 1,
          "name": "Food"
      },
      "tax_code_id": 1
  }
}
```
- Simple API structure for recieving list of taxes
- Requests is handled by `handler` package using Golang Web HTTP library `Gin` for recieving POST requests and returning response (V in MVC Pattern) 
- Serialization are handled in `controller` package for `database` and `handler` (C in MVC Pattern)

### Database Documentation
```Golang
// Tax model
type Tax struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name" gorm:"UNIQUE"`
	Price     int       `json:"price"`

	TaxCode   TaxCode `json:"tax_code" gorm:"foreignkey:TaxCodeID"`
	TaxCodeID uint    `json:"tax_code_id"`
}

// TaxCode model for Tax
type TaxCode struct {
	Code uint   `json:"code" gorm:"primary_key"`
	Name string `json:"name"`
}
```

- `Tax` Object represents the Tax from the user, with a foreign key to a Tax Code
- `Tax Code` are made into seperate table for cases if it needs an additional tax codes
- Using Golang ORM library `Gorm`, Database objects will be handled by `database` package (M in MVC pattern)
- Data migration are handler in `migration` packages, also by using `Gorm` for auto migration on modified tables