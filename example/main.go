package main

import (
	"encoding/json"
	"net/http"

	"github.com/sjkaliski/pinned"
)

var (
	users = []*User{
		{
			ID:        1,
			Email:     "foo@bar.com",
			Name:      "foo",
			CreatedAt: 1257894000,
		},
	}
)

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// If no version supplied, default to latest.
	v, err := vm.Parse(r)
	if err == pinned.ErrNoVersionSupplied {
		v = vm.Latest()
	} else if err != nil {
		panic(err)
	}

	m, err := vm.Apply(v, users[0])
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/users", getUsersHandler)
	http.ListenAndServe(":8080", nil)
}
