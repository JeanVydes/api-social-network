package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserPublicData(c *gin.Context) {
	accountID, accountIDParameterIncluded := c.GetQuery("id")
	username, usernameParameterIncluded := c.GetQuery("username")

	if !accountIDParameterIncluded && !usernameParameterIncluded {
		JSON(c, http.StatusBadRequest, false, "Missing id or username parameter", nil)
		return
	}

	if username == "" && accountID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing id or username parameter", nil)
		return
	}

	var user User
	found := false

	if username != "" && accountID == "" {
		user, found = GetUserByID(accountID)
	} else if username == "" && accountID != "" {
		user, found = GetUserByID(accountID)
	} else if username != "" && accountID != "" {
		user, found = GetUserByID(accountID)
		if !found {
			user, found = GetUserByUsername(username)
		}
	} else {
		JSON(c, http.StatusBadRequest, false, "Missing id or username parameter", nil)
		return
	}

	if !found {
		JSON(c, http.StatusNotFound, false, "User not found", nil)
		return
	}

	publicUser := PublicUser{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Username:       user.Username,
		Birthday:       user.Birthday,
		Followers:      len(user.Followers),
		Following:      len(user.Following),
		ProfilePicture: user.ProfilePicture,
		Verified:       user.Verified,
	}

	JSON(c, http.StatusOK, true, "Users has been retrieved", publicUser)
}

func UpdateUserInformation(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	firstName, firstNameParameterIncluded := c.GetQuery("first_name")
	lastName, lastNameParameterIncluded := c.GetQuery("last_name")
	email, emailParameterIncluded := c.GetQuery("email")

	if firstNameParameterIncluded && firstName != "" && len(firstName) > 100 {
		JSON(c, http.StatusBadRequest, false, "First name must be less than 100 characters", nil)
		return
	}

	if lastNameParameterIncluded && lastName != "" && len(lastName) > 100 {
		JSON(c, http.StatusBadRequest, false, "Last name must be less than 100 characters", nil)
		return
	}

	if emailParameterIncluded && email != "" && len(email) > 100 {
		JSON(c, http.StatusBadRequest, false, "Email must be less than 100 characters", nil)
		return
	}

	if emailParameterIncluded && email != "" && !isAValidEmail(email) {
		JSON(c, http.StatusBadRequest, false, "Email is not valid", nil)
		return
	}

	userUsingEmail, found := GetUserByEmailAddress(email)
	if emailParameterIncluded && email != "" && (found || userUsingEmail.ID != "") {
		JSON(c, http.StatusBadRequest, false, "Email is already used", nil)
		return
	}

	_, err = UpdateUser(accountID, firstName, lastName, email)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Profile information has been updated", nil)
}

func FollowUser(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	userID, accountIDParameterIncluded := c.GetQuery("id")
	if !accountIDParameterIncluded || userID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing id parameter", nil)
		return
	}

	user, found := GetUserByID(userID)
	if !found {
		JSON(c, http.StatusNotFound, false, "User not found", nil)
		return
	}

	if user.ID == accountID {
		JSON(c, http.StatusBadRequest, false, "You can't follow yourself", nil)
		return
	}

	if user.Followers[accountID].AccountID == accountID {
		JSON(c, http.StatusBadRequest, false, "You already follow this user", nil)
		return
	}

	_, err := AddFollower(accountID, user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Error while following user", nil)
		return
	}

	_, err = AddFollowing(accountID, user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Error while following user", nil)
		return
	}

	JSON(c, http.StatusOK, true, "User has been followed", nil)
}

func UnfollowUser(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	userID, accountIDParameterIncluded := c.GetQuery("id")
	if !accountIDParameterIncluded || userID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing id parameter", nil)
		return
	}

	user, found := GetUserByID(userID)
	if !found {
		JSON(c, http.StatusNotFound, false, "User not found", nil)
		return
	}

	if user.ID == accountID {
		JSON(c, http.StatusBadRequest, false, "You can't unfollow yourself", nil)
		return
	}

	if user.Followers[accountID].AccountID == "" {
		JSON(c, http.StatusBadRequest, false, "You are not following this user", nil)
		return
	}

	_, err := RemoveFollower(accountID, user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Error while unfollowing user", nil)
		return
	}

	_, err = RemoveFollowing(accountID, user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Error while unfollowing user", nil)
		return
	}

	JSON(c, http.StatusOK, true, "User has been unfollowed", nil)
}

func GetSearchUsers(c *gin.Context) {
	searchText, searchTextParameterIncluded := c.GetQuery("search_text")
	if !searchTextParameterIncluded {
		JSON(c, http.StatusNotFound, false, "Search text not provided", nil)
		return
	}

	users, err := SearchUsers(searchText)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Error while trying to find users", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Users has been retrieved", users)
}

func BlockAccount(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	userID, accountIDParameterIncluded := c.GetQuery("id")
	if !accountIDParameterIncluded || userID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing id parameter", nil)
		return
	}

	user, found := GetUserByID(userID)
	if !found {
		JSON(c, http.StatusNotFound, false, "User not found", nil)
		return
	}

	me, found := GetUserByID(accountID)
	if !found {
		JSON(c, http.StatusNotFound, false, "You haven't found, sign in again", nil)
		return
	}

	if user.ID == me.ID {
		JSON(c, http.StatusBadRequest, false, "You can't block yourself", nil)
		return
	}

	if me.BlockedAccounts[user.ID].AccountID == user.ID {
		JSON(c, http.StatusBadRequest, false, "You already blocked this user", nil)
		return
	}

	_, err := BlockAccountByID(me.ID, user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Error while blocking user", nil)
		return
	}

	_, err = RemoveFollower(me.ID, user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Error while unfollowing user", nil)
		return
	}

	_, err = RemoveFollowing(me.ID, user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Error while unfollowing user", nil)
		return
	}

	JSON(c, http.StatusOK, true, "User has been blocked", nil)
}

func UnblockAccount(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	userID, accountIDParameterIncluded := c.GetQuery("id")
	if !accountIDParameterIncluded || userID == "" {
		JSON(c, http.StatusBadRequest, false, "Missing id parameter", nil)
		return
	}

	user, found := GetUserByID(userID)
	if !found {
		JSON(c, http.StatusNotFound, false, "User not found", nil)
		return
	}

	me, found := GetUserByID(accountID)
	if !found {
		JSON(c, http.StatusNotFound, false, "You haven't found, sign in again", nil)
		return
	}

	if user.ID == me.ID {
		JSON(c, http.StatusBadRequest, false, "You can't unblock yourself", nil)
		return
	}

	if me.BlockedAccounts[user.ID].AccountID != user.ID {
		JSON(c, http.StatusBadRequest, false, "You haven't blocked this user", nil)
		return
	}

	_, err := UnblockAccountByID(me.ID, user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Error while unblocking user", nil)
		return
	}

	JSON(c, http.StatusOK, true, "User has been unblocked", nil)
}

// This need a super hyper mega fix to link people with number phone
// suggest that people and watch the history of following
// too people that are followed by the users that X person follow
func GetPeopleSuggestions(c *gin.Context) {
	accountIDInterface, exists := c.Get("accountID")
	if !exists {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized: Header", nil)
		return
	}

	accountID := accountIDInterface.(string)

	users, err := GetPeopleSuggestionsByAccountID(accountID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Error while trying to find people", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Users has been retrieved", users)
}
