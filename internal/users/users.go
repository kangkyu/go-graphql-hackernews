package users

import (
	"database/sql"
	"log"

	database "github.com/kangkyu/hackernews/internal/pkg/db/mysql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	Password string
}

func (user *User) Create() {
	stmt, err := database.Db.Prepare("INSERT INTO Users (Username, Password) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := HashPassword(user.Password)
	_, err = stmt.Exec(user.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (user *User) Authenticate() bool {
	stmt, err := database.Db.Prepare("SELECT Password from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(user.Username)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return CheckPasswordHash(user.Password, hashedPassword)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserIdByUsername will be used in the auth middleware
func GetUserIdByUsername(username string) (int, error) {
	stmt, err := database.Db.Prepare("SELECT ID FROM Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}
	return Id, nil
}
