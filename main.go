package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"
	"token-encrypt/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 定义配置文件解析后的结构
type User struct {
	Id                      primitive.ObjectID `bson:"_id"       json:"_id"`
	Name                    string             `bson:"name"       json:"name"`
	Email                   string             `bson:"email"      json:"email"`
	Bio                     string             `bson:"bio"        json:"bio"`
	AvatarId                string             `bson:"avatar_id"  json:"avatar_id"`
	PlatformToken           string             `bson:"token"      json:"token"`
	PlatformUserId          string             `bson:"uid"        json:"uid"`
	PlatformUserNamespaceId string             `bson:"nid"        json:"nid"`
	Follower                []string           `bson:"follower"   json:"follower"`
	Following               []string           `bson:"following"  json:"following"`
	Version                 int                `bson:"version"    json:"version"`
}

func main() {

	eHelper, _ := utils.NewSymmetricEncryption("xx", "")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := client.Database("test").Collection("user")

	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cur, _ := collection.Find(ctx, bson.M{})

	for cur.Next(ctx) {
		var result User
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		token := result.PlatformToken
		uname := result.Name
		etoken, err := eHelper.Encrypt([]byte(token))
		if err != nil {
			return
		}
		newtoken := hex.EncodeToString(etoken)

		r, err := collection.UpdateOne(
			ctx, bson.M{
				"name": uname,
			},
			bson.M{
				"$set": bson.M{"token": newtoken},
			},
		)

		fmt.Printf("a: %v\n", r)
		if err != nil {
			fmt.Println("aaaa")
		}

	}

}
