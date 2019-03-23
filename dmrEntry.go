package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type DmrEntry struct {
	Id          int
	Call        string
	Name        string
	City        string
	State       string
	CountryCode string
}

func (e DmrEntry) String() string {
	return fmt.Sprintf("ID: %d\tCall: %s\tName: %s\tCity: %s\tState: %s\tCountry: %s", e.Id, e.Call, e.Name, e.City, e.State, e.CountryCode)
}

func getUsersCSVurl() string {
	result := ""
	doc, err := goquery.NewDocument(USR_URL)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a.readmore").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		result = USR_URL + link
	})
	return result
}

func trim(s string) string {
	return strings.Trim(s, "\" ")
}

func parseDmrEntry(line string) (DmrEntry, error) {
	var cc string

	rec := strings.Split(line, ",")

	id, err := strconv.Atoi(rec[0])
	if err != nil {
		return DmrEntry{}, errors.New("Wrong DMR ID")
	}

	cc = CODES[strings.Trim(rec[5], "\"")]

	return DmrEntry{
		id,
		trim(rec[1]),
		trim(rec[2]),
		trim(rec[3]),
		trim(rec[4]),
		cc,
	}, nil
}

func parseCsvUserData(data string) []DmrEntry {
	log.Println("Parsing the given data...\n")
	var result []DmrEntry
	var x = 0

	for _, line := range strings.Split(data, "\n")[1:] {
		x += 1
		entry, err := parseDmrEntry(line)
		if err == nil {
			result = append(result, entry)
		}
	}
	log.Printf("DEBUG: line count: %d", x)
	return result
}
