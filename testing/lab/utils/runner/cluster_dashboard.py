# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2026-present, Ukama Inc.
"""Offline presentation of saved cluster reports; no collection or verdicts."""
import html
import json
from pathlib import Path
import re

ASSETS = Path(__file__).with_name("dashboard")


def metric(item, name, phase="during", key="peak"):
    return item.get("statistics", {}).get(name, {}).get(phase, {}).get(key)


def total(item, name):
    return item.get("totals", {}).get(name, {}).get("during")


def number(value, digits=1):
    return "N/A" if value is None else f"{value:.{digits}f}"


def is_mesh(service):
    return service.get("service", "").split("/", 2)[-1].startswith(
        "ukama-mesh-node-")


def baseline_deficit(service, start):
    desired = dict(service.get("series", {}).get("desired", []))
    pairs = [(value, desired[t]) for t, value in
             service.get("series", {}).get("replicas", [])
             if t < start and t in desired and desired[t] > 0]
    return any(value < target for value, target in pairs) if pairs else None


def findings(report):
    rows = []
    start = report["windows"]["during"][0]

    def add(service, message, priority, before=None):
        context = ("present before run" if before is True else
                   "baseline below marker" if before is False else
                   "baseline unknown / event observed during run")
        rows.append({"service": service["service"], "message": message,
                     "context": context, "priority": priority})

    for service in report.get("services", []):
        resource = service.get("resources", {})
        for name, label in (("oom_kills", "OOMKilled terminations"),
                            ("oom_events", "OOM events")):
            count = total(resource, name)
            if count:
                add(service, f"At least {count:g} observed {label}", 0)
        if service.get("degraded_sample_seconds"):
            seconds = service["degraded_sample_seconds"]
            add(service, f"Replica deficit: {seconds:g} sampled seconds",
                1, baseline_deficit(service, start))
        if service.get("restart_increases"):
            count = service["restart_increases"]
            add(service, f"At least {count:g} observed restarts", 2)
        for name, label, marker in (
                ("cpu_throttle_pct", "Throttled periods", 10),
                ("memory_limit_pct", "Total memory / container limit", 90),
                ("cpu_limit_pct", "CPU / container limit", 90)):
            peak = metric(resource, name)
            before = metric(resource, name, "before", "mean")
            after = metric(resource, name, "after", "last")
            if peak is not None and peak >= marker:
                message = (f"{label}: {number(before)}% before, "
                           f"{number(peak)}% peak, {number(after)}% after")
                if name == "memory_limit_pct":
                    message += "; total includes cache"
                add(service, message, 3, None if before is None else
                    before >= marker)
    return sorted(rows, key=lambda row: (row["priority"], row["service"]))


def coverage(report):
    containers = report.get("containers", [])

    def present(item, name):
        return metric(item, name) is not None

    observed = [c for c in containers if present(c, "cpu_mcpu") or
                present(c, "memory_working_mib")]
    return {"containers": len(observed), **{
        name: sum(present(c, name) for c in observed)
        for name in ("cpu_limit_pct", "memory_limit_pct",
                     "cpu_throttle_pct", "memory_psi_pct")}}


def summary_text(report, rows, covered):
    run = report.get("run", {})
    result = report.get("results", {})
    services = report.get("services", [])
    mesh = sum(is_mesh(s) for s in services)
    queries = report.get("queries", {})
    ok = sum(q.get("status") == "OK" for q in queries.values())
    text = ["Ukama cluster KPI report", "", f"Batch: {report['batch']}",
            f"Suite: {run.get('suite', 'N/A')}",
            "Workers requested / launched: "
            f"{run.get('workers_requested', 'N/A')} / "
            f"{run.get('workers_launched', 'N/A')}",
            "PASS / FAIL / SKIP: " + " / ".join(
                str(result.get(k, 0)) for k in ("PASS", "FAIL", "SKIP")),
            f"Executed pass rate: {number(report.get('pass_rate'))}%",
            f"Wall time: {number(report.get('wall_seconds'))} s",
            f"Peak concurrency: {report.get('peak_concurrency', 'N/A')}",
            f"Cleanup warnings/failures: {report.get('cleanup_warnings', 0)}",
            f"Workload groups: {len(services)-mesh} other; {mesh} mesh-node",
            "Recovery: " + ("complete" if report.get("recovery_complete")
                             else "partial; recovery is not yet established"),
            "", "Query coverage", f"{ok}/{len(queries)} queries have data.",
            "Container CPU: "
            f"{covered['cpu_limit_pct']}/{covered['containers']} observed "
            "containers have usable CPU / limit data.",
            "Container memory: "
            f"{covered['memory_limit_pct']}/{covered['containers']} observed "
            "containers have usable memory / limit data."]
    for name, item in queries.items():
        if item.get("status") != "OK":
            detail = " ".join(str(item.get("error") or "no samples").split())
            text.append(f"  {name}: {item.get('status')} — {detail[:250]}")
    text += ["", "Findings (observations, not scenario verdicts)"]
    for row in rows[:12]:
        text += [f"  {row['service']}",
                 f"    {row['message']} [{row['context']}]"]
    if not rows:
        text.append("  No configured markers triggered in observed data; "
                    "check missing coverage.")
    if len(rows) > 12:
        text.append(f"  {len(rows)-12} more findings in report.html.")
    core = [s for s in services if not is_mesh(s)]
    for name, label in (("cpu_mcpu", "CPU mCPU"),
                        ("memory_mib", "Memory MiB")):
        ranked = sorted([s for s in core if metric(s, name) is not None],
                        key=lambda s: metric(s, name), reverse=True)[:8]
        text += ["", label + " — top 8 non-mesh peaks",
                 "  before average -> during peak -> after last"]
        for item in ranked:
            values = [metric(item, name, p, k) for p, k in
                      (("before", "mean"), ("during", "peak"),
                       ("after", "last"))]
            text.append("  " + item["service"] + ": " +
                        " -> ".join(number(v) for v in values))
    text += ["", "Scenarios", "Open report.html for scenario search, "
             "events, aligned charts and run comparison.",
             "Full measurements remain in report.json, timeline.csv and "
             "resources.csv.", "", "Interpretation",
             "Markers: throttled periods >=10%; CPU or total memory / "
             "limit >=90%. They do not change PASS/FAIL.",
             *report.get("notes", [])]
    return "\n".join(text) + "\n"


def write_outputs(directory, report):
    directory = Path(directory)
    rows, covered = findings(report), coverage(report)
    summary = summary_text(report, rows, covered)
    (directory / "report.txt").write_text(summary)
    # Escape the HTML parser's closing-script sequence, including file imports.
    data = json.dumps({"report": report, "findings": rows,
                       "coverage": covered}, separators=(",", ":"),
                      ensure_ascii=True, allow_nan=False)
    for char, escaped in (("<", "\\u003c"), (">", "\\u003e"),
                          ("&", "\\u0026")):
        data = data.replace(char, escaped)
    page = (ASSETS / "report.html").read_text()
    replacements = {"STYLE": (ASSETS / "report.css").read_text(),
                    "SCRIPT": (ASSETS / "report.js").read_text(),
                    "DATA": data, "SUMMARY": html.escape(summary)}
    # One substitution pass: user-provided values cannot become placeholders.
    page = re.sub(r"@@(STYLE|SCRIPT|DATA|SUMMARY)@@",
                  lambda match: replacements[match[1]], page)
    (directory / "report.html").write_text(page)
