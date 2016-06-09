package main

import (
	"net/http"
	"fmt"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/memap"
)

func main() {


	store := &memap.MmapKVStore{}
	store.Open()
	defer store.Close()
	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Fprint(w, r.Form.Get("key"))
		fmt.Fprint(w, store.Get(r.Form.Get("key")))
	})

	http.HandleFunc("/put/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		store.Put(r.Form.Get("key"), r.Form.Get("value"))
	})

	http.ListenAndServe(":2233", nil)
}
