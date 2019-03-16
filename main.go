package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"
)

const (
	VER       = "0.2"
	USR_URL   = "https://ham-digital.org/user_by_call.php"
	RPT_ULR   = "http://przemienniki.net/export/rxf.xml"
	Spinner   = "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
	USR_CACHE = ".usr.cache"
	THRESHOLD = "24h"
)

var quit = make(chan struct{})
var debug = false

func printVersion() {
	fmt.Println("gpcps\t", VER)
}

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

	// Main flags
	var debug = flag.Bool("d", false, "Debug mode")
	var ver = flag.Bool("v", false, "Print version")
	var cc = flag.Bool("pc", false, "Just print country codes")
	var pretty = flag.Bool("p", false, "Print pretty table")

	// Fitler flags
	var f_cc = flag.String("c", "", "Filter entries Country Code")
	var f_nm = flag.String("n", "", "Filter entries by Name")
	var f_cl = flag.String("s", "", "Filter entries by Call Sign")

	// Format
	var f_format = flag.String(
		"f",
		"{{.id}},{{.call}},{{.name}},{{.cc}}",
		"Format of the output lines")

	var pass bool
	var w *tabwriter.Writer

	flag.Parse()

	// Disable logs if not enabled explicitly
	if !*debug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	if *ver {
		printVersion()
		return
	}

	if *cc {
		printCountries()
		return
	}

	cc_rgx := rgxInit(f_cc, "country")
	nm_rgx := rgxInit(f_nm, "name")
	cl_rgx := rgxInit(f_cl, "call sign")
	raw_format := *f_format
	if *pretty { // we will print pretty table -> replace comas with \t
		raw_format = strings.ReplaceAll(raw_format, ",", "\t") + "\t"
	}
	format := template.Must(template.New("").Parse(raw_format))

	data, err := getRawData(USR_CACHE, USR_URL)
	if err != nil {
		log.Fatalf("Cannot retrieve data from %s", USR_URL)
	}

	entries := parseCsvUserData(data)
	filtered := 0

	if *pretty {
		w = tabwriter.NewWriter(os.Stdout, 3, 0, 1, ' ', tabwriter.AlignRight)
	}

	// filter entries
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
