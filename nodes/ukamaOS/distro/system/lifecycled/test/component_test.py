#!/usr/bin/env python3

"""Component test for lifecycle.d using fake starter and notify services."""

import http.client
import json
import os
import socket
import subprocess
import tempfile
import threading
import time
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer


class HarnessState:
    def __init__(self):
        self.lock = threading.Lock()
        self.aggregate = "pending"
        self.aggregate_reason = "applications are starting"
        self.config_phase = "awaiting_configuration"
        self.config_request_id = ""
        self.events = []

    def starter_status(self):
        with self.lock:
            entry = {
                "space": "services",
                "name": "configd",
                "service": "config",
                "state": (
                    "pending"
                    if self.config_phase == "configuration_in_progress"
                    else "faulty"
                    if self.config_phase == "configuration_failed"
                    else "ready"
                ),
                "httpStatus": (
                    202
                    if self.config_phase == "configuration_in_progress"
                    else 503
                    if self.config_phase == "configuration_failed"
                    else 200
                ),
                "reason": self.config_phase,
                "checkedAt": int(time.time()),
            }
            if self.config_request_id:
                entry["requestId"] = self.config_request_id

            return {
                "spaces": [],
                "starterd": {
                    "readiness": {
                        "enabled": True,
                        "state": self.aggregate,
                        "reason": self.aggregate_reason,
                        "apps": [entry],
                    }
                },
            }

    def set_starter(self, aggregate, config_phase, request_id=""):
        with self.lock:
            self.aggregate = aggregate
            self.aggregate_reason = aggregate
            self.config_phase = config_phase
            self.config_request_id = request_id

    def add_event(self, event):
        with self.lock:
            self.events.append(event)

    def event_values(self):
        with self.lock:
            return [event.get("value") for event in self.events]


class HarnessHandler(BaseHTTPRequestHandler):
    state = None
    role = None

    def log_message(self, _format, *_args):
        return

    def send_json(self, status, body):
        encoded = json.dumps(body).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(encoded)))
        self.end_headers()
        self.wfile.write(encoded)

    def do_GET(self):
        if self.role == "starter" and self.path == "/v1/status":
            self.send_json(200, self.state.starter_status())
            return

        self.send_json(404, {"error": "not found"})

    def do_POST(self):
        length = int(self.headers.get("Content-Length", "0"))
        body = self.rfile.read(length)

        if self.role == "notify" and self.path == "/v1/event/lifecycle":
            self.state.add_event(json.loads(body.decode("utf-8")))
            self.send_json(202, {"status": "accepted"})
            return

        self.send_json(404, {"error": "not found"})


def start_server(role, state):
    handler = type(
        f"{role.title()}Handler",
        (HarnessHandler,),
        {"state": state, "role": role},
    )
    server = ThreadingHTTPServer(("127.0.0.1", 0), handler)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    return server


def unused_port():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.bind(("127.0.0.1", 0))
    port = sock.getsockname()[1]
    sock.close()
    return port


def request(port, method, path, body=None):
    connection = http.client.HTTPConnection("127.0.0.1", port, timeout=2)
    encoded = json.dumps(body) if body is not None else None
    headers = {"Content-Type": "application/json"} if body is not None else {}
    connection.request(method, path, body=encoded, headers=headers)
    response = connection.getresponse()
    payload = response.read().decode("utf-8")
    connection.close()
    return response.status, json.loads(payload) if payload else {}


def wait_http(port, process, timeout=5):
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        if process.poll() is not None:
            output = process.stdout.read() if process.stdout else ""
            raise AssertionError(f"lifecycle.d exited early:\n{output}")
        try:
            status, _ = request(port, "GET", "/v1/ping")
            if status == 200:
                return
        except (ConnectionError, OSError, TimeoutError):
            pass
        time.sleep(0.05)
    raise AssertionError("lifecycle.d did not start its HTTP service")


def wait_state(port, expected, timeout=5):
    deadline = time.monotonic() + timeout
    last = None
    while time.monotonic() < deadline:
        status, body = request(port, "GET", "/v1/status")
        if status == 200:
            last = body.get("state")
            if last == expected:
                return body
        time.sleep(0.05)
    raise AssertionError(f"expected state {expected}, last state was {last}")


def wait_config_received(port, timeout=5):
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        _, body = request(port, "GET", "/v1/status")
        if body.get("configuration", {}).get("received"):
            return
        time.sleep(0.05)
    raise AssertionError("lifecycle did not observe configuration in progress")


def assert_subsequence(values, expected):
    cursor = 0
    for value in values:
        if cursor < len(expected) and value == expected[cursor]:
            cursor += 1
    if cursor != len(expected):
        raise AssertionError(
            f"event sequence missing; expected {expected}, received {values}"
        )


def run():
    binary = os.environ.get("LIFECYCLED_BIN")
    if not binary:
        raise SystemExit("LIFECYCLED_BIN is required")

    binary = os.path.abspath(binary)
    state = HarnessState()
    starter = start_server("starter", state)
    notify = start_server("notify", state)
    lifecycle_port = unused_port()

    with tempfile.TemporaryDirectory(prefix="lifecycled-test-") as temp_dir:
        environment = os.environ.copy()
        environment.update(
            {
                "LIFECYCLED_HTTP_ADDR": "127.0.0.1",
                "LIFECYCLED_HTTP_PORT": str(lifecycle_port),
                "LIFECYCLED_STARTER_HOST": "127.0.0.1",
                "LIFECYCLED_STARTER_PORT": str(starter.server_port),
                "LIFECYCLED_NOTIFY_HOST": "127.0.0.1",
                "LIFECYCLED_NOTIFY_PORT": str(notify.server_port),
                "LIFECYCLED_STATE_FILE": os.path.join(temp_dir, "state"),
                "LIFECYCLED_CHECKIN_TIMEOUT_SEC": "1",
                "LIFECYCLED_CONFIG_TIMEOUT_SEC": "1",
                "LIFECYCLED_STARTER_UNAVAILABLE_TIMEOUT_SEC": "2",
                "LIFECYCLED_POLL_INTERVAL_MS": "100",
                "LIFECYCLED_REQUEST_TIMEOUT_SEC": "1",
                "LIFECYCLED_LOG_LEVEL": "error",
            }
        )

        process = subprocess.Popen(
            [binary],
            env=environment,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
        )

        try:
            wait_http(lifecycle_port, process)
            wait_state(lifecycle_port, "STARTING")

            status, _ = request(
                lifecycle_port,
                "POST",
                "/v1/check-in",
                {"bootResult": "ready"},
            )
            assert status == 202
            wait_state(lifecycle_port, "CHECKING_IN")

            deadline = time.monotonic() + 3
            while time.monotonic() < deadline:
                status, gate = request(lifecycle_port, "GET", "/v1/gate")
                if status == 200 and gate.get("proceed"):
                    break
                time.sleep(0.05)
            else:
                raise AssertionError("check-in gate did not open")

            state.set_starter("ready", "awaiting_configuration")
            wait_state(lifecycle_port, "READY")

            time.sleep(1.2)
            wait_state(lifecycle_port, "READY")

            status, _ = request(
                lifecycle_port,
                "POST",
                "/v1/configure",
                {"requestId": "no-config", "assignmentId": "site-1"},
            )
            assert status == 202
            wait_state(lifecycle_port, "CONFIGURING")

            status, _ = request(
                lifecycle_port,
                "POST",
                "/v1/configure",
                {"requestId": "no-config", "assignmentId": "site-1"},
            )
            assert status == 200
            wait_state(lifecycle_port, "OPERATIONAL")

            status, _ = request(
                lifecycle_port,
                "POST",
                "/v1/configure",
                {"requestId": "with-config", "assignmentId": "site-1"},
            )
            assert status == 202
            wait_state(lifecycle_port, "CONFIGURING")

            state.set_starter(
                "pending", "configuration_in_progress", "with-config"
            )
            wait_config_received(lifecycle_port)
            state.set_starter("ready", "configuration_applied", "with-config")
            wait_state(lifecycle_port, "OPERATIONAL")

            status, _ = request(
                lifecycle_port,
                "POST",
                "/v1/configure",
                {"requestId": "bad-config", "assignmentId": "site-1"},
            )
            assert status == 202
            state.set_starter("faulty", "configuration_failed", "bad-config")
            wait_state(lifecycle_port, "FAULTY")

            expected = [
                "STARTING",
                "CHECKING_IN",
                "READY",
                "CONFIGURING",
                "OPERATIONAL",
                "CONFIGURING",
                "OPERATIONAL",
                "CONFIGURING",
                "FAULTY",
            ]
            deadline = time.monotonic() + 3
            while time.monotonic() < deadline:
                try:
                    assert_subsequence(state.event_values(), expected)
                    break
                except AssertionError:
                    time.sleep(0.05)
            else:
                assert_subsequence(state.event_values(), expected)

            print("PASS: lifecycle component flow")
        finally:
            process.terminate()
            try:
                process.wait(timeout=3)
            except subprocess.TimeoutExpired:
                process.kill()
                process.wait(timeout=3)
            starter.shutdown()
            notify.shutdown()
            starter.server_close()
            notify.server_close()


if __name__ == "__main__":
    run()
