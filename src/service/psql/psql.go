package main

import (
	"database/sql"
	_ "engine/pq"
	"fmt"
	"sync"
	"time"
)

var (
	db  *sql.DB
	mux sync.Mutex
)

var userTableSql string = `
    create table if not exists user_profile
    (
        uid serial,
        name varchar(20) not null,
        created varchar(20) not null,
        primary key(uid)
    )
`

func init() {
	mux.Lock()
	defer mux.Unlock()

	// check
	if db != nil {
		return
	}

	// open
	mysqldb, err := sql.Open("postgres", "host=192.168.1.138 port=5432 user=postgres password=admin dbname=t2m sslmode=disable")
	checkErr(err)

	// new db
	db = mysqldb

	// create database table
	_, err = db.Exec(userTableSql)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic("psql err: " + err.Error())
	}
}

func main() {
	// insert
	insertSql := `insert into user_profile values(nextval('user_profile_uid_seq'),$1,$2)`
	stmt, err := db.Prepare(insertSql)
	checkErr(err)

	_, err = stmt.Exec("viney", time.Now().Format("2006-01-02 15:04:05"))
	checkErr(err)

	// update
	updateSql := `update user_profile set name=$1 where name=$2`
	stmt, err = db.Prepare(updateSql)
	checkErr(err)

    res, err := stmt.Exec("viney.chow", "viney")
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("exec update,rows affected: " + fmt.Sprint(affect))

	// select
	querySql := `select * from user_profile where name=$1`
	rows, err := db.Query(querySql, "viney.chow")

	type user struct {
		uid     int
		name    string
		created string
	}

	var u = &user{}
	for rows.Next() {
		err = rows.Scan(
			&u.uid,
			&u.name,
			&u.created)
		checkErr(err)
	}

	fmt.Println(*u)

	// delete
	deleteSql := `delete from user_profile where name=$1`
	stmt, err = db.Prepare(deleteSql)
	checkErr(err)

	res, err = stmt.Exec("viney.chow")
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println("exec delete,rows affected: " + fmt.Sprint(affect))
}
