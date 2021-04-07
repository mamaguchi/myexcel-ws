package main

import (
	"time"
    "errors"
    "fmt"
	"log"
	"strings"
    "github.com/tealeg/xlsx"
	"pg_service/db"
	"myexcel/excel"
)

/* tealeg xlsx package tutorial */
func cellVisitor(c *xlsx.Cell) error {
    value, err := c.FormattedValue()
    if err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println("Cell value:", value)
    }
    return err
}

func rowVisitor(r *xlsx.Row) error {
    return r.ForEachCell(cellVisitor)
}

func rowStuff() {
    filename := "samplefile.xlsx"
    wb, err := xlsx.OpenFile(filename)
    if err != nil {
        panic(err)
    }
    sh, ok := wb.Sheet["Sample"]
    if !ok {
        panic(errors.New("Sheet not found"))
    }
    fmt.Println("Max row is", sh.MaxRow)
    sh.ForEachRow(rowVisitor)
}

/* My Code */
func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

func rowVisitorLinelisting(r *xlsx.Row) error {
	if r.GetCoordinate() == 0 { return nil }

	bil, err := r.GetCell(xlsx.ColLettersToIndex("A")).Int()
	check(err)
	notifydt := r.GetCell(xlsx.ColLettersToIndex("AD")).String()
    name, err := r.GetCell(xlsx.ColLettersToIndex("E")).FormattedValue()
	check(err)
    ident, err := r.GetCell(xlsx.ColLettersToIndex("F")).FormattedValue()
	check(err)
    tel, err := r.GetCell(xlsx.ColLettersToIndex("K")).FormattedValue()
	check(err)
    address, err := r.GetCell(xlsx.ColLettersToIndex("J")).FormattedValue()
	check(err)
	state, err := r.GetCell(xlsx.ColLettersToIndex("L")).FormattedValue()
	check(err)

	ident = strings.ReplaceAll(ident, "-", "")

	// fmt.Printf("%v, %v, %v, %v\n", name, ident, tel, address)
	newPerson := excel.NewPersonIn{
		Bil: bil,
		Name: name,
		Ident: ident,
		Tel: tel,
		Address: address,
		State: state,
		Notifydt: notifydt,
	}	

	db.CheckDbConn()
	err = excel.AddNewPerson(db.Conn, newPerson)
    check(err) 

	return nil
}

func parseLinelisting(excelFile string) {	
	wb, err := xlsx.OpenFile(excelFile)
	if err != nil {
		panic(err)
	}
	sh, ok := wb.Sheet["LINELISTING"]
	if !ok {
		panic(errors.New("Sheet 'LINELISTING' not found"))
	}
	fmt.Println("Max row is", sh.MaxRow)
	sh.ForEachRow(rowVisitorLinelisting)
}

func main() {
    /* tealeg xlsx package tutorial */
    // rowStuff()

	/* INIT DATABASE CONNECTION */
	db.Open()
	defer db.Close()
	
	start := time.Now()

	// filename := "Masterlist Test.xlsx"
	filename := "Masterlist dari March 2021.xlsx"
	parseLinelisting(filename)   

	execDuration := time.Since(start) 	
	fmt.Println("Execution time: ", execDuration)
}