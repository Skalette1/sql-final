package main

import (
	"database/sql"
	"fmt"

)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	// верните идентификатор последней добавленной записи
	res, err := s.db.Exec("INSERT INTO parcel  VALUES (:Address, :Client, :CreatedAt, :Number, :Status)",
	sql.Named("Address", p.Address),
	sql.Named("Client", p.Client),
	sql.Named("CreatedAt", p.CreatedAt),
	sql.Named("Number", p.Number),
	sql.Named("Status", p.Status))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	query := "SELECT address, client, createdAt, number, status FROM parcel WHERE number = ?"
	row := s.db.QueryRow(query, number)
	var (
		address string
		client int
		createdAt string
		status string
	)
	err := row.Scan(&address, &client, &createdAt, &number, &status)
	if err != nil {
		fmt.Println(err)
		return Parcel{}, err
	}
	// здесь из таблицы должна вернуться только одна строка
	// заполните объект Parcel данными из таблицы
	p := Parcel{
		Address: address,
		Client:  client,
		CreatedAt: createdAt,
		Status: status,
		Number: number,
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	query := "SELECT address, client, createdAt, number, status FROM parcel WHERE client = ?"
	row, err := s.db.Query(query, client)
	if err != nil {
		fmt.Println(err)
	}
	// здесь из таблицы может вернуться несколько строк
	if err != nil {
		fmt.Println(err)
	}
	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for row.Next() {
		var (
		address string
		client int
		createdAt string
		number int
		status string
	)
	err := row.Scan(&address, &client, &createdAt,&number, &status)
	if err != nil {
		fmt.Println(err)
	}
	p := Parcel{
			Address:   address,
			Client:    client,
			CreatedAt: createdAt,
			Number:    number,
			Status:    status,
	}
	res = append(res, p)
	}
	if err = row.Err(); err != nil {
		fmt.Println(err)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE id = :number",
	sql.Named("status", status),
	sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	query := "SELECT status FROM parcel WHERE number = ?"
	row := s.db.QueryRow(query, number)
	var status string
	err := row.Scan(&status)
	if err != nil {
		return err
	}
	if status == ParcelStatusRegistered {
		_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE id = :number",
		sql.Named("address", address),
		sql.Named("number", number))
		if err != nil {
			fmt.Println(err)
			return err
		} 
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	query := "SELECT status FROM parcel WHERE number = ?"
	row := s.db.QueryRow(query, number)
	var status string
	err := row.Scan(&status)
	if err != nil {
		return err
	}
	if status == ParcelStatusRegistered {
		_, err := s.db.Exec("DELETE FROM parcel WHERE id = :number", sql.Named("id", number))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
