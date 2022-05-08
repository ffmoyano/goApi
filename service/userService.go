package service

import (
	"database/sql"
	"github.com/ffmoyano/goApi/database"
	"github.com/ffmoyano/goApi/model"
)

// FindAllUsers makes a query to the DB and returns all registered users in the app
func FindAllUsers() ([]model.User, error) {
	db := database.Get()
	var users []model.User
	var result *sql.Rows
	var innerResult *sql.Rows

	// queries to the db
	if result, err = db.Query(
		"Select id, name, username, password, email, active from user"); err != nil {
		return users, err
	} else {
		// generates a new user for each result row and appends it to the users list
		for result.Next() {
			var user model.User
			var roles []model.Role
			if err = result.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email,
				&user.Active); err != nil {
				return users, err
			}
			// for each user row it makes a query to retrieve the role, and appends it to the user
			if innerResult, err = db.Query(
				"select r.id, r.name from role r inner join user_role ur on r.id = ur.id_role where ur.id_user = ?", user.ID); err != nil {
				return users, err
			}

			for innerResult.Next() {
				var role model.Role
				if err = innerResult.Scan(&role.ID, &role.Name); err != nil {
					return users, err
				}
				roles = append(roles, role)

			}
			user.Roles = roles
			users = append(users, user)
		}
		// closes the inner resultset
		if innerResult != nil {
			if err = innerResult.Close(); err != nil {
				return users, err
			}
		}

		// closes the outer resultset
		if result != nil {
			if err = result.Close(); err != nil {
				return users, err
			}
		}

	}
	return users, nil
}
