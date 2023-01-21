package main

import (
	"context"

	measurementapi "github.com/hieutrtr/influxdb-poc/services/measurements/api"
	measurementdb "github.com/hieutrtr/influxdb-poc/services/measurements/db"
	"github.com/hieutrtr/influxdb-poc/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conf, err := utils.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MongoURL))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	collection := client.Database(conf.MongoDB).Collection(conf.MongoCollection)
	store := measurementdb.NewMeasurementStore(collection)
	server := measurementapi.NewServer(store, &measurementapi.Credentials{
		Username: conf.BasicAuthUsername,
		Password: conf.BasicAuthPassword,
	})
	err = server.Run(conf.ServerAddress)
	if err != nil {
		panic(err)
	}

}
