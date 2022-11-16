package main

import (
	"fmt"
	"net/http"
)

type databaseV2 map[string]dollars

func (db databaseV2) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		_, err := fmt.Fprintf(w, "%s: %s \n", item, price)
		if err != nil {
			return
		}
	}
}

func (db databaseV2) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

func main() {
	db := databaseV2{"shoes": 50, "socks": 5}
	//mux := http.NewServeMux()
	//mux.Handle("/list", http.HandlerFunc(db.list))
	//mux.Handle("/price", http.HandlerFunc(db.price))
	//err := http.ListenAndServe("localhost:8080", mux)

	//简写
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/list", db.price)
	err := http.ListenAndServe("localhost:8080", nil)

	if err != nil {
		return
	}

}
