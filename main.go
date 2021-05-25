package main

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

type PassedData struct {
	PostData []Board
	Target   string
	Value    string
	PageList []string
	Page     string
}

var (
	tpl    *template.Template
	gormDB *gorm.DB
	//go:embed web
	staticContent embed.FS
)

const (
	MaxPerPage = 5
)

func init() {
	// tpl = template.Must(template.ParseGlob("web/templates/*.gohtml"))
	tpl = template.Must(template.ParseFS(staticContent, "web/templates/*"))
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
	http.HandleFunc("/delete/", delete)
	// http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.Handle("/web/", http.FileServer(http.FS(staticContent)))

	fmt.Println("Listening ... !")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/delete/")
	gormDB.Delete(&Board{}, id)

	http.Redirect(w, r, "/board", http.StatusSeeOther)
}

func write(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		title := r.PostFormValue("title")
		author := r.PostFormValue("author")
		content := r.PostFormValue("content")

		newPost := Board{Title: title, Author: author, Content: content}
		gormDB.Create(&newPost)

		http.Redirect(w, r, "/", http.StatusSeeOther)

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
			pg := paginator.New(adapter.NewGORMAdapter(q), MaxPerPage)
			pg.SetPage(pageInt)

			if err := pg.Results(&b); err != nil {
				panic(err)
			}
			pgNums, _ := pg.PageNums()
			pageSlice := getPageList(page, pgNums)

			temp := PassedData{
				PostData: b,
				Target:   target,
				Value:    keyword,
				PageList: pageSlice,
				Page:     page,
			}

			tpl.ExecuteTemplate(w, "board.gohtml", temp)
			return
		case "author":
			q := gormDB.Where("author LIKE ?", fmt.Sprintf("%%%s%%", keyword)).Find(&b)
			pg := paginator.New(adapter.NewGORMAdapter(q), MaxPerPage)
			pg.SetPage(pageInt)

			if err := pg.Results(&b); err != nil {
				panic(err)
			}
			pgNums, _ := pg.PageNums()
			pageSlice := getPageList(page, pgNums)

			temp := PassedData{
				PostData: b,
				Target:   target,
				Value:    keyword,
				PageList: pageSlice,
				Page:     page,
			}

			tpl.ExecuteTemplate(w, "board.gohtml", temp)
			return
		}
	}

	q := gormDB.Order("id desc").Find(&b)
	pg := paginator.New(adapter.NewGORMAdapter(q), MaxPerPage)

	pg.SetPage(pageInt)

	if err := pg.Results(&b); err != nil {
		panic(err)
	}

	pgNums, _ := pg.PageNums()
	pageSlice := getPageList(page, pgNums)

	temp := PassedData{
		PostData: b,
		PageList: pageSlice,
		Page:     page,
	}

	tpl.ExecuteTemplate(w, "board.gohtml", temp)
}

func post(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	var b Board
	gormDB.First(&b, id)

	tpl.ExecuteTemplate(w, "post.gohtml", b)
}

func getPageList(p string, limit int) []string {
	page, _ := strconv.Atoi(p)
	var result []string

	for i := page - 2; i <= page+2; i++ {
		if i > 0 && i <= limit {
			result = append(result, strconv.Itoa(i))
		}
	}
	return result
}
