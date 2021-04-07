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
	Bil int 
	Notifydt string
	Name string 
	Ident string 
	Tel string 
	Address string 
	State string
}

func AddNewPerson(conn *pgxpool.Pool, npi NewPersonIn) error {
	sql :=
		`insert into wbk.linelisting
		( bil, notifydt, name, ident, tel, address, state )
		values
		( $1, $2, $3, $4, $5, $6, $7 )`
		
	_, err := conn.Exec(context.Background(), sql, 
		npi.Bil, npi.Notifydt, npi.Name, 
		npi.Ident, npi.Tel, npi.Address, npi.State)
	if err != nil {
		return err 
	}
	return nil 
}

func GetPersons(conn *pgxpool.Pool) (Linelisting, error) {
	sql :=
		`select bil, address, state
		   from wbk.linelisting`

	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}

	var linelisting Linelisting  
	for rows.Next() {
		var bil int 
		var address string 
		var state string 

		err = rows.Scan(&bil, &address, &state)
		if err != nil {
			return nil, err 
		}

		person := Person{
			Bil: bil,
			Address: address,
			State: state,
		}
		linelisting.Persons = append(linelisting.Persons, person)
	}
	return linelisting, err
}