package mongorepo

import (
	"context"
	"log"
	"web-app-test/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepo struct {
	Client  *mongo.Client
	hostUrl string
}

func NewDB(hostUrl string) *mongoRepo {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(hostUrl).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return &mongoRepo{hostUrl: hostUrl, Client: client}
}

func (repo *mongoRepo) Insert(obj entity.RequestObj) {
	_, err := repo.Client.Database("core").Collection("request").InsertOne(context.TODO(), obj)
	if err != nil {
		log.Println(err)
	}
}
