#!/usr/bin/env python3
"""End-to-end starter readiness test for a running notifyd."""

import http.client
import json
import threading
import time
from http.server import BaseHTTPRequestHandler, HTTPServer

NOTIFY_HOST = "127.0.0.1"
NOTIFY_PORT = 18009
BACKEND_HOST = "127.0.0.1"
BACKEND_PORT = 18300
SERVICE = "starter"

received = []


class BackendHandler(BaseHTTPRequestHandler):

    def do_POST(self):
        length = int(self.headers.get("Content-Length", "0"))
        body = self.rfile.read(length)
        received.append((self.path, json.loads(body)))
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(b'{"status":"received"}')

    def log_message(self, format_string, *args):
        del format_string, args


def send_event(state, severity):
    payload = json.dumps({
        "service_name": SERVICE,
        "severity": severity,
        "time": int(time.time()),
        "module": "node",
        "name": state,
        "value": state,
        "units": "",
        "details": f"starter readiness changed to {state}",
    })
    headers = {"Content-Type": "application/json"}
    conn = http.client.HTTPConnection(NOTIFY_HOST, NOTIFY_PORT, timeout=5)
    conn.request("POST", f"/v1/event/{SERVICE}", payload, headers)
    response = conn.getresponse()
    response_body = response.read().decode("utf-8", errors="replace")
    conn.close()
    assert response.status == 202, (
        f"notifyd rejected {state}: HTTP {response.status} "
        f"body={response_body!r}"
    )


def main():
    backend = HTTPServer((BACKEND_HOST, BACKEND_PORT), BackendHandler)
    thread = threading.Thread(target=backend.serve_forever, daemon=True)
    thread.start()

    try:
        send_event("ready", "low")
        send_event("fault", "high")
    finally:
        backend.shutdown()
        backend.server_close()
        thread.join(timeout=2)

    assert len(received) == 2, f"backend received {len(received)} events"

    expected = [
        ("ready", "low", 8700),
        ("fault", "high", 8100),
    ]
    for (path, event), (state, severity, status) in zip(received, expected):
        assert path == "/node/v1/notify", path
        assert event["service_name"] == SERVICE, event
        assert event["type"] == "event", event
        assert event["severity"] == severity, event
        assert event["status"] == status, event
        assert event["details"]["name"] == state, event
        assert event["details"]["value"] == state, event

    print("PASS: notifyd accepted ready and fault with HTTP 202")
    print("PASS: backend received matching ready/fault names and values")


if __name__ == "__main__":
    main()
