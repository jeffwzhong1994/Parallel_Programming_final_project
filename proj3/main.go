package main

import (
	"flag"
	"fmt"
	"os"
	"proj3/parallel"
	"proj3/sequential"
	"proj3/workstealing"
)

func main() {
	seq := flag.Bool("seq", false, "Run the sequential implementation")
	chanParallel := flag.Bool("chan", false, "Run the parallel implementation using channels")
	workStealing := flag.Bool("workstealing", false, "Run the parallel implementation using work-stealing")
	numThreads := flag.Int("threads", 1, "Number of threads to use for parallel implementations")
	flag.Parse()

	baseURL := "https://www.scrapingcourse.com/ecommerce/page/%d/"
	totalPages := 12
	urls := make([]string, totalPages)
	for i := 1; i <= totalPages; i++ {
		urls[i-1] = fmt.Sprintf(baseURL, i)
	}

	switch {
	case *seq:
		fmt.Println("Running sequential implementation")
		sequential.Run(urls)
	case *chanParallel:
		fmt.Println("Running parallel implementation using channels")
		parallel.Run(urls, *numThreads)
	case *workStealing:
		fmt.Println("Running parallel implementation using work-stealing")
		workstealing.StartWorkStealing(urls, *numThreads)
	default:
		fmt.Println("Please specify a valid flag: --seq, --chan, or --workstealing")
		os.Exit(1)
	}
}
