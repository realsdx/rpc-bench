import matplotlib.pyplot as plt
import numpy as np
from scipy.ndimage import gaussian_filter1d
import seaborn as sns

def read_latencies(filename):
    with open(filename, 'r') as f:
        return [float(line.strip()) for line in f]

def microseconds_to_milliseconds(latencies):
    return [latency / 1000 for latency in latencies]

def smooth_data(data, sigma=25):
    return gaussian_filter1d(data, sigma)

def calculate_statistics(latencies):
    return {
        'avg': np.mean(latencies),
        'median': np.median(latencies),
        'p50': np.percentile(latencies, 50),
        'p90': np.percentile(latencies, 90),
        'p99': np.percentile(latencies, 99)
    }

def plot_latencies(grpc_latencies, ttrpc_latencies):
    plt.figure(figsize=(16, 10))
    sns.set_style("whitegrid")
    sns.set_palette("deep")
    
    # Smooth data
    grpc_smooth = smooth_data(grpc_latencies)
    ttrpc_smooth = smooth_data(ttrpc_latencies)
    
    # Create x-axis values
    x_grpc = range(len(grpc_smooth))
    x_ttrpc = range(len(ttrpc_smooth))
    
    # Plot smoothed lines
    sns.lineplot(x=x_grpc, y=grpc_smooth, label='gRPC', linewidth=2)
    sns.lineplot(x=x_ttrpc, y=ttrpc_smooth, label='ttrpc', linewidth=2)
    
    # Calculate statistics
    grpc_stats = calculate_statistics(grpc_latencies)
    ttrpc_stats = calculate_statistics(ttrpc_latencies)
    
    # Create a text box with statistics
    stats_text = (
        "gRPC Statistics (ms):\n"
        f"Avg: {grpc_stats['avg']:.2f}\n"
        f"Median: {grpc_stats['median']:.2f}\n"
        f"P50: {grpc_stats['p50']:.2f}\n"
        f"P90: {grpc_stats['p90']:.2f}\n"
        f"P99: {grpc_stats['p99']:.2f}"
    )
    stats_text2 = (
        "ttrpc Statistics (ms):\n"
        f"Avg: {ttrpc_stats['avg']:.2f}\n"
        f"Median: {ttrpc_stats['median']:.2f}\n"
        f"P50: {ttrpc_stats['p50']:.2f}\n"
        f"P90: {ttrpc_stats['p90']:.2f}\n"
        f"P99: {ttrpc_stats['p99']:.2f}"
    )
    
    # Add text box to the plot
    plt.text(0.80, 0.98, stats_text, transform=plt.gca().transAxes, fontsize=12,
             verticalalignment='top', bbox=dict(boxstyle='round,pad=0.5', facecolor='white', alpha=0.8))
    # Add text box to the plot
    plt.text(0.65, 0.98, stats_text2, transform=plt.gca().transAxes, fontsize=12,
             verticalalignment='top', bbox=dict(boxstyle='round,pad=0.5', facecolor='white', alpha=0.8))
    
    # Customize the plot
    plt.xlabel('Request Number', fontsize=12)
    plt.ylabel('Latency (ms)', fontsize=12)
    plt.title('Latency Comparison: gRPC vs ttrpc', fontsize=16, fontweight='bold')
    plt.legend(fontsize=10, loc='upper right')
    
    # Customize ticks
    plt.tick_params(axis='both', which='major', labelsize=10)
    
    # Improve layout
    plt.tight_layout()
    
    # Save and show the plot
    plt.savefig('output/latency_comparison.png', dpi=300, bbox_inches='tight')
    print('Latency comparison plot saved to output/latency_comparison.png')

def main(limit=-1):
    grpc_latencies = read_latencies('output/grpc_latencies.txt')
    ttrpc_latencies = read_latencies('output/ttrpc_latencies.txt')
    
    grpc_latencies_ms = microseconds_to_milliseconds(grpc_latencies[:limit])
    ttrpc_latencies_ms = microseconds_to_milliseconds(ttrpc_latencies[:limit])
    
    plot_latencies(grpc_latencies_ms, ttrpc_latencies_ms)

if __name__ == "__main__":
    import sys

    if len(sys.argv) > 1:
        main(int(sys.argv[1]))
    else:
        main()