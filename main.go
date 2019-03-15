package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"
)

const (
	URL = "https://ham-digital.org/user_by_call.php"
	// Spinner = "◐◓◑◒"
	Spinner   = "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
	CACHE     = ".cache"
	THRESHOLD = "24h"
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

func writeCache(body []byte) error {
	if len(body) == 0 {
		return errors.New("Nothing to write")
	}

	err := ioutil.WriteFile(CACHE, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readCacheIfValid() []byte {
	fileStat, err := os.Stat(CACHE)
	if err != nil {
		return nil
	}
	how_old := time.Since(fileStat.ModTime())
	threshold, _ := time.ParseDuration(THRESHOLD)
	if how_old >= threshold {
		return nil
	}
	data, err := ioutil.ReadFile(CACHE)
	if err != nil {
		log.Printf("Cannot read cache file!")
		return nil
	}
	return data
}

/**
 * Retrieve the body of the page of given URL
 */
func getRawData(url string) (string, error) {
	data := readCacheIfValid()
	if data == nil {
		go spin(fmt.Sprintf("Getting the data from %s...", url))
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		writeCache(data)
		close(quit)
	} else {
		log.Println("Using cached data...")
	}
	fmt.Println()
	return string(data), nil
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

	var cc = flag.Bool("pc", false, "Just print country codes")
	var pretty = flag.Bool("p", false, "Print pretty table")
	var f_cc = flag.String("c", "", "Filter entries Country Code")
	var f_nm = flag.String("n", "", "Filter entries by Name")
	var f_cl = flag.String("s", "", "Filter entries by Call Sign")
	var f_format = flag.String("f", "{{.id}},{{.call}},{{.name}},{{.cc}}", "Format of the output lines")
	var pass bool
	var w *tabwriter.Writer

	flag.Parse()

	if *cc {
		printCountries()
		return
	}

	cc_rgx := rgxInit(f_cc, "country")
	nm_rgx := rgxInit(f_nm, "name")
	cl_rgx := rgxInit(f_cl, "call sign")
	raw_format := *f_format
	if *pretty {
		raw_format = strings.ReplaceAll(raw_format, ",", "\t") + "\t"
	}
	format := template.Must(template.New("").Parse(raw_format))

	data, err := getRawData(URL)
	if err != nil {
		log.Fatalf("Cannot retrieve data from %s", URL)
	}

	entries := parseCsvUserData(data)
	filtered := 0

	if *pretty {
		w = tabwriter.NewWriter(os.Stdout, 3, 0, 1, ' ', tabwriter.AlignRight)
	}

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

			if *pretty {
				fmt.Fprintln(w, buf.String())
			} else {
				fmt.Println(buf.String())
			}
			filtered++
		}
	}
	if *pretty {
		w.Flush()
	}
	fmt.Println()
	log.Printf("Listed %d records out of %d\n", filtered, len(entries))
}
