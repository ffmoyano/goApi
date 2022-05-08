package main

import (
	controller "github.com/ffmoyano/goApi/controller"
	"net/http"
)

func addHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/signup", isValidMethod(controller.SignUp, Post))
	mux.HandleFunc("/login", isValidMethod(controller.Login, Post))
	mux.HandleFunc("/users", isValidMethod(controller.GetUsers, Get))
}
