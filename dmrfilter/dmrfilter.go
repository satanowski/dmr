package main

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "os"
    "flag"
    "strings"
)


func main() {
    cntrPtr := flag.String("country", "", "filter by country")
    callPtr := flag.String("call", "", "filter by callsign")
    flag.Parse()

    reader := csv.NewReader(bufio.NewReader(os.Stdin))
    reader.Comma =';'
    pass := false
    for {
        row, error := reader.Read()
        pass = true
        if error == io.EOF {
            break
        } else if error != nil {
            log.Fatal(error)
        }
        if len(*cntrPtr) > 0 {
            pass = pass && strings.Compare(*cntrPtr, row[3]) == 0
        }

        if len(*callPtr) > 0 {
            pass = pass && strings.Contains(row[0], *callPtr)
        }

        if pass {
            fmt.Println(strings.Join(row, ","))
        }
    }
}
