package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AddCommentReply(c *gin.Context) {
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

	replyContent, replyContentParameterIncluded := c.GetQuery("content")
	if !replyContentParameterIncluded || replyContent == "" {
		JSON(c, http.StatusBadRequest, false, "Missing 'rontent' parameter", nil)
		return
	}

	if len(replyContent) > 150 || len(replyContent) <= 0 {
		JSON(c, http.StatusInternalServerError, false, "The maximum length of a reply is 150 characters", nil)
	}

	comment, found := GetCommentByID(commentID)
	if comment.ID == "" || !found {
		JSON(c, http.StatusBadRequest, false, "Comment not found", nil)
		return
	}

	if len(comment.Replies) > 99 {
		JSON(c, http.StatusInternalServerError, false, "The maximum number of replies is 100", nil)
		return
	}

	replyID := RandomAccountID()
	commentReply := PostCommentReply{
		ID:       fmt.Sprint(replyID),
		AuthorID: accountID,
		Content:  replyContent,
		Date:     time.Now().Unix(),
	}

	added, err := AddReplyToComment(comment.ID, commentReply)
	if err != nil || !added {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Reply added succesfully.", commentReply)
}

func RemoveCommentReply(c *gin.Context) {
	_, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	commentID, commentIDParameterIncluded := c.GetQuery("comment_id")
	if !commentIDParameterIncluded || commentID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing 'comment_id' parameter", nil)
		return
	}

	replyID, replyIDParameterIncluded := c.GetQuery("reply_id")
	if !replyIDParameterIncluded || replyID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing 'reply_id' parameter", nil)
		return
	}

	comment, found := GetCommentByID(commentID)
	if comment.ID == "" || !found {
		JSON(c, http.StatusBadRequest, false, "Comment not found", nil)
		return
	}

	if comment.Replies[replyID].ID == "" {
		JSON(c, http.StatusBadRequest, false, "Reply in the comment not found", nil)
		return
	}

	added, err := RemoveReplyFromComment(comment.ID, replyID)
	if err != nil || !added {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Reply deleted succesfully.", nil)
}
