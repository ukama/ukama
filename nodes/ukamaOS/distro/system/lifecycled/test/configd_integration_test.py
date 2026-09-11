#!/usr/bin/env python3
"""Exercise lifecycle drop 2 with a real config.d drop 1 executable.

CONFIGD_BIN uses the drop-1 test build's CONFIGD_TEST_PORT adapter, or a
production binary listening on CONFIGD_PORT. LIFECYCLED_BIN is required.
"""
import json
import os
from pathlib import Path
import subprocess
import tempfile

from component_test import Peers, request, server, unused_port, wait

peers = Peers()
peers.notify_available = True
starter = server("starter", peers)
notify = server("notify", peers)
config_port = int(os.environ.get("CONFIGD_PORT", unused_port()))
lifecycle_port = unused_port()
processes = []

with tempfile.TemporaryDirectory(prefix="lifecycle-configd-") as directory:
    env = dict(os.environ, ENV_CONFIG_DEBUG_MODE="1",
               CONFIGD_TEST_PORT=str(config_port),
               LIFECYCLED_HTTP_PORT=str(lifecycle_port),
               LIFECYCLED_STARTER_PORT=str(starter.server_port),
               LIFECYCLED_CONFIG_PORT=str(config_port),
               LIFECYCLED_NOTIFY_PORT=str(notify.server_port),
               LIFECYCLED_STATE_FILE=str(Path(directory) / "lifecycle.json"),
               LIFECYCLED_CHECKIN_TIMEOUT_SEC="1", LIFECYCLED_POLL_INTERVAL_MS="50")
    try:
        processes.append(subprocess.Popen(
            [os.environ["CONFIGD_BIN"], "--state-file", str(Path(directory) / "config.status")],
            env=env, stdout=subprocess.DEVNULL))
        wait(lambda: request(config_port, path="/v1/config/status")[0] == 200, "configd startup")
        processes.append(subprocess.Popen([os.environ["LIFECYCLED_BIN"]],
                                         env=env, stdout=subprocess.DEVNULL))
        wait(lambda: request(lifecycle_port)[0] == 200, "lifecycle startup")
        assert request(lifecycle_port, "POST", "/v1/check-in", {"bootResult": "ready"})[0] == 202
        wait(lambda: request(lifecycle_port)[1]["state"] == "READY", "READY")
        command = dict(mode="NOCONFIG", requestId="assignment-real")
        assert request(config_port, "POST", "/v1/config", command)[0] == 200
        wait(lambda: request(lifecycle_port)[1]["state"] == "OPERATIONAL", "OPERATIONAL")
        first_sequence = request(lifecycle_port)[1]["sequence"]
        assert request(config_port, "POST", "/v1/config", command)[0] == 200
        wait(lambda: request(lifecycle_port)[1]["sequence"] > first_sequence, "confirmation retry")
        assert request(lifecycle_port)[1]["state"] == "OPERATIONAL"
        wait(lambda: not request(lifecycle_port)[1]["notificationPending"], "event delivery")
        with peers.lock:
            values = [item["value"] for item in peers.events]
            assert values == ["INIT", "READY", "CONFIGURING", "OPERATIONAL", "OPERATIONAL"], values
            last = json.loads(peers.events[-1]["details"])
            assert last["requestId"] == "assignment-real" and last["configGeneration"] == 2
        print("PASS: real configd NOCONFIG, lifecycle transitions, and repeated correlated completion")
    finally:
        for process in reversed(processes):
            process.kill()
            process.wait(timeout=5)
        for peer in (starter, notify):
            peer.shutdown()
            peer.server_close()
