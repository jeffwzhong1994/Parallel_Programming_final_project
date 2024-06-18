package workstealing

import (
	"fmt"
	"proj3/parse"
	"proj3/util"
	"proj3/fetch"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Worker(id int, deques []*Deque, results chan<- []parse.Product, done <-chan struct{}, emptyTasks chan<- int, imgWg *sync.WaitGroup, taskProcessed *sync.WaitGroup) {
	defer wg.Done()
	myDeque := deques[id]

	imageFolder := "scraped_image"
	os.MkdirAll(imageFolder, os.ModePerm)

	taskCount := 0
	productCount := 0

	for {
		select {
		case <-done:
			return
		default:
			task, ok := myDeque.PopBottom()
			if !ok {
				// Attempt to steal from another deque
				stolen := false
				for i := 0; i < len(deques); i++ {
					if i == id {
						continue
					}
					time.Sleep(100 * time.Millisecond) // Delay to observe steps
					task, ok = deques[i].PopTop()
					if ok {
						stolen = true
						break
					}
				}
				if !stolen {
					emptyTasks <- id // Signal that this worker has no tasks left
					continue
				}
			}

			taskProcessed.Add(1)
			taskCount++
			html, err := fetch.Fetch(task)
			if err != nil {
				fmt.Printf("Worker %d: Failed to fetch %s: %s\n", id, task, err)
				taskProcessed.Done()
				continue
			}
			products := parse.ParseProducts(html)
			productCount += len(products)
			results <- products

			// Download images in parallel
			for _, product := range products {
				imgWg.Add(1)
				go func(imageURL string) {
					defer imgWg.Done()
					err := util.DownloadImage(imageURL, filepath.Join(imageFolder, filepath.Base(imageURL)))
					if err != nil {
						fmt.Printf("Worker %d: Failed to download image %s: %s\n", id, imageURL, err)
					}
				}(product.Image)
			}
			taskProcessed.Done()
		}
	}
}

func StartWorkStealing(urls []string, numWorkers int) {
	start := time.Now()
	deques := make([]*Deque, numWorkers)
	for i := range deques {
		deques[i] = &Deque{}
	}

	done := make(chan struct{})
	results := make(chan []parse.Product, numWorkers)
	emptyTasks := make(chan int, numWorkers)
	var imgWg sync.WaitGroup
	var taskProcessed sync.WaitGroup
	var once sync.Once

	for i, url := range urls {
		deques[i%numWorkers].PushBottom(url)
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go Worker(i, deques, results, done, emptyTasks, &imgWg, &taskProcessed)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allProducts []parse.Product
	allDone := false
	emptyCount := 0
	for !allDone {
		select {
		case products := <-results:
			if products != nil {
				allProducts = append(allProducts, products...)
			}
		case <-emptyTasks:
			emptyCount++
			if emptyCount == numWorkers {
				once.Do(func() {
					close(done)
				})
				allDone = true
			}
		}
	}

	imgWg.Wait()         // Ensure all image downloads are complete
	taskProcessed.Wait() // Ensure all tasks are processed

	util.SaveToCSV(allProducts)
	duration := time.Since(start)
	fmt.Printf("Work-stealing implementation took %s\n", duration)
	fmt.Printf("Total products scraped: %d\n", len(allProducts))
}
