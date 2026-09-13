#!/usr/bin/env python3
"""Reject incomplete/duplicate/unassigned scenario results before publishing DONE."""
import json, pathlib, sys

def scenario_name(path):
    """Match repository-relative shard paths to suite-relative report names."""
    for suite in ("p0", "resilience"):
        prefix = f"scenarios/{suite}/"
        if path.startswith(prefix):
            return path[len(prefix):]
    return path


expected = [scenario_name(p) for p in pathlib.Path(sys.argv[1]).read_text().splitlines() if p]
try:
    report = json.loads(pathlib.Path(sys.argv[2]).read_text())
    rows = report["results"]
    actual = [r["scenario"] for r in rows]
    valid = (sorted(actual) == sorted(expected) and report["total"] == len(expected)
             and all(r["outcome"] in {"PASS", "FAIL", "SKIP"} for r in rows)
             and all(report[key] == sum(r["outcome"] == value for r in rows)
                     for key, value in (("passed", "PASS"), ("failed", "FAIL"), ("skipped", "SKIP"))))
except (ValueError, KeyError, TypeError, OSError):
    valid = False
if not valid:
    sys.exit("Worker report does not account for every assigned scenario exactly once")
