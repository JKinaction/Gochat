package model

import (
	_ "github.com/lib/pq"
)

type message struct {
	User string `json:"user"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
}
type primes struct {
	User   string `json:"user"`
	Msg    string `json:"msg"`
	ToUser string `json:"touser"`
	Time   string `json:"time"`
}

var GlobalmsgList = []message{{
	User: "Mod",
	Msg:  "Welcome to the global chat.",
	Time: "[00:00:00] ",
}}

func AddPriMsgtoDB(user string, msg string, touser string, time string) (*primes, error) {
	created := primes{}

	row := db.QueryRow(
		`INSERT INTO primess (username,msg,touser,time) VALUES ($1,$2,$3,$4) RETURNING username, msg, touser, time`,
		user, msg, touser, time,
	)

	err := row.Scan(&created.User, &created.Msg, &created.ToUser, &created.Time)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
func AddMsgtoDB(user string, msg string, time string) (*message, error) {
	created := message{}

	row := db.QueryRow(
		`INSERT INTO messages (username,msg,time) VALUES ($1,$2,$3) RETURNING username, msg, time`,
		user, msg, time,
	)

	err := row.Scan(&created.User, &created.Msg, &created.Time)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
func GetPriMsgsDB(touser string, user string) ([]primes, error) {

	rows, err := db.Query(
		`SELECT * FROM primess where (username = $1 and touser= $2) or (username = $2 and touser=$1) `, touser, user,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	msgs := make([]primes, 0, 10)

	for rows.Next() {
		m := primes{}
		err = rows.Scan(&m.User, &m.Msg, &m.ToUser, &m.Time)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}

	return msgs, nil
}
func GetAllMsgsDB() ([]message, error) {

	rows, err := db.Query(
		`SELECT * FROM messages`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	msgs := make([]message, 0, 10)

	for rows.Next() {
		m := message{}

		err = rows.Scan(&m.User, &m.Msg, &m.Time)

		if err != nil {
			return nil, err
		}

		msgs = append(msgs, m)
	}

	return msgs, nil
}
