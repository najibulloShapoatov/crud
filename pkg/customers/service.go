package customers

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"errors"
	"log"
	"time"
)

//ErrNotFound ...
var ErrNotFound = errors.New("item not found")

//ErrInternal ...
var ErrInternal = errors.New("internal error")

//Service ..
type Service struct {
	db *pgxpool.Pool
}

//NewService ..
func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

//Customer ...
type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

//All ....
func (s *Service) All(ctx context.Context) (cs []*Customer, err error) {
 
	//это наш sql запрос
	sqlStatement := `select * from customers`

	rows, err := s.db.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)
		if err != nil {
			log.Println(err)
		}
		cs = append(cs, item)
	}

	return cs, nil
}

//AllActive ....
func (s *Service) AllActive(ctx context.Context) (cs []*Customer, err error) {

	//это наш sql запрос
	sqlStatement := `select * from customers where active=true`

	rows, err := s.db.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)
		if err != nil {
			log.Println(err)
		}
		cs = append(cs, item)
	}

	return cs, nil
}

//ByID ...
func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	//это наш sql запрос
	sqlStatement := `select * from customers where id=$1`
	//выполняем запрос к базу
	err := s.db.QueryRow(ctx, sqlStatement, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)

	//если sql нам не вернул резултат
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	//проверим ошибку если во время выполнения запросы было какая то ошибка то вернем InternalError
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

//ChangeActive ...
func (s *Service) ChangeActive(ctx context.Context, id int64, active bool) (*Customer, error) {
	item := &Customer{}

	//это наш sql запрос
	sqlStatement := `update customers set active=$2 where id=$1 returning *`
	//выполняем запрос к базу
	err := s.db.QueryRow(ctx, sqlStatement, id, active).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)
	//если sql нам не вернул резултат
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	//проверим ошибку если во время выполнения запросы было какая то ошибка то вернем InternalError
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

//Delete ...
func (s *Service) Delete(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	//это наш sql запрос
	sqlStatement := `delete from customers  where id=$1 returning *`
	//выполняем запрос к базу
	err := s.db.QueryRow(ctx, sqlStatement, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)

	//если sql нам не вернул резултат
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	//проверим ошибку если во время выполнения запросы было какая то ошибка то вернем InternalError
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

//Save ...
func (s *Service) Save(ctx context.Context, customer *Customer) (c *Customer, err error) {

	//обявляем пустую структуру
	item := &Customer{}

	//если id равно то сделаем инцерт (тоест создаем и веренем только что созданный клиент)
	if customer.ID == 0 {

		//это наш sql запрос
		sqlStatement := `insert into customers(name, phone) values($1, $2) returning *`

		//выполняем запрос к базу
		err = s.db.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)

	} else { //если нет обновляем и вернем обновленный

		//это наш sql запрос
		sqlStatement := `update customers set name=$1, phone=$2 where id=$3 returning *`
		//выполняем запрос к базу
		err = s.db.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone, customer.ID).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
	}

	//проверим ошибку если во время выполнения запросы было какая то ошибка то вернем InternalError
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}
