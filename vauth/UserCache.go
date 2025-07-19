package vauth

import "vauth/models"

type UserCache interface {
	Get(identifier string) (*models.User, bool)
	Set(user *models.User)
	Delete(identifier string)
}
