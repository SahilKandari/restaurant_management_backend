package controllers

import (
	"database/sql"
	"restaurant-management/database"
)

var Db *sql.DB

func InitControllers() {
	Db = database.Client
}
