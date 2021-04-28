package dnslite

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Conf *MongoClientConfig
	C    *mongo.Collection
}

type MongoClientConfig struct {
	DB   string
	Coll string
	Addr string
}

func NewMongoClient(conf *MongoClientConfig) *MongoClient {
	client, err := mongo.NewClient(options.Client().ApplyURI(conf.Addr))
	if err != nil {
		log.Println("new client error:", conf.Addr, err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Println("connect error:", err)
		return nil
	}
	d := client.Database(conf.DB)
	c := d.Collection(conf.Coll)
	return &MongoClient{
		Conf: conf,
		C:    c,
	}
}

func (mc *MongoClient) Insert(data interface{}) error {
	ctx := context.Background()
	_, err := mc.C.InsertOne(ctx, data)
	if err != nil {
		log.Println("insert data error:", data, err)
		return err
	}
	return nil
}

func (mc *MongoClient) Find(filter, ret interface{}) error {
	ctx := context.Background()
	c, err := mc.C.Find(ctx, filter)
	if err != nil {
		log.Println("find error:", err)
		return err
	}
	err = c.All(ctx, ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.Println("find one error:", filter, err)
		return err
	}
	return nil
}

func (mc *MongoClient) Del(filter interface{}) error {
	ctx := context.Background()
	_, err := mc.C.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("delete one error:", filter, err)
		return err
	}
	return nil
}
