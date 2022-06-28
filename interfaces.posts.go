package main

type Post struct {
	ID        string                 `json:"_id" bson:"_id"`
	AuthorID  string                 `json:"author_id" bson:"author_id"`
	Title     string                 `json:"title" bson:"title"`
	Content   string                 `json:"content" bson:"content"`
	CreatedAt int64                  `json:"created_at" bson:"created_at"`
	Likes     map[string]GenericLike `json:"likes" bson:"likes"`
}

type SavedPost struct {
	ID   string `json:"post_id" bson:"post_id"`
	Date int64  `json:"date" bson:"date"`
}
