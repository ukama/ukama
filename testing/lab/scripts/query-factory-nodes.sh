#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

FACTORY_URL="${FACTORY_URL:-http://localhost:8070}"
CSV="${CSV:-}"

need_cmd() {
    command -v "$1" >/dev/null 2>&1 || {
        echo "error: missing required command: $1" >&2
        exit 1
    }
}

need_cmd curl
need_cmd python3

tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT

curl -sS \
    -H "accept: application/json" \
    "${FACTORY_URL}/v1/nodefactory/nodes?isProvisioned=false" >"$tmp"

python3 - "$tmp" "$CSV" <<'PY'
import csv
import json
import sys

path = sys.argv[1]
csv_path = sys.argv[2] if len(sys.argv) > 2 else ""

with open(path) as f:
    data = json.load(f)

def find_list(x):
    if isinstance(x, list):
        return x
    if isinstance(x, dict):
        for k in ("nodes", "Nodes", "data", "Data", "items", "Items", "results", "Results"):
            if k in x:
                r = find_list(x[k])
                if r is not None:
                    return r
    return None

nodes = find_list(data)
if nodes is None:
    print("error: could not find node list in response", file=sys.stderr)
    print(json.dumps(data, indent=2)[:4000], file=sys.stderr)
    sys.exit(1)

def get(n, *keys):
    for k in keys:
        if isinstance(n, dict) and k in n and n[k] not in (None, ""):
            return n[k]
    return ""

def norm_type(t, node_id):
    t = str(t).lower()
    node_id = str(node_id).lower()

    if t in ("tnode", "tower", "tower_node") or "-tnode-" in node_id:
        return "tnode"
    if t in ("cnode", "controller", "controller_node") or "-cnode-" in node_id:
        return "cnode"
    if t in ("anode", "amplifier", "access", "amplifier_node") or "-anode-" in node_id:
        return "anode"

    return t

def derived(node_id, kind):
    return node_id.replace("-tnode-", f"-{kind}-")

by_id = {}
tnodes = {}
cnode_count = 0
anode_count = 0

for n in nodes:
    nid = str(get(n, "id", "Id", "nodeId", "NodeId", "node_id"))
    if not nid:
        continue

    ntype = norm_type(
        get(n, "type", "Type", "nodeType", "NodeType", "node_type"),
        nid,
    )

    rec = {
        "id": nid,
        "type": ntype,
        "raw": n,
    }

    by_id[nid] = rec

    if ntype == "tnode":
        tnodes[nid] = rec
    elif ntype == "cnode":
        cnode_count += 1
    elif ntype == "anode":
        anode_count += 1

complete = []
incomplete = []

for tid in sorted(tnodes):
    cid = derived(tid, "cnode")
    aid = derived(tid, "anode")

    has_c = cid in by_id and by_id[cid]["type"] == "cnode"
    has_a = aid in by_id and by_id[aid]["type"] == "anode"

    if has_c and has_a:
        complete.append((tid, cid, aid))
    else:
        incomplete.append((tid, has_c, has_a))

print("factory bundle audit")
print(f"  total nodes: {len(by_id)}")
print(f"  tnodes:      {len(tnodes)}")
print(f"  cnodes:      {cnode_count}")
print(f"  anodes:      {anode_count}")
print(f"  complete:    {len(complete)}")
print(f"  incomplete:  {len(incomplete)}")
print()

if complete:
    print("complete bundles:")
    for tid, cid, aid in complete[:25]:
        print(f"  {tid}  cnode={cid}  anode={aid}")
    if len(complete) > 25:
        print(f"  ... {len(complete) - 25} more")
    print()

if incomplete:
    print("incomplete bundles:")
    for tid, has_c, has_a in incomplete[:50]:
        print(f"  {tid}  cnode={str(has_c).lower()}  anode={str(has_a).lower()}")
    if len(incomplete) > 50:
        print(f"  ... {len(incomplete) - 50} more")
    print()

if csv_path:
    print("csv verification:")
    missing = 0
    bad_bundle_id = 0

    with open(csv_path, newline="") as f:
        r = csv.DictReader(f)
        for row in r:
            tid = row.get("tnode", "")
            cid = row.get("cnode", "")
            aid = row.get("anode", "")

            ok = True

            if tid not in by_id:
                print(f"  missing tnode: {tid}")
                missing += 1
                ok = False
            if cid not in by_id:
                print(f"  missing cnode: {cid}")
                missing += 1
                ok = False
            if aid not in by_id:
                print(f"  missing anode: {aid}")
                missing += 1
                ok = False

            expected_cid = derived(tid, "cnode")
            expected_aid = derived(tid, "anode")

            if cid != expected_cid or aid != expected_aid:
                print(f"  bad derived ids set={row.get('set')}:")
                print(f"    tnode={tid}")
                print(f"    cnode={cid} expected={expected_cid}")
                print(f"    anode={aid} expected={expected_aid}")
                bad_bundle_id += 1
                ok = False

            if ok:
                if by_id[tid]["type"] != "tnode":
                    print(f"  wrong type for tnode id: {tid} type={by_id[tid]['type']}")
                    bad_bundle_id += 1
                if by_id[cid]["type"] != "cnode":
                    print(f"  wrong type for cnode id: {cid} type={by_id[cid]['type']}")
                    bad_bundle_id += 1
                if by_id[aid]["type"] != "anode":
                    print(f"  wrong type for anode id: {aid} type={by_id[aid]['type']}")
                    bad_bundle_id += 1

    print(f"  missing:       {missing}")
    print(f"  bad_bundle_id: {bad_bundle_id}")
PY
