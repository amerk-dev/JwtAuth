package main

import (
	"JwtAuth/internal/handler"
	"JwtAuth/internal/repo"
	"net/http"
)

func main() {
	db, err := repo.InitDB()
	if err != nil {
		panic(err)
	}
	pgDB := repo.NewDB(db)
	server := handler.NewServer(pgDB)

	// дада, можно заморочиться и сделать нормально, вынести "/api/v1/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hello World!")) })
	http.HandleFunc("/auth/get-token", server.AccessMethod)
	http.HandleFunc("/auth/refresh", server.RefreshHandler)

	http.ListenAndServe(":8080", nil)
}
