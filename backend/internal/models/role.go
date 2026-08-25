package models

// Role represents the access level assigned to a User.
type Role string

const (
	// RoleAdmin grants full access to all resources.
	RoleAdmin Role = "admin"
	// RoleLogistics grants access to manage items, loan requests, and organisations.
	RoleLogistics Role = "logistics"
	// RoleStudent grants access to browse items and manage the user's own loan requests.
	RoleStudent Role = "student"
)
