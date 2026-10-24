#!/usr/bin/env python3
"""Delete only mesh resources belonging to a completed lab run.

SPDX-License-Identifier: MPL-2.0
Copyright (c) 2026-present, Ukama Inc.
"""
import argparse
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import time


KINDS = ("deployment", "replicaset", "pod", "service")
STATE_KEYS = {"FACTORY_NODE_ID", "TNODE_ID", "CNODE_ID", "ANODE_ID"}


def node_ids(directory):
    """Use the run's persisted claims, never the shared Factory inventory."""
    nodes = set()
    for subdir in ("runtime-nodes", "runtime-sites"):
        for path in sorted((directory / subdir).glob("*.env")):
            for line in path.read_text().splitlines():
                key, sep, value = line.partition("=")
                if not sep or key not in STATE_KEYS or not value:
                    continue
                if not re.fullmatch(r"[a-z0-9][a-z0-9-]*", value):
                    raise RuntimeError(f"invalid node ID in {path}: {key}")
                nodes.add(value)
    return sorted(nodes)


def owned_resources(items, nodes):
    names = {"ukama-mesh-node-" + node for node in nodes}
    owned = []
    for item in items:
        kind = item.get("kind", "").lower()
        name = item.get("metadata", {}).get("name", "")
        if kind not in KINDS:
            continue
        # Match the complete factory ID, with a boundary before any suffix.
        # This also includes generated deployment/service names and pod hashes.
        if any(name == base or name.startswith(base + "-") for base in names):
            owned.append((kind, name))
    return sorted(owned, key=lambda value: (KINDS.index(value[0]), value[1]))


def cleanup(directory, dry_run=False):
    nodes = node_ids(directory)
    if not nodes:
        print("cleanup-mesh: no recorded virtual nodes")
        return
    kubectl = os.environ.get("ULAB_KUBECTL") or shutil.which("kubectl")
    if not kubectl:
        candidate = Path.home() / "kubectl"
        kubectl = str(candidate) if candidate.is_file() else ""
    if not kubectl or not shutil.which(kubectl):
        raise RuntimeError("kubectl unavailable; set ULAB_KUBECTL")
    namespace = os.environ.get("ULAB_MESH_NAMESPACE", "ukama-messaging")

    def kube(*args):
        command = [kubectl, *args, "-n", namespace,
                   "--request-timeout=20s"]
        result = subprocess.run(command, capture_output=True, text=True,
                                timeout=45)
        if result.returncode:
            detail = result.stderr.strip() or result.stdout.strip()
            raise RuntimeError(f"kubectl {args[0]} failed: {detail}")
        return result.stdout

    def remaining():
        data = json.loads(kube("get", "deploy,rs,pod,svc", "-o", "json"))
        if not isinstance(data.get("items"), list):
            raise RuntimeError("kubectl returned no resource list")
        return owned_resources(data["items"], nodes)

    print(f"cleanup-mesh: namespace={namespace} nodes={len(nodes)}")
    resources = remaining()
    if dry_run:
        for kind, name in resources:
            print(f"cleanup-mesh: would delete {kind}/{name}")
        return

    # Stop controllers before hard-deleting their pods. Explicit ReplicaSet
    # removal closes the background garbage-collection/replacement-pod race.
    for attempt in range(3):
        if not resources:
            print("cleanup-mesh: verified no mesh resources remain")
            return
        for kind in KINDS:
            selected = [f"{kind}/{name}" for k, name in resources if k == kind]
            # Keep command size bounded for scenarios containing many sites.
            for offset in range(0, len(selected), 50):
                group = selected[offset:offset + 50]
                print("cleanup-mesh: delete " + " ".join(group), flush=True)
                options = ["--ignore-not-found=true"]
                if kind == "pod":
                    options += ["--force", "--grace-period=0", "--wait=false"]
                else:
                    options += ["--wait=true", "--timeout=30s"]
                    if kind in ("deployment", "replicaset"):
                        options += ["--cascade=background"]
                kube("delete", *group, *options)
        resources = remaining()
        if resources and attempt < 2:
            time.sleep(1)
    if resources:
        names = ", ".join(f"{kind}/{name}" for kind, name in resources)
        raise RuntimeError("mesh resources still present: " + names)
    print("cleanup-mesh: verified no mesh resources remain")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("run_dir", type=Path)
    parser.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()
    if not args.run_dir.is_dir():
        parser.error("run directory does not exist")
    try:
        cleanup(args.run_dir, args.dry_run)
    except (OSError, ValueError, RuntimeError,
            subprocess.TimeoutExpired) as exc:
        print(f"cleanup-mesh: ERROR: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
