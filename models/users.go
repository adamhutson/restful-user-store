package models

// User comment
type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserID    string `json:"userid"`
}

//AllUsers comment
func AllUsers() ([]User, error) {
	users := make([]User, 0)

	stmt := "SELECT FirstName, LastName, UserID FROM users"
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.FirstName, &user.LastName, &user.UserID); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

//OneUser comment
func OneUser(userid string) (User, error) {
	user := User{}

	stmt := "SELECT FirstName, LastName, UserID FROM users WHERE UserID = $1"

	rows, err := db.Query(stmt, userid)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.FirstName, &user.LastName, &user.UserID); err != nil {
			return user, err
		}
	}
	if err = rows.Err(); err != nil {
		return user, err
	}

	return user, nil
}

//CreateUser comment
func CreateUser(user User) (int64, error) {
	stmt := "INSERT INTO users (FirstName, LastName, UserID) VALUES ($1, $2, $3) ON CONFLICT (UserID) DO NOTHING"
	result, err := db.Exec(stmt, user.FirstName, user.LastName, user.UserID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// UpdateUser comment
func UpdateUser(user User) (int64, error) {
	stmt := "UPDATE users SET FirstName = $1, LastName = $2 WHERE UserID = $3"
	result, err := db.Exec(stmt, user.FirstName, user.LastName, user.UserID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// DeleteUser comment
func DeleteUser(userid string) (int64, error) {
	stmt := "DELETE FROM users WHERE UserID = $1"

	result, err := db.Exec(stmt, userid)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
