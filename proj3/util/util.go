package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"proj3/parse"
)

// DownloadImage downloads an image from a URL and saves it to the specified filepath, overwriting if it already exists.
func DownloadImage(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Open the file for writing, which will overwrite the existing file if it exists
	out, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// SaveToCSV saves the products to a CSV file.
func SaveToCSV(products []parse.Product) {
	file, err := os.Create("products.csv")
	if err != nil {
		fmt.Println("Failed to create output CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	headers := []string{"url", "image", "name", "price"}
	writer.Write(headers)

	for _, product := range products {
		record := []string{product.URL, product.Image, product.Name, product.Price}
		writer.Write(record)
	}
	writer.Flush()
}
