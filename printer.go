package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

var output string
var lastData = time.Now()

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	go printWhenDone()
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		// Pages are separated via \f or line feed
		//fmt.Println(scanner.Text())
		output += fmt.Sprintf("%+v\n", scanner.Text())
		lastData = time.Now()
	}
}

func printWhenDone() {
	for {
		if time.Since(lastData) > 5*time.Second && output != "" {
			makePdf()
		}
		time.Sleep(2 * time.Second)
	}
}

func delete_empty (s []string) []string {
    var r []string
    for _, str := range s {
        if str != "" {
            r = append(r, str)
        }
    }
    return r
}

func makePdf() {
	o := output
	output = ""
	outfile := "nil"
	m := pdf.NewMaroto(consts.Landscape, consts.Letter)
	scanner := bufio.NewScanner(strings.NewReader(o))
	format := props.Text{
		Size:   9,
		Family: consts.Courier,
		Top:    0,
	}
	firstRun := true
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "****A") && strings.Contains(line, "START") {
			
			parts := strings.Split(line, "  ")
			parts = delete_empty(parts)
			var jobname, timestamp string
			for i, part := range parts {
				if len(part) == 0 {
					continue
				}
				if i == 2 {
					if parts[2] == "JOB" {
						jobname = parts[4]
					} else {
						jobname = fmt.Sprintf("%+v%+v", parts[4], parts[2])
					}
				}
				if len(part) > 12 {
					ts, err := time.Parse("03.04.05 PM 02 Jan 06", strings.Trim(part, " "))
					if err != nil {
						continue
					}
					timestamp = ts.Format("2006-01-02_15-04-05")
					break
				}
				
			}
			

			outfile = fmt.Sprintf("%+v-JOB_%+v", timestamp, jobname)
		}
		page := m.GetCurrentPage()
		if strings.HasPrefix(line, "\f") && !firstRun {
			for {
				if m.GetCurrentPage() != page {
					break
				}
				m.Row(3, func() {
					m.Col(10, func() {
						m.Text("\n", format)
					})
				})
			}
		}
		firstRun = false
		m.Row(3, func() {
			m.Col(12, func() {
				m.Text(line, format)
			})
		})
	}
	err := m.OutputFileAndClose(fmt.Sprintf("%+v.pdf", outfile))
	if err != nil {
		fmt.Printf("Error closing pdf\n")
		return
	}
}
