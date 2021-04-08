package main

import (
	"os"
	"flag"
	"time"
    "errors"
    "fmt"
	"log"
	"strings"
    "github.com/tealeg/xlsx"
	"pg_service/db"
	"myexcel/data"
	"github.com/kr/pretty"
)

var (
	mode = flag.String("mode", "", `
	Please specify a running mode for this program:
	1. parse_excel	-	read excel data into the database
	2. geocode	-	geocode the address in the database`)	
)

func usageAndExit() {
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

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
	newPerson := data.NewPersonIn{
		Bil: bil,
		Name: name,
		Ident: ident,
		Tel: tel,
		Address: address,
		State: state,
		Notifydt: notifydt,
	}	

	db.CheckDbConn()
	err = data.AddNewPerson(db.Conn, newPerson)
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

func geocode() {
	GMAP_API_KEY := os.Getenv("GMAP_API_KEY")
	if GMAP_API_KEY == "" {
		log.Fatalf("Google maps api key not found!")
	}
	client, err := maps.NewClient(maps.WithAPIKey(GMAP_API_KEY))
	check(err)

	db.CheckDbConn()
	persons, err := data.GetPersons(db.Conn)
	check(err)
	// pretty.Println(persons)

	for _, person := range persons {
		address := fmt.Sprintf("%s, %s", person.Address, person.State)	
		r := &maps.GeocodingRequest{
			Address:  address,
			Language: "en",
			Region:   "my",
		}

		resp, err := client.Geocode(context.Background(), r)
		// 'err' will not be nil if Geocode status != "OK" && != "ZERO_RESULTS".
		if err != nil {
			errFirstSplit := strings.Split(err.Error(), "-")[0]
			errSecondSplit := strings.Split(errFirstSplit, ":")[1]
			errSecondSplit = strings.TrimSpace(errSecondSplit)
			invldGeocodedPersonAddr := GeocodedPersonAddrIn{
				Bil: person.Bil,
				GeocodeStatus: errSecondSplit,
			}
			db.CheckDbConn()
			data.UpdateInvldPersonGeocodedAddr(db.Conn, invldGeocodedPersonAddr)
			continue
		}

		// Now we check for "ZERO_RESULTS" Geocode status.
		if len(resp) == 0 {
			invldGeocodedPersonAddr := GeocodedPersonAddrIn{
				Bil: person.Bil,
				GeocodeStatus: "ZERO_RESULTS",
			}
			db.CheckDbConn()
			data.UpdateInvldPersonGeocodedAddr(db.Conn, invldGeocodedPersonAddr)
			continue
		}

		// Here we print the Geocode result.
		// pretty.Println(resp)
		// fmt.Println(resp[0].FormattedAddress)
		// fmt.Printf("Longitude: %v\n", resp[0].Geometry.Location.Lng)
		// fmt.Printf("Latitude: %v\n", resp[0].Geometry.Location.Lat)
		geocodedPersonAddr := GeocodedPersonAddrIn{
			Bil: person.Bil,
			Lon: resp[0].Geometry.Location.Lng,
			Lat: resp[0].Geometry.Location.Lat,
			FormattedAddr: resp[0].FormattedAddress,
			GeocodeStatus: "OK",
		}
		db.CheckDbConn()
		data.UpdatePersonGeocodedAddr(db.Conn, geocodedPersonAddr)
	}		
}

func main() {
	/* INIT POSTGRESQL DATABASE CONNECTION */
	db.Open()
	defer db.Close()
	
	start := time.Now()
	flag.Parse()

	if *mode == "parse_excel" {
		fmt.Println("Parsing excel data into database...")
		filename := "Masterlist dari March 2021.xlsx"
		parseLinelisting(filename) 
	} else if *mode == "geocode" {
		fmt.Println("Geocoding the address in database...")
		geocode()
	} else {
		usageAndExit()
	}

	execDuration := time.Since(start) 
	fmt.Println("Done")	
	fmt.Println("Execution time: ", execDuration)
}