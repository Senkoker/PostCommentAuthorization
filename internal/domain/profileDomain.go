package domain

import (
	"VK_posts/internal/logger"
	"VK_posts/internal/models"
)

type ProfileDomain struct {
	PostgresProfileInterface PostgresProfileInterface
	SelectelS3Interface      SelectelS3Interface
}

type PostgresProfileInterface interface {
	PGFillUserProfile(profile models.ProfileFill) (string, error)
	PGGetUserProfile(username string) (models.ProfileFill, error)
}

func NewProfileDomain(postgresProfileInterface PostgresProfileInterface, s3Interface SelectelS3Interface) *ProfileDomain {
	return &ProfileDomain{PostgresProfileInterface: postgresProfileInterface, SelectelS3Interface: s3Interface}
}
func (d *ProfileDomain) FillUserProfile(profile models.ProfileFill) (string, error) {
	op := "fillUserProfile"
	localLogger := logger.GetLogger().With("op", op)
	img_url, err := d.SelectelS3Interface.SendImage(profile.Image)
	if err != nil {
		localLogger.Error("Failed to send image to selectel", "err", err)
		img_url = "some_url"
	}
	profile.ImgURL = img_url
	id, err := d.PostgresProfileInterface.PGFillUserProfile(profile)
	if err != nil {
		localLogger.Error("pgFillUserProfile", "err", err.Error())
		return "", err
	}
	return id, nil
}
func (d *ProfileDomain) GetUserProfile(username string) (models.ProfileFill, error) {
	op := "getUserProfile"
	localLogger := logger.GetLogger().With("op", op)
	profile, err := d.PostgresProfileInterface.PGGetUserProfile(username)
	if err != nil {
		localLogger.Error("pgGetUserProfile", "err", err.Error())
		return models.ProfileFill{}, err
	}
	return profile, nil
}
