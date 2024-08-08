package main

import (
	"context"
	"fmt"
	gohttp "net/http"

	"github.com/ISSuh/sos/internal/http"
)

func m1(next gohttp.HandlerFunc) gohttp.HandlerFunc {
	return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		fmt.Println("start: m1")
		next.ServeHTTP(w, r)
		fmt.Println("end: m1")

	})
}

func m2(next gohttp.HandlerFunc) gohttp.HandlerFunc {
	return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		fmt.Println("start: m2")
		next.ServeHTTP(w, r)
		fmt.Println("end: m2")
	})
}

func handler(w gohttp.ResponseWriter, r *gohttp.Request) {
	fmt.Println("Hello World")
}

type myHandler func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request)

func h2(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
	fmt.Println("Hello World")
}

func wrpper(my myHandler) gohttp.HandlerFunc {
	return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		my(, w, r)
	})
}

func main() {
	// mux := gohttp.NewServeMux()

	// mux.HandleFunc("/", m1(m2(handler)))
	// gohttp.ListenAndServe(":33114", mux)

	s := http.NewServer()
	s.Use(m1)
	s.Use(m2)

	s.Mux("/", handler)
	s.Run(":33114")
}
