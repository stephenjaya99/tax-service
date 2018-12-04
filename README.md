## Tax Service API Documentation
Simple Tax Service API. This API uses :
- Golang (Gin Web Framework, Gorm for ORM, Govendor for handling dependencies)
- Docker for Containerization (Web API & Database)
- Postgres for Database
- Using Model-Controller-Handler pattern for handling API calls
- Utilize Golang's Interfaces for Abstraction
- Unit Tests

### Usages
- `docker-compose up` will start pulling and creating container for database, web api, and migration 
- `make test` will run unit test and calculate code coverage
- Web API will available at `localhost:8000`

### API Calls
### `GET /tax` 
#### Response
```json
{
    "Meta": {
        "Code": 200,
        "Message": "Success",
        "Error": ""
    },
    "Body": {
        "TaxDetails": [
            {
                "Name": "Entertaiment",
                "TaxCode": 3,
                "TaxType": "Entertainment",
                "Refundable": false,
                "Price": 150,
                "TaxFee": 0.5,
                "TotalPrice": 150.5
            },
            {
                "Name": "Tobacco",
                "TaxCode": 2,
                "TaxType": "Tobacco",
                "Refundable": false,
                "Price": 1000,
                "TaxFee": 30,
                "TotalPrice": 1030
            },
            {
                "Name": "Burger",
                "TaxCode": 1,
                "TaxType": "Food",
                "Refundable": true,
                "Price": 1000,
                "TaxFee": 100,
                "TotalPrice": 1100
            }
        ],
        "PriceSubTotal": 2150,
        "TaxSubTotal": 130.5,
        "GrandTotal": 2280.5
    }
}
```

### `POST /tax`
#### Request
```JSON
{
	"name":"Burger",
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
        "id": 3,
        "created_at": "2018-12-04T17:45:39.751766Z",
        "name": "Burger",
        "price": 1000,
        "tax_code": 1
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
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	TaxCode   uint      `json:"tax_code"`
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
- Data migration are handler in `migration` packages, also using `Gorm` for auto migration on modified tables

### Misc
- Spend 4 Days for creating the API, mostly spent on ~learning~ debugging Golang and Docker :v
- Excited to learn more about best practices using Docker & Dockerfile
- Gorm Postgres for some weird reason won't create related objects (Will create using default values) :")
- Workout for above is by not using Gorm foreign key definition, instead just create an additional column for references