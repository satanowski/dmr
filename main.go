package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	URL = "https://ham-digital.org/user_by_call.php"
	// Spinner = "◐◓◑◒"
	Spinner = "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
)

type DmrEntry struct {
	Id          int
	Call        string
	Name        string
	CountryCode string
}

var quit = make(chan struct{})

func spin(prefix string) {
	for {
		for _, c := range Spinner {
			select {
			case <-quit:
				return
			default:
				fmt.Printf("\r%s %c   ", prefix, c)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (e DmrEntry) String() string {
	return fmt.Sprintf("ID: %d\tCall: %s\tName: %s\tCountry: %s", e.Id, e.Call, e.Name, e.CountryCode)
}

/**
 * Retrieve the body of the page of given URL
 */
func getRawData(url string) (string, error) {
	go spin(fmt.Sprintf("Getting the data from %s...", url))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	close(quit)
	fmt.Println()
	return string(body), nil
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
	log.Println("Parsing the given data...")
	var result []DmrEntry

	for _, line := range strings.Split(data, "\n") {
		entry, err := parseDmrEntry(line)
		if err == nil {
			result = append(result, entry)
		}
	}
	return result
}

func printCountries() {
	var codes []string
	var names = map[string]string{}

	for name := range CODES {
		if names[CODES[name]] == "" {
			codes = append(codes, CODES[name])
		}
		names[CODES[name]] = name
	}

	fmt.Println("Countries:")
	sort.Strings(codes)
	for _, code := range codes {
		fmt.Printf("%s\t%s\n", code, names[code])
	}
}

func rgxInit(expr *string, kind string) *regexp.Regexp {
	if *expr == "" {
		return nil
	}
	reg, err := regexp.Compile(*expr)
	if err != nil {
		*expr = ""
		log.Fatalf("Incorrect regexp \"%s\" for %s!", *expr, kind)
	}
	return reg
}

func main() {
	type val_map map[string]interface{}

	var cc = flag.Bool("c", false, "Just print country codes")
	var f_cc = flag.String("cc", "", "Filter entries Country Code")
	var f_nm = flag.String("name", "", "Filter entries by Name")
	var f_cl = flag.String("call", "", "Filter entries by Call Sign")
	var f_format = flag.String("f", "{{.id}},{{.call}},{{.name}},{{.cc}}", "Format of the output lines")
	var pass bool

	flag.Parse()

	if *cc {
		printCountries()
		return
	}

	cc_rgx := rgxInit(f_cc, "country")
	nm_rgx := rgxInit(f_nm, "name")
	cl_rgx := rgxInit(f_cl, "call sign")
	format := template.Must(template.New("").Parse(*f_format))

	data, err := getRawData(URL)
	if err != nil {
		log.Fatalf("Cannot retrieve data from %s", URL)
	}

	entries := parseCsvUserData(data)
	for _, entry := range entries {
		pass = true
		if *f_cc != "" {
			pass = pass && cc_rgx.Match([]byte(entry.CountryCode))
		}

		if *f_nm != "" {
			pass = pass && nm_rgx.Match([]byte(entry.Name))
		}

		if *f_cl != "" {
			pass = pass && cl_rgx.Match([]byte(entry.Call))
		}

		if pass {
			buf := bytes.Buffer{}
			format.Execute(&buf, val_map{
				"id":   entry.Id,
				"call": entry.Call,
				"name": entry.Name,
				"cc":   entry.CountryCode,
			})
			fmt.Println(buf.String())
		}
	}
	log.Printf("Retrieved %d records\n", len(entries))
}
