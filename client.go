// sean at shanghai
// 2021
// dnslite

package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	client, err := mongo.NewClient(options.Client().
		ApplyURI(conf.Addr).
		SetConnectTimeout(time.Second * 2).
		SetHeartbeatInterval(time.Second * 10).
		SetSocketTimeout(time.Second * 2).
		SetServerSelectionTimeout(time.Second * 2))
	if err != nil {
		log.Println("new client error:", conf.Addr, err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

func (mc *MongoClient) Insert(data DNSRecord) error {
	ctx := context.Background()
	data.Id = primitive.NewObjectID().Hex()
	_, err := mc.C.InsertOne(ctx, data)
	if err != nil {
		log.Println("insert data error:", data, err)
		return err
	}
	return nil
}

func (mc *MongoClient) Find(name string, tp uint16) ([]DNSRecord, error) {
	ctx := context.Background()
	filter := bson.M{"name": name, "type": tp}
	c, err := mc.C.Find(ctx, filter)
	if err != nil {
		log.Println("find error:", err)
		return nil, err
	}
	var ret []DNSRecord
	err = c.All(ctx, &ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Println("find one error:", filter, err)
		return nil, err
	}
	return ret, nil
}

func (mc *MongoClient) Del(name string, tp uint16) error {
	ctx := context.Background()
	filter := bson.M{"name": name, "type": tp}
	_, err := mc.C.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("delete one error:", filter, err)
		return err
	}
	return nil
}
