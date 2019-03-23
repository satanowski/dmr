package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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
		strings.Trim(rec[1], "\""),
		strings.Trim(rec[2], "\""),
		strings.Trim(rec[3], "\""),
		strings.Trim(rec[4], "\""),
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
