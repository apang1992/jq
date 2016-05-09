# jq
--------------------
json query with golang
## Introduction
--------------------
 You can query a json with a string like this : __".data.productList[1].productInfo.salePrice"__  
 The idea comes from a C project : https://github.com/stedolan/jq  
 You are free to use it. However ,it's your own risk to use it in production!  
 Have Fun!  

### Following is a sample program.
```go 
package main

import (
	"fmt"
	"os"

	"github.com/apang1992/jq"
)

var data = `
  {
    "result": "1",
    "errorCode": 0,
    "data": {
      "productList": [
        {
          "productID": 100,
          "type": "simple",
          "productInfo": {
            "productID": 100,
            "productName": "Apple 500g",
            "salePrice": 118
          }
        },
        {
          "productID": 101,
          "type": "complext",
          "productInfo": {
            "productID": 101,
            "productName": "Pear 500g",
            "salePrice": 130
          }
        }
      ]
    }
  }
`

func main() {
	price, err := jq.JsonQuery([]byte(data), ".data.productList[1].productInfo.salePrice")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println("the price is:", string(price))
	}
}



//you will get :the price is: 130
```
