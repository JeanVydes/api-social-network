package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func AddReplyToComment(commentID string, reply PostCommentReply) (bool, error) {
	command := bson.M{
		"$set": bson.M{
			"replies." + reply.ID: reply,
		},
	}

	filter := bson.M{"_id": commentID}

	_, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(CommentsCollection).UpdateOne(context.TODO(), filter, command)
	if err != nil {
		return false, err
	}

	return true, nil
}

func RemoveReplyFromComment(commentID string, replyID string) (bool, error) {
	command := bson.M{
		"$unset": bson.M{
			"replies." + replyID: 1,
		},
	}

	filter := bson.M{"_id": commentID}

	_, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(CommentsCollection).UpdateOne(context.TODO(), filter, command)
	if err != nil {
		return false, err
	}

	return true, nil
}
