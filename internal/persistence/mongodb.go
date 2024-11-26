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

	"github.com/ISSuh/sos/internal/apm"
	"github.com/ISSuh/sos/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	engin *mongo.Database
}

func ConnectMongoDB(c context.Context, dbConfig config.Database) (*MongoDB, error) {
	logger := options.
		Logger().
		SetComponentLevel(options.LogComponentCommand, logLevel(dbConfig.LogLevel))

	credential := options.Credential{
		Username: dbConfig.Credentials.Username,
		Password: dbConfig.Credentials.Password,
	}

	options :=
		options.Client().
			ApplyURI(dbConfig.Host).
			SetAuth(credential).
			SetLoggerOptions(logger)

	options.Monitor = apm.WrapDatabaseMonitor()

	client, err := mongo.Connect(c, options)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	err = client.Ping(c, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	db := client.Database(dbConfig.DatabaseName)
	return &MongoDB{
		engin: db,
	}, nil
}

func CloseMongoDB(c context.Context, db *MongoDB) error {
	client := db.engin.Client()
	return client.Disconnect(c)
}

func logLevel(level string) options.LogLevel {
	switch level {
	case "info":
		return options.LogLevelInfo
	case "debug":
		return options.LogLevelDebug
	default:
		return options.LogLevelInfo
	}
}

func (d *MongoDB) Collection(collection string) (*mongo.Collection, error) {
	switch {
	case d.engin == nil:
		return nil, fmt.Errorf("database engine is nil")
	case collection == "":
		return nil, fmt.Errorf("collection name is empty")
	}

	return d.engin.Collection(collection), nil
}
