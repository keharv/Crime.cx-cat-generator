package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generateMD5Hash(filename string) string {
	file, err := os.Open(filename)
	CheckError(err)
	defer file.Close()

	// get hash
	hash := md5.New()
	_, err = io.Copy(hash, file)
	CheckError(err)
	// convert from bytes to string
	md5String := hex.EncodeToString(hash.Sum(nil)[:])
	return md5String
}

// generate random string with length of n and possibilities of 62 * n
func generateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func scrapeLinks(url string) {
	c := colly.NewCollector(
		// max  depth because it will go on forever if not
		colly.MaxDepth(10),
	)
	// rate limit to be nice. dont want to get blacklisted or banned
	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "*",
		// Set a delay between requests to these domains
		Delay: 300 * time.Millisecond,
		// Add an additional random delay
		RandomDelay: 300 * time.Millisecond,
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		addLink(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.Visit(url)
}

func scrape(domain string, outputDir string) {
	fmt.Println("Checking", domain)

	c := colly.NewCollector(
		// max  depth because it will go on forever if not
		colly.MaxDepth(10),
	)
	// rate limit to be nice. dont want to get blacklisted or banned
	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "*",
		// Set a delay between requests to these domains
		Delay: 300 * time.Millisecond,
		// Add an additional random delay
		RandomDelay: 300 * time.Millisecond,
	})

	imageDownloader := c.Clone()
	var link_id = getLinkID(domain)
	fmt.Println("done getting link id")
	// download images
	imageDownloader.OnResponse(func(r *colly.Response) {
		requestedURL := r.Request.URL.String()
		filename := generateRandomString(100)
		dir := outputDir + "/"
		var ext string
		if filepath.Ext(requestedURL) == "" {
			contentType := r.Headers.Get("Content-Type")
			ext = "." + strings.Split(contentType, "/")[1]
		} else {
			//make sure there are no invalid chars
			if strings.Contains(filepath.Ext(requestedURL), "?") {
				ext = strings.Split(filepath.Ext(requestedURL), "?")[0]
			} else {
				ext = filepath.Ext(requestedURL)
			}
		}
		final_filename := dir + filename + ext
		// save image
		err := r.Save(final_filename)
		if err != nil {
			log.Println("Error downloading picture.")
			log.Println(err)
			return
		}

		// check md5 of file
		md5_hash := generateMD5Hash(final_filename)
		// add to database
		if !addImage(filename+ext, md5_hash, link_id) {
			err := os.Remove(final_filename)
			if err != nil {
				fmt.Println("Error duplicate deleting file", final_filename)
			}
		}
	})

	// find all images
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		imageDownloader.Visit(e.Request.AbsoluteURL(e.Attr("src")))
	})
	scrapeLinks(domain)
	fmt.Println("Scraping for images.")
	c.Visit(domain)

}
