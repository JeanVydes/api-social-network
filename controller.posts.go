package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// TODO

// Add Attachment support
func NewPost(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	title, titleParameterIncluded := c.GetQuery("title")
	if !titleParameterIncluded {
		JSON(c, http.StatusBadRequest, false, "Missing title parameter", nil)
		return
	}

	content, contentParameterIncluded := c.GetQuery("content")
	if !contentParameterIncluded {
		JSON(c, http.StatusBadRequest, false, "Missing content parameter", nil)
		return
	}

	if title == "" || content == "" {
		JSON(c, http.StatusBadRequest, false, "Missing title or content", nil)
		return
	}

	if len(title) > 100 {
		JSON(c, http.StatusBadRequest, false, "Title must be less than 100 characters", nil)
		return
	}

	if len(content) > 1024 {
		JSON(c, http.StatusBadRequest, false, "Content is too long", nil)
		return
	}

	randomNumber := RandomAccountID()
	postID, err := GenerateToken(fmt.Sprint(randomNumber))

	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	post := Post{
		ID:        postID,
		AuthorID:  accountID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().Unix(),
		Likes:     map[string]GenericLike{},
	}

	_, err = InsertDocument(PostsCollection, post)
	if err != nil {
		fmt.Println(err)
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Post created succesfully.", post)
}

func DeletePost(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	postID, postIDParameterIncluded := c.GetQuery("post_id")
	if !postIDParameterIncluded {
		JSON(c, http.StatusBadRequest, false, "Missing post_id parameter", nil)
		return
	}

	if postID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing post_id parameter", nil)
		return
	}

	post, found := GetPostByID(postID)
	if !found {
		JSON(c, http.StatusBadRequest, false, "Post not found", nil)
		return
	}

	if post.AuthorID != accountID {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Not Permission", nil)
		return
	}

	deleted, err := DeletePostByID(postID)
	if err != nil || !deleted {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	deleted, err = DeleteAllCommentsFromPost(postID)
	if err != nil || !deleted {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Post deleted succesfully.", nil)
}

func UserPosts(c *gin.Context) {
	accountID, accountIDParameterIncluded := c.GetQuery("author_id")
	if !accountIDParameterIncluded || accountID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing 'author_id' parameter", nil)
		return
	}

	limit, limitParameterProvided := c.GetQuery("limit")
	if !limitParameterProvided {
		JSON(c, http.StatusInternalServerError, false, "Query 'limit' not provided", nil)
		return
	}

	orderBy, orderByParameterProvided := c.GetQuery("order_by")
	if !orderByParameterProvided {
		JSON(c, http.StatusInternalServerError, false, "Query 'order_by not provided", nil)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Query 'limit' is not a number", nil)
		return
	}

	if limitInt < 0 {
		JSON(c, http.StatusInternalServerError, false, "The minimum range to fetch posts is 0", nil)
		return
	}

	if limitInt > 1000 {
		JSON(c, http.StatusInternalServerError, false, "The maximum range to fetch posts is 1k", nil)
		return
	}

	var order bson.M
	switch orderBy {
	case "newest":
		order = bson.M{"created_at": -1}
	case "oldest":
		order = bson.M{"created_at": 1}
	case "none":
		order = bson.M{}
	case "":
		order = bson.M{}
	default:
		JSON(c, http.StatusInternalServerError, false, "Query 'order_by' is not valid", nil)
		return
	}

	_, found := GetUserByID(accountID)
	if !found {
		JSON(c, http.StatusNotFound, false, "User not found", nil)
		return
	}

	posts, err := GetUserPosts(accountID, int64(limitInt), order)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Posts fetched succesfully.", posts)
}

func LikePost(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	postID, postIDParameterIncluded := c.GetQuery("post_id")
	if !postIDParameterIncluded {
		JSON(c, http.StatusBadRequest, false, "Missing post_id parameter", nil)
		return
	}

	if postID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing post_id parameter", nil)
		return
	}

	post, found := GetPostByID(postID)
	if !found {
		JSON(c, http.StatusBadRequest, false, "Post not found", nil)
		return
	}

	if post.Likes[accountID].AccountID != "" {
		JSON(c, http.StatusBadRequest, false, "You already liked this post", nil)
		return
	}

	_, err := AddLikeToPost(postID, accountID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Post liked succesfully.", nil)
}

func UnlikePost(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	postID, postIDParameterIncluded := c.GetQuery("post_id")
	if !postIDParameterIncluded || postID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing post_id parameter", nil)
		return
	}

	post, found := GetPostByID(postID)
	if !found {
		JSON(c, http.StatusBadRequest, false, "Post not found", nil)
		return
	}

	if post.Likes[accountID].AccountID == "" {
		JSON(c, http.StatusBadRequest, false, "You haven't liked this post", nil)
		return
	}

	_, err := RemoveLikeFromPost(postID, accountID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Post unliked succesfully.", nil)
}

func SavePost(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	postID, postIDParameterIncluded := c.GetQuery("post_id")
	if !postIDParameterIncluded || postID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing post_id parameter", nil)
		return
	}

	post, found := GetPostByID(postID)
	if !found {
		JSON(c, http.StatusBadRequest, false, "Post not found", nil)
		return
	}

	user, found := GetUserByID(accountID)
	if !found {
		JSON(c, http.StatusNotFound, false, "User not found", nil)
		return
	}

	if user.SavedPosts[postID].ID != "" {
		JSON(c, http.StatusBadRequest, false, "You already saved this post", nil)
		return
	}

	_, err := SavePostByID(accountID, post.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Post saved succesfully.", nil)
}

func UnsavePost(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	postID, postIDParameterIncluded := c.GetQuery("post_id")
	if !postIDParameterIncluded || postID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing post_id parameter", nil)
		return
	}

	post, found := GetPostByID(postID)
	if !found {
		JSON(c, http.StatusBadRequest, false, "Post not found", nil)
		return
	}

	user, found := GetUserByID(accountID)
	if !found {
		JSON(c, http.StatusNotFound, false, "User not found", nil)
		return
	}

	if user.SavedPosts[postID].ID != "" {
		JSON(c, http.StatusBadRequest, false, "You haven't saved this post", nil)
		return
	}

	_, err := DeleteSavedPostByID(accountID, post.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Post unsaved succesfully.", nil)
}