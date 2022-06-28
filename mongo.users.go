package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserByID(accountID string) (User, bool) {
	var user User

	filter := bson.M{"_id": accountID}
	_ = MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).FindOne(context.TODO(), filter).Decode(&user)

	if user.ID == "" {
		return user, false
	}

	return user, true
}

func GetUserByEmailAddress(emailAddress string) (User, bool) {
	var user User

	filter := bson.M{"email": emailAddress}
	err = MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return user, false
	}

	if user.ID == "" {
		return user, false
	}

	return user, true
}

func GetUserByUsername(username string) (User, bool) {
	var user User

	filter := bson.M{"username": username}
	err = MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return user, false
	}

	if user.ID == "" {
		return user, false
	}

	return user, true
}

func SearchUsers(searchText string) ([]PublicUser, error) {
	var users []PublicUser
	filter := bson.M{
		"$or": []interface{}{
			bson.M{
				"first_name": bson.M{
					"$regex": primitive.Regex{Pattern: fmt.Sprintf("^%s.*", searchText), Options: "i"},
				},
			},
			bson.M{
				"last_name": bson.M{
					"$regex": primitive.Regex{Pattern: fmt.Sprintf("^%s.*", searchText), Options: "i"},
				},
			},
			bson.M{
				"username": bson.M{
					"$regex": primitive.Regex{Pattern: fmt.Sprintf("^%s.*", searchText), Options: "i"},
				},
			},
		},
	}

	options := options.Find()
	options.SetLimit(25)
	options.SetSort(bson.D{})

	cur, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			continue
		}

		publicUser := PublicUser{
			ID:             user.ID,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Username:       user.Username,
			Birthday:       user.Birthday,
			ProfilePicture: user.ProfilePicture,
			Verified:       user.Verified,
		}

		users = append(users, publicUser)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}

func UpdateUser(accountID string, firstName string, lastName string, email string) (*mongo.UpdateResult, error) {
	newData := bson.M{}

	if len(firstName) > 2 {
		newData["firstname"] = firstName
	}

	if len(lastName) > 2 {
		newData["lastname"] = lastName
	}

	if email != "" {
		newData["email"] = email
	}

	filter := bson.M{"_id": accountID}
	command := bson.M{
		"$set": newData,
	}

	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).UpdateOne(context.TODO(), filter, command)

	return result, err
}

func AddFollower(transmitterID string, targetID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": targetID}
	command := bson.M{
		"$set": bson.M{
			"followers." + transmitterID: Follow{
				AccountID: transmitterID,
				Date:      time.Now().Unix(),
			},
		},
	}

	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).UpdateOne(context.TODO(), filter, command)

	return result, err
}

func AddFollowing(transmitterID string, targetID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": transmitterID}
	command := bson.M{
		"$set": bson.M{
			"following." + targetID: Follow{
				AccountID: targetID,
				Date:      time.Now().Unix(),
			},
		},
	}

	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).UpdateOne(context.TODO(), filter, command)

	return result, err
}

func RemoveFollower(transmitterID string, targetID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": targetID}
	command := bson.M{
		"$unset": bson.M{
			"followers." + transmitterID: 1,
		},
	}

	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).UpdateOne(context.TODO(), filter, command)

	return result, err
}

func RemoveFollowing(transmitterID string, targetID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": transmitterID}
	command := bson.M{
		"$unset": bson.M{
			"following." + targetID: 1,
		},
	}

	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).UpdateOne(context.TODO(), filter, command)

	return result, err
}

func BlockAccountByID(accountID string, toBlock string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": accountID}
	command := bson.M{
		"$set": bson.M{
			"blocked_users." + toBlock: BlockedAccount{
				AccountID: toBlock,
				Date:      time.Now().Unix(),
			},
		},
	}

	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).UpdateOne(context.TODO(), filter, command)

	return result, err
}

func UnblockAccountByID(accountID string, toUnblock string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": accountID}
	command := bson.M{
		"$unset": bson.M{
			"blocked_users." + toUnblock: 1,
		},
	}

	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).UpdateOne(context.TODO(), filter, command)

	return result, err
}

func GetPeopleSuggestionsByAccountID(accountID string) ([]PublicUser, error) {
	var users []PublicUser

	user, found := GetUserByID(accountID)
	if !found {
		return nil, errors.New("User not found")
	}

	var filterOmitAlreadyFollowing []interface{}
	var filterOfUsersFollowings []interface{}
	for _, follow := range user.Following {
		filterOfUsersFollowings = append(filterOfUsersFollowings, bson.M{
			"followers." + follow.AccountID: bson.M{
				"$exists": true,
			},
		})

		filterOmitAlreadyFollowing = append(filterOmitAlreadyFollowing, bson.M{
			"_id": bson.M{
				"$ne": follow.AccountID,
			},
		})
	}

	for _, blockedAccount := range user.BlockedAccounts {
		filterOmitAlreadyFollowing = append(filterOmitAlreadyFollowing, bson.M{
			"_id": bson.M{
				"$ne": blockedAccount.AccountID,
			},
		})
	}

	filter := bson.M{
		"$or":  filterOfUsersFollowings,
		"$and": filterOmitAlreadyFollowing,
	}

	options := options.Find()
	options.SetLimit(25)
	options.SetSort(bson.D{})

	cur, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			continue
		}

		publicUser := PublicUser{
			ID:             user.ID,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Username:       user.Username,
			ProfilePicture: user.ProfilePicture,
			Verified:       user.Verified,
		}

		users = append(users, publicUser)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}

func NewSessionLog(accountID string, IPAddress string, UserAgent string) (*mongo.UpdateResult, error) {
	sessionLog := SessionLog{
		IPAddress: IPAddress,
		UserAgent: UserAgent,
		Date:      time.Now().Unix(),
	}

	command := bson.M{
		"$push": bson.M{
			"session_logs": sessionLog,
		},
	}

	filter := bson.M{"_id": accountID}

	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection("users").UpdateOne(context.TODO(), filter, command)

	return result, err
}
