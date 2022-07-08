package model

import (
	_ "github.com/lib/pq"
)

type USER struct {
	Username string `json:"user"`
	Pwd      string `json:"pwd"`
}

func AddUsertoDB(username string, pwd string) (*USER, error) {
	created := USER{}

	row := db.QueryRow(
		`INSERT INTO users (username,pwd) VALUES ($1,$2) RETURNING username, pwd`,
		username, pwd,
	)

	err := row.Scan(&created.Username, &created.Pwd)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
func GetPwd(username string) (*USER, error) {
	var user1 USER
	rows, err := db.Query("SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user1.Username, &user1.Pwd)
		if err != nil {
			return nil, err
		}
	}
	return &user1, nil

}
func GetAllUsers() ([]USER, error) {
	var users []USER
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		u := USER{}

		err = rows.Scan(&u.Username, &u.Pwd)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}
	return users, nil
}
