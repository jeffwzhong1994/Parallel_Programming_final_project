package sequential

import (
	"fmt"
	"proj3/parse"
	"proj3/util"
	"proj3/fetch"
	"path/filepath"
	"time"
)

func Run(urls []string) {
	start := time.Now()
	var allProducts []parse.Product

	for _, url := range urls {
		html, err := fetch.Fetch(url)
		if err != nil {
			fmt.Printf("Failed to fetch %s: %s\n", url, err)
			continue
		}
		products := parse.ParseProducts(html)
		allProducts = append(allProducts, products...)

		for _, product := range products {
			err := util.DownloadImage(product.Image, filepath.Join("scraped_image", filepath.Base(product.Image)))
			if err != nil {
				fmt.Printf("Failed to download image %s: %s\n", product.Image, err)
			}
		}
	}

	util.SaveToCSV(allProducts)
	duration := time.Since(start)
	fmt.Printf("Sequential implementation took %s\n", duration)
}
