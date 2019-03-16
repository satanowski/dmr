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
	CountryCode string
}

func (e DmrEntry) String() string {
	return fmt.Sprintf("ID: %d\tCall: %s\tName: %s\tCountry: %s", e.Id, e.Call, e.Name, e.CountryCode)
}

func parseDmrEntry(line string) (DmrEntry, error) {
	var cc string

	rec := strings.Split(line, ";")
	if len(rec) < 5 {
		return DmrEntry{}, errors.New("Wrong number of fields in DMR entry")
	}

	id, err := strconv.Atoi(rec[2])
	if err != nil {
		return DmrEntry{}, errors.New("Wrong DMR ID")
	}

	cc = CODES[rec[4]]
	if cc == "" {
		return DmrEntry{}, errors.New("No such country")
	}

	return DmrEntry{id, rec[1], rec[3], cc}, nil
}

func parseCsvUserData(data string) []DmrEntry {
	log.Println("Parsing the given data...\n")
	var result []DmrEntry

	for _, line := range strings.Split(data, "\n") {
		entry, err := parseDmrEntry(line)
		if err == nil {
			result = append(result, entry)
		}
	}
	return result
}
