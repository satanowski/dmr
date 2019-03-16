package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

/**
 * Retrieve the body of the page of given URL
 */
func getRawData(cacheFile, url string) (string, error) {
	data := readCacheIfValid(cacheFile)
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
		writeCache(cacheFile, data)
		close(quit)
	} else {
		log.Println("Using cached data...")
	}
	fmt.Println()
	return string(data), nil
}

func writeCache(cacheFile string, body []byte) error {
	if len(body) == 0 {
		return errors.New("Nothing to write")
	}

	err := ioutil.WriteFile(cacheFile, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readCacheIfValid(cacheFile string) []byte {
	fileStat, err := os.Stat(cacheFile)
	if err != nil {
		return nil
	}
	how_old := time.Since(fileStat.ModTime())
	threshold, _ := time.ParseDuration(THRESHOLD)
	if how_old >= threshold {
		return nil
	}
	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		log.Printf("Cannot read cache file!")
		return nil
	}
	return data
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
