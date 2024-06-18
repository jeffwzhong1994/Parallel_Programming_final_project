package parallel

import (
	"fmt"
	"proj3/parse"
	"proj3/util"
	"proj3/fetch"
	"path/filepath"
	"time"
)

func Run(urls []string, numWorkers int) {
	start := time.Now()

	urlChan := make(chan string, len(urls))
	results := make(chan []parse.Product, len(urls))
	doneChan := make(chan struct{}, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go worker(i, urlChan, results, doneChan)
	}

	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	// Wait for all workers to complete
	for i := 0; i < numWorkers; i++ {
		<-doneChan
	}
	close(results)

	var allProducts []parse.Product
	for products := range results {
		allProducts = append(allProducts, products...)
	}

	util.SaveToCSV(allProducts)
	duration := time.Since(start)
	fmt.Printf("Parallel implementation using channels took %s\n", duration)
}

func worker(id int, urls <-chan string, results chan<- []parse.Product, doneChan chan<- struct{}) {
	for url := range urls {
		html, err := fetch.Fetch(url)
		if err != nil {
			fmt.Printf("Worker %d: Failed to fetch %s: %s\n", id, url, err)
			continue
		}
		products := parse.ParseProducts(html)
		results <- products

		imgDone := make(chan struct{})
		go downloadImages(products, imgDone, id)

		<-imgDone // Wait for image downloads to complete
	}
	doneChan <- struct{}{}
}

func downloadImages(products []parse.Product, done chan<- struct{}, id int) {
	var imgCount = len(products)
	imgDone := make(chan struct{}, imgCount)

	for _, product := range products {
		go func(imageURL string) {
			defer func() { imgDone <- struct{}{} }()
			err := util.DownloadImage(imageURL, filepath.Join("scraped_image", filepath.Base(imageURL)))
			if err != nil {
				fmt.Printf("Worker %d: Failed to download image %s: %s\n", id, imageURL, err)
			}
		}(product.Image)
	}

	// Wait for all image downloads to complete
	for i := 0; i < imgCount; i++ {
		<-imgDone
	}
	close(done)
}
