package dnslite

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DNSRecord struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Type  uint16             `json:"type" bson:"type"`
	TTL   uint32             `json:"ttl" bson:"ttl"`
	View  string             `json:"view" bson:"view"`
	A     string             `json:"a,omitempty" bson:"a,omitempty"`
	AAAA  string             `json:"aaaa,omitempty" bson:"aaaa,omitempty"`
	NS    string             `json:"ns,omitempty" bson:"ns,omitempty"`
	CNAME string             `json:"cname,omitempty" bson:"cname,omitempty"`
	Txt   string             `json:"txt,omitempty" bson:"txt,omitempty"`
	CAA   string             `json:"caa,omitempty" bson:"caa,omitempty"`
	Tag   string             `json:"tag,omitempty" bson:"tag,omitempty"`
	Flag  uint8              `json:"flag,omitempty" bson:"flag,omitempty"`
	Code  int                `json:"code,omitempty" bson:"code,omitempty"`
	Msg   string             `json:"msg,omitempty" bson:"msg,omitempty"`
}

type ReadRecordArgs struct {
	Name string
	Type uint16
}

type RecordManager struct {
	DB *MongoClient
}

func NewRecordManager(conf *MongoClientConfig) *RecordManager {
	var rm RecordManager
	m := NewMongoClient(conf)
	if m == nil {
		return nil
	}
	rm.DB = m
	return &rm
}

func (rm *RecordManager) FindRecords(args ReadRecordArgs) ([]DNSRecord, error) {
	var ret []DNSRecord
	err := rm.DB.Find(bson.M{"name": args.Name, "type": args.Type}, &ret)
	if err != nil {
		log.Println("find records error:", args, err)
		return nil, err
	}
	return ret, nil
}
