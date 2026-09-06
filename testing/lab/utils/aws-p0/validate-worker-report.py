#!/usr/bin/env python3
"""Reject incomplete/duplicate/unassigned scenario results before publishing DONE."""
import json, pathlib, sys
expected = [p[len("scenarios/p0/"):] if p.startswith("scenarios/p0/") else p for p in pathlib.Path(sys.argv[1]).read_text().splitlines() if p]
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
