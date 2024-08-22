package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"gobooks/internal/cli"
	"gobooks/internal/service"
	"gobooks/internal/web"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  mySQLURL := os.Getenv("MYSQL_URL")
	// Conexão com o banco de dados SQLite3
	db, err := sql.Open("mysql", mySQLURL)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Inicializando o serviço
	bookService := service.NewBookService(db)

	// Inicializando os handlers
	bookHandlers := web.NewBookHandlers(bookService)

	// Verifica se o CLI foi chamado
	if len(os.Args) > 1 && (os.Args[1] == "search" || os.Args[1] == "simulate"){
		bookCLI := cli.NewBookCLI(bookService)
		bookCLI.Run()
		return
	}

	// Criando o roteador com o novo servidor
	router := http.NewServeMux()

	// Configurando as rotas RESTful
	router.HandleFunc("GET /books", bookHandlers.GetBooks)
	router.HandleFunc("POST /books", bookHandlers.CreateBook)
	router.HandleFunc("GET /books/{id}", bookHandlers.GetBookByID)
	router.HandleFunc("PUT /books/{id}", bookHandlers.UpdateBook)
	router.HandleFunc("DELETE /books/{id}", bookHandlers.DeleteBook)

	// Iniciando o servidor
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}