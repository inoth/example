package main

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	C_DbConstr = ""
)

var (
	DBC  *DBConnect
	once sync.Once
)

type DBConnect struct {
	Mg *mongo.Client
}

func InitDb() {
	once.Do(func() {
		DBC = &DBConnect{
			Mg: setConnect(),
		}
	})
}

func setConnect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 连接池
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(C_DbConstr).SetMaxPoolSize(20))
	if err != nil {
		log.Error(err)
	}
	return client
}

func main() {
	InitDb()
	mg := NewMgo("database", "table")
	mg.Count(map[string]string{"aaa": "aaa"})
}

type Mgo struct {
	database   string
	collection string
}

func NewMgo(database, collection string) *Mgo {
	return &Mgo{
		database,
		collection,
	}
}

func (m *Mgo) GetCollection() (*mongo.Collection, error) {
	client := DBC.Mg
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "")
	}
	return col, nil
}

func (m *Mgo) InsertOne(value interface{}) bool {
	client := DBC.Mg
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return false
	}
	_, err = col.InsertOne(context.TODO(), value)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (m *Mgo) InsertMany(values []interface{}) bool {
	client := DBC.Mg
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return false
	}
	_, err = col.InsertMany(context.TODO(), values)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (m *Mgo) FindOne(filter interface{}, res interface{}) bool {
	client := DBC.Mg
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return false
	}
	err = col.FindOne(context.TODO(), filter).Decode(res)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (m *Mgo) FindMany(skip, limit int64, filter, sort interface{}) ([]bson.M, bool) {
	client := DBC.Mg
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return nil, false
	}
	var findoptions *options.FindOptions
	if limit > 0 {
		findoptions = &options.FindOptions{}
		findoptions.SetSkip(limit * (skip - 1))
		findoptions.SetLimit(limit)
		findoptions.SetSort(sort)
	}
	cur, err := col.Find(context.Background(), filter, findoptions)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	r := make([]bson.M, 0)
	ctx := context.TODO()
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var tmp bson.M
		if err = cur.Decode(&tmp); err != nil {
			log.Fatal(err)
		}
		r = append(r, tmp)
	}
	return r, true
}

func (m *Mgo) Aggregate(pipeline interface{}) ([]map[string]interface{}, bool) {
	client := DBC.Mg
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return nil, false
	}
	opts := options.Aggregate()
	cur, err := col.Aggregate(context.Background(), pipeline, opts)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	ctx := context.TODO()
	r := make([]map[string]interface{}, 0)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var tmp map[string]interface{}
		if err = cur.Decode(&tmp); err != nil {
			log.Fatal(err)
		}
		r = append(r, tmp)
	}
	return r, true
}

func (m *Mgo) Count(filter interface{}) int {
	client := DBC.Mg
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return 0
	}
	c, err := col.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0
	}
	return int(c)
}

func (m *Mgo) UpdateOne(filter, update interface{}) bool {
	client := DBC.Mg
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return false
	}
	_, err = col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}
