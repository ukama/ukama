#!/usr/bin/env python3
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2026-present, Ukama Inc.
"""Collect existing Prometheus history for a scenario batch (stdlib only)."""
import argparse
import hashlib
import json
import math
import os
from pathlib import Path
import re
import shlex
import subprocess
import sys
import time
import urllib.error
import urllib.parse
import urllib.request

import cluster_resources
from cluster_report import build_report, write_report, compare

DEFAULT_URL = "https://prometheus.udev.ukama.com"
PREFIX = "ULAB_RUNNER_METRICS_"


def save(path, value):
    path = Path(path)
    temporary = path.with_suffix(path.suffix + ".tmp")
    temporary.write_text(json.dumps(value, indent=2, allow_nan=False) + "\n")
    temporary.chmod(0o600)
    temporary.replace(path)


def read(path):
    return json.loads(Path(path).read_text())


def number(name, default, minimum, maximum):
    value = int(os.environ.get(PREFIX + name, default))
    if not minimum <= value <= maximum:
        raise ValueError(f"{PREFIX}{name} must be {minimum}..{maximum}")
    return value


def configuration():
    url = os.environ.get("PROMETHEUS_URL", DEFAULT_URL).rstrip("/")
    parsed = urllib.parse.urlsplit(url)
    if (parsed.scheme not in ("http", "https") or not parsed.hostname
            or parsed.username or parsed.password or parsed.query
            or parsed.fragment):
        raise ValueError("PROMETHEUS_URL must be an HTTP(S) base URL")
    service = os.environ.get(PREFIX + "SERVICE", ".*")
    try:
        re.compile(service)
    except re.error as exc:
        raise ValueError(f"invalid service regex: {exc}") from exc
    return {
        "url": url.removesuffix("/api/v1"),
        "namespace": os.environ.get(PREFIX + "NAMESPACE", "ukama-.*"),
        "service": service,
        "before": number("BEFORE", 600, 0, 86400),
        "after": number("AFTER", 600, 0, 86400),
        "step": number("STEP", 30, 1, 3600),
        "timeout": number("TIMEOUT", 15, 1, 120),
    }


def begin(batch, metadata=None):
    directory = Path(batch) / "metrics"
    directory.mkdir(mode=0o700, exist_ok=True)
    path = directory / "config.json"
    if not path.exists():
        save(path, {"config": configuration(), "created_at": time.time(),
                    "run": metadata or {}})


def safe_begin(batch, metadata=None):
    try:
        begin(batch, metadata)
    except Exception as exc:
        # Optional observation must never change the batch verdict.
        print(f"Metrics warning: {exc}", file=sys.stderr, flush=True)


def resolve_report(batch, row):
    original = Path(row.get("report") or "missing-report.json")
    if original.is_file():
        return original
    local = batch / original
    if local.is_file():
        return local
    # Reports may have been moved together after collection from AWS.
    for marker in ("workers", "runs"):
        if marker in original.parts:
            index = original.parts.index(marker)
            local = batch.joinpath(*original.parts[index:])
            if local.is_file():
                return local
    return batch / "runs" / row.get("run_id", "missing") / "report.json"


def scenario_results(batch):
    path = batch / "batch-report.json"
    summary = read(path) if path.exists() else {}
    rows = []
    for entry in summary.get("results", []):
        row = dict(entry)
        row.setdefault("worker", "local")
        try:
            report = read(resolve_report(batch, row))
            row.update({key: report.get(key) for key in
                        ("started_at", "ended_at", "cleanup")})
            if not row.get("failure_reason"):
                failures = [r.get("detail", r.get("name", ""))
                            for r in report.get("results", [])
                            if r.get("state") == "FAIL"]
                row["failure_reason"] = "; ".join(failures)
            start, end = row.get("started_at"), row.get("ended_at")
            if (not isinstance(start, (int, float))
                    or not isinstance(end, (int, float))
                    or not math.isfinite(start) or not math.isfinite(end)
                    or start <= 0 or end < start):
                raise ValueError("missing or invalid scenario timestamps")
        except (OSError, ValueError, TypeError) as exc:
            row["started_at"] = row["ended_at"] = None
            row["timing_error"] = str(exc)
        rows.append(row)
    return summary, rows


def queries(config, duration):
    ns = 'namespace=~' + json.dumps(config["namespace"])
    pod = ns + ',pod!=""'
    container = pod + ',container!="",container!="POD"'
    deployment = lambda metric: (
        f"max by (namespace, deployment) ({metric}{{{ns}}})")
    available = ("kube_deployment_status_condition{" + ns
                 + ',condition="Available",status="true"}')
    rate = lambda metric, labels=pod: (
        f"sum by (namespace, pod) (rate({metric}{{{labels}}}[5m]))")
    network = lambda metric: (
        "sum by (namespace, pod) (max by (namespace, pod, interface) "
        f"(rate({metric}{{{pod}}}[5m])))")
    values = {
        "available": ("range", "max by (namespace, deployment) ("
                      + available + ")"),
        "uptime": ("instant", "max by (namespace, deployment) "
                   f"(avg_over_time({available}[{duration}s])) * 100"),
        "availability_hour": ("range", "max by (namespace, deployment) "
                              f"(avg_over_time({available}[1h])) * 100"),
        "replicas": ("range", deployment(
            "kube_deployment_status_replicas_available")),
        "desired": ("range", deployment("kube_deployment_spec_replicas")),
        "pod_age": ("range", f"time() - kube_pod_start_time{{{pod}}}"),
        "grpc_rps": ("range", "sum by (namespace, pod, grpc_code) "
                     f"(rate(grpc_server_handled_total{{{pod}}}[5m]))"),
        "db_rps": ("range", rate("ukama_db_requests_total")),
        "cpu_mcpu": ("range", "sum by (namespace, pod) (max by "
                     "(namespace, pod, container) (rate("
                     "container_cpu_usage_seconds_total{"
                     + container + "}[5m]))) * 1000"),
        "memory_mib": ("range", "sum by (namespace, pod) (max by "
                       "(namespace, pod, container) ("
                       "container_memory_working_set_bytes{"
                       + container + "})) / 1048576"),
        "rx_bps": ("range", network("container_network_receive_bytes_total")),
        "tx_bps": ("range", network("container_network_transmit_bytes_total")),
        "restarts_hour": ("range", "sum by (namespace, pod) (max by "
                          "(namespace, pod, container, uid) (increase("
                          "kube_pod_container_status_restarts_total{"
                          + pod + "}[1h])))"),
        "termination_reason": ("range", "max by (namespace, pod, "
                               "container, reason, uid) ("
                               "kube_pod_container_status_last_terminated_"
                               "reason{"
                               + pod + "} == 1)"),
        # Raw counters avoid adding overlapping one-hour restart windows.
        "restart_counter": ("range", "max by (namespace, pod, container, "
                            "uid) (kube_pod_container_status_restarts_total{"
                            + pod + "})"),
        "pod_owner": ("range", "max by (namespace, pod, owner_kind, "
                      "owner_name) (kube_pod_owner{" + ns
                      + ',owner_is_controller="true"})'),
        "replicaset_owner": ("range", "max by (namespace, replicaset, "
                             "owner_kind, owner_name) (kube_replicaset_owner{"
                             + ns + ',owner_is_controller="true"})'),
    }
    values.update(cluster_resources.queries(ns))
    return values


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, *args, **kwargs):
        return None


class Prometheus:
    def __init__(self, config):
        self.config = config
        self.opener = urllib.request.build_opener(NoRedirect())

    def request(self, endpoint, parameters):
        url = self.config["url"] + "/api/v1/" + endpoint
        query = urllib.parse.urlencode(parameters)
        headers = {"Accept": "application/json"}
        token = os.environ.get("PROMETHEUS_TOKEN")
        if token:
            headers["Authorization"] = "Bearer " + token
        request = urllib.request.Request(url + "?" + query, headers=headers)
        try:
            with self.opener.open(request,
                                  timeout=self.config["timeout"]) as response:
                result = json.load(response)
        except urllib.error.HTTPError as exc:
            body = exc.read(4096).decode("utf-8", errors="replace")
            try:
                payload = json.loads(body)
                detail = payload.get("error", body)
            except (ValueError, AttributeError):
                detail = re.sub(r"<[^>]*>", " ", body)
            detail = " ".join(str(detail).split())
            if token:
                detail = detail.replace(token, "[redacted]")
            raise ValueError(f"HTTP {exc.code} {exc.reason}: "
                             + detail[:500]) from None
        if result.get("status") != "success":
            raise ValueError(result.get("error", "Prometheus query failed"))
        data = result.get("data", {})
        expected = "matrix" if endpoint == "query_range" else "vector"
        if (data.get("resultType") != expected
                or not isinstance(data.get("result"), list)):
            raise ValueError("unexpected Prometheus response shape")
        return result

    def fetch(self, kind, query, start, end, step, at):
        if kind == "instant":
            return self.request("query", {"query": query, "time": at})
        merged, warnings, infos = {}, [], []
        cursor = start
        while cursor <= end:
            last = min(end, cursor + 9999 * step)
            result = self.request("query_range", {
                "query": query, "start": cursor, "end": last, "step": step})
            if result.get("data", {}).get("resultType") != "matrix":
                raise ValueError("expected a Prometheus matrix")
            warnings.extend(result.get("warnings", []))
            infos.extend(result.get("infos", []))
            for series in result["data"]["result"]:
                key = json.dumps(series["metric"], sort_keys=True)
                target = merged.setdefault(key, {
                    "metric": series["metric"], "values": []})
                target["values"].extend(series.get("values", []))
            cursor += 10000 * step
        return {"status": "success", "warnings": warnings, "infos": infos,
                "data": {"resultType": "matrix",
                         "result": list(merged.values())}}


def collect(batch, wait=False, partial=False):
    batch = Path(batch).resolve()
    begin(batch)
    directory = batch / "metrics"
    saved = read(directory / "config.json")
    config = saved["config"]
    summary, scenarios = scenario_results(batch)
    timed = [s for s in scenarios if s.get("started_at") is not None
             and s.get("outcome") != "SKIP"]
    if not timed:
        raise ValueError("no executed scenario has usable start/end times")
    start = min(s["started_at"] for s in timed)
    end = max(s["ended_at"] for s in timed)
    partial = bool(partial or summary.get("unfinished")
                   or summary.get("infrastructure_errors")
                   or any(s.get("timing_error") for s in scenarios
                          if s.get("outcome") != "SKIP"))
    target = end + config["after"]
    if wait and not partial:
        remaining = target - time.time()
        if remaining > 0:
            print(f"Metrics: observing recovery for {math.ceil(remaining)}s; "
                  "Ctrl-C collects the available window.", flush=True)
        try:
            while time.time() < target:
                time.sleep(max(0, min(30, target - time.time())))
        except KeyboardInterrupt:
            print("\nMetrics: collecting the available recovery window.",
                  flush=True)
    observed_end = min(target, time.time())
    if end > observed_end:
        raise ValueError("scenario times are ahead of the collector clock")
    query_start = math.floor((start - config["before"]) / config["step"])
    query_start *= config["step"]
    query_end = max(query_start, observed_end)
    duration = max(1, math.ceil(end - start))
    client = Prometheus(config)
    raw = {"generated_at": time.time(), "config": config, "start": start,
           "end": end, "observed_end": observed_end,
           "requested_end": target, "partial_batch": partial,
           "run": saved.get("run", {}), "queries": {}}
    state_path = batch / "state.json"
    state = read(state_path) if state_path.exists() else {}
    raw["run"].update({key: state[key] for key in
                       ("suite", "ami", "source_sha256") if key in state})
    if state.get("workers"):
        raw["run"]["workers_launched"] = sum(
            bool(w.get("instance_id")) for w in state["workers"].values())
    planned = queries(config, duration)
    raw["queries"] = {
        name: {"query": query, "kind": kind, "error": "not collected"}
        for name, (kind, query) in planned.items()}
    for name, (kind, query) in planned.items():
        item = {"query": query, "kind": kind}
        try:
            item["response"] = client.fetch(
                kind, query, query_start, query_end, config["step"], end)
        except (OSError, ValueError, urllib.error.URLError) as exc:
            item["error"] = str(exc)
        except KeyboardInterrupt:
            item["error"] = "collection interrupted"
            raw["queries"][name] = item
            break
        raw["queries"][name] = item
        save(directory / "raw.json", raw)
    save(directory / "raw.json", raw)
    report = build_report(batch.name, raw, summary, scenarios)
    save(directory / "report.json", report)
    write_report(directory, report)
    (directory / "error.json").unlink(missing_ok=True)
    unavailable = sum(q["status"] != "OK" for q in report["queries"].values())
    print(f"Metrics: {len(report['services'])} observed services; "
          f"{unavailable}/{len(planned)} queries unavailable or empty",
          flush=True)
    for name, status in report["queries"].items():
        if status["status"] != "OK":
            detail = status["error"] or "no usable samples in this window"
            print(f"Metrics warning: {name}: {detail}",
                  file=sys.stderr, flush=True)
    if not report["services"]:
        print("Metrics warning: no service KPIs collected; see Query coverage "
              "in the report. Scenario results are unchanged.",
              file=sys.stderr, flush=True)
    print(f"Metrics: {directory / 'report.txt'}", flush=True)
    print(f"Metrics: {directory / 'report.html'}", flush=True)
    return report


def open_dashboard(batch):
    report = (Path(batch) / "metrics/report.html").resolve()
    print(f"Dashboard: {report}", flush=True)
    try:
        if not report.is_file():
            raise OSError("report.html was not generated")
        command = shlex.split(
            os.environ.get("ULAB_DASHBOARD_COMMAND") or "xdg-open")
        if not command:
            raise ValueError("ULAB_DASHBOARD_COMMAND must name a program")
        log = report.with_name("dashboard-open.log")
        with log.open("w") as output:
            process = subprocess.Popen(
                [*command, report.as_uri()], stdin=subprocess.DEVNULL,
                stdout=output, stderr=subprocess.STDOUT,
                start_new_session=True)
        try:
            status = process.wait(timeout=2)
        except subprocess.TimeoutExpired:
            # A browser may stay running until its window is closed.
            return
        if status:
            raise OSError(f"browser command exited {status}; see {log}")
    except (OSError, ValueError) as exc:
        print(f"Dashboard warning: {exc}; open {report} manually.",
              file=sys.stderr, flush=True)


def safe_collect(batch, wait=True, partial=False, dashboard=False):
    if not (Path(batch) / "metrics/config.json").exists():
        return
    try:
        collect(batch, wait, partial)
        if dashboard:
            open_dashboard(batch)
    except Exception as exc:
        # Optional observation must never change the batch verdict.
        print(f"Metrics warning: {exc}", file=sys.stderr, flush=True)
        try:
            save(Path(batch) / "metrics/error.json",
                 {"error": str(exc), "generated_at": time.time()})
        except OSError:
            pass
    except KeyboardInterrupt:
        print("\nMetrics collection interrupted; scenario results unchanged.",
              file=sys.stderr, flush=True)


def manifest(paths):
    values = []
    for name in paths:
        path = Path(name)
        values.append({"scenario": str(path),
                       "sha256": hashlib.sha256(path.read_bytes()).hexdigest()})
    return values


def main():
    parser = argparse.ArgumentParser(
        description=__doc__,
        epilog="PROMETHEUS_URL sets the server; PROMETHEUS_TOKEN is optional. "
        + PREFIX + "NAMESPACE/SERVICE set filters (ukama-.* / .*). "
        + PREFIX + "BEFORE/AFTER/STEP/TIMEOUT are seconds (600/600/30/15). "
        "Default rate window is 5m. Collection failures never change verdicts.")
    sub = parser.add_subparsers(dest="action", required=True)
    start = sub.add_parser("start", help="save optional batch configuration")
    start.add_argument("batch", type=Path)
    start.add_argument("--suite", default="")
    start.add_argument("--mode", default="local")
    start.add_argument("--binary", type=Path)
    start.add_argument("--scenarios", nargs="*", default=[])
    fetch = sub.add_parser("collect", help="collect or refresh batch history")
    fetch.add_argument("batch", type=Path)
    fetch.add_argument("--wait", action="store_true",
                       help="wait for the configured recovery window")
    fetch.add_argument("--partial", action="store_true",
                       help="batch was interrupted or is still running")
    fetch.add_argument("--dashboard", action="store_true",
                       help="open the dashboard after collection")
    render = sub.add_parser("render", help="rebuild HTML/text from saved JSON")
    render.add_argument("batch", type=Path)
    diff = sub.add_parser("compare", help="compare saved metrics reports")
    diff.add_argument("batches", nargs="+", type=Path)
    diff.add_argument("--output", type=Path)
    args = parser.parse_args()
    try:
        if args.action == "start":
            metadata = {"suite": args.suite, "mode": args.mode,
                        "workers_requested": 1,
                        "selected": manifest(args.scenarios)}
            if args.binary and args.binary.is_file():
                metadata["binary_sha256"] = hashlib.sha256(
                    args.binary.read_bytes()).hexdigest()
            begin(args.batch, metadata)
        elif args.action == "collect":
            collect(args.batch, args.wait, args.partial)
            if args.dashboard:
                open_dashboard(args.batch)
        elif args.action == "render":
            from cluster_dashboard import write_outputs
            directory = args.batch / "metrics"
            write_outputs(directory, read(directory / "report.json"))
            print(f"Metrics: {directory / 'report.html'}")
            print(f"Metrics: {directory / 'report.txt'}")
        else:
            result = compare(args.batches)
            if args.output:
                args.output.write_text(result)
            else:
                print(result, end="")
    except (OSError, ValueError, KeyError, TypeError) as exc:
        print(f"Metrics warning: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
