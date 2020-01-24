package repository

import (
	"database/sql"
	"libretaxi/objects"
	"log"
)

type Repository struct {
	db *sql.DB
}

func (repo *Repository) FindUser(userId int64) *objects.User {
	user := &objects.User{}

	rows, err := repo.db.Query(`select "userId", "menuId" from users where "userId"=$1 limit 1`, userId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cnt := 0
	for rows.Next() {
		cnt++
		rows.Scan(&user.UserId, &user.MenuId)
	}

	if cnt == 0 {
		return nil
	}

	return user
}

func (repo *Repository) SaveUser(user *objects.User) {
	// Upsert syntax: https://stackoverflow.com/questions/1109061/insert-on-duplicate-update-in-postgresql
	// Geo populate syntax: https://gis.stackexchange.com/questions/145007/creating-geometry-from-lat-lon-in-table-using-postgis/145009
	_, err := repo.db.Query(`INSERT INTO users ("userId", "menuId")
		VALUES ($1, $2)
		ON CONFLICT ("userId") DO UPDATE
		  SET "menuId" = $2`, user.UserId, user.MenuId)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Saved user")
	}
}

func NewRepository(db *sql.DB) *Repository {
	repo := &Repository{db: db}
	return repo
}