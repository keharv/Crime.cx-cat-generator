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

func generateMD5Hash(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// get hash
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Fatal(err)
	}
	// convert from bytes to string
	md5String := hex.EncodeToString(hash.Sum(nil)[:])
	return md5String
}

func generateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	imageCollector := c.Clone()
	imageDownloader := c.Clone()

	// download images
	imageDownloader.OnResponse(func(r *colly.Response) {
		requestedURL := r.Request.URL.String()
		filename := generateRandomString(10)
		dir := outputDir + "/"
		fmt.Println(requestedURL)
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
		err := r.Save(final_filename)
		if err != nil {
			log.Println("Error downloading picture.")
			log.Println(err)
			return
		}

		// check md5 of file
		md5_hash := generateMD5Hash(final_filename)
		log.Println("MD5:", md5_hash)
	})

	// find all images
	imageCollector.OnHTML("img[src]", func(e *colly.HTMLElement) {
		foundImages := e.Request.AbsoluteURL(e.Attr("src"))
		imageDownloader.Visit(foundImages)
	})
	// find all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		scrapedLinks := e.Request.AbsoluteURL(e.Attr("href"))
		imageCollector.Visit(scrapedLinks)
		c.Visit(scrapedLinks)
	})

	c.Visit(domain)

}
