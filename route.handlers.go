package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

var (
	API   *gin.RouterGroup
	Auth  *gin.RouterGroup
	Users *gin.RouterGroup
	Posts *gin.RouterGroup
)

func SetRouters() {
	API = App.Group("/api")

	AuthenticationController()
	UsersController()
	PostsController()
}

func AuthenticationController() {
	Auth = API.Group("/auth")
	Auth.POST("/signin", AuthSignIn)
	Auth.POST("/signup", AuthSignUp)
	Auth.POST("signout", AuthSignOut)
	Auth.GET("/request/data", AuthRequestUserData)
}

func UsersController() {
	Users = API.Group("/users")
	Users.GET("/", GetSearchUsers)

	Users.GET("/user", GetUserPublicData)
	Users.PUT("/user", TokenMiddlware(), UpdateUserInformation)

	Users.GET("/posts", UserPosts)
	Users.GET("/suggestions", TokenMiddlware(), GetPeopleSuggestions)

	Users.PUT("/follow", TokenMiddlware(), FollowUser)
	Users.DELETE("/follow", TokenMiddlware(), UnfollowUser)

	Users.PUT("/block", TokenMiddlware(), BlockAccount)
	Users.DELETE("/block", TokenMiddlware(), UnblockAccount)
}

func PostsController() {
	Posts = API.Group("/posts")

	Posts.PUT("/", TokenMiddlware(), NewPost)
	Posts.DELETE("/", TokenMiddlware(), DeletePost)

	Posts.PUT("/like", TokenMiddlware(), LikePost)
	Posts.DELETE("/like", TokenMiddlware(), UnlikePost)

	Posts.PUT("/save", TokenMiddlware(), SavePost)
	Posts.DELETE("/save", TokenMiddlware(), UnsavePost)

	Posts.GET("/comment", GetCommentsByPost)
	Posts.PUT("/comment", TokenMiddlware(), AddComment)
	Posts.DELETE("/comment", TokenMiddlware(), RemoveComment)

	Posts.PUT("/comment/like", TokenMiddlware(), LikeComment)
	Posts.DELETE("/comment/like", TokenMiddlware(), UnlikeComment)

	Posts.PUT("/comment/reply", TokenMiddlware(), AddCommentReply)
	Posts.DELETE("/comment/reply", TokenMiddlware(), RemoveCommentReply)
}

func TokenMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Auth-Token")
		if token == "" {
			Abort(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		session := SessionTokens[token]
		if session.Token == "" || session.AccountID == "" {
			Abort(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		c.Set("accountID", session.AccountID)
		c.Next()
	}
}