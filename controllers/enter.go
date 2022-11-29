package controllers

import (
	"enian_blog/controllers/adminApi"
	"enian_blog/controllers/profileApi"
	"enian_blog/controllers/views"
)

type ControllersGroup struct {
	View       views.View
	AdminApi   adminApi.AdminApi
	ProfileApi profileApi.ProfileApi
}

var Controllers = ControllersGroup{}
