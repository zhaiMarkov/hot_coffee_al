package usecase

import "hot-coffee/internal/dal"

type Application struct {
	Repository dal.DataRepository
}

func NewApplication(repoObject dal.DataRepository) *Application {
	return &Application{Repository: repoObject}
}
