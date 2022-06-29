package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"newfinalproject/db"
	"newfinalproject/handler"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.UserPq, db.Password, db.Dbname)
	db.Db, db.Err = sql.Open("postgres", psqlInfo)
	defer db.Db.Close()
	if db.Err != nil {
		panic(db.Err)
	}

	db.Err = db.Db.Ping()
	if db.Err != nil {
		panic(db.Err)
	}
	fmt.Println("Done")
	r := mux.NewRouter()
	r.HandleFunc("/users/register", handler.PostRegister).Methods("POST")
	r.HandleFunc("/users/login", handler.UserLogin).Methods("POST")
	r.HandleFunc("/users", handler.UserUpdate).Methods("PUT")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
