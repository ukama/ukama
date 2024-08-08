import requests
import json
import time
import statistics
import argparse

# Configuration
SERVICE_NAME = "deviced"
WEB_SERVICE_ENDPOINT = "http://localhost:18009/v1/alert/deviced"
FORWARD_ENDPOINT = "http://localhost:18300/node/v1/notify"

# JSON payload
payload = {
    "serviceName": SERVICE_NAME,
    "severity": "high",
    "time": int(time.time()),
    "module": "none",
    "name": "node",
    "value": "reboot",
    "units": "",
    "details": "Node about to reboot"
}

headers = {
    'Content-Type': 'application/json'
}

def send_notification(num_requests):
    times = []
    for _ in range(num_requests):
        start_time = time.time()
        response = requests.post(WEB_SERVICE_ENDPOINT,
                                 headers=headers,
                                 data=json.dumps(payload))
        end_time = time.time()
        times.append(end_time - start_time)
    
    return times

def forward_notification(num_requests):
    times = []
    for _ in range(num_requests):
        start_time = time.time()
        response = requests.post(FORWARD_ENDPOINT,
                                 headers=headers, data=json.dumps(payload))
        end_time = time.time()
        times.append(end_time - start_time)
    
    return times

def calculate_statistics(times):
    avg_time = statistics.mean(times)
    max_time = max(times)
    min_time = min(times)
    return avg_time, max_time, min_time

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Send notifications and measure times.')
    parser.add_argument('num_requests', type=int, help='Number of requests to send')
    
    args = parser.parse_args()
    num_requests = args.num_requests

    print(f"Sending {num_requests} notifications to web service...")
    web_service_times = send_notification(num_requests)
    print(f"Sending {num_requests} notifications to forward endpoint...")
    forward_times = forward_notification(num_requests)

    web_service_avg, web_service_max, web_service_min = calculate_statistics(web_service_times)
    forward_avg, forward_max, forward_min = calculate_statistics(forward_times)

    print("\nStatistics for web service:")
    print(f"Average time: {web_service_avg:.6f} seconds")
    print(f"Max time: {web_service_max:.6f} seconds")
    print(f"Min time: {web_service_min:.6f} seconds")

    print("\nStatistics for forward endpoint:")
    print(f"Average time: {forward_avg:.6f} seconds")
    print(f"Max time: {forward_max:.6f} seconds")
    print(f"Min time: {forward_min:.6f} seconds")


