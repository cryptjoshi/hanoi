package common 

import (
	"os"
	"fmt"
	"encoding/base64"
)

// var PG_API_KEY = "9dc857f4-2225-45ef-bf0f-665bcf7d4a1b"  
// var PG_API_KEY= "31d3cc58-4e34-4dc4-9c45-b8abe6a1b0d2"
var SECRET_KEY = os.Getenv("PASSWORD_SECRET")
var pg_prod_code = os.Getenv("PG_PRODUCT_ID")
var OPERATOR_CODE = "sunshinetest"//"sunshinetest",
var SECRET_API_KEY = os.Getenv("PG_API_KEY")
var PG_PROD_CODE= os.Getenv("PG_PRODUCT_ID")
var PG_PROD_URL = os.Getenv("PG_API_URL") //"https://prod_md.9977997.com"
var PG_API_URL = "https://test.ambsuperapi.com"


func CreateBasicAuthHeader(operatorCode, secretKey string) string {
    // Combine the operator code and secret key
    credentials := fmt.Sprintf("%s:%s", operatorCode, secretKey)

    // Encode the credentials to base64
    encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

    // Return the Basic Auth header
    return "basic " + encodedCredentials
}