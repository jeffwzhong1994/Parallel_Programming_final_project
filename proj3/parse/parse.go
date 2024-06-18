package parse

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Product struct {
	URL   string
	Image string
	Name  string
	Price string
}

func ParseProducts(html string) []Product {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("Error creating document from reader:", err)
		return nil
	}

	var products []Product
	doc.Find("li.product").Each(func(i int, s *goquery.Selection) {
		product := Product{
			URL:   s.Find("a").AttrOr("href", ""),
			Image: s.Find("img").AttrOr("src", ""),
			Name:  s.Find("h2").Text(),
			Price: s.Find(".price").Text(),
		}
		products = append(products, product)
	})

	return products
}

func SaveToCSV(products []Product) {
	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"URL", "Image", "Name", "Price"}
	writer.Write(headers)

	for _, product := range products {
		record := []string{product.URL, product.Image, product.Name, product.Price}
		writer.Write(record)
	}
}
