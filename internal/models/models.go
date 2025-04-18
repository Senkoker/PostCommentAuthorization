package models

import (
	"mime/multipart"
	"time"
)

type InterestingPost struct {
	PostsIDs []string `json:"posts_ids"`
	Hashtags string   `json:"hashtags"`
}
type NewPost struct {
	AuthorId string
	Hashtags []string
	Content  string
	Img      *multipart.FileHeader
	ImgUrl   string
	Private  bool
}

type Post struct {
	PostID       string    `json:"post_id"`
	ImgPersonURL string    `json:"img_person_url"`
	Author       string    `json:"author"`
	AuthorId     string    `json:"author_id"`
	TagsIds      []string  `json:"tags_ids"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	Watched      int       `json:"watched"`
	Likes        int       `json:"likes"`
}

type NewComment struct {
	ReplyTo   string    `json:"repy_to"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Main      bool      `json:"main"`
}
type UserInfo struct {
	FirstName  string
	SecondName string
	ImgUrl     string
}

type MainComment struct {
	PostID       string    `json:"post_id"`
	AuthorID     string    `json:"author_id"`
	AuthorName   string    `json:"author"`
	AuthorImgUrl string    `json:"author_img"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
}

type Comment struct {
	MainCommentID string    `json:"main_comment_id"`
	AuthorID      string    `json:"author_id"`
	AuthorName    string    `json:"author"`
	AuthorImgUrl  string    `json:"author_img"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
}
