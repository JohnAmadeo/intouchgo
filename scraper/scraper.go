package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/johnamadeo/intouchgo/models"
	"golang.org/x/net/html"
)

const (
	CTHomePage      = "http://www.ctinmateinfo.state.ct.us/"
	CTSearchForm    = "#frmSearchOp"
	CTLastNameInput = "#frmSearchOp tr:nth-of-type(5) td:nth-of-type(2) input"
	CTSubmitButton  = "#submit1"
	CTInmateTable   = "table[summary='Result.']"

	AlphabetSize = 26
)

func findInmatesByLastName(letter string, html *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(CTHomePage),
		chromedp.WaitVisible(CTSearchForm, chromedp.NodeVisible),
		chromedp.SendKeys(CTLastNameInput, letter),
		chromedp.Click(CTSubmitButton, chromedp.NodeVisible),
		// Wait for all the table rows to load; assumes the network connection
		// is fast enough that this will occur after 1 second
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(CTInmateTable, chromedp.NodeVisible),
		chromedp.OuterHTML(CTInmateTable, html, chromedp.NodeVisible),
	}
}

func formatName(rawName string) (string, string) {
	names := strings.SplitN(rawName, ",", 2)

	style := func(name string) string {
		return strings.Title(strings.ToLower(strings.TrimSpace(name)))
	}

	firstName := style(names[1])
	lastName := style(names[0])
	return firstName, lastName
}

func extractInmatesFromHTML(html string) ([]models.Inmate, error) {
	inmates := []models.Inmate{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return inmates, err
	}

	trs := doc.Find("tr")

	for _, tr := range trs.Nodes {
		tds := nodeToSelection(tr).Find("td")

		if len(tds.Nodes) == 4 {
			var inmateNumber, firstName, lastName, dateOfBirth, facility string
			for i, td := range tds.Nodes {
				text := nodeToSelection(td).Text()
				switch i {
				case 0:
					inmateNumber = text
				case 1:
					firstName, lastName = formatName(text)
				case 2:
					dateOfBirth = text
				case 3:
					facility = text
				}
			}

			inmates = append(inmates, models.Inmate{
				InmateNumber: inmateNumber,
				FirstName:    firstName,
				LastName:     lastName,
				DateOfBirth:  dateOfBirth,
				Facility:     facility,
			})
		}
	}

	return inmates, nil
}

func nodeToSelection(node *html.Node) *goquery.Selection {
	return &goquery.Selection{
		Nodes: []*html.Node{node},
	}
}

func printInmates(inmates []models.Inmate) {
	if len(inmates) > 0 {
		fmt.Println(inmates[0].LastName[:1], " : ", len(inmates))
	}
}

func getInmatesByLastName(
	ctxt context.Context,
	chrome *chromedp.CDP,
	letter string,
) []models.Inmate {
	fmt.Println("Scraping all inmates whose last name start with " + letter)

	var html string
	err := chrome.Run(ctxt, findInmatesByLastName(letter, &html))
	if err != nil {
		fmt.Println(err)
		return []models.Inmate{}
	}

	inmates, err := extractInmatesFromHTML(html)
	if err != nil {
		fmt.Println(err)
		return []models.Inmate{}
	}

	return inmates
}

func main() {
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	chrome, err := chromedp.New(ctxt /*chromedp.WithLog(log.Printf)*/)
	if err != nil {
		log.Fatal(err)
	}

	inmates := []models.Inmate{}
	for i := 65; i < 65+AlphabetSize; i++ {
		letter := string(i)
		letterInmates := getInmatesByLastName(ctxt, chrome, letter)
		printInmates(letterInmates)
		inmates = append(inmates, letterInmates...)
	}

	fmt.Println("All: ", len(inmates))

	err = chrome.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	err = chrome.Wait()
	if err != nil {
		log.Fatal(err)
	}

}
