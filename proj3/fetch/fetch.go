package fetch

import (
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Fetch(url string) (string, error) {
	c := colly.NewCollector(
		colly.Async(true),
	)

	var htmlContent string
	var fetchErr error

	c.OnHTML("html", func(e *colly.HTMLElement) {
		htmlContent, _ = e.DOM.Html()
	})

	c.OnError(func(r *colly.Response, err error) {
		fetchErr = fmt.Errorf("failed to fetch URL %s: %v", url, err)
	})

	c.Visit(url)
	c.Wait()

	return htmlContent, fetchErr
}

func SaveImage(imageURL, dir string) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to fetch image %s: %v", imageURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch image %s: %s", imageURL, resp.Status)
	}

	fileName := filepath.Base(imageURL)
	filePath := filepath.Join(dir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create image file %s: %v", filePath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save image %s: %v", filePath, err)
	}

	return nil
}
