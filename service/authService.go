package service

import (
	"context"
	"database/sql"
	"github.com/ffmoyano/goApi/database"
	"github.com/ffmoyano/goApi/model"
	"time"
)

var err error
var tx *sql.Tx

// SignUp receives a request with the user data and, if data is valid, generates a new user of the app
func SignUp(user model.SignUpRequest) error {

	db := database.Get()
	var result sql.Result

	// we use context to start a transaction because we may need to rollback in case of failure as there are multiple
	// queries that make changes to the table in this function
	ctx := context.Background()
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}

	// insert user
	if result, err = tx.Exec("insert into user (name, username, password, email) "+
		" values (?, ?, ?, ?)", user.Name, user.Username, user.Password, user.Email); err != nil {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			return rollBackErr
		}
	} else {
		// if user is successfully inserted we insert its role too
		var lastId, _ = result.LastInsertId()
		if _, err = tx.Exec("insert into user_role (id_user, id_role) values (?, 1)", lastId); err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				return rollBackErr
			}
		} else {
			// if everything is correct we commit transaction
			if err = tx.Commit(); err != nil {
				return err
			}
		}
	}

	return nil
}

// FindByUsername queries the DB and returns the user with the provided username, or empty user if there is none
func FindByUsername(username string) (model.User, error) {

	db := database.Get()
	var user model.User
	var result *sql.Rows
	var innerResult *sql.Rows

	if result, err = db.Query(
		"Select id, name, username, password, email, active from user where username = ? ", username); err != nil {
		return user, err
	}

	for result.Next() {
		var roles []model.Role
		if err = result.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email,
			&user.Active); err != nil {
			return user, err
		}
		// for each retrieved user we recover its roles
		if innerResult, err = db.Query(
			"select r.id, r.name from role r inner join user_role ur on r.id = ur.id_role where ur.id_user = ? ", user.ID); err != nil {
			return user, err
		}

		for innerResult.Next() {
			var role model.Role
			if err = innerResult.Scan(&role.ID, &role.Name); err != nil {
				return user, err
			}
			roles = append(roles, role)
		}
		user.Roles = roles
	}
	if innerResult != nil {
		if err = innerResult.Close(); err != nil {
			return user, err
		}
	}
	if result != nil {
		if err = result.Close(); err != nil {
			return user, err
		}
	}

	return user, nil
}

func InsertRefreshToken(user model.User, tokenResponse model.TokenResponse) error {
	db := database.Get()
	expirationTime := time.Now().Add(30 * 24 * time.Hour)

	if _, err = db.Exec("delete from token WHERE id_user = ?", user.ID); err != nil {
		return err
	}

	if _, err = db.Exec("INSERT INTO token (id_user, expiration, refresh_token) VALUES (?, ?, ?);", user.ID, expirationTime, tokenResponse.RefreshToken); err != nil {
		return err
	}
	return nil
}
