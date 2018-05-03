package data

import (
	"database/sql"
	"fmt"
	"hatchery/models"
	"log"
	"os"
	"time"

	// fallback for the sql package
	_ "github.com/lib/pq"
)

var eggs []models.Egg
var index = 0
var db *sql.DB
var pw, dbname, user string

//TODO: add function to fetch into memory
func InitializeRepository() {
	//fetch initial data, start timer for further fetching etc
	pw = os.Getenv("PGPASSWORD")
	dbname = os.Getenv("PGDBNAME")
	user = os.Getenv("PGUSER")

	connect()

	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS $1.eggs (
		id SERIAL PRIMARY KEY,
		author character varying(255) NOT NULL,
		target character varying(255) NOT NULL,
		message character varying(255),
		picture text,
		layed timestamp with time zone NOT NULL,
		hatchtime timestamp with time zone NOT NULL,
		)`)

	if err != nil {
		log.Fatal(err)
		fmt.Print("I broke")
	}

	_, execErr := stmt.Exec(dbname)

	if execErr != nil {
		log.Fatal(execErr)
		fmt.Print("I broke")
	}

	fmt.Print("We made the table!")

	db.Close()
}

func connect() {
	connStr := "user=" + user + " password=" + pw + " dbname=" + dbname + " sslmode=verify-full"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

//StoreEgg will put an egg into either memory or db storage depending on how far out it should hatch
func StoreEgg(egg models.Egg) {

	egg.Id = index
	index++
	// db insert
	if egg.HatchTime.After(time.Now().Add(time.Minute * 15)) {
		connect()
		db.Query(fmt.Sprintf("INSERT INTO eggs VALUES (%s, %s, %s, %s, %s, %s);", egg.Author, egg.Target, egg.Message, egg.Picture, egg.Layed, egg.HatchTime))

		db.Close()
	} else {
		// in memory insert
		for i, item := range eggs {
			if egg.HatchTime.Before(item.HatchTime) {
				eggs = append(eggs[:i], append([]models.Egg{egg}, eggs[i:]...)...)
				return
			}
		}
		eggs = append(eggs, egg)
	}
}

//RetrieveEgg gets an egg by Id
func RetrieveEgg(id int) models.Egg {
	//do nothing atm
	var egg models.Egg
	return egg
}

//RetrieveEggs gets all eggs
func RetrieveEggs() []models.Egg {
	if eggs == nil {
		return []models.Egg{}
	}
	return eggs
}

//RemoveEgg removes an egg by Id
func RemoveEgg(id int) {
	eggs = eggs[1:]
}

//RemoveEggs removes n eggs from memory
func RemoveEggs(number int) {
	eggs = eggs[number:]
}
