// @APIVersion 1.0.0
// @Title API-Satya
// @Description Simple API for generating and parsing Authentication tokens and retrieving IMEIS
// @Contact leonardogcsoares93@gmail.com
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"api-satya/database"
	"api-satya/models"
	"api-satya/restiful"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type key int

const (
	secret  string = "secret" // secret key used for parsing JWT
	jwtKey  key    = 0
	imeiKey key    = 1
)

func main() {

	router := httprouter.New()

	router.Handler("POST", "/SayHello", restiful.Handle(
		AuthKeyGenerator,
		SayHello,
	))
	router.Handler("GET", "/GetHello", restiful.Handle(
		GetHello,
	))

	log.Fatal(http.ListenAndServe(":8080", router))
}

// AuthKeyGenerator is used to generate a JWT given an IMEI number
func AuthKeyGenerator(w http.ResponseWriter, r *http.Request) error {
	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)
	context.Set(r, imeiKey, u.Imei)

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["id"] = u.Imei
	// token.Claims["iat"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(time.Second * 3600 * 24).Unix()
	jwtString, err := token.SignedString([]byte(secret))

	// In case of nil error, save token and IMEI to Database
	if err == nil {
		fmt.Fprint(w, "\n\n\n"+jwtString)
		context.Set(r, jwtKey, jwtString)
	}

	return nil
}

// SayHello is a handler for generating an authentication token
// @Title SayHello
// @Description Generates and returns authentication code for given IMEI number.
// @Success 200 {object} string "Success"
// @Failure 401 {object} string "Access denied"
// @Failure 404 {object} string "Not Found"
// @router /SayHello [post]
func SayHello(w http.ResponseWriter, r *http.Request) error {
	// u := models.User{}
	// json.NewDecoder(r.Body).Decode(&u)

	imei := context.Get(r, imeiKey)
	jwtString := context.Get(r, jwtKey)
	// fmt.Println(jwtString.(string))
	database.SaveIDToDb(jwtString.(string), imei.(string))

	fmt.Fprint(w, imei.(string))
	return nil
}

// GetHello is the handler for retrieving the IMEI for a given Authentication token
// @Title GetHello
// @Description Retrieves IMEI number for given Authentication code
// @Success 200 {object} string "Success"
// @Failure 401 {object} string "Access denied"
// @Failure 404 {object} string "Not Found"
// @router /GetHello [get]
func GetHello(w http.ResponseWriter, r *http.Request) error {

	jwtString := r.Header.Get("authorization")
	// token, err := jwt.Parse(jwtString, KeyFn)
	token, err := jwt.Parse(jwtString, KeyFn)

	// if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
	// 	return nil
	// }

	// Given a valid token and no error, retrieve the ID from the database
	if err == nil && token.Valid {
		fmt.Fprint(w, "Is Valid token.\n")
		// fmt.Fprint(w, token)
		id := database.GetIDFromDb(jwtString)
		fmt.Fprintln(w, "Id for Auth token: "+id)
	}
	return nil
}

// KeyFn is a simple callback that returns the secret key for parsing an Auth token
func KeyFn(token *jwt.Token) (interface{}, error) {
	return []byte(secret), nil
}
