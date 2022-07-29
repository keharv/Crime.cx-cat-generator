package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load(".env")
	CheckError(err)
	// connect to postgresql DB
	var (
		host     = os.Getenv("DATABASE_HOST")
		port     = os.Getenv("DATABASE_PORT")
		user     = os.Getenv("DATABASE_USER")
		password = os.Getenv("DATABASE_PASSWORD")
		dbname   = os.Getenv("DATABASE_NAME")
	)

	// connection string
	psqlconn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected to database!")

	// insert tables if not exists
	insertImgTable := `CREATE TABLE IF NOT EXISTS images (id SERIAL PRIMARY KEY, filename VARCHAR(255) UNIQUE, md5 VARCHAR(255) UNIQUE, link_id INTEGER, time_downloaded TIMESTAMP)`
	_, err = db.Exec(insertImgTable)
	CheckError(err)
	insertLinkTable := `CREATE TABLE IF NOT EXISTS links (id SERIAL PRIMARY KEY, link VARCHAR(512) UNIQUE, scraped BOOLEAN, time_inserted TIMESTAMP, time_scraping_started TIMESTAMP, status_id INTEGER)`
	_, err = db.Exec(insertLinkTable)
	CheckError(err)

}

func main() {
	outputDir := flag.String("output", "scraped", "Output directory")
	flag.Parse()

	// get links in the database
	var link string

	for {
		link = getLinks()
		if link == "" {
			fmt.Println("No links in the database! Waiting 5 minutes for links to be added...")
			time.Sleep(5 * time.Minute)
		} else {
			// scrape for more links first
			scrapeLinks(link)
			// scrape for images
			scrape(link, *outputDir)
			doneWithLink(link)
		}
	}

}
