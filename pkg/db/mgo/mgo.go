package mgo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"im-service/conf"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// 初始化
func Open() {
	var err error
	addr := conf.TomlConfig.MongoConfig(gin.Mode()).Addr
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(addr).SetMinPoolSize(30).SetMaxPoolSize(50)) // 连接池
	if err != nil {
		panic(err)
	}
}

type mgo struct {
	database   string
	collection string
}

func NewMgo(collection string) *mgo {
	return &mgo{
		conf.TomlConfig.MongoConfig(gin.Mode()).Database,
		collection,
	}
}

//  查询单个
func (m *mgo) FindOne(key string, value interface{}) *mongo.SingleResult {
	collection, _ := client.Database(m.database).Collection(m.collection).Clone()
	// collection.
	filter := bson.D{{key, value}}
	singleResult := collection.FindOne(context.TODO(), filter)
	return singleResult
}

// 查询多个
func (m *mgo) GetAll() (*mongo.Cursor, error) {
	collection := client.Database(m.database).Collection(m.collection)
	// 将选项传递给Find()
	ops := options.Find()
	ops.SetLimit(0)
	return collection.Find(context.TODO(), bson.D{{}}, ops)
}

// 保存对话记录
func Save(coll string, log interface{}) *mongo.InsertOneResult {
	return NewMgo(coll).InsertOne(log)
}

// 插入单个
func (m *mgo) InsertOne(value interface{}) *mongo.InsertOneResult {
	//client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	insertResult, err := collection.InsertOne(context.TODO(), value)
	if err != nil {
		fmt.Println(err)
	}
	return insertResult
}

// 查询集合里有多少数据
func (m *mgo) CollectionCount() (string, int64) {
	//client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	name := collection.Name()
	size, _ := collection.EstimatedDocumentCount(context.TODO())
	return name, size
}

// 按选项查询集合 Skip 跳过 Limit 读取数量 sort 1 ，-1 . 1 为最初时间读取 ， -1 为最新时间读取
func (m *mgo) CollectionDocuments(Skip, Limit int64, sort int) *mongo.Cursor {
	//client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	SORT := bson.D{{"_id", sort}} // filter := bson.D{{key,value}}
	filter := bson.D{{}}
	findOptions := options.Find().SetSort(SORT).SetLimit(Limit).SetSkip(Skip)
	// findOptions.SetLimit(i)
	temp, _ := collection.Find(context.Background(), filter, findOptions)
	return temp
}

// 获取集合创建时间和编号
func (m *mgo) ParsingId(result string) (time.Time, uint64) {
	temp1 := result[:8]
	timestamp, _ := strconv.ParseInt(temp1, 16, 64)
	dateTime := time.Unix(timestamp, 0) // 这是截获情报时间 时间格式 2019-04-24 09:23:39 +0800 CST
	temp2 := result[18:]
	count, _ := strconv.ParseUint(temp2, 16, 64) // 截获情报的编号
	return dateTime, count
}

// 删除文章和查询文章
func (m *mgo) DeleteAndFind(key string, value interface{}) (int64, *mongo.SingleResult) {
	//client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	filter := bson.D{{key, value}}
	singleResult := collection.FindOne(context.TODO(), filter)
	DeleteResult, err := collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		fmt.Println("删除时出现错误，你删不掉的~")
	}
	return DeleteResult.DeletedCount, singleResult
}

// 删除文章
func (m *mgo) Delete(key string, value interface{}) int64 {
	//client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	filter := bson.D{{key, value}}
	count, err := collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		fmt.Println(err)
	}
	return count.DeletedCount

}

// 删除多个
func (m *mgo) DeleteMany(key string, value interface{}) int64 {
	//client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	filter := bson.D{{key, value}}

	count, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}
	return count.DeletedCount
}
