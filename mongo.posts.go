package main

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPostByID(postID string) (Post, bool) {
	var post Post

	filter := bson.M{"_id": postID}
	_ = MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(PostsCollection).FindOne(context.TODO(), filter).Decode(&post)

	if post.ID == "" {
		return post, false
	}

	return post, true
}

func DeletePostByID(postID string) (bool, error) {
	filter := bson.M{"_id": postID}
	_, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(PostsCollection).DeleteOne(context.TODO(), filter)

	if err != nil {
		return false, err
	}

	return true, err
}

func AddLikeToPost(postID string, accountID string) (*mongo.UpdateResult, error) {
	command := bson.M{
		"$set": bson.M{
			"likes." + accountID: GenericLike{
				AccountID: accountID,
				Date:      time.Now().Unix(),
			},
		},
	}

	filter := bson.M{"_id": postID}

	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(PostsCollection).UpdateOne(context.TODO(), filter, command)

	return result, err
}

func RemoveLikeFromPost(postID string, accountID string) (*mongo.UpdateResult, error) {
	command := bson.M{
		"$unset": bson.M{
			"likes." + accountID: 1,
		},
	}

	filter := bson.M{"_id": postID}

	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(PostsCollection).UpdateOne(context.TODO(), filter, command)

	return result, err
}

func GetUserPosts(accountID string, limit int64, order interface{}) ([]Post, error) {
	var posts []Post
	filter := bson.M{
		"author_id": accountID,
	}

	options := options.Find()
	options.SetLimit(limit)
	options.SetSort(order)

	cur, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(PostsCollection).Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var post Post
		err := cur.Decode(&post)
		if err != nil {
			continue
		}

		posts = append(posts, post)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func SavePostByID(accountID string, postID string) (bool, error) {
	command := bson.M{
		"$set": bson.M{
			"saved_posts." + postID: SavedPost{
				ID:   postID,
				Date: time.Now().Unix(),
			},
		},
	}

	filter := bson.M{"_id": accountID}

	_, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).UpdateOne(context.TODO(), filter, command)

	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteSavedPostByID(accountID string, postID string) (bool, error) {
	command := bson.M{
		"$unset": bson.M{
			"saved_posts." + postID: 1,
		},
	}

	filter := bson.M{"_id": accountID}

	_, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).UpdateOne(context.TODO(), filter, command)

	if err != nil {
		return false, err
	}

	return true, nil
}
