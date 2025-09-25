package repository

import "tasklist/db"

type Repository struct {
	UserRepo User
	TaskRepo Task
	database *db.Database
}

func NewRepository(database *db.Database) *Repository {
	return &Repository{
		UserRepo: NewUserRepo(database),
		TaskRepo: NewTaskRepo(database),
	}
}
