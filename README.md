
# Final Project: Parallel Web Scraping System

**Author:** Wenlue Jeff Zhong

**Project Date:** 05/23/2024

## Overview

This project demonstrates the implementation of a parallel web scraping system designed to scrape product data from an e-commerce demo site. The system includes three versions:

1. Sequential Implementation
2. Parallel Implementation using Channels
3. Work-Stealing Implementation

## Usage

To run the benchmark script and reproduce the results, use the following command:

```bash
python benchmark/benchmark.py
```

This script will execute each implementation (sequential, parallel channel-based, and work-stealing) 10 times, calculate the average execution time, and generate performance analysis.

### Specific Implementation Execution

For examining specific implementations:

**Sequential Implementation:**

```bash
./scraper --seq
```

**Parallel Implementation using Channels:**

```bash
./scraper --chan --threads <numThreads>
```

**Work-Stealing Implementation:**

```bash
./scraper --workstealing --threads <numThreads>
```

## Final Report

The final report summarizing the results and analysis can be found at:

```bash
benchmark/final_report.pdf
```

## Project Structure

- `main.go`: The main entry point for the program.
- `sequential/sequential.go`: Contains the sequential implementation.
- `parallel/parallel.go`: Contains the parallel implementation using channels.
- `workstealing/workstealing.go`: Contains the work-stealing implementation.
- `util/util.go`: Contains utility functions for downloading images and saving to CSV.
- `benchmark/benchmark.py`: The benchmark script for running and analyzing the implementations.
- `benchmark/speedup_graph.png`: The generated speedup graph.
- `benchmark/benchmark_results.txt`: The raw results of the benchmark runs.
- `benchmark/final_report.pdf`: The final report in PDF format.

## System Requirements

- **Go:** Ensure you have Go installed on your system.
- **Python:** For running the benchmark script and generating graphs.

