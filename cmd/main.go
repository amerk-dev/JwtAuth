package main

import (
	_ "JwtAuth/docs"
	"JwtAuth/internal/handler"
	"JwtAuth/internal/repo"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// @title JWT Auth Service
// @version 1.0
// @description Тестовое задание на позицию Junior Backend Developer
// @host localhost:8080
// @BasePath /
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

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.ListenAndServe(":8080", nil)
}
