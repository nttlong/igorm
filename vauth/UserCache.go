package unvsauth

import "unvs-auth/models"

type UserCache interface {
	Get(identifier string) (*models.User, bool)
	Set(user *models.User)
	Delete(identifier string)
}
