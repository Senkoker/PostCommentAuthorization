package models

type CreatePostResponse struct {
	Id string `json:"id"`
}
type GetPostResponse struct {
	Posts []Post `json:"posts"`
}
type CreateCommentResponse struct {
	Id string `json:"id"`
}
type GetMainCommentResponse struct {
	MainComments []MainComment `json:"main_comments"`
}

type GetCommentResponse struct {
	Comments []Comment `json:"comments"`
}
