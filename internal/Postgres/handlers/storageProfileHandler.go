package handlers

import (
	"VK_posts/internal/models"
	"VK_posts/pkg/Postgres"
	"context"
	"fmt"
)

type StorageProfileHandler struct {
	storage *Postgres.Storage
}

func NewProfileHandler(storage *Postgres.Storage) *StorageProfileHandler {
	return &StorageProfileHandler{storage}
}

func (s *StorageProfileHandler) PGFillUserProfile(profile models.ProfileFill) (string, error) {
	fmt.Println(profile)
	stmt, err := s.storage.Db.Prepare("INSERT INTO users_info(user_id,first_name,second_name,img_url,birth_date,education,country,city) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) returning id")
	defer stmt.Close()
	if err != nil {
		return "", fmt.Errorf("Problem preparing storage FillUserPRofile: %w", err)
	}
	var id string
	ctx := context.Background()
	row := stmt.QueryRowContext(ctx, profile.UserID, profile.FirstName, profile.SecondName, profile.ImgURL, profile.BirthDate, profile.Education, profile.Country, profile.City)
	err = row.Err()
	if err != nil {
		return "", fmt.Errorf("Problem to executedata FillUserPRofile: %w", err)
	}
	err = row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("Problem to scan FillUserPRofile: %w", err)
	}
	return id, nil
}

func (s *StorageProfileHandler) PGGetUserProfile(username string) (models.ProfileFill, error) {
	stmt, err := s.storage.Db.Prepare("SELECT user_id,first_name,second_name,img_url,birth_date,education,country,city FROM users_info where user_id=$1")
	defer stmt.Close()
	if err != nil {
		return models.ProfileFill{}, fmt.Errorf("Problem preparing storage PGGetUserPRofile: %w", err)
	}
	ctx := context.Background()
	row := stmt.QueryRowContext(ctx, username)
	err = row.Err()
	if err != nil {
		return models.ProfileFill{}, fmt.Errorf("Problem to return data PGGetUserPRofile: %w", err)
	}
	var profile models.ProfileFill
	err = row.Scan(&profile.UserID, &profile.FirstName, &profile.SecondName, &profile.ImgURL, &profile.BirthDate, &profile.Education, &profile.Country, &profile.City)
	if err != nil {
		return models.ProfileFill{}, fmt.Errorf("Problem to scan PGGetUserPRofile: %w", err)
	}
	return profile, nil

}
