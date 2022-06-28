package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthSignUp(c *gin.Context) {
	firstName, firstNameParameterIncluded := c.GetQuery("first_name")
	lastName, lastNameParameterIncluded := c.GetQuery("last_name")
	userName, userNameParameterIncluded := c.GetQuery("username")
	emailAddress, emailAddressParameterIncluded := c.GetQuery("email_address")
	emailAddressConfirmation, emailAddressConfirmationParameterIncluded := c.GetQuery("email_address_confirmation")
	password, passwordParameterIncluded := c.GetQuery("password")
	birthdayDay, birthdayDayParameterIncluded := c.GetQuery("birthday_day")
	birthdayMonth, birthdayMonthParameterIncluded := c.GetQuery("birthday_month")
	birthdayYear, birthdayYearParameterIncluded := c.GetQuery("birthday_year")

	if !(firstNameParameterIncluded || lastNameParameterIncluded || userNameParameterIncluded || emailAddressParameterIncluded || emailAddressConfirmationParameterIncluded || passwordParameterIncluded || birthdayDayParameterIncluded || birthdayMonthParameterIncluded || birthdayYearParameterIncluded) {
		JSON(c, http.StatusBadRequest, false, "Missing required parameters", nil)
		return
	}

	if len(firstName) < 2 || len(firstName) > 32 {
		JSON(c, http.StatusBadRequest, false, "First name must be between 2 and 32 characters", nil)
		return
	}

	if len(lastName) < 2 || len(lastName) > 32 {
		JSON(c, http.StatusBadRequest, false, "Last name must be between 2 and 32 characters", nil)
		return
	}

	if len(userName) < 2 || len(userName) > 32 {
		JSON(c, http.StatusBadRequest, false, "Username must be between 2 and 32 characters", nil)
		return
	}

	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	if !usernameRegex.MatchString(userName) {
		JSON(c, http.StatusBadRequest, false, "Username must only contain letters, numbers, and underscores", nil)
		return
	}

	_, accountFoundWithUsername := GetUserByUsername(userName)
	if accountFoundWithUsername {
		JSON(c, http.StatusBadRequest, false, "Username already in use", nil)
		return
	}

	if emailAddress == "" || len(emailAddress) > 254 {
		JSON(c, http.StatusBadRequest, false, "Invalid email address", nil)
		return
	}

	_, accountFoundWithEmail := GetUserByEmailAddress(emailAddress)
	if accountFoundWithEmail {
		JSON(c, http.StatusBadRequest, false, "Email address already in use", nil)
		return
	}

	emailRegex := regexp.MustCompile(`(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)`)
	if !emailRegex.MatchString(emailAddress) {
		JSON(c, http.StatusBadRequest, false, "Invalid email address", nil)
		return
	}

	if emailAddress != emailAddressConfirmation {
		JSON(c, http.StatusBadRequest, false, "Email addresses do not match", nil)
		return
	}

	if len(password) < 4 {
		JSON(c, http.StatusBadRequest, false, "Password must be at least 4 characters", nil)
		return
	}

	birthdayDayInt, err := strconv.Atoi(birthdayDay)
	if birthdayDayInt < 1 || birthdayDayInt > 31 || err != nil {
		JSON(c, http.StatusBadRequest, false, "Invalid birthday day", nil)
		return
	}

	birthdayMonthInt, err := strconv.Atoi(birthdayMonth)
	if birthdayMonthInt < 1 || birthdayMonthInt > 12 || err != nil {
		JSON(c, http.StatusBadRequest, false, "Invalid birthday month", nil)
		return
	}

	birthdayYearInt, err := strconv.Atoi(birthdayYear)
	if birthdayYearInt < 1900 || birthdayYearInt > 2100 || err != nil {
		JSON(c, http.StatusBadRequest, false, "Invalid birthday year", nil)
		return
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Could not hash password", nil)
		return
	}

	birthdayTimestamp, err := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", birthdayYear, birthdayMonth, birthdayDay))
	if err != nil {
		fmt.Println(err)
		JSON(c, http.StatusInternalServerError, false, "Could not parse birthday. (format: 2006-01-02)", nil)
		return
	}

	accountID := RandomAccountID()
	user := User{
		ID:              fmt.Sprintf("%d", accountID),
		FirstName:       firstName,
		LastName:        lastName,
		Username:        userName,
		Email:           emailAddress,
		SignUpDate:      time.Now().Unix(),
		Password:        hashedPassword,
		Verified:        false,
		Followers:       map[string]Follow{},
		Following:       map[string]Follow{},
		Likes:           []UserLike{},
		BlockedAccounts: map[string]BlockedAccount{},
		SavedPosts:      map[string]SavedPost{},
		Birthday: Birthday{
			Timestamp: birthdayTimestamp.Unix(),
		},
		Preferences: Preferences{
			PublicLikes: false,
		},
	}

	sessionLog := SessionLog{
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Date:      time.Now().Unix(),
	}

	user.SessionLogs = append(user.SessionLogs, sessionLog)

	_, err = InsertDocument("users", user)
	if err != nil {
		JSON(c, http.StatusOK, false, "An internal error has been generated, retry later", nil)
		return
	}

	token, err := AssignToken(user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Signed up succesfully.", Map{
		"token": token,
	})
}

func AuthSignIn(c *gin.Context) {
	emailAddress, emailAddressParameterIncluded := c.GetQuery("email_address")
	password, passwordParameterIncluded := c.GetQuery("password")

	if !(emailAddressParameterIncluded && passwordParameterIncluded) {
		JSON(c, http.StatusBadRequest, false, "Missing required parameters", nil)
		return
	}

	if emailAddress == "" || len(emailAddress) > 254 {
		JSON(c, http.StatusBadRequest, false, "Invalid email address", nil)
		return
	}

	emailRegex := "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"
	if matched, err := regexp.MatchString(emailRegex, emailAddress); !matched || err != nil {
		JSON(c, http.StatusBadRequest, false, "Invalid email address", nil)
		return
	}

	if len(password) < 4 {
		JSON(c, http.StatusBadRequest, false, "Password must be at least 4 characters", nil)
		return
	}

	user, found := GetUserByEmailAddress(emailAddress)

	if !found || !CheckPasswordHash(password, user.Password) {
		JSON(c, http.StatusBadRequest, false, "Invalid email address or password", nil)
		return
	}

	_, err := NewSessionLog(user.ID, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	token, err := AssignToken(user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Signed in succesfully.", Map{
		"token": token,
	})
}

func AuthRequestUserData(c *gin.Context) {
	token := c.GetHeader("X-Auth-Token")
	if token == "" {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	session := SessionTokens[token]
	if session.Token == "" || session.AccountID == "" {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	user, found := GetUserByID(session.AccountID)
	if !found {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	publicUser := User{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Username:       user.Username,
		Email:          user.Email,
		Birthday:       user.Birthday,
		SignUpDate:     user.SignUpDate,
		Verified:       user.Verified,
		ProfilePicture: user.ProfilePicture,
		Preferences:    user.Preferences,
	}

	JSON(c, http.StatusOK, true, "Requested user data succesfully.", publicUser)
}

func AuthSignOut(c *gin.Context) {
	token := c.GetHeader("X-Auth-Token")
	if token == "" {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	session := SessionTokens[token]
	if session.Token == "" || session.AccountID == "" {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	RemoveToken(token)

	JSON(c, http.StatusOK, true, "Signed out succesfully.", nil)
}
