package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/vcraescu/go-paginator/v2"
	"github.com/vcraescu/go-paginator/v2/adapter"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Board struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Author    string
	Content   string
}

var (
	tpl    *template.Template
	gormDB *gorm.DB
)

func init() {
	tpl = template.Must(template.ParseGlob("web/templates/*.gohtml"))
}

func main() {

	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True", user, password, host, database)
	mysqlDB, err := sql.Open("mysql", connectionString)
	defer mysqlDB.Close()

	gormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: mysqlDB,
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	gormDB.AutoMigrate(&Board{})

	http.HandleFunc("/", index)
	http.HandleFunc("/write", write)
	http.HandleFunc("/board/", board)
	http.HandleFunc("/post/", post)
	http.HandleFunc("/test", test)
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	fmt.Println("Listening ... !")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func test(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/board?target=title&v=writing", http.StatusAccepted)
	return

}

func write(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		title := r.PostFormValue("title")
		author := r.PostFormValue("author")
		content := r.PostFormValue("content")

		newPost := Board{Title: title, Author: author, Content: content}
		gormDB.Create(&newPost)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		// http.Redirect(w, r, "/", http.StatusCreated)

		return
	}

	tpl.ExecuteTemplate(w, "write.gohtml", nil)
}

func board(w http.ResponseWriter, r *http.Request) {
	var b []Board

	page := r.FormValue("page")
	if page == "" {
		page = "1"
	}
	pageInt, _ := strconv.Atoi(page)

	if keyword := r.FormValue("v"); keyword != "" {
		target := r.FormValue("target")

		switch target {
		case "title":
			q := gormDB.Where("title LIKE ?", fmt.Sprintf("%%%s%%", keyword)).Find(&b)
			pg := paginator.New(adapter.NewGORMAdapter(q), 3)
			pg.SetPage(pageInt)

			if err := pg.Results(&b); err != nil {
				panic(err)
			}

			tpl.ExecuteTemplate(w, "board.gohtml", b)
			return
		case "author":
			gormDB.Where("author LIKE ?", fmt.Sprintf("%%%s%%", keyword)).Find(&b)
			tpl.ExecuteTemplate(w, "board.gohtml", b)
			return
		}

	}

	q := gormDB.Order("id desc").Limit(10).Offset(0).Find(&b)
	pg := paginator.New(adapter.NewGORMAdapter(q), 3)

	pg.SetPage(pageInt)

	if err := pg.Results(&b); err != nil {
		panic(err)
	}

	fmt.Println("page:", pageInt, "\nboards:", b)
	// gormDB.Limit(10).Offset(0).Find(&b)
	// gormDB.Select("id", "title", "author").Find(&b)
	// gormDB.Find(&b)
	// gormDB.First(&b)

	tpl.ExecuteTemplate(w, "board.gohtml", b)
}

func post(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	var b Board
	gormDB.First(&b, id)
	// gormDB.First(&b, "id = ?", id)

	tpl.ExecuteTemplate(w, "post.gohtml", b)
}
