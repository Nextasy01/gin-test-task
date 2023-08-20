package db

import (
	"fmt"

	"github.com/Nextasy01/gin-test-task/entity"
)

type PhoneRepository interface {
	SavePhone(phone entity.Phone) error
	UpdatePhone(phone_id, number, description string, is_fax bool) error
	DeletePhone(phone_id string) error
	GetPhone(number string) ([]*entity.Phone, error)
	CheckPhone(number string) error
}

func NewPhoneRepository(db *Database) PhoneRepository {
	return db
}

func (db *Database) SavePhone(phone entity.Phone) error {
	var err error

	_, err = db.conn.Exec("INSERT INTO phone(id,phone,description,user_id,is_fax) VALUES(?,?,?,?,?)", phone.ID, phone.Number, phone.Description, phone.UserId, phone.IsFax)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetPhone(number string) ([]*entity.Phone, error) {
	var phones []*entity.Phone
	query := "SELECT user_id, phone, description, is_fax FROM phone WHERE phone LIKE ?"
	rows, err := db.conn.Query(query, number+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := entity.NewPhone()

		if err := rows.Scan(&p.UserId, &p.Number, &p.Description, &p.IsFax); err != nil {
			break
		}

		phones = append(phones, p)
	}

	return phones, nil

}

func (db *Database) CheckPhone(number string) error {
	phone := entity.NewPhone()
	query := "SELECT id FROM phone WHERE phone = ?"
	err := db.conn.QueryRow(query, number).Scan(&phone.ID)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdatePhone(phone_id, number, description string, is_fax bool) error {
	var err error

	_, err = db.conn.Exec("UPDATE phone SET phone = ?, description = ?, is_fax = ? WHERE id=?", number, description, is_fax, phone_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeletePhone(phone_id string) error {
	var err error

	res, err := db.conn.Exec("DELETE FROM phone WHERE id=?", phone_id)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("phone with id '%s' does not exist", phone_id)
	}
	return nil
}
