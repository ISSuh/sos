// MIT License

// Copyright (c) 2024 ISSuh

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	middlewares []MiddlewareFunc
	router      *mux.Router
}

func NewServer() Server {
	return Server{
		router: mux.NewRouter(),
	}
}

func (s *Server) Use(m MiddlewareFunc) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) Mux(pattern string, method string, handler http.HandlerFunc) {
	h := http.HandlerFunc(handler)
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		h = s.middlewares[i](h)
	}

	s.router.Handle(pattern, h).Methods(method)
}

func (s *Server) MuxAll(routeList RouteList) {
	for _, item := range routeList {
		h := http.HandlerFunc(item.Handler)

		for i := len(item.Middlewares) - 1; i >= 0; i-- {
			h = item.Middlewares[i](h)
		}

		for i := len(s.middlewares) - 1; i >= 0; i-- {
			h = s.middlewares[i](h)
		}

		s.router.Handle(item.URL, h).Methods(item.Method)
	}
}

func (s *Server) Run(address string) error {
	return http.ListenAndServe(address, s.router)
}
