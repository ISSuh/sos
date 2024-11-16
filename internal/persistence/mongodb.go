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

package persistence

import (
	"context"
	"fmt"

	"github.com/ISSuh/sos/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	config config.Database
	engin  *mongo.Client
}

func NewMongoDB(dbConfig config.Database) (*MongoDB, error) {
	return &MongoDB{
		config: dbConfig,
		engin:  nil,
	}, nil
}

func (p *MongoDB) Connect(c context.Context) (*DB, error) {
	credential := options.Credential{
		Username: p.config.Credentials.Username,
		Password: p.config.Credentials.Password,
	}

	options :=
		options.Client().
			ApplyURI(p.config.Host).
			SetAuth(credential)

	client, err := mongo.Connect(c, options)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	err = client.Ping(c, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	db := client.Database(p.config.DatabaseName)
	return &DB{
		p.engin: db,
	}, nil
}

func Close(c context.Context, db *DB) error {
	client := db.engin.Client()
	return client.Disconnect(c)
}

func (d *DB) Collection(collection string) *mongo.Collection {
	return d.engin.Collection(collection)
}
