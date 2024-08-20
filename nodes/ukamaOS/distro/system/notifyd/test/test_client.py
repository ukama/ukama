import http.client
import json
import time
import statistics
import argparse

# Configuration
SERVICE_NAME = "deviced"
WEB_SERVICE_ENDPOINT_HOST = "localhost"
WEB_SERVICE_ENDPOINT_PATH = "/v1/alert/deviced"
WEB_SERVICE_PORT = 18009

# JSON payload
payload = json.dumps({
    "service_name": SERVICE_NAME,
    "severity": "high",
    "time": int(time.time()),
    "module": "none",
    "name": "node",
    "value": "reboot",
    "units": "",
    "details": "Node about to reboot"
})

# Headers including User-Agent
headers = {
    'Content-Type': 'application/json',
    'User-Agent': 'PythonTestClient/1.0'
}

def send_notification_httpclient(num_requests):
    times = []
    conn = http.client.HTTPConnection(WEB_SERVICE_ENDPOINT_HOST, WEB_SERVICE_PORT)

    for _ in range(num_requests):
        start_time = time.time()

        # Debugging: print the headers to ensure they are set correctly
        print(f"Sending request with headers: {headers}")

        conn.request("POST", WEB_SERVICE_ENDPOINT_PATH, payload, headers)
        response = conn.getresponse()

        # Debugging: print the status and reason to ensure the server is responding
        print(f"Response status: {response.status}, reason: {response.reason}")

        end_time = time.time()
        times.append(end_time - start_time)

        # Read and close the response to free up the connection
        response.read()

    conn.close()
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
    web_service_times = send_notification_httpclient(num_requests)

    web_service_avg, web_service_max, web_service_min = calculate_statistics(web_service_times)

    print("\nStatistics for web service:")
    print(f"Average time: {web_service_avg:.6f} seconds")
    print(f"Max time: {web_service_max:.6f} seconds")
    print(f"Min time: {web_service_min:.6f} seconds")
