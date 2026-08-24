#!/usr/bin/env python3
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.
#

import argparse
import json
import threading
import time
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer


class LifecycleState:
    def __init__(self, gate_delay):
        self.gate_delay = gate_delay
        self.checkins = 0
        self.boot_result = None
        self.checked_in_at = None
        self.lock = threading.Lock()

    def check_in(self, boot_result):
        with self.lock:
            self.checkins += 1
            self.boot_result = boot_result
            if self.checked_in_at is None:
                self.checked_in_at = time.monotonic()

    def snapshot(self):
        with self.lock:
            gate_open = (
                self.checked_in_at is not None
                and time.monotonic() - self.checked_in_at >= self.gate_delay
            )
            return {
                "checkins": self.checkins,
                "bootResult": self.boot_result,
                "gateOpen": gate_open,
            }


class Handler(BaseHTTPRequestHandler):
    state = None

    def reply(self, status, body):
        encoded = json.dumps(body).encode()
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(encoded)))
        self.end_headers()
        self.wfile.write(encoded)

    def do_POST(self):
        if self.path != "/v1/check-in":
            self.reply(404, {"error": "not found"})
            return

        try:
            length = int(self.headers.get("Content-Length", "0"))
            body = json.loads(self.rfile.read(length) or b"{}")
        except (ValueError, json.JSONDecodeError):
            self.reply(400, {"error": "invalid json"})
            return

        boot_result = body.get("bootResult")
        if boot_result not in ("ready", "degraded"):
            self.reply(400, {"error": "invalid bootResult"})
            return

        self.state.check_in(boot_result)
        self.reply(202, {"status": "accepted"})

    def do_GET(self):
        snapshot = self.state.snapshot()

        if self.path == "/v1/ping":
            self.reply(200, {"status": "ok"})
        elif self.path == "/v1/gate":
            self.reply(
                200 if snapshot["gateOpen"] else 202,
                {"proceed": snapshot["gateOpen"]},
            )
        elif self.path == "/test/status":
            self.reply(200, snapshot)
        else:
            self.reply(404, {"error": "not found"})

    def log_message(self, fmt, *args):
        print("mock_lifecycle:", fmt % args, flush=True)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, required=True)
    parser.add_argument("--gate-delay", type=float, default=0.0)
    args = parser.parse_args()

    Handler.state = LifecycleState(args.gate_delay)
    with ThreadingHTTPServer((args.host, args.port), Handler) as server:
        server.serve_forever()


if __name__ == "__main__":
    main()
