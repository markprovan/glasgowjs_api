package main

import (
	"encoding/json"
	"net/http"

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
	m.Get("/authors", AuthorsIndex)
	m.Get("/authors/:id", AuthorsShow)
	m.Post("/posts", PostsCreate)
	m.Options("/posts", PostsOptions)
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
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

func PostsCreate(req *http.Request, r render.Render, db *sqlx.DB) {
	var postJSON PostJSON
	err := json.NewDecoder(req.Body).Decode(&postJSON)
	PanicIf(err)

	post := postJSON.Post
	dbsql, err := db.Exec("insert into posts (title, body, author_id) values (?, ?, 1)", post.Title, post.Body)
	PanicIf(err)
	id, err := dbsql.LastInsertId()
	PanicIf(err)
	post.Id = id
	r.JSON(200, map[string]interface{}{"post": post})
}

func PostsOptions(r render.Render, db *sqlx.DB, res http.ResponseWriter) {
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	res.Header().Set("Access-Control-Allow-Credentials", "true")
	res.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

func PostsIndex(r render.Render, db *sqlx.DB) {
	posts := []Post{}
	err := db.Select(&posts, "SELECT * FROM posts ORDER BY id DESC")
	PanicIf(err)
	r.JSON(200, map[string]interface{}{"posts": posts})
}

func PostsShow(r render.Render, db *sqlx.DB, params martini.Params) {
	post := Post{}
	err := db.Get(&post, "SELECT * FROM posts WHERE id = $1", params["id"])
	if err != nil {
		r.JSON(404, nil)
	} else {
		r.JSON(200, map[string]interface{}{"post": post})
	}
}

func AuthorsIndex(r render.Render, db *sqlx.DB) {
	authors := []Author{}
	err := db.Select(&authors, "SELECT * FROM authors ORDER BY id DESC")
	PanicIf(err)
	r.JSON(200, map[string]interface{}{"authors": authors})
}

func AuthorsShow(r render.Render, db *sqlx.DB, params martini.Params) {
	author := Author{}
	err := db.Get(&author, "SELECT * FROM authors WHERE id = $1", params["id"])
	if err != nil {
		r.JSON(404, nil)
	} else {
		r.JSON(200, map[string]interface{}{"author": author})
	}
}

type Post struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"db:"author_id"`
}

type Author struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type AuthorJSON struct {
	Author Author `json:"author"`
}

type PostJSON struct {
	Post Post `json:"post"`
}
