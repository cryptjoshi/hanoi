package common 

import (
	"os"
	"fmt"
	"encoding/base64"
)

// var PG_API_KEY = "9dc857f4-2225-45ef-bf0f-665bcf7d4a1b"  
// var PG_API_KEY= "31d3cc58-4e34-4dc4-9c45-b8abe6a1b0d2"
var SECRET_KEY = os.Getenv("PASSWORD_SECRET")
var PG_prod_code = os.Getenv("PG_PRODUCT_ID")
var OPERATOR_CODE = "sunshinepgthb"// os.Getenv("PG_API_USER")//"sunshinepgthb"//"sunshinetest"//"sunshinetest",
var SECRET_API_KEY =  "9dc857f4-2225-45ef-bf0f-665bcf7d4a1b"//os.Getenv("PG_API_KEY")
var PG_PROD_CODE= os.Getenv("PG_PRODUCT_ID")
var PG_PROD_URL = os.Getenv("PG_API_URL") // "https://api.hentory.io"
var PG_API_URL =  "https://api.hentory.io"//os.Getenv("PG_API_URL") //"https://test.ambsuperapi.com"

var EF_SECRET_KEY= os.Getenv("EF_SECRET_KEY") //"1g1bb3" //"456Ayb"  //stagging
var EF_OPERATOR_CODE = "E293" //os.Getenv("EF_OPERATOR")

//var EF_API_URL = os.Getenv("INFINITY_STAG_URL") //"https://swmd.6633663.com/"
//var OPERATOR_CODE = os.Getenv("INFINITY_OPERATOR_CODE") || "E293"
var INFINITY_PROD_URL  = os.Getenv("INFINITY_PROD_URL") // "https://prod_md.9977997.com"
var INFINITY_STAG_URL = "https://swmd.6633663.com"  //os.Getenv("INFINITY_STAG_URL") // "https://stag_md.9977997.com"
//var DEVELOPMENT_SECRET_KEY = os.Getenv("DEVELOPMENT_SECRET_KEY")    || "1g1bb3"  //staging
//var PRODUCTION_SECRET_KEY= os.Getenv("PRODUCTION_SECRET_KEY")  || "456Ayb" //product

var USER_FIX = os.Getenv("USER_FIX") 
var PASS_FIX = os.Getenv("PASS_FIX") 

func CreateBasicAuthHeader(operatorCode, secretKey string) string {
    // Combine the operator code and secret key
    credentials := fmt.Sprintf("%s:%s", operatorCode, secretKey)

    // Encode the credentials to base64
    encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

    // Return the Basic Auth header
    return "basic " + encodedCredentials
}