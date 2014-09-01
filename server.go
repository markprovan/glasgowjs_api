package main

import (
	"github.com/codegangsta/martini-contrib/cors"
	"github.com/go-martini/martini"
	"github.com/jmoiron/sqlx"
	"github.com/martini-contrib/render"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Map(SetupDB())
	m.Get("/posts", PostsIndex)
	m.Get("/posts/:id", PostsShow)
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	m.Run()
}

func SetupDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "./blog.db")
	PanicIf(err)
	return db
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func PostsIndex(r render.Render, db *sqlx.DB) {
	posts := []Post{}
	err := db.Select(&posts, "SELECT * FROM posts")
	PanicIf(err)
	r.JSON(200, map[string]interface{}{"posts": posts})
}

func PostsShow(r render.Render, db *sqlx.DB, params martini.Params) {
	post := Post{}
	err := db.Get(&post, "SELECT * FROM posts WHERE id = $1", params["id"])
	if err != nil {
		r.JSON(404, nil)
	} else {
		r.JSON(200, post)
	}
}

type Post struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"db:"author_id"`
}
