#!/usr/bin/env python3
"""Run lifecycle.d with HTTP peers for starter, configd and notifyd."""
import http.client
import json
import os
from pathlib import Path
import socket
import subprocess
import tempfile
import threading
import time
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer


class Peers:
    def __init__(self):
        self.lock = threading.Lock()
        self.ready = "ready"
        self.config_available = True
        self.notify_available = False
        self.events = []
        self.config = dict(schemaVersion=1, mode="NONE", phase="awaiting",
                           requestId="", generation=0, revision=0, error="")

    def decision(self, generation):
        with self.lock:
            self.config.update(mode="NOCONFIG", phase="completed",
                               requestId="assignment-1", generation=generation)


class Handler(BaseHTTPRequestHandler):
    def log_message(self, *_args):
        pass

    def reply(self, code, body):
        data = json.dumps(body).encode()
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    def do_GET(self):
        with self.peers.lock:
            if self.role == "starter" and self.path == "/v1/status":
                self.reply(200, {"starterd": {"readiness": {
                    "state": self.peers.ready, "reason": self.peers.ready}}})
            elif self.role == "config" and self.path == "/v1/config/status":
                self.reply(200 if self.peers.config_available else 503, self.peers.config)
            else:
                self.reply(404, {})

    def do_POST(self):
        body = self.rfile.read(int(self.headers.get("Content-Length", 0)))
        with self.peers.lock:
            if self.role == "notify" and self.path == "/v1/event/lifecycle":
                if self.peers.notify_available:
                    event = json.loads(body)
                    event["metadata"] = json.loads(event["details"])
                    self.peers.events.append(event)
                    self.reply(202, {})
                else:
                    self.reply(503, {})
            else:
                self.reply(404, {})


def server(role, peers):
    handler = type(role, (Handler,), dict(role=role, peers=peers))
    result = ThreadingHTTPServer(("127.0.0.1", 0), handler)
    threading.Thread(target=result.serve_forever, daemon=True).start()
    return result


def unused_port():
    with socket.socket() as sock:
        sock.bind(("127.0.0.1", 0))
        return sock.getsockname()[1]


def request(port, method="GET", path="/v1/status", body=None):
    connection = http.client.HTTPConnection("127.0.0.1", port, timeout=3)
    data = None if body is None else json.dumps(body)
    connection.request(method, path, data, {"Content-Type": "application/json"})
    response = connection.getresponse()
    result = response.status, json.loads(response.read())
    connection.close()
    return result


def wait(check, message, timeout=8):
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        try:
            if check():
                return
        except (OSError, ValueError):
            pass
        time.sleep(0.05)
    raise AssertionError(message)


def run():
    peers = Peers()
    servers = {role: server(role, peers) for role in ("starter", "config", "notify")}
    port = unused_port()
    binary = str(Path(os.environ["LIFECYCLED_BIN"]).resolve())
    process = None

    with tempfile.TemporaryDirectory(prefix="lifecycle-http-") as directory:
        checkpoint = Path(directory) / "checkpoint.json"
        boot_file = Path(directory) / "boot-id"
        boot_file.write_text("boot-1\n")
        environment = dict(os.environ)
        environment.update(LIFECYCLED_HTTP_PORT=str(port),
                           LIFECYCLED_STATE_FILE=str(checkpoint),
                           LIFECYCLED_BOOT_ID_FILE=str(boot_file),
                           LIFECYCLED_CHECKIN_TIMEOUT_SEC="1",
                           LIFECYCLED_STARTER_UNAVAILABLE_TIMEOUT_SEC="1",
                           LIFECYCLED_CONFIG_UNAVAILABLE_TIMEOUT_SEC="1",
                           LIFECYCLED_POLL_INTERVAL_MS="50",
                           LIFECYCLED_REQUEST_TIMEOUT_SEC="1")
        for role, peer in servers.items():
            environment[f"LIFECYCLED_{role.upper()}_PORT"] = str(peer.server_port)

        def start():
            nonlocal process
            process = subprocess.Popen([binary], env=environment,
                                       stdout=subprocess.DEVNULL)
            wait(lambda: request(port)[0] == 200, "HTTP startup")

        def stop():
            nonlocal process
            if process:
                process.kill()
                process.wait(timeout=5)
                process = None

        def state(value):
            wait(lambda: request(port)[1]["state"] == value, f"state {value}")

        def event_count(value, boot="boot-1"):
            with peers.lock:
                return sum(event["value"] == value and event["metadata"]["bootId"] == boot
                           for event in peers.events)

        try:
            start()
            assert request(port, "POST", "/v1/check-in", {"bootResult": "ready"})[0] == 202
            state("READY")
            time.sleep(1.2)
            assert request(port)[1]["state"] == "READY"
            assert request(port, "POST", "/v1/configure", {"requestId": "wrong"})[0] == 409
            peers.decision(1)
            state("OPERATIONAL")
            saved = json.loads(checkpoint.read_text())
            assert [item["state"] for item in saved["events"]] == [0, 2, 3, 4]
            stop()
            start()
            with peers.lock:
                peers.notify_available = True
            wait(lambda: event_count("OPERATIONAL") >= 1, "queued completion after crash")
            with peers.lock:
                assert [event["value"] for event in peers.events[:4]] == [
                    "INIT", "READY", "CONFIGURING", "OPERATIONAL"]
                sequence = [event["metadata"]["sequence"] for event in peers.events[:4]]
                assert sequence == sorted(set(sequence))
            print("PASS: READY waits; legacy configure blocked; queue survives crash and notify outage")

            wait(lambda: not request(port)[1]["notificationPending"], "drain queue")
            before = event_count("OPERATIONAL")
            peers.decision(2)
            wait(lambda: event_count("OPERATIONAL") == before + 1, "fresh repeated confirmation")
            assert request(port)[1]["state"] == "OPERATIONAL"
            time.sleep(0.2)
            assert event_count("OPERATIONAL") == before + 1
            with peers.lock:
                peers.ready = "pending"
            peers.decision(3)
            time.sleep(0.3)
            assert event_count("OPERATIONAL") == before + 1
            with peers.lock:
                peers.ready = "ready"
            wait(lambda: event_count("OPERATIONAL") == before + 2, "readiness-gated retry")
            print("PASS: repeated confirmation is fresh, idempotent across polls, and readiness gated")

            stop()
            with peers.lock:
                peers.ready = "pending"
            start()
            before = event_count("OPERATIONAL")
            time.sleep(0.3)
            assert event_count("OPERATIONAL") == before
            with peers.lock:
                peers.ready = "ready"
            wait(lambda: event_count("OPERATIONAL") > before, "same-boot crash revalidation")
            print("PASS: daemon restart revalidates before a fresh Operational observation")

            with peers.lock:
                peers.config_available = False
            state("FAULTY")
            with peers.lock:
                peers.config_available = True
            state("OPERATIONAL")
            print("PASS: configd outage and recovery")

            stop()
            boot_file.write_text("boot-2\n")
            start()
            assert request(port, "POST", "/v1/check-in", {"bootResult": "ready"})[0] == 202
            state("OPERATIONAL")
            wait(lambda: event_count("OPERATIONAL", "boot-2") == 1, "new boot completion")
            with peers.lock:
                values = [event["value"] for event in peers.events
                          if event["metadata"]["bootId"] == "boot-2"]
            assert values == ["INIT", "READY", "CONFIGURING", "OPERATIONAL"], values
            print("PASS: new boot consumes the retained decision through the complete flow")
        finally:
            stop()
            for peer in servers.values():
                peer.shutdown()
                peer.server_close()


if __name__ == "__main__":
    run()
