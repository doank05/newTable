package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Customer struct (Model) ...
type Favorit struct {
	id      string `json:"id"`
	makanan string `json:"makanan"`
	minuman string `json:"minuman"`
	hewan   string `json:"hewan"`
	benda   string `json:"benda"`
}

// Get all orders

func getFavorit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var favorits []Favorit

	sql := `SELECT
				id,
				IFNULL(makanan,'') makanan,
				IFNULL(minuman,'') minuman,
				IFNULL(hewan,'') hewan,
				IFNULL(benda,'') benda
			FROM favorit`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var favorit Favorit
		err := result.Scan(&favorit.id, &favorit.makanan, &favorit.minuman, &favorit.hewan, &favorit.benda)

		if err != nil {
			panic(err.Error())
		}
		favorits = append(favorits, favorit)
	}

	json.NewEncoder(w).Encode(favorits)
}

func getFavorit2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var favorits []Favorit
	params := mux.Vars(r)

	sql := `SELECT
				id,
				IFNULL(makanan,'') makanan,
				IFNULL(minuman,'') minuman,
				IFNULL(hewan,'') hewan,
				IFNULL(benda,'') benda
			FROM favorit`

	result, err := db.Query(sql, params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var favorit Favorit

	for result.Next() {

		err := result.Scan(&favorit.id, &favorit.makanan, &favorit.minuman, &favorit.hewan, &favorit.benda)

		if err != nil {
			panic(err.Error())
		}

		favorits = append(favorits, favorit)
	}

	json.NewEncoder(w).Encode(favorits)
}

func createFavorit(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		id := r.FormValue("id")
		makanan := r.FormValue("makanan")
		minuman := r.FormValue("minuman")
		hewan := r.FormValue("hewan")
		benda := r.FormValue("benda")

		stmt, err := db.Prepare("INSERT INTO favorit (id,makanan,minuman,hewan,benda) VALUES (?,?,?,?,?)")

		_, err = stmt.Exec(id, makanan, minuman, hewan, benda)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}

func updateFavorit(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newMakanan := r.FormValue("makanan")

		stmt, err := db.Prepare("UPDATE favorit SET makanan = ? WHERE id = ?")

		_, err = stmt.Exec(newMakanan, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Favorit with id = %s was updated", params["id"])
	}
}

func deleteFavorit(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM favorit WHERE id = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Favorit with ID = %s was deleted", params["id"])
}

// Main function
func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/northwind")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/favorit", getFavorit).Methods("GET")
	r.HandleFunc("/favorit/{id}", getFavorit2).Methods("GET")
	r.HandleFunc("/favorit", createFavorit).Methods("POST")
	r.HandleFunc("/favorit/{id}", updateFavorit).Methods("PUT")
	r.HandleFunc("/favorit/{id}", deleteFavorit).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
