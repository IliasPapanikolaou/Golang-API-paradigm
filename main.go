package main

import (
	"log"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	store, err := newPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.init(); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%+v\n", store)

	server := NewApiServer(":3000", store)
	server.Run()
}
