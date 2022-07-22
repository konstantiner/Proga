package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Items struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Items{}
var showItems = Items{}

func main() {
	handleFunc()
}

func show_items(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventory")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Выборка данных
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `items` WHERE `id` = '%s' ", vars["id"]))
	if err != nil {
		panic(err)
	}

	showItems = Items{} //очищаем список
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

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create_items).Methods("GET")
	rtr.HandleFunc("/save_item", save_items).Methods("POST")
	rtr.HandleFunc("/show_item/{id:[0-9]+}", show_items).Methods("GET")
	rtr.HandleFunc("/delete_item/{id:[0-9]+}", delete_items)
	rtr.HandleFunc("/edit/{id:[0-9]+}", EditItemPage).Methods("GET")
	rtr.HandleFunc("/edit/{id:[0-9]+}", EditItemHandler).Methods("POST")
	rtr.HandleFunc("/write_off_item/{id:[0-9]+}", WriteOffItem).Methods("POST", "GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":1212", nil)
}

func WriteOffItem(w http.ResponseWriter, r *http.Request) { //Списание
	vars := mux.Vars(r)

	db, _ := sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventory")
	defer db.Close()
	//Выбор данных по id
	_, _ = db.Query("SELECT title, anons, full_text FROM `items` WHERE `id` = '%s' ", vars["id"])

	prod := showItems
	id := prod.Id
	title := prod.Title
	anons := prod.Anons
	full_text := prod.FullText
	//fmt.Println(id, title, anons, full_text)

	//Добавление данных
	_, err := db.Exec("INSERT INTO writeoffitem (id, title, anons, full_text) VALUES(?, ?, ?, ?)", id, title, anons, full_text)
	if err != nil {
		panic(err)
	}
	//defer insert.Close()

	//Удаление данных
	delete, err := db.Query(fmt.Sprintf("DELETE FROM `items` WHERE `id` = '%s' ", vars["id"]))
	if err != nil {
		panic(err)
	}
	defer delete.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditItemPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, _ := sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventory")

	defer db.Close()

	//Выбор данных по id
	row := db.QueryRow(fmt.Sprintf("SELECT * FROM `items` WHERE `id` = '%s' ", vars["id"]))
	prod := showItems
	err := row.Scan(&prod.Id, &prod.Title, &prod.Anons, &prod.FullText)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		t, err := template.ParseFiles("templates/edit.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		t.ExecuteTemplate(w, "edit", prod)
	}
}

func EditItemHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	id := r.FormValue("id")
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	db, err := sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventory")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//  fmt.Println(id, title, anons, full_text)

	_, err = db.Exec("UPDATE items set title=?, anons=?, full_text=? WHERE id=?", title, anons, full_text, id)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	}

	// update, err := db.Query(fmt.Sprintf("UPDATE `items` set `title`=%s, `anons`=%s, `full_text`=%s WHERE `id`=%s", title, anons, full_text, id))
	//  if err != nil {
	//    http.Error(w, http.StatusText(404), http.StatusNotFound)
	//  }
	//  defer update.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func delete_items(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventory")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Удаление данных
	delete, err := db.Query(fmt.Sprintf("DELETE FROM `items` WHERE `id` = '%s' ", vars["id"]))
	if err != nil {
		panic(err)
	}
	defer delete.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventory")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Выборка данных
	res, err := db.Query("SELECT * FROM `items` ")
	if err != nil {
		panic(err)
	}

	posts = []Items{} //очищаем список
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

func create_items(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)

}

func save_items(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {

		db, err := sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventory")
		if err != nil {
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
