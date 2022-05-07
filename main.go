package main

import (
  "fmt"
  "net/http"
	"html/template"

  "github.com/gorilla/mux"

  "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Items struct {
  Id uint16
  Title, Anons, FullText string
}

var posts = []Items{}
var showItems = Items{}

func main()  {
  handleFunc()
}

func show_items (w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)

  t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")

	if err != nil{
		fmt.Fprintf(w, err.Error())
	}

  db, err := sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventory")
  if err != nil{
    panic(err)
  }
  defer db.Close()

  //Выборка данных
  res, err := db.Query(fmt.Sprintf("SELECT * FROM `items` WHERE `id` = '%s' ", vars["id"]))
  if err!= nil {
    panic(err)
  }

  showItems = Items{}  //очищаем список
  for res.Next() {
    var post Items
    err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
    if err != nil {
      panic(err)
    }

    showItems = post
  }

  t.ExecuteTemplate(w, "show", showItems)

}

func handleFunc()  {
  rtr := mux.NewRouter()
  rtr.HandleFunc("/", index).Methods("GET")
  rtr.HandleFunc("/create", create).Methods("GET")
  rtr.HandleFunc("/save_article", save_article).Methods("POST")
  rtr.HandleFunc("/post/{id:[0-9]+}", show_items).Methods("GET")

  http.Handle("/", rtr)
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  http.ListenAndServe(":1212", nil)
}

func index(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil{
		fmt.Fprintf(w, err.Error())
	}

  db, err := sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventory")
  if err != nil{
    panic(err)
  }
  defer db.Close()

  //Выборка данных
  res, err := db.Query("SELECT * FROM `items` ")
  if err!= nil {
    panic(err)
  }

  posts = []Items{}  //очищаем список
  for res.Next() {
    var post Items
    err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
    if err != nil {
      panic(err)
    }

    posts = append(posts, post)
  }

  t.ExecuteTemplate(w, "index", posts)

}

func create(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil{
		fmt.Fprintf(w, err.Error())
	}

  t.ExecuteTemplate(w, "create", nil)

}

func save_article(w http.ResponseWriter, r *http.Request){
  title := r.FormValue("title")
  anons := r.FormValue("anons")
  full_text := r.FormValue("full_text")

  if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {

  db, err := sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventory")
  if err != nil{
    panic(err)
  }
  defer db.Close()

  //Добавление данных
  insert, err := db.Query(fmt.Sprintf("INSERT INTO `items` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, full_text))
  if err != nil {
    panic(err)
  }
  defer insert.Close()

  http.Redirect(w, r, "/", http.StatusSeeOther)
}
}
