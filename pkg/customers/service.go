package customers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	//ErrNotFound ...
	ErrNotFound = errors.New("item not found")
	//ErrInternal ...
	ErrInternal = errors.New("internal error")
	//ErrTokenNotFound ...
	ErrTokenNotFound = errors.New("token not found")
	//ErrNoSuchUser ...
	ErrNoSuchUser = errors.New("no such user")
	//ErrInvalidPassword ..
	ErrInvalidPassword = errors.New("invalid password")
	//ErrPhoneUsed ...
	ErrPhoneUsed = errors.New("phone alredy registered")
	//ErrTokenExpired ...
	ErrTokenExpired = errors.New("token expired")
)

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
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

//Product ...
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Qty   int    `json:"qty"`
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
		sqlStatement := `insert into customers(name, phone, password) values($1, $2, $3) returning *`
		//выполняем запрос к базу
		err = s.db.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone, customer.Password).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Password,
			&item.Active,
			&item.Created)

	} else { //если нет обновляем и вернем обновленный

		//это наш sql запрос
		sqlStatement := `update customers set name=$1, phone=$2, password=$3 where id=$4 returning *`
		//выполняем запрос к базу
		err = s.db.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone, customer.Password, customer.ID).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Password,
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

//Token .... метод для генерации токена
func (s *Service) Token(ctx context.Context, phone, password string) (string, error) {

	//обявляем переменную хеш и парол
	var hash string
	var id int64
	//выполняем запрос и извелекаем ид и хеш пароля
	err := s.db.QueryRow(ctx, "select id, password from customers where phone = $1", phone).Scan(&id, &hash)
	//если ничего не получили вернем ErrNoSuchUser
	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}
	//если другая ошибка то вернем ErrInternal
	if err != nil {
		return "", ErrInternal
	}
	//проверим хеш с представленным паролем
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}

	//генерируем токен
	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}

	token := hex.EncodeToString(buffer)
	_, err = s.db.Exec(ctx, "insert into customers_tokens(token, customer_id) values($1, $2)", token, id)
	if err != nil {
		return "", ErrInternal
	}

	return token, nil

}

//Products ...
func (s *Service) Products(ctx context.Context) ([]*Product, error) {

	items := make([]*Product, 0)

	sqlStatement := `select id, name, price, qty from products where active = true order by id limit 500`
	rows, err := s.db.Query(ctx, sqlStatement)

	if err != nil {
		if err == pgx.ErrNoRows {
			return items, nil
		}
		return nil, ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		item := &Product{}
		err = rows.Scan(&item.ID, &item.Name, &item.Price, &item.Qty)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

//IDByToken ....
func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	sqlStatement := `select customer_id from customers_tokens where token = $1`
	err := s.db.QueryRow(ctx, sqlStatement, token).Scan(&id)

	if err != nil {

		if err == pgx.ErrNoRows {
			return 0, nil
		}

		return 0, ErrInternal
	}

	return id, nil
}
