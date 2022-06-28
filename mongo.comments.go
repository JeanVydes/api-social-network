package main

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCommentsByPostID(postID string, limit int64, order interface{}) ([]PostComment, error) {
	var comments []PostComment
	filter := bson.M{
		"post_id": postID,
	}

	options := options.Find()
	options.SetLimit(limit)
	options.SetSort(order)

	cur, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(CommentsCollection).Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var comment PostComment
		err := cur.Decode(&comment)
		if err != nil { 
			continue
		}
		
		comments = append(comments, comment)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func GetCommentByID(commentID string) (PostComment, bool) {
	var comment PostComment
	filter := bson.M{"_id": commentID}
	_ = MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(CommentsCollection).FindOne(context.TODO(), filter).Decode(&comment)

	if comment.ID == "" {
		return comment, false
	}

	return comment, true
}

func DeleteCommentByID(commentID string) (bool, error) {
	filter := bson.M{"_id": commentID}
	_, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(CommentsCollection).DeleteOne(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteAllCommentsFromPost(postID string) (bool, error) {
	filter := bson.M{"_id": postID}
	_, err = MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(PostsCollection).DeleteMany(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return true, nil
}

func LikeCommentByID(accountID string, commentID string) (bool, error) {
	command := bson.M{
		"$set": bson.M{
			"likes" + "." + accountID: GenericLike{
				AccountID: accountID,
				Date: time.Now().Unix(),
			},
		},
	}

	filter := bson.M{"_id": commentID}

	_, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(CommentsCollection).UpdateOne(context.TODO(), filter, command)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UnlikeCommentByID(accountID string, commentID string) (bool, error) {
	command := bson.M{
		"$unset": bson.M{
			"likes" + "." + accountID: 1,
		},
	}

	filter := bson.M{"_id": commentID}

	_, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(CommentsCollection).UpdateOne(context.TODO(), filter, command)
	if err != nil {
		return false, err
	}

	return true, nil
}