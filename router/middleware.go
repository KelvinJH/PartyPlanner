package router

import (
	"fmt"
	"net/http"
)

func Authenticated(handler http.HandlerFunc) http.HandlerFunc {
	fmt.Println("Inside authenticated middleware returning" ) 
	return func(w http.ResponseWriter, r *http.Request) {

		handler(w, r)
	}
}