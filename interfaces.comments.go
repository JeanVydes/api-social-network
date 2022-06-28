package main

type PostCommentReply struct {
	ID        string `json:"id" bson:"id"`
	CommentID string `json:"comment_id" bson:"comment_id"`
	AuthorID  string `json:"author_id" bson:"author_id"`
	Content   string `json:"content" bson:"content"`
	Date      int64  `json:"date" bson:"date"`
}

type PostComment struct {
	ID       string                      `json:"_id" bson:"_id"`
	PostID   string                      `json:"post_id" bson:"post_id"`
	AuthorID string                      `json:"author_id" bson:"author_id"`
	Content  string                      `json:"content" bson:"content"`
	Date     int64                       `json:"date" bson:"date"`
	Likes    map[string]GenericLike      `json:"likes" bson:"likes"`
	Replies  map[string]PostCommentReply `json:"replies" bson:"replies"`
}
