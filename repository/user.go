package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"crud-product/model"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	DB *sql.DB
}

type ErrorResponse struct {
	Err string
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &User{
		DB: db,
	}
}

func (u *User) FindOne(ctx context.Context, email, password string) (model.User, error) {
	query := `
			SELECT 
				user_id,
				name,
				email,
				password,
			    role
			FROM 
				user
			WHERE
				email = ?`

	user := model.User{}
	err := u.DB.QueryRowContext(ctx, query, email).Scan(
		&user.Id, &user.Name, &user.Email,
		&user.Password, &user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("data not found %s", err.Error())
		}
		return user, err
	}

	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return user, errors.New("invalid password")
	}

	tk := &model.Token{
		UserID: user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		return user, err
	}

	user.Token = tokenString

	return user, err
}

func (u *User) Store(ctx context.Context, user model.User) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(pass)

	query := `
				INSERT INTO user 
					(name, email, password, gender, role)
				VALUES
					(?, ?, ?, ?, ?)
			`

	_, err = u.DB.ExecContext(ctx, query,
		user.Name, user.Email, user.Password, user.Gender, user.Role)

	if err != nil {
		return err
	}

	return nil
}
