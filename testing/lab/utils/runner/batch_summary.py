#!/usr/bin/env python3
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

"""Category totals and text summary shared by local and AWS batch reports."""
import argparse
import json
from pathlib import Path


CATEGORY_LABELS = {
    "billing": "Billing",
    "console": "Console",
    "console-kpi": "Console KPI",
    "data-package": "Data package",
    "node": "Node",
    "sim": "SIM",
    "site": "Site",
    "software-update": "Software update",
    "ue": "UE",
    "usage": "Usage",
}


def totals(passed, failed):
    executed = passed + failed
    return {
        "executed": executed,
        "passed": passed,
        "failed": failed,
        "pass_rate": round(100 * passed / executed, 1) if executed else None,
        "fail_rate": round(100 * failed / executed, 1) if executed else None,
    }


def summarize(results):
    categories = {}
    failed_scenarios = []
    for result in results:
        outcome = result.get("outcome")
        if outcome not in ("PASS", "FAIL"):
            continue
        category = result.get("category") or result["scenario"].split("/", 1)[0]
        counts = categories.setdefault(category, {"PASS": 0, "FAIL": 0})
        counts[outcome] += 1
        if outcome == "FAIL":
            failed_scenarios.append(result["scenario"])
    rows = [
        {
            "category": category,
            "label": CATEGORY_LABELS.get(category, category.replace("-", " ").capitalize()),
            **totals(counts["PASS"], counts["FAIL"]),
        }
        for category, counts in sorted(categories.items())
    ]
    return {
        **totals(sum(row["passed"] for row in rows), sum(row["failed"] for row in rows)),
        "categories": rows,
        "failed_scenarios": sorted(failed_scenarios),
    }


def format_summary(summary, suite="p0"):
    def rate(value):
        return f"{value:.1f}%" if value is not None else "N/A"

    width = max([20, *[len(row["label"]) for row in summary["categories"]]])
    header = (f"{'CATEGORY':<{width}} {'EXECUTED':>9} {'PASS':>6} {'FAIL':>6} "
              f"{'PASS RATE':>11} {'FAIL RATE':>11}")
    separator = "-" * len(header)

    def row(label, values):
        return (f"{label:<{width}} {values['executed']:>9} {values['passed']:>6} "
                f"{values['failed']:>6} {rate(values['pass_rate']):>11} "
                f"{rate(values['fail_rate']):>11}")

    lines = [
        f"{suite.upper()} RESULTS — SKIPS EXCLUDED", "",
        f"Executed : {summary['executed']}",
        f"PASS     : {summary['passed']:>3}  ({rate(summary['pass_rate'])})",
        f"FAIL     : {summary['failed']:>3}  ({rate(summary['fail_rate'])})",
        "", header, separator,
        *[row(item["label"], item) for item in summary["categories"]],
        separator, row("OVERALL", summary), "", "Failed:",
    ]
    lines.extend(f"{index}. {scenario}" for index, scenario in
                 enumerate(summary["failed_scenarios"], 1))
    if not summary["failed_scenarios"]:
        lines.append("None")
    return "\n".join(lines)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("report", type=Path, help="existing batch-report.json")
    parser.add_argument("--suite", default="p0")
    args = parser.parse_args()
    report = json.loads(args.report.read_text(encoding="utf-8"))
    print(format_summary(summarize(report["results"]), args.suite))
