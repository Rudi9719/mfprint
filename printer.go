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

func makePdf() {
	o := output
	output = ""
	jobname := "nil"
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
			jobname = fmt.Sprintf("%+v-JOB_%+v", parts[4], strings.ReplaceAll(parts[11], " ", "_"))
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
	err := m.OutputFileAndClose(fmt.Sprintf("%+v.pdf", jobname))
	if err != nil {
		fmt.Printf("Error closing pdf\n")
		return
	}
}
