package models

import (
	"mime/multipart"
)

type ProfileFill struct {
	UserID     string                `json:"user_id,omitempty"`
	FirstName  string                `json:"first_name" json:"first_name,omitempty"`
	SecondName string                `json:"second_name" json:"second_name,omitempty"`
	Image      *multipart.FileHeader `json:"image_url,omitempty"`
	ImgURL     string                `json:"img_url,omitempty"`
	BirthDate  string                `json:"birth_date" json:"birth_date"`
	Education  string                `json:"education" json:"education,omitempty"`
	Country    string                `json:"country" json:"country,omitempty"`
	City       string                `json:"city" json:"city,omitempty"`
	Friends    []string              `json:"friends" json:"friends,omitempty"`
}
type ProfileResponse struct {
	Id string `json:"id"`
}
