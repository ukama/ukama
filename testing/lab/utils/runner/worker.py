#!/usr/bin/env python3
import json
import os
from pathlib import Path
import shutil
import signal
import subprocess
import sys
import threading
import time
from common import AWS, RunnerError, atomic_json, now, read_json, run, sha256
from config import ROOT
from artifacts import check_compressed_image, extract
import network

RESULTS = ROOT / "results"
STATUS_LOCK = threading.RLock()


def progress(config):
    status = RESULTS / "progress.tsv"
    result = {}
    if status.exists():
        for line in status.read_text().splitlines():
            parts = line.split("\t", 1)
            if len(parts) == 2:
                result[parts[0]] = parts[1]
    for key in ("completed", "passed", "failed", "skipped", "total"):
        result[key] = int(result.get(key, 0))
    return result


def publish_status(aws, config, phase, error="", stage=""):
    status = {"phase": phase, "updated_at": now(), "heartbeat_epoch": time.time(),
              "worker": config["worker"], "progress": progress(config), "error": error, "stage": stage}
    # Foreground readiness and the heartbeat can publish at the same time.
    # Serialize both the atomic replacement and its S3 upload.
    with STATUS_LOCK:
        atomic_json(RESULTS / "status.json", status)
        aws.copy(RESULTS / "status.json", result_uri(config) + "/status.json")


def result_uri(config):
    return f"s3://{config['bucket']}/results/{config['batch']}/{config['worker']}"


def check_stop(aws, config):
    value = aws.optional_json(config["bucket"], f"control/{config['batch']}/stop.json")
    if value:
        atomic_json(ROOT / "stop.json", value)
    return bool(value)


def upload_results(aws, config):
    aws.sync(RESULTS, result_uri(config))


def load_images(image_file, info, min_free_gb):
    """Materialize and verify the tar, then load it into this worker's Podman."""
    reserve = min_free_gb * 1024**3
    # Space for the plain tar plus unpacked image layers, with the run reserve
    # left available. The compressed download is already accounted for by df.
    if shutil.disk_usage(ROOT).free < 2 * info.get("tar_size", 0) + reserve:
        raise RunnerError("insufficient worker disk for image archive, image layers and MIN_FREE_DISK_GB")
    raw = ROOT / "images.tar"
    check_compressed_image(image_file, info, raw)
    print(f"Loading {info['tar_size']:,} verified tar bytes into worker Podman...", flush=True)
    result = run(["podman", "load", "--input", raw], check=False, timeout=7200)
    output = (result.stdout or "") + (result.stderr or "")
    (RESULTS / "image-load.log").write_text(output)
    if output:
        print(output, flush=True)
    if result.returncode:
        raise RunnerError(f"Podman image load exited {result.returncode}: {output[-4000:].strip()}")
    for tag, expected in info["images"].items():
        loaded = json.loads(run(["podman", "image", "inspect", tag]).stdout)[0]
        if loaded["Id"].removeprefix("sha256:") != expected:
            raise RunnerError(f"loaded image does not match controller image: {tag}")
    raw.unlink()
    image_file.unlink()


def bounded_retry(callback, seconds, label):
    deadline = time.monotonic() + seconds
    while True:
        try:
            return callback()
        except Exception as exc:
            if time.monotonic() >= deadline:
                raise RunnerError(f"{label} failed: {exc}") from exc
            print(f"{label}: retrying after {exc}", flush=True)
            time.sleep(10)


def before(config):
    aws = AWS(config["config"]["REGION"])
    if bounded_retry(lambda: check_stop(aws, config), 180, "batch control"):
        return 10
    disk = shutil.disk_usage(ROOT)
    if disk.free < config["config"]["MIN_FREE_DISK_GB"] * 1024**3:
        raise RunnerError("worker disk is below MIN_FREE_DISK_GB; no further scenarios started")
    bounded_retry(lambda: network.check(config["config"], config["environment"]), 180, "worker connectivity")
    return 0


def after(config, rc):
    aws = AWS(config["config"]["REGION"])
    if config["fail_fast"] and rc:
        stop = ROOT / "stop.json"
        atomic_json(stop, {"worker": config["worker"], "reason": "scenario exit code was nonzero", "created_at": now()})
        bounded_retry(lambda: aws.copy(stop, f"s3://{config['bucket']}/control/{config['batch']}/stop.json"), 180, "fail-fast notification")
    if rc:
        network.snapshot(RESULTS / "infrastructure")
    bounded_retry(lambda: upload_results(aws, config), config["config"]["UPLOAD_RETRY_MINUTES"] * 60, "scenario result upload")
    return 0


def worker(config):
    aws = AWS(config["config"]["REGION"])
    phase = {"value": "PREPARING", "error": "", "stage": "source-download"}
    def stage(value):
        phase["stage"] = value
        print(f"Worker stage: {value}", flush=True)
    RESULTS.mkdir(parents=True, exist_ok=True)
    stop = threading.Event()
    def heartbeat():
        while not stop.is_set():
            try:
                with STATUS_LOCK:
                    publish_status(aws, config, phase["value"], phase["error"], phase["stage"])
                for name in ("bootstrap.log", "worker.log"):
                    source = RESULTS / name
                    if source.exists():
                        # Small live tail; complete originals are uploaded at scenario boundaries.
                        with source.open("rb") as stream:
                            stream.seek(max(0, source.stat().st_size - 65536))
                            tail = ROOT / (name + ".tail")
                            tail.write_bytes(stream.read())
                        aws.copy(tail, result_uri(config) + "/live/" + name)
            except Exception as exc:
                print(f"heartbeat publication: {exc}", flush=True)
            stop.wait(30)
    publisher = threading.Thread(target=heartbeat, daemon=True)
    publisher.start()
    lab_rc = None
    error = ""
    try:
        prefix = f"s3://{config['bucket']}/inputs/{config['batch']}"
        archive = ROOT / "source.tar.gz"
        aws.copy(prefix + "/source.tar.gz", archive)
        if sha256(archive) != config["source_sha256"]:
            raise RunnerError("source download checksum mismatch")
        extract(archive, ROOT)
        archive.unlink()
        stage("vpn-and-dns")
        if config["config"]["NETWORK_MODE"] == "vpn":
            aws.copy(prefix + "/" + config["worker"] + ".ovpn", ROOT / "client.ovpn")
            (ROOT / "client.ovpn").chmod(0o600)
        # Set the Podman backend and subnet pools before its first invocation.
        network.configure(config["config"])
        stage("image-download")
        image_file = ROOT / "images.tar.gz"
        aws.copy(f"s3://{config['bucket']}/{config['image_key']}", image_file)
        stage("image-load")
        load_images(image_file, config["image_info"], config["config"]["MIN_FREE_DISK_GB"])
        stage("runtime-check")
        env = {**os.environ, **config["environment"], "HOME": "/root",
               "UKAMA_REPO": str(ROOT / "repo"), "LAB_BIN": str(ROOT / "runtime/lab/run"),
               "P0_RUNS_DIR": str(RESULTS / "batch"), "P0_STATUS_FILE": str(RESULTS / "progress.tsv"),
               "SCENARIO_ROOT": str(ROOT / "lab" / config["scenario_root"]),
               "ULAB_RUNNER_HOOK": str(ROOT / "hook.sh")}
        # AWS mode does not use optional kubectl mesh-pod cleanup.
        # Also remove controller AWS credentials and development library paths.
        for key in ("LD_LIBRARY_PATH", "LD_PRELOAD", "AWS_PROFILE", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN", "ULAB_KUBECTL", "KUBECONFIG"):
            env.pop(key, None)
        run([ROOT / "lab/scripts/node-host-control.sh", "check", "ulab-probe"], env=env)
        # Exercise this binary/loader without contacting any backend or running a scenario.
        run([env["LAB_BIN"], "list-checks"], env=env, cwd=ROOT / "lab")
        stage("connectivity")
        bounded_retry(lambda: network.check(config["config"], env), 180, "startup connectivity")
        (ROOT / "hook.sh").write_text(f'#!/bin/sh\nexec /usr/bin/python3 {ROOT}/bootstrap/worker.py hook "$@"\n')
        (ROOT / "hook.sh").chmod(0o700)
        selected = ROOT / "scenarios.txt"
        selected.write_text("\n".join(config["scenarios"]) + "\n")
        command = [ROOT / "lab/utils/run-scenarios.sh", config["suite"], "--mode", "local",
                   "--scenario-list", selected, "--batch-id", config["batch_id"], "--factory-nodes", "0"]
        if config["fail_fast"]:
            command.append("--fail-fast")
        # A durable infrastructure-ready marker lets the controller admit the
        # rest of the workers even if this worker has already started scenarios.
        # This is independent of any scenario's verdict or backend business data.
        with STATUS_LOCK:
            phase["value"] = "READY"
            publish_status(aws, config, "READY", stage="ready")
        ready = RESULTS / "ready.json"
        atomic_json(ready, {"worker": config["worker"], "batch": config["batch"], "ready_at": now(),
                            "image_sha256": config["image_info"]["sha256"]})
        bounded_retry(lambda: aws.copy(ready, result_uri(config) + "/ready.json"), 180, "infrastructure readiness upload")
        phase["value"] = "RUNNING"
        stage("scenario-execution")
        with (RESULTS / "worker.log").open("a") as log:
            process = subprocess.Popen([str(x) for x in command], cwd=ROOT / "lab", env=env, stdout=log, stderr=subprocess.STDOUT)
            lab_rc = process.wait()
        report = RESULTS / "batch" / config["batch_id"] / "batch-report.json"
        if lab_rc == 70 or not report.is_file():
            raise RunnerError("local execution could not finish normally; see worker.log and uploaded original scenario reports")
        # A failed scenario gives local exit 1 but is a completed worker, not an infrastructure error.
        count = len(read_json(report).get("results", []))
        if count == len(config["scenarios"]):
            phase["value"] = "DONE"
        elif config["fail_fast"] and ((ROOT / "stop.json").exists() or lab_rc != 0):
            phase["value"] = "STOPPED"
        else:
            raise RunnerError("worker ended with unfinished assignments")
    except Exception as exc:
        error = f"{phase['stage']}: {exc}"
        phase.update(value="INFRA_ERROR", error=error)
        print(f"Worker infrastructure error: {error}", flush=True)
    finally:
        network.snapshot(RESULTS / "infrastructure")
        stop.set()
        publisher.join(timeout=180)
    # Upload complete.json LAST, and only after all original results are durable.
    bounded_retry(lambda: upload_results(aws, config), config["config"]["UPLOAD_RETRY_MINUTES"] * 60, "final result upload")
    publish_status(aws, config, phase["value"], error, phase["stage"])
    marker = ROOT / "complete.json"
    atomic_json(marker, {"phase": phase["value"], "error": error, "lab_exit": lab_rc,
                         "worker": config["worker"], "finished_at": now(), "progress": progress(config), "stage": phase["stage"]})
    aws.copy(marker, result_uri(config) + "/complete.json")
    run(["shutdown", "-h", "now"], check=False)
    return 0


if __name__ == "__main__":
    try:
        config = read_json(ROOT / "worker.json")
        if len(sys.argv) > 1 and sys.argv[1] == "hook":
            action = sys.argv[2]
            sys.exit(before(config) if action == "before" else after(config, int(sys.argv[3])))
        sys.exit(worker(config))
    except Exception as exc:
        print(f"worker: {exc}", file=sys.stderr)
        sys.exit(70)
