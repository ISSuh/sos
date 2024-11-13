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

package rpc

import (
	"fmt"
	"net"

	"github.com/ISSuh/sos/internal/validation"
	"google.golang.org/grpc"
)

type Server struct {
	engine    Engine
	registers []RegisterFunc
}

func NewServer() Server {
	return Server{
		engine: Engine{
			Server: grpc.NewServer(),
		},
		registers: make([]RegisterFunc, 0),
	}
}

func (s *Server) Regist(functions []RegisterFunc) {
	s.registers = append(s.registers, functions...)
}

func (s *Server) Run(address string) error {
	switch {
	case validation.IsEmpty(address):
		return fmt.Errorf("address is empty")
	case len(s.registers) == 0:
		return fmt.Errorf("register functions is empty")
	}

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	for _, f := range s.registers {
		f(&s.engine)
	}

	if err := s.engine.Serve(l); err != nil {
		return err
	}
	return nil
}
