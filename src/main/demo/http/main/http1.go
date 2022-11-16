package main

import (
	"fmt"
	"net/http"
)

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database map[string]dollars

//v1 简单版
//func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	for item, price := range db {
//		_, err := fmt.Fprintf(w, "%s: %s \n", item, price)
//		if err != nil {
//			return
//		}
//	}
//}

//v2 路由版
func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	switch path {
	case "/list": //查询全部
		for item, price := range db {
			_, err := fmt.Fprintf(w, "%s: %s \n", item, price)
			if err != nil {
				return
			}
		}
	case "/price": //参数查询单个商品
		itemArg := req.URL.Query().Get("item")
		price, ok := db[itemArg]
		if !ok {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "no such item: %q\n", itemArg)
			return
		}
		fmt.Fprintf(w, "%s: %s \n", itemArg, price)
	default:
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such path: %s\n", path)
	}
}

func main() {
	db := database{"shoes": 200, "socks": 10}
	http.ListenAndServe("localhost:8080", db)
}
