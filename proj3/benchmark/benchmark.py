import subprocess
import os
import pandas as pd
import matplotlib.pyplot as plt
import time

# Define paths
project_root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
scraper_executable = os.path.join(project_root, 'scraper')

# Define the number of threads and the number of runs
threads = [2, 4, 6, 8, 12]
runs = 10
implementations = ['--seq', '--chan', '--workstealing']

# Prepare output file
output_file = os.path.join(project_root, 'benchmark', 'benchmark_results.txt')
if os.path.exists(output_file):
    os.remove(output_file)

# Function to run benchmarks
def run_benchmark(implementation, num_threads=None):
    for run in range(runs):
        if implementation == '--seq':
            cmd = [scraper_executable, implementation]
        else:
            cmd = [scraper_executable, implementation, '--threads', str(num_threads)]
        
        start_time = time.time()
        subprocess.run(cmd, check=True)
        end_time = time.time()
        
        elapsed_time = end_time - start_time
        with open(output_file, 'a') as f:
            if num_threads:
                f.write(f'{implementation}, {num_threads}, {run+1}, {elapsed_time:.3f}\n')
            else:
                f.write(f'{implementation}, 1, {run+1}, {elapsed_time:.3f}\n')

# Run benchmarks for each implementation and thread count
for implementation in implementations:
    if implementation == '--seq':
        run_benchmark(implementation)
    else:
        for thread in threads:
            run_benchmark(implementation, thread)

# Load the benchmark results into a DataFrame
data = pd.read_csv(output_file, header=None)
data.columns = ['Implementation', 'Thread Count', 'Run Number', 'Elapsed Time']

# Filter to get only the average times
avg_data = data.groupby(['Implementation', 'Thread Count'])['Elapsed Time'].mean().reset_index()

# Separate the data into sequential and parallel results
seq_data = avg_data[avg_data['Implementation'] == '--seq'].set_index('Thread Count')['Elapsed Time']
parallel_data = avg_data[avg_data['Implementation'] != '--seq']

# Calculate speedup
parallel_data['Speedup'] = parallel_data.apply(lambda x: seq_data.at[1] / x['Elapsed Time'], axis=1)

# Plot the speedup graph
plt.figure(figsize=(10, 5))
for impl in parallel_data['Implementation'].unique():
    subset = parallel_data[parallel_data['Implementation'] == impl]
    plt.plot(subset['Thread Count'], subset['Speedup'], marker='o', label=impl)

plt.title('Speedup Graph for Various Implementations')
plt.xlabel('Number of Threads')
plt.ylabel('Speedup (Relative to Sequential)')
plt.legend(title='Implementation')
plt.grid(True)
plt.savefig(os.path.join(project_root, 'benchmark', 'speedup_graph.png'))
#plt.show()
