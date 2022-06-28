package main

type GenericLike struct {
	AccountID string `json:"account_id" bson:"account_id"`
	Date      int64  `json:"date" bson:"date"`
}

type Message struct {
	ExitCode int         `json:"exit_code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
}

type Session struct {
	Token     string `json:"token"`
	AccountID string `json:"account_id"`
}
