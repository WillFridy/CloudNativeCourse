package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := database{data: map[string]dollars{"shoes": 50, "socks": 10}}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))

}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database struct { // converted database to a struct
	sync.Mutex
	data map[string]dollars
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	item := req.URL.Query().Get("item")
	if price, ok := db.data[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	newItem := req.URL.Query().Get("item") // get item and price
	newPrice := req.URL.Query().Get("price")
	tempPrice, _ := strconv.ParseFloat(newPrice, 32)
	updatedPrice := dollars(tempPrice)

	for item := range db.data { // check to see if item is already made
		if newItem != item {
			db.data[newItem] = updatedPrice
		}
	}

	for item, price := range db.data { // print updated list
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}

}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	item := req.URL.Query().Get("item") // get item and price
	newPrice := req.URL.Query().Get("price")
	tempPrice, _ := strconv.ParseFloat(newPrice, 32) // convert price from string to float
	updatedPrice := dollars(tempPrice)               // convert to type dollars
	db.data[item] = updatedPrice

	for item, price := range db.data { // print updated list
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}

}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	delItem := req.URL.Query().Get("item")
	for item := range db.data { // check for item in the list
		if item == delItem {
			delete(db.data, item) // delete item from list
		}
	}
	for item, price := range db.data { // print updated list
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}

}
