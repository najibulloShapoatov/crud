package security

import (
	"log"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)



//Service ...
type Service struct {
	db *pgxpool.Pool
}

//NewService ..
func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}



//Auth ...
func (s *Service) Auth(login, password string) bool{

	//это наш sql запрос
	sqlStatement := `select login, password from managers where login=$1 and password=$2`

	//выполняем запрос к базу
	err := s.db.QueryRow(context.Background(), sqlStatement, login, password).Scan(&login, &password)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}