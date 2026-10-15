# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2026-present, Ukama Inc.
"""Summarize observed metrics without inventing missing samples or causes."""
import csv
from datetime import datetime, timezone
import html
import json
import math
from pathlib import Path
import re
from collections import defaultdict

METRICS = {
    "available": "Deployment available (0/1)",
    "uptime": "Run availability (%)",
    "availability_hour": "Rolling 1h availability (%)",
    "replicas": "Available replicas", "desired": "Desired replicas",
    "pod_age": "Youngest observed pod age (s)",
    "grpc_rps": "Handled gRPC calls/s", "db_rps": "DB requests/s",
    "cpu_mcpu": "CPU (mCPU)", "memory_mib": "Memory (MiB)",
    "rx_bps": "Network receive (bytes/s)",
    "tx_bps": "Network transmit (bytes/s)",
    "restarts_hour": "Rolling 1h restart increase",
}
NOTES = [
    "All times are UTC. Scenario windows include setup and cleanup.",
    "Scenario end times are not exact failure times. Overlap is correlation, "
    "not proof that a scenario caused a service change.",
    "CPU, request and network rates use a 5m lookback. Peaks are sampled "
    "peaks of those averages. Query step is not the scrape interval.",
    "N/A means no usable samples, not zero. Aggregates include only observed "
    "pods; partial scrape coverage can understate totals.",
    "Availability percentages use observed samples. Missing scrapes are not "
    "counted as available or unavailable. Available is not API health.",
    "Restart counts are observed counter increases, a lower bound. Restarts "
    "before the first sample or in deleted, unsampled pods may be missed.",
    "Termination reasons are last-reason samples and can predate this run. "
    "Rolling 1h restart values are never summed over time.",
    "Pod age is not container uptime. CPU limits/throttling, latency, database "
    "connections, queue depth and worker/VPN health are not measured here.",
    "Background traffic is included. Service grouping uses Kubernetes owner "
    "metrics; unresolved owners remain separate pods, without name guessing.",
]


def utc(value):
    if value is None:
        return "N/A"
    return datetime.fromtimestamp(value, timezone.utc).isoformat(
        timespec="seconds").replace("+00:00", "Z")


def points(series):
    values = series.get("values", [])
    if "value" in series:
        values = [series["value"]]
    result = []
    for timestamp, value in values:
        try:
            timestamp, value = float(timestamp), float(value)
            if math.isfinite(timestamp) and math.isfinite(value):
                result.append((timestamp, value))
        except (ValueError, TypeError):
            continue
    return sorted(result)


def phase(timestamp, raw):
    if timestamp < raw["start"]:
        return "before"
    if timestamp <= raw["end"]:
        return "during"
    return "after"


def statistics(values):
    numbers = [v for _, v in values]
    if not numbers:
        return {"mean": None, "peak": None, "last": None, "samples": 0}
    return {"mean": sum(numbers) / len(numbers), "peak": max(numbers),
            "last": numbers[-1], "samples": len(numbers)}


def owner_index(series, name):
    owners = defaultdict(set)
    for item in series:
        labels = item["metric"]
        if any(value == 1 for _, value in points(item)):
            owners[(labels.get("namespace"), labels.get(name))].add(
                (labels.get("owner_kind"), labels.get("owner_name")))
    return {key: next(iter(value)) for key, value in owners.items()
            if len(value) == 1}


def service_key(labels, pod_owners, rs_owners):
    namespace = labels.get("namespace", "unknown")
    if labels.get("deployment"):
        return (namespace, "Deployment", labels["deployment"])
    pod = labels.get("pod", "unknown")
    kind, name = pod_owners.get((namespace, pod), ("Pod", pod))
    if kind == "ReplicaSet":
        kind, name = rs_owners.get((namespace, name), (kind, name))
    return namespace, kind or "Pod", name or pod


def overlap(scenarios, timestamp):
    return [s["run_id"] for s in scenarios
            if s.get("started_at") is not None
            and s.get("outcome") != "SKIP"
            and s["started_at"] <= timestamp < s["ended_at"]]


def concurrency(scenarios):
    changes = defaultdict(int)
    for row in scenarios:
        start, end = row.get("started_at"), row.get("ended_at")
        if start is not None and end > start and row.get("outcome") != "SKIP":
            changes[start] += 1
            changes[end] -= 1
    current, peak, timeline = 0, 0, []
    for timestamp, delta in sorted(changes.items()):
        current += delta
        peak = max(peak, current)
        timeline.append([timestamp, current])
    return peak, timeline


def build_report(batch, raw, batch_summary, scenarios):
    responses, query_status = {}, {}
    for name, item in raw["queries"].items():
        response = item.get("response", {})
        series = response.get("data", {}).get("result", [])
        responses[name] = series
        samples = sum(len(points(s)) for s in series)
        status = "ERROR" if item.get("error") else (
            "OK" if samples else "NO_DATA")
        query_status[name] = {
            "status": status, "series": len(series), "samples": samples,
            "error": item.get("error"),
            "warnings": response.get("warnings", []),
            "infos": response.get("infos", []),
        }
    pod_owners = owner_index(responses.get("pod_owner", []), "pod")
    rs_owners = owner_index(responses.get("replicaset_owner", []), "replicaset")
    selected = re.compile(raw["config"]["service"])
    services, events = {}, []
    step = raw["config"]["step"]

    def get_service(labels):
        key = service_key(labels, pod_owners, rs_owners)
        if not selected.fullmatch(key[2]):
            return None
        identity = "/".join(key)
        if identity not in services:
            services[identity] = {
                "service": identity, "pods": set(), "series": {},
                "restart_increases": None, "termination_reasons": set(),
            }
        service = services[identity]
        if labels.get("pod"):
            service["pods"].add(labels["pod"])
        return service

    for name, series in responses.items():
        if name == "replicaset_owner":
            continue
        for item in series:
            service = get_service(item["metric"])
            if service is None:
                continue
            values = [(t, v) for t, v in points(item)
                      if raw["start"] - raw["config"]["before"] <= t
                      <= raw["observed_end"]]
            if name == "restart_counter":
                previous = None
                for timestamp, value in values:
                    if (phase(timestamp, raw) == "during" and previous
                            and previous[0] >= raw["start"]):
                        # A reset starts a new counter; do not subtract it.
                        delta = value - previous[1]
                        delta = delta if delta >= 0 else value
                        old = service["restart_increases"] or 0
                        service["restart_increases"] = old + delta
                        if delta:
                            events.append({"time": timestamp,
                                           "service": service["service"],
                                           "event": "restart increase "
                                           f"{delta:g}",
                                           "pod": item["metric"].get("pod")})
                    previous = (timestamp, value)
            elif name == "termination_reason":
                if any(phase(t, raw) == "during" and v == 1
                       for t, v in values):
                    reason = item["metric"].get("reason", "unknown")
                    service["termination_reasons"].add(reason)
            elif name in METRICS:
                target = service["series"].setdefault(name, {})
                for timestamp, value in values:
                    if name == "pod_age":
                        target[timestamp] = min(
                            target.get(timestamp, value), value)
                    else:
                        target[timestamp] = target.get(timestamp, 0) + value

    for service in services.values():
        service["pods"] = sorted(service["pods"])
        service["termination_reasons"] = sorted(
            service["termination_reasons"])
        service["statistics"] = {}
        for metric, values in service["series"].items():
            ordered = sorted(values.items())
            service["series"][metric] = ordered
            service["statistics"][metric] = {
                name: statistics([(t, v) for t, v in ordered
                                  if phase(t, raw) == name])
                for name in ("before", "during", "after")}
        desired = dict(service["series"].get("desired", []))
        degraded, paired, previous, minimum = 0, 0, None, None
        for timestamp, available in service["series"].get("replicas", []):
            if phase(timestamp, raw) != "during" or timestamp not in desired:
                continue
            deficit = available < desired[timestamp]
            paired += 1
            if desired[timestamp] > 0:
                ratio = available / desired[timestamp]
                if minimum is None or ratio < minimum[0]:
                    minimum = (ratio, available, desired[timestamp])
            if deficit:
                degraded += min(step, max(0, raw["end"] - timestamp))
            if deficit != previous and (deficit or previous is not None):
                events.append({"time": timestamp,
                               "service": service["service"],
                               "event": "replicas "
                               f"{available:g}/{desired[timestamp]:g}"})
            previous = deficit
        service["degraded_sample_seconds"] = degraded if paired else None
        service["lowest_replicas"] = list(minimum[1:]) if minimum else None
    for event in events:
        event["overlapping_runs"] = overlap(scenarios, event["time"])
    peak, activity = concurrency(scenarios)
    outcomes = {key: sum(s.get("outcome") == key for s in scenarios)
                for key in ("PASS", "FAIL", "SKIP")}
    executed = outcomes["PASS"] + outcomes["FAIL"]
    windows = {"before": [raw["start"] - raw["config"]["before"],
                          raw["start"]],
               "during": [raw["start"], raw["end"]],
               "after": [raw["end"], raw["observed_end"]]}
    return {
        "version": 1, "batch": batch, "generated_at": raw["generated_at"],
        "run": raw.get("run", {}), "config": raw["config"],
        "windows": windows, "requested_end": raw["requested_end"],
        "partial_batch": raw["partial_batch"],
        "recovery_complete": (raw["observed_end"] >= raw["requested_end"]
                              and not raw["partial_batch"]),
        "results": outcomes,
        "selected_count": batch_summary.get("selected",
            len(raw.get("run", {}).get("selected", [])) or len(scenarios)),
        "pass_rate": 100 * outcomes["PASS"] / executed if executed else None,
        "wall_seconds": raw["end"] - raw["start"],
        "peak_concurrency": peak, "activity": activity,
        "cleanup_warnings": sum(s.get("cleanup") in ("warning", "failed")
                                for s in scenarios),
        "missing_timestamps": sum(s.get("started_at") is None
                                  for s in scenarios),
        "unfinished": batch_summary.get("unfinished", []),
        "infrastructure_errors": batch_summary.get("infrastructure_errors", []),
        "queries": query_status, "scenarios": scenarios,
        "services": [services[k] for k in sorted(services)],
        "events": sorted(events, key=lambda e: e["time"]), "notes": NOTES,
    }


def value(number, digits=1):
    return "N/A" if number is None else f"{number:.{digits}f}"


def statistic(service, metric, phase_name, key):
    return service.get("statistics", {}).get(metric, {}).get(
        phase_name, {}).get(key)


def table(headers, rows):
    rows = [[str(v).replace("\n", " ") for v in row] for row in rows]
    widths = [max([len(h), *[len(row[i]) for row in rows]])
              for i, h in enumerate(headers)]
    line = lambda row: "  ".join(str(v).ljust(w)
                                 for v, w in zip(row, widths))
    return "\n".join([line(headers), line(["-" * w for w in widths]),
                      *[line(row) for row in rows]])


def report_sections(report):
    results = report["results"]
    run = report["run"]
    details = [
        ["Batch", report["batch"]],
        ["Suite / mode", f"{run.get('suite', 'N/A')} / "
         f"{run.get('mode', 'N/A')}"],
        ["Workers requested / launched",
         str(run.get("workers_requested", "N/A"))
         + " / " + str(run.get("workers_launched", "N/A"))],
        ["Scenarios selected / completed",
         f"{report['selected_count']} / {len(report['scenarios'])}"],
        ["Peak concurrent scenarios", report["peak_concurrency"]],
        ["Execution wall time (s)", value(report["wall_seconds"])],
        ["PASS / FAIL / SKIP", " / ".join(str(results[k])
                                          for k in ("PASS", "FAIL", "SKIP"))],
        ["Executed pass rate (%)", value(report["pass_rate"])],
        ["Cleanup warnings / failures", report["cleanup_warnings"]],
        ["Missing scenario timestamps", report["missing_timestamps"]],
        ["Unfinished / infrastructure errors",
         f"{len(report['unfinished'])} / "
         f"{len(report['infrastructure_errors'])}"],
        ["Recovery window", "complete" if report["recovery_complete"]
         else "partial; rerun collect later to extend"],
        ["Batch observation", "partial" if report["partial_batch"]
         else "completed report collection"],
    ]
    for name, bounds in report["windows"].items():
        details.append([name.title() + " (UTC)", " to ".join(map(utc, bounds))])
    sections = [("Run", ["Field", "Value"], details)]
    rows = []
    for service in report["services"]:
        cpu = [statistic(service, "cpu_mcpu", p, k)
               for p, k in (("before", "mean"), ("during", "peak"))]
        memory = [statistic(service, "memory_mib", p, k)
                  for p, k in (("before", "mean"), ("during", "peak"),
                               ("after", "last"))]
        replicas = service["lowest_replicas"]
        rows.append([service["service"], " -> ".join(map(value, cpu)),
                     " -> ".join(map(value, memory)),
                     value(service["restart_increases"], 0),
                     "/".join(value(v, 0) for v in replicas)
                     if replicas else "N/A",
                     value(service["degraded_sample_seconds"], 0)])
    sections.append(("All observed services", ["Service", "CPU base->peak mCPU",
                     "Memory base->peak->after MiB", "Restarts >=",
                     "Min avail/desired", "Degraded sampled s"], rows))
    rows = [[utc(e["time"]), e["service"], e["event"],
             ", ".join(e["overlapping_runs"]) or "none"]
            for e in report["events"]]
    sections.append(("Observed service events", ["UTC", "Service", "Event",
                     "Overlapping scenario run IDs"], rows))
    rows = [[s.get("scenario", s.get("run_id", "unknown")),
             s.get("worker", "local"), s.get("outcome", "N/A"),
             utc(s.get("started_at")), utc(s.get("ended_at")),
             s.get("cleanup", "N/A"), s.get("failure_reason") or ""]
            for s in report["scenarios"]]
    sections.append(("Scenarios", ["Scenario", "Worker", "Outcome", "Start",
                     "End", "Cleanup", "Failure detail"], rows))
    rows = [[name, q["status"], q["series"], q["samples"],
             q["error"] or "; ".join(q["warnings"] + q["infos"])]
            for name, q in report["queries"].items()]
    sections.append(("Query coverage", ["Query", "Status", "Series", "Samples",
                     "Detail"], rows))
    for service in report["services"]:
        rows = []
        for metric, label in METRICS.items():
            rows.append([label, *[value(statistic(service, metric, p, k))
                                 for p, k in (("before", "mean"),
                                              ("during", "mean"),
                                              ("during", "peak"),
                                              ("after", "last"))]])
        rows.append(["Observed last termination reasons",
                     ", ".join(service["termination_reasons"]) or "N/A",
                     "", "", ""])
        sections.append((service["service"], ["Metric", "Before avg",
                         "During avg", "During peak", "After last"], rows))
    return sections


def svg_chart(series, windows, step):
    if not series:
        return "<p>No samples.</p>"
    start, end = windows["before"][0], windows["after"][1]
    width, height = 720, 100
    maximum = max(v for _, v in series) or 1
    x = lambda t: width * (t - start) / max(1, end - start)
    y = lambda v: height - 5 - (height - 10) * v / maximum
    segments, current, previous = [], [], None
    for timestamp, sample in series:
        if timestamp < start or timestamp > end:
            continue
        if previous is not None and timestamp - previous > step * 1.5:
            segments.append(current)
            current = []
        current.append(f"{x(timestamp):.1f},{y(sample):.1f}")
        previous = timestamp
    segments.append(current)
    left, right = map(x, windows["during"])
    return (f'<svg viewBox="0 0 {width} {height}" role="img" '
            'aria-label="Metric over the observation window">'
            f'<rect x="{left:.1f}" y="0" width="{right-left:.1f}" '
            f'height="{height}" fill="#e8eefc"/>'
            + "".join('<polyline fill="none" stroke="#2859a7" '
                      f'stroke-width="2" points="{" ".join(part)}"/>'
                      for part in segments if part) + '</svg>')


def write_report(directory, report):
    sections = report_sections(report)
    text = "Ukama cluster KPI report\n\n" + "\n\n".join(
        title + "\n" + table(headers, rows)
        for title, headers, rows in sections)
    text += "\n\nInterpretation\n" + "\n".join(report["notes"]) + "\n"
    (directory / "report.txt").write_text(text)
    escape = lambda value: html.escape(str(value))
    body = ['<!doctype html><html lang="en"><meta charset="utf-8">',
            '<meta name="viewport" '
            'content="width=device-width,initial-scale=1">',
            '<title>Ukama cluster KPI report</title><style>',
            'body{font:14px system-ui;color:#182333;margin:32px}',
            'table{border-collapse:collapse;width:100%;margin:16px 0 32px}',
            'td,th{padding:8px;border:1px solid #d7dee7;text-align:left}',
            'th{background:#edf2f8}td{overflow-wrap:anywhere}',
            'svg{width:100%;max-width:900px;border-bottom:1px solid #ddd}',
            'details{margin:20px 0}h1,h2{color:#22436d}</style><body>',
            '<h1>Ukama cluster KPI report</h1>']
    for title, headers, rows in sections:
        body.append('<h2>' + escape(title) + '</h2><table><thead><tr>')
        body.extend('<th>' + escape(h) + '</th>' for h in headers)
        body.append('</tr></thead><tbody>')
        for row in rows:
            body.append('<tr>' + ''.join('<td>' + escape(v) + '</td>'
                                        for v in row) + '</tr>')
        body.append('</tbody></table>')
    body.append('<h2>Metric history</h2><p>Shaded: scenario execution. '
                'Left: baseline. Right: recovery. '
                'Each plot has its own scale.</p>')
    for service in report["services"]:
        body.append('<details><summary>' + escape(service["service"])
                    + '</summary>')
        for metric in ("cpu_mcpu", "memory_mib", "grpc_rps", "db_rps",
                       "rx_bps", "tx_bps", "replicas", "desired"):
            series = service["series"].get(metric, [])
            peak = max((v for _, v in series), default=None)
            body.append('<p>' + escape(METRICS[metric]) + ' — plot peak '
                        + escape(value(peak)) + '</p>')
            body.append(svg_chart(series, report["windows"],
                                  report["config"]["step"]))
        body.append('</details>')
    body.append('<h2>Interpretation</h2><ul>')
    body.extend('<li>' + escape(note) + '</li>' for note in report["notes"])
    body.append('</ul></body></html>')
    (directory / "report.html").write_text("\n".join(body))
    with (directory / "timeline.csv").open("w", newline="") as stream:
        writer = csv.writer(stream)
        writer.writerow(["utc", "epoch", "phase", "active_scenarios",
                         "service", "metric", "value"])
        for service in report["services"]:
            for metric, series in service["series"].items():
                for timestamp, sample in series:
                    start, end = report["windows"]["during"]
                    name = ("before" if timestamp < start else
                            "during" if timestamp <= end else "after")
                    writer.writerow([utc(timestamp), timestamp, name,
                                     len(overlap(report["scenarios"],
                                                 timestamp)),
                                     service["service"], metric, sample])


def compare(batches):
    if len(batches) < 2:
        raise ValueError("compare requires at least two batches")
    reports = [json.loads((p / "metrics/report.json").read_text())
               for p in batches]
    headers = ["Measurement", *[r["batch"] for r in reports]]
    rows = [["Execution wall seconds", *[value(r["wall_seconds"])
                                        for r in reports]],
            ["Selected scenarios", *[r["selected_count"] for r in reports]],
            ["PASS / FAIL / SKIP", *[
                "/".join(str(r["results"][k]) for k in ("PASS", "FAIL", "SKIP"))
                for r in reports]],
            ["Cleanup warnings", *[r["cleanup_warnings"] for r in reports]],
            ["Peak concurrency", *[r["peak_concurrency"] for r in reports]],
            ["Executed pass rate (%)",
             *[value(r["pass_rate"]) for r in reports]],
            ["Namespace filter", *[r["config"]["namespace"] for r in reports]],
            ["Service filter", *[r["config"]["service"] for r in reports]],
            ["Step (s)", *[r["config"]["step"] for r in reports]],
            ["Workers requested", *[r["run"].get("workers_requested", "N/A")
                                     for r in reports]],
            ["Source/binary SHA256", *[r["run"].get("source_sha256",
                                      r["run"].get("binary_sha256", "N/A"))
                                      for r in reports]]]
    lookup = [{s["service"]: s for s in r["services"]} for r in reports]
    identities = sorted(set().union(*(set(items) for items in lookup)))
    for identity in identities:
        for key, label in (("restart_increases", "observed restarts"),
                           ("degraded_sample_seconds", "degraded sampled s")):
            rows.append([identity + " " + label,
                         *[value(items.get(identity, {}).get(key))
                           for items in lookup]])
        for metric in ("cpu_mcpu", "memory_mib", "grpc_rps", "db_rps"):
            rows.append([identity + " peak " + METRICS[metric],
                         *[value(statistic(items.get(identity, {}), metric,
                                           "during", "peak"))
                           for items in lookup]])
    return ("Ukama cluster KPI comparison\n"
            "Runs may differ in scenario selection, timing, cluster state and "
            "background traffic. Check each report's coverage.\n\n"
            + table(headers, rows) + "\n")
