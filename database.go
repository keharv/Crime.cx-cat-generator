package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getLinks() string {
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

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	CheckError(err)

	// get links in the database
	var website string
	rows, err := db.Query("SELECT link FROM links WHERE status_id = 1 LIMIT 1")
	CheckError(err)
	for rows.Next() {
		var siteURL string
		err = rows.Scan(&siteURL)
		CheckError(err)
		website = siteURL

	}
	if website != "" {
		// tell everyone we are working on this url
		_, err = tx.ExecContext(ctx, "UPDATE links SET status_id = 2 WHERE link = $1", website)
		if err != nil {
			tx.Rollback()
			website = ""
		} else {
			time_scraping_started := time.Now()
			_, err = tx.ExecContext(ctx, "UPDATE links SET time_scraping_started = $1 WHERE link = $2", time_scraping_started, website)
			if err != nil {
				tx.Rollback()
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Finished getting links.", website)
	}
	return website
}

func doneWithLink(link string) {
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

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	CheckError(err)

	// get links in the database
	_, err = tx.ExecContext(ctx, "UPDATE links SET status_id = 3 WHERE link = $1", link)
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.ExecContext(ctx, "UPDATE links SET scraped = true WHERE link = $1", link)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Telling everyone we are done with the url.")
	}
}

func addLink(link string) {
	if link != "" {
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

		// get links in the database
		time_inserted := time.Now()
		updateStmt := `INSERT INTO links (link, time_inserted, status_id) VALUES ($1, $2, 1)`
		_, err = db.Exec(updateStmt, link, time_inserted)
		if err != nil {
			fmt.Println("Error adding link to database.")
		} else {
			fmt.Println("Successfully added link")
		}
	}
}

func getLinkID(link string) int {
	fmt.Println("Getting link ID")
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

	// get links in the database
	var link_id int
	query := `SELECT id FROM links WHERE link = $1`
	rows, err := db.Query(query, link)
	CheckError(err)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		CheckError(err)
		link_id = id

	}
	return link_id
}

// add image info to db
func addImage(filename string, file_hash string, link_id int) bool {
	err := godotenv.Load(".env")
	fmt.Println("Adding image to database.")
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

	// get links in the database
	time_downloaded := time.Now()
	updateStmt := `INSERT INTO images (filename, md5, link_id, time_downloaded) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(updateStmt, filename, file_hash, link_id, time_downloaded)
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return true
	}

}
