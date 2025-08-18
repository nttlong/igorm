package controllers

import "wx"

type UserClaimsHandler struct {
	// The UserClaimsHandler handles user claims related operations.
	// It is responsible for validating, parsing, and managing user claims.
	User wx.UserClaims
}
