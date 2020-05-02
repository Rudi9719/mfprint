package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

var output string
var jobname string

func main() {
	arguments := os.Args
	if len(arguments) < 1 {
		fmt.Println("Please supply a text file to convert.")
		return
	}
	jobname = strings.Replace(arguments[0], ".txt", "", -1)
	outbytes, err := ioutil.ReadFile(arguments[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	output = string(outbytes)
	makePdf()
}

func makePdf() {
	o := output
	output = ""
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
			parts := strings.Split(line, " ")
			jobname = fmt.Sprintf("%+v-JOB%+v", parts[8], parts[6])
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
