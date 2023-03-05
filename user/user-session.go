package user

import (
	"controller"
	"web/view"
)

type UserSessionFactory struct {
	controller.Persistance
}

type UserSession struct {
	*controller.FlashcardsController
	*view.HttpView
}

type UserContext struct {
	Id    string
	Email string
}

func (u *UserSessionFactory) Create(userContext UserContext) *UserSession {
	store := u.Persistance.Create("db", userContext.Id)
	view := view.CreateHttpView(userContext.Email)

	return &UserSession{
		FlashcardsController: controller.CreateFlashcardsController(view, store),
		HttpView:             view,
	}
}
