#!/usr/bin/env python3
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

from __future__ import annotations

import argparse
import json
import math
import os
import pty
import signal
import sys
import threading
import time

from http.server import BaseHTTPRequestHandler
from http.server import ThreadingHTTPServer
from typing import Optional


SERVICE_NAME = "controlleremu"

DEF_LISTEN_ADDR = "0.0.0.0"
DEF_HTTP_PORT = 18089
DEF_SERIAL_LINK = "/tmp/victron-tty"
DEF_READY_FILE = "/tmp/victron-emu.ready"
DEF_INTERVAL_SEC = 1.0

PRODUCT_ID = "0xA042"
FIRMWARE = "159"
SERIAL = "HQ2242ABCDE"

g_stop = False
g_stats_lock = threading.Lock()
g_stats = {
    "started": False,
    "serial_path": "",
    "serial_link": "",
    "frame_count": 0,
    "last_frame_ms": 0,
    "battery_voltage_v": 0.0,
    "battery_current_a": 0.0,
    "pv_voltage_v": 0.0,
    "pv_power_w": 0,
    "temperature_c": 0,
    "charge_state": "unknown",
    "error": "0",
}


def read_version() -> str:
    version = os.getenv("CONTROLLEREMU_VERSION")

    if version and version.strip():
        return version.strip()

    version = os.getenv("VERSION")

    if version and version.strip():
        return version.strip()

    return "unknown"


VERSION = read_version()


def now_ms() -> int:
    return int(time.time() * 1000)


def safe_unlink(path: str) -> None:
    try:
        os.unlink(path)
    except FileNotFoundError:
        pass


def ensure_parent(path: str) -> None:
    parent = os.path.dirname(path)

    if parent:
        os.makedirs(parent, exist_ok=True)


def publish_serial(slave_path: str,
                   serial_link: Optional[str],
                   ready_file: Optional[str]) -> None:
    if serial_link:
        ensure_parent(serial_link)

        if os.path.lexists(serial_link):
            safe_unlink(serial_link)

        os.symlink(slave_path, serial_link)

    if ready_file:
        ensure_parent(ready_file)

        with open(ready_file, "w", encoding="utf-8") as file_desc:
            file_desc.write(slave_path + "\n")


def cleanup_serial(serial_link: Optional[str],
                   ready_file: Optional[str]) -> None:
    if serial_link and os.path.islink(serial_link):
        safe_unlink(serial_link)

    if ready_file and os.path.exists(ready_file):
        safe_unlink(ready_file)


def charge_state_label(value: int) -> str:
    labels = {
        0: "off",
        1: "low-power",
        2: "fault",
        3: "bulk",
        4: "absorption",
        5: "float",
        6: "storage",
    }

    return labels.get(value, "unknown")


def build_frame(fields: list[tuple[str, str]]) -> bytes:
    body = b"".join(
        f"{key}\t{value}\r\n".encode("ascii")
        for key, value in fields
    )

    prefix = body + b"Checksum\t"
    checksum = (-sum(prefix)) % 256

    return prefix + bytes([checksum]) + b"\r\n"


def build_sample_frame(elapsed_sec: float) -> tuple[bytes, dict[str, object]]:
    cycle_sec = 180.0
    phase = (elapsed_sec % cycle_sec) / cycle_sec
    sun = math.sin(math.pi * phase)

    if sun < 0.0:
        sun = 0.0

    pv_power_w = int(round(980.0 * sun))
    battery_voltage_v = 49.2 + (7.0 * sun)
    battery_current_a = pv_power_w / battery_voltage_v

    pv_voltage_v = 0.0
    if pv_power_w > 0:
        pv_voltage_v = battery_voltage_v + 8.0 + (8.0 * sun)

    temp_c = int(round(24.0 + (12.0 * sun)))

    charge_state = 0
    mppt = 0

    if pv_power_w > 220:
        charge_state = 5
        mppt = 2
    elif pv_power_w > 30:
        charge_state = 3
        mppt = 2

    yield_today = 5.6 * sun

    frame = build_frame([
        ("PID", PRODUCT_ID),
        ("FW", FIRMWARE),
        ("SER#", SERIAL),
        ("V", str(int(round(battery_voltage_v * 1000.0)))),
        ("I", str(int(round(battery_current_a * 1000.0)))),
        ("VPV", str(int(round(pv_voltage_v * 1000.0)))),
        ("PPV", str(pv_power_w)),
        ("CS", str(charge_state)),
        ("MPPT", str(mppt)),
        ("OR", "0"),
        ("ERR", "0"),
        ("LOAD", "ON"),
        ("IL", "0"),
        ("H19", f"{yield_today:.2f}"),
        ("H20", f"{1248.4 + yield_today:.2f}"),
        ("H21", str(pv_power_w)),
        ("H22", "5.20"),
        ("H23", "910"),
        ("T", str(temp_c)),
    ])

    sample = {
        "battery_voltage_v": battery_voltage_v,
        "battery_current_a": battery_current_a,
        "pv_voltage_v": pv_voltage_v,
        "pv_power_w": pv_power_w,
        "temperature_c": temp_c,
        "charge_state": charge_state_label(charge_state),
        "error": "0",
    }

    return frame, sample


class ApiHandler(BaseHTTPRequestHandler):
    server_version = SERVICE_NAME
    sys_version = ""

    def log_message(self, fmt: str, *args: object) -> None:
        return

    def send_text(self, status: int, body: str) -> None:
        data = body.encode("utf-8")

        self.send_response(status)
        self.send_header("Content-Type", "text/plain")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    def send_json(self, status: int, body: dict[str, object]) -> None:
        data = json.dumps(body, indent=2).encode("utf-8")

        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    def do_GET(self) -> None:
        if self.path == "/v1/ping":
            self.send_text(200, "OK\n")
            return

        if self.path == "/v1/version":
            self.send_text(200, VERSION + "\n")
            return

        if self.path == "/v1/status":
            with g_stats_lock:
                stats = dict(g_stats)

            self.send_json(200, {
                "service": SERVICE_NAME,
                "status": "running" if stats["started"] else "starting",
                "version": VERSION,
                "stats": stats,
            })
            return

        if self.path == "/v1/metrics":
            with g_stats_lock:
                stats = dict(g_stats)

            self.send_json(200, {
                "timestamp_ms": now_ms(),
                "metrics": [
                    {
                        "name": "controlleremu_frames_total",
                        "value": stats["frame_count"],
                        "unit": "count",
                    },
                    {
                        "name": "controlleremu_pv_power_watts",
                        "value": stats["pv_power_w"],
                        "unit": "W",
                    },
                    {
                        "name": "controlleremu_battery_voltage",
                        "value": stats["battery_voltage_v"],
                        "unit": "V",
                    },
                    {
                        "name": "controlleremu_temperature_c",
                        "value": stats["temperature_c"],
                        "unit": "C",
                    },
                ],
            })
            return

        self.send_json(404, {
            "error": "not-found",
            "path": self.path,
        })


def start_http_server(addr: str, port: int) -> ThreadingHTTPServer:
    server = ThreadingHTTPServer((addr, port), ApiHandler)

    thread = threading.Thread(target=server.serve_forever)
    thread.daemon = True
    thread.start()

    return server


def run_emulator(args: argparse.Namespace) -> None:
    global g_stop

    master_fd = -1
    slave_fd = -1

    try:
        master_fd, slave_fd = pty.openpty()
        slave_path = os.ttyname(slave_fd)

        publish_serial(slave_path, args.serial_link, args.ready_file)

        with g_stats_lock:
            g_stats["started"] = True
            g_stats["serial_path"] = slave_path
            g_stats["serial_link"] = args.serial_link or ""

        print(f"{SERVICE_NAME}: serial path {slave_path}", flush=True)
        print(f"{SERVICE_NAME}: serial link {args.serial_link}", flush=True)

        start_sec = time.time()

        while not g_stop:
            elapsed_sec = time.time() - start_sec
            frame, sample = build_sample_frame(elapsed_sec)

            try:
                os.write(master_fd, frame)
            except OSError:
                time.sleep(args.interval)
                continue

            with g_stats_lock:
                g_stats["frame_count"] += 1
                g_stats["last_frame_ms"] = now_ms()
                g_stats.update(sample)

            time.sleep(args.interval)

    finally:
        cleanup_serial(args.serial_link, args.ready_file)

        if master_fd >= 0:
            try:
                os.close(master_fd)
            except OSError:
                pass

        if slave_fd >= 0:
            try:
                os.close(slave_fd)
            except OSError:
                pass


def handle_signal(_signum: int, _frame: object) -> None:
    global g_stop

    g_stop = True


def getenv_int(name: str, default: int) -> int:
    value = os.getenv(name)

    if not value:
        return default

    try:
        return int(value)
    except ValueError:
        return default


def getenv_float(name: str, default: float) -> float:
    value = os.getenv(name)

    if not value:
        return default

    try:
        return float(value)
    except ValueError:
        return default


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="VE.Direct emulator for controller.d",
    )

    parser.add_argument(
        "--listen-addr",
        default=os.getenv("CONTROLLEREMU_LISTEN_ADDR", DEF_LISTEN_ADDR),
        help="HTTP listen address",
    )

    parser.add_argument(
        "--http-port",
        type=int,
        default=getenv_int("CONTROLLEREMU_HTTP_PORT", DEF_HTTP_PORT),
        help="HTTP listen port",
    )

    parser.add_argument(
        "--serial-link",
        default=os.getenv("CONTROLLEREMU_SERIAL_LINK", DEF_SERIAL_LINK),
        help="Stable serial symlink used by controller.d",
    )

    parser.add_argument(
        "--ready-file",
        default=os.getenv("CONTROLLEREMU_READY_FILE", DEF_READY_FILE),
        help="File written when PTY is ready",
    )

    parser.add_argument(
        "--interval",
        type=float,
        default=getenv_float("CONTROLLEREMU_INTERVAL_SEC",
                             DEF_INTERVAL_SEC),
        help="Seconds between VE.Direct frames",
    )

    return parser.parse_args()


def main() -> int:
    signal.signal(signal.SIGINT, handle_signal)
    signal.signal(signal.SIGTERM, handle_signal)

    args = parse_args()
    httpd = start_http_server(args.listen_addr, args.http_port)

    print(
        f"{SERVICE_NAME}: http {args.listen_addr}:{args.http_port}",
        flush=True,
    )
    print(f"{SERVICE_NAME}: version {VERSION}", flush=True)

    try:
        run_emulator(args)
    finally:
        httpd.shutdown()
        httpd.server_close()

    return 0


if __name__ == "__main__":
    sys.exit(main())
