package main

import (
	"fmt"
	forum "forum/backend"
	"forum/database"
	"log"
	"net/http"
)

func main() {
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/main.html")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		forum.HandleWebSocket(w, r, db)
	})

	port := "8090"
	fmt.Printf("Listening on port %v\n", port)
	fmt.Println("server started . . .")
	fmt.Println("ctrl(cmd) + click: http://localhost:8090/")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
