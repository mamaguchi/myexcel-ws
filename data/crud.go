package data

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

type Person struct{
	Bil int 
	Address string 
	State string
}

type Linelisting struct{
	Persons []Person
}

type GeocodedPersonAddrIn struct{
	Bil int 
	Lon float64
	Lat float64 
	FormattedAddr string 
	GeocodeStatus string 
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

func GetPersons(conn *pgxpool.Pool) (linelisting Linelisting, err error) {
	sql :=
		`select bil, address, state
		   from wbk.linelisting`

	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		return 
	}

	for rows.Next() {
		var bil int 
		var address string 
		var state string 

		err = rows.Scan(&bil, &address, &state)
		if err != nil {
			return 
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

func UpdatePersonGeocodedAddr(conn *pgxpool.Pool, gpai GeocodedPersonAddrIn) error {
	sql := 
		`update wbk.linelisting
		   set lon=$1, lat=$2, 
		     formatted_address=$3,
			 geocode_status=$4
		   where bil=$5`
		   
	_, err := conn.Exec(context.Background(), sql, 
		gpai.Lon, gpai.Lat, gpai.FormattedAddr,
		gpai.GeocodeStatus, gpai.Bil)
	if err != nil {
		return err
	}
	return nil
}

func UpdateInvldPersonGeocodedAddr(conn *pgxpool.Pool, gpai GeocodedPersonAddrIn) error {
	sql := 
		`update wbk.linelisting
		   set geocode_status=$1
		   where bil=$2`
		   
	_, err := conn.Exec(context.Background(), sql, 
		gpai.GeocodeStatus, gpai.Bil)
	if err != nil {
		return err
	}
	return nil
}