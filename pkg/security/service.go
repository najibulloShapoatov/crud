package security

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

//Service сервис авторизации
type Service struct {
	db *pgxpool.Pool
}

//NewService создаем новый сервис авторизации
func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

//Auth в этом мметоде проверяем логин и парол если они верны то вернем true если нет false
func (s *Service) Auth(login, password string) bool {

	//это наш sql запрос
	sqlStatement := `select login, password from managers where login=$1 and password=$2`

	//выполняем запрос к базу
	err := s.db.QueryRow(context.Background(), sqlStatement, login, password).Scan(&login, &password)
	//если при выполнения запроса к базе  получили ошибку то печатаем ошибку и вернем false
	if err != nil {
		log.Print(err)
		return false
	}
	//если все хорошо (тоест такой пользовател есть) то вернем true
	return true
}
