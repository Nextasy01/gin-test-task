package db

import (
	"github.com/Nextasy01/gin-test-task/entity"
	"github.com/Nextasy01/gin-test-task/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	SaveUser(user entity.User) error
	// GetUserByID(uid string) (*entity.User, error)
	GetUsersByName(name string) ([]*entity.User, error)
	LoginCheck(name, password string) (string, string, error)
}

func NewUserRepository(db *Database) UserRepository {
	return db
}

func (db *Database) SaveUser(user entity.User) error {
	var err error

	user.Password, err = HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec("INSERT INTO users(id,email,name,age,password) VALUES(?,?,?,?,?)", user.ID, user.Email, user.Name, user.Age, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) LoginCheck(email, password string) (string, string, error) {
	u := entity.User{}

	err := db.conn.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&u.ID, &u.Email, &u.Name, &u.Age, &u.Password, &u.CreatedAt)
	if err != nil {
		return "", "", err
	}

	if err := VerifyPassword(password, u.Password); err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", "", err
	}

	token, err := utils.GenerateToken(u.ID)

	if err != nil {
		return "", "", err
	}

	return token, u.ID.String(), nil
}

func (db *Database) GetUsersByName(name string) ([]*entity.User, error) {
	var users []*entity.User
	query := "SELECT id, name, age FROM users WHERE name LIKE ?"
	rows, err := db.conn.Query(query, name+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := entity.NewUser()
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			break
		}
		users = append(users, u)
	}

	return users, err
}

func VerifyPassword(password, hashedPassword string) error {
	_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	password = string(hashedPassword)

	return password, nil
}
