package main

type User struct {
	ID              string                    `json:"id" bson:"_id"`
	FirstName       string                    `json:"first_name" bson:"first_name"`
	LastName        string                    `json:"last_name" bson:"last_name"`
	Username        string                    `json:"username" bson:"username"`
	Email           string                    `json:"email" bson:"email"`
	Birthday        Birthday                  `json:"birthday" bson:"birthday"`
	SignUpDate      int64                     `json:"sign_up_date" bson:"sign_up_date"`
	Password        []byte                    `json:"password" bson:"password"`
	Verified        bool                      `json:"verified" bson:"verified"`
	ProfilePicture  []byte                    `json:"profile_picture" bson:"profile_picture"`
	Followers       map[string]Follow         `json:"followers" bson:"followers"`
	Following       map[string]Follow         `json:"following" bson:"following"`
	Likes           []UserLike                `json:"likes" bson:"likes"`
	Preferences     Preferences               `json:"preferences" bson:"preferences"`
	SessionLogs     []SessionLog              `json:"session_logs" bson:"session_logs"`
	BlockedAccounts map[string]BlockedAccount `json:"blocked_users" bson:"blocked_users"`
	SavedPosts      map[string]SavedPost      `json:"saved_posts" bson:"saved_posts"`
}

type PublicUser struct {
	ID             string   `json:"id" bson:"_id"`
	FirstName      string   `json:"first_name" bson:"first_name"`
	LastName       string   `json:"last_name" bson:"last_name"`
	Username       string   `json:"username" bson:"username"`
	Birthday       Birthday `json:"birthday" bson:"birthday"`
	Posts          int      `json:"posts" bson:"posts"`
	Followers      int      `json:"followers" bson:"followers"`
	Following      int      `json:"following" bson:"following"`
	ProfilePicture []byte   `json:"profile_picture" bson:"profile_picture"`
	Verified       bool     `json:"verified" bson:"verified"`
}

type Follow struct {
	AccountID string `json:"account_id" bson:"account_id"`
	Date      int64  `json:"date" bson:"date"`
}

type Birthday struct {
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
}

type Preferences struct {
	PublicLikes bool `json:"public_likes" bson:"public_likes"`
}

type SessionLog struct {
	IPAddress string `json:"ip_address" bson:"ip_address"`
	UserAgent string `json:"user_agent" bson:"user_agent"`
	Date      int64  `json:"date" bson:"date"`
}

type UserLike struct {
	PostID string `json:"post_id" bson:"post_id"`
	Date   int64  `json:"date" bson:"date"`
}

type BlockedAccount struct {
	AccountID string `json:"account_id" bson:"account_id"`
	Date      int64  `json:"date" bson:"date"`
}
