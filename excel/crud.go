package excel

import (
    // "net/http"
    // "encoding/json"  
    // "errors"
    // "fmt"
    "context"
    // "github.com/jackc/pgx"
    "github.com/jackc/pgx/pgxpool"
    // "github.com/jackc/pgconn"
    // "pg_service/db"
    // "pg_service/util"
    // "pg-service//auth"
)


type NewPersonIn struct{
	Name string 
	Ident string 
	Tel string 
	Address string 
}

func AddNewPerson(conn *pgxpool.Pool, npi NewPersonIn) error {
	sql :=
		`insert into wbk.linelisting
		( name, ident, tel, address )
		values
		( $1, $2, $3, $4 )`
		
	_, err := conn.Exec(context.Background(), sql, 
		npi.Name, npi.Ident, npi.Tel, npi.Address)
	if err != nil {
		return err 
	}
	return nil 
}