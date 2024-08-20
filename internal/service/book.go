package service

import "database/sql"

type Book struct {
	ID int
	Title string
	Author string
	Genre string
}

type BookService struct {
	db * sql.DB
}

func NewBookService(db *sql.DB) *BookService{
	return &BookService{db: db}
}

func (s *BookService) CreateBook (book *Book) error {
	query := "INSERT INTO books (title, author, genre) VALUES (?, ?, ?)"
	result, err := s.db.Exec(query, book.Title, book.Author, book.Genre)
	
	if err != nil {
		return err
	}

	lastInsertID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	book.ID = int(lastInsertID)
	return nil
}

func (s *BookService) GetBooks() ([]Book, error) {
	query := "Select id, title, author, genre from books"
	rows, err := s.db.Query(query)

	if err != nil {
			return nil, err
	}

	var books []Book 
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
		if err != nil {
			return nil, err
		}
		books = append(books, book )
	}
	return books, nil
}

func (s *BookService) GetBookByID(id int) (*Book, error){
	query := "select id, title, author, genre from books where id = ?"
	rows := s.db.QueryRow(query, id)

	var book Book
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
	if err != nil{
		return nil, err
	}
	return &book, nil
}

// UpdateBook atualiza as informações de um livro no banco de dados.
func (s *BookService) UpdateBook(book *Book) error {
	query := "UPDATE books SET title = ?, author = ?, genre = ? WHERE id = ?"
	_, err := s.db.Exec(query, book.Title, book.Author, book.Genre, book.ID)
	return err
}

// DeleteBook deleta um livro do banco de dados.
func (s *BookService) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id = ?"
	_, err := s.db.Exec(query, id)
	return err
}

// SearchBooksByName busca livros pelo nome (título) no banco de dados.
func (s *BookService) SearchBooksByName(name string) ([]Book, error) {
	query := "SELECT id, title, author, genre FROM books WHERE title LIKE ?"
	rows, err := s.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}