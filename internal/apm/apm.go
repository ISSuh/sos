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

package apm

import (
	"os"
	"sync"

	"github.com/ISSuh/sos/internal/config"
	"go.elastic.co/apm"
)

type agent struct {
	tracer *apm.Tracer
	isInit bool
}

func (a *agent) IsUsingAPM() bool {
	return isInitialized
}

var a *agent
var once sync.Once
var isInitialized bool

func Initialize(config config.APM) error {
	os.Setenv("ELASTIC_APM_SERVER_URL", config.Host)
	os.Setenv("ELASTIC_APM_SERVICE_NAME", config.ServiceName)
	os.Setenv("ELASTIC_APM_SERVICE_VERSION", config.ServiceVersion)

	once.Do(func() {
		a = &agent{
			tracer: apm.DefaultTracer,
			isInit: true,
		}
	})

	isInitialized = true
	return nil
}

func Tracer() *apm.Tracer {
	if a == nil {
		return nil
	}
	return a.tracer
}
