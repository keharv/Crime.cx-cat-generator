package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

func main() {
	websites := []string{}
	listFile := flag.String("list", "websites.txt", "List of websites to scrape.")
	flag.Parse()

	// read wordlist of websites to scrape
	readFile, err := os.Open(*listFile)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	// add to website array
	for fileScanner.Scan() {
		websites = append(websites, fileScanner.Text())
		scrape(fileScanner.Text())
	}
	readFile.Close()

}
