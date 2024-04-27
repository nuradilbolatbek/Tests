package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	const url = "https://hypeauditor.com/top-instagram-all-russia/"
	data, err := parser(url)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if len(data) == 0 {
		log.Fatal("error")
	}

	err = createFile(data, "instagram_influencers.csv")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

}

func parser(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	var records [][]string

	doc.Find(".container .table .row").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 {
			return true
		}
		if i > 50 {
			return false
		}

		var record []string
		//a := 0
		s.Find(".row-cell").Each(func(j int, cell *goquery.Selection) {
			if j > 6 {
				return
			}
			text := cleaner(strings.TrimSpace(cell.Text()))
			//fmt.Printf(text)
			record = append(record, text)
			//a++
		})

		fmt.Printf("rank #%d: %+v\n", i, record)

		records = append(records, record)
		return true
	})

	//fmt.Printf("total: %d\n", len(records))

	return records, nil
}

func createFile(records [][]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("can't create: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Rank", "About", "Category", "Subscribers", "Authentic", "Engagement", "Engagement(avg."} // Update these headers as per actual data
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error: %w", err)
	}

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error: %w", err)
		}
	}

	return nil
}
func cleaner(text string) string {
	str := strings.Replace(text, "\u200c", "", -1)
	str = strings.Replace(str, "\uF336", "  name:", -1)
	str = strings.Replace(str, "\uF16D", "", -1)

	return str
}
