package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CommentPost(c *gin.Context) {
	postID, postIDParameterIncluded := c.GetQuery("post_id")
	if !postIDParameterIncluded {
		JSON(c, http.StatusBadRequest, false, "Missing post_id parameter", nil)
		return
	}

	if postID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing post_id parameter", nil)
		return
	}

	commentContent, commentParameterIncluded := c.GetQuery("comment_content")
	if !commentParameterIncluded {
		JSON(c, http.StatusBadRequest, false, "Missing comment parameter", nil)
		return
	}

	if commentContent == "" {
		JSON(c, http.StatusBadRequest, false, "Missing comment parameter", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Comment added succesfully.", nil)
}

func AddComment(c *gin.Context) {
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

	commentContent, commentContentParameterIncluded := c.GetQuery("content")
	if !commentContentParameterIncluded {
		JSON(c, http.StatusBadRequest, false, "Missing content parameter", nil)
		return
	}

	if commentContent == "" {
		JSON(c, http.StatusBadRequest, false, "Missing content parameter", nil)
		return
	}

	if len(commentContent) > 150 {
		JSON(c, http.StatusBadRequest, false, "Comment content too long", nil)
		return
	}

	post, found := GetPostByID(postID)
	if post.ID == "" || !found {
		JSON(c, http.StatusBadRequest, false, "Post not found", nil)
		return
	}

	commentID := RandomAccountID()
	comment := PostComment{
		ID:       fmt.Sprint(commentID),
		PostID:   postID,
		AuthorID: accountID,
		Content:  commentContent,
		Date:     time.Now().Unix(),
		Likes:    map[string]GenericLike{},
		Replies:  map[string]PostCommentReply{},
	}

	_, err := InsertDocument(CommentsCollection, comment)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Comment created succesfully.", comment)
}

func RemoveComment(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	commentID, commentIDParameterIncluded := c.GetQuery("comment_id")
	if !commentIDParameterIncluded {
		JSON(c, http.StatusBadRequest, false, "Missing comment_id parameter", nil)
		return
	}

	if commentID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing comment_id parameter", nil)
		return
	}

	comment, found := GetCommentByID(commentID)
	if !found {
		JSON(c, http.StatusBadRequest, false, "Comment not found", nil)
		return
	}

	if comment.AuthorID != accountID {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	deleted, err := DeleteCommentByID(commentID)
	if err != nil || !deleted {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Comment deleted succesfully.", nil)
}

func GetCommentsByPost(c *gin.Context) {
	postID, postIDParameterIncluded := c.GetQuery("post_id")
	if !postIDParameterIncluded || postID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing 'post_id' parameter", nil)
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
		JSON(c, http.StatusInternalServerError, false, "The minimum range to fetch comments is 0", nil)
		return
	}

	if limitInt > 1000 {
		JSON(c, http.StatusInternalServerError, false, "The maximum range to fetch comments is 1k", nil)
		return
	}

	var order bson.M
	switch orderBy {
	case "newest":
		order = bson.M{"date": -1}
	case "oldest":
		order = bson.M{"date": 1}
	case "none":
		order = bson.M{}
	case "":
		order = bson.M{}
	default:
		JSON(c, http.StatusInternalServerError, false, "Query 'order_by' is not valid", nil)
		return
	}

	posts, err := GetCommentsByPostID(postID, int64(limitInt), order)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Comments fetched succesfully.", posts)
}

func LikeComment(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	commentID, commentIDParameterIncluded := c.GetQuery("comment_id")
	if !commentIDParameterIncluded || commentID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing 'comment_id' parameter", nil)
		return
	}

	comment, found := GetCommentByID(commentID)
	if comment.ID == "" || !found {
		JSON(c, http.StatusBadRequest, false, "Comment not found", nil)
		return
	}

	if comment.Likes[accountID].AccountID != "" {
		JSON(c, http.StatusBadRequest, false, "You already liked this comment", nil)
		return
	}

	liked, err := LikeCommentByID(accountID, commentID)
	if err != nil || !liked {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Comment liked succesfully.", nil)
}

func UnlikeComment(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)


	commentID, commentIDParameterIncluded := c.GetQuery("comment_id")
	if !commentIDParameterIncluded || commentID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing 'comment_id' parameter", nil)
		return
	}

	comment, found := GetCommentByID(commentID)
	if comment.ID == "" || !found {
		JSON(c, http.StatusBadRequest, false, "Comment not found", nil)
		return
	}

	if comment.Likes[accountID].AccountID == "" {
		JSON(c, http.StatusBadRequest, false, "You haven't liked this comment", nil)
		return
	}

	liked, err := UnlikeCommentByID(accountID, commentID)
	if err != nil || !liked {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Comment unliked succesfully.", nil)
}
