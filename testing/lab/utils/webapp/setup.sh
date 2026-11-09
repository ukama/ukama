#!/usr/bin/env bash
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2026-present, Ukama Inc.

webapp_node_install() {
    python3 - "$1" <<'PY'
import hashlib
import os
from pathlib import Path
import platform
import re
import shutil
import subprocess
import sys
import tempfile
import urllib.request

target = Path(sys.argv[1])
system = {"Linux": "linux", "Darwin": "darwin"}.get(platform.system())
arch = {"x86_64": "x64", "aarch64": "arm64", "arm64": "arm64"}.get(platform.machine())
if not system or not arch:
    sys.exit("webapp: automatic Node installation supports Linux/macOS x64 and arm64")
base = "https://nodejs.org/dist/latest-v22.x/"
try:
    with urllib.request.urlopen(base + "SHASUMS256.txt", timeout=60) as response:
        checksums = response.read().decode()
    match = re.search(r"^([a-f0-9]{64})\s+(node-v22\.\d+\.\d+-" + system + "-" + arch + r"\.tar\.gz)$", checksums, re.M)
    if not match:
        raise ValueError("Node 22 archive is absent from the official checksums")
    digest, name = match.groups()
    with tempfile.TemporaryDirectory(prefix="node-", dir=target.parent) as temp:
        archive = Path(temp) / name
        checksum = hashlib.sha256()
        with urllib.request.urlopen(base + name, timeout=60) as response, archive.open("wb") as stream:
            while block := response.read(1024 * 1024):
                checksum.update(block)
                stream.write(block)
        if checksum.hexdigest() != digest:
            raise ValueError("Node download checksum mismatch")
        extracted = Path(temp) / "runtime"
        extracted.mkdir()
        subprocess.run(["tar", "--no-same-owner", "-xzf", str(archive), "--strip-components=1", "-C", str(extracted)], check=True)
        subprocess.run([str(extracted / "bin/node"), "--version"], check=True)
        if target.exists():
            shutil.rmtree(target)
        extracted.rename(target)
except (OSError, ValueError, subprocess.CalledProcessError) as exc:
    sys.exit(f"webapp: cannot install private Node runtime: {exc}")
PY
}

webapp_fingerprint() {
    python3 - "$@" <<'PY'
import hashlib
from pathlib import Path
import sys
digest = hashlib.sha256()
for name in sys.argv[1:]:
    root = Path(name)
    paths = sorted(p for p in root.rglob("*") if p.is_file()) if root.is_dir() else [root]
    for path in paths:
        digest.update(str(path).encode() + b"\0" + path.read_bytes() + b"\0")
print(digest.hexdigest())
PY
}

webapp_browser_check() {
    (cd adapters/webapp && "$WEBAPP_NODE" --input-type=module - "$1" <<'JS'
import {chromium, firefox, webkit} from 'playwright';
const name = process.argv[2];
const executablePath = process.env[`ULAB_WEBAPP_${name.toUpperCase()}_EXECUTABLE_PATH`] || process.env.ULAB_WEBAPP_EXECUTABLE_PATH;
try {
    const browser = await ({chromium, firefox, webkit})[name].launch({headless: true, timeout: 15000, ...(executablePath ? {executablePath} : {})});
    await browser.close();
} catch { process.exit(1); }
JS
    )
}

webapp_setup() {
    local cache browser major fingerprint stamp node_dir
    cache="${XDG_CACHE_HOME:-$HOME/.cache}/ukama-lab/webapp"
    browser="${ULAB_WEBAPP_BROWSER:-chromium}"
    case "$browser" in chromium|firefox|webkit) ;; *) echo 'error: invalid ULAB_WEBAPP_BROWSER' >&2; return 1 ;; esac
    case "${ULAB_WEBAPP_HEADLESS:-true}" in true|false) ;; *) echo 'error: ULAB_WEBAPP_HEADLESS must be true or false' >&2; return 1 ;; esac
    mkdir -p "$cache" || return 1
    # Serialize writes to the dependency tree when two batches start together.
    exec {WEBAPP_SETUP_LOCK}>"$cache/setup.lock" || return 1
    flock "$WEBAPP_SETUP_LOCK" || return 1
    major="$(node -p 'process.versions.node.split(".")[0]' 2>/dev/null)" || major=0
    if [[ ! "$major" =~ ^[0-9]+$ ]] || ((major < 22)) || ! command -v npm >/dev/null 2>&1; then
        node_dir="$cache/node"
        if [[ ! -x "$node_dir/bin/node" ]]; then
            echo 'WEBAPP setup: installing private Node 22 runtime'
            webapp_node_install "$node_dir" || return 1
        fi
        export PATH="$node_dir/bin:$PATH"
    fi
    WEBAPP_NODE="$(command -v node)" || return 1
    export WEBAPP_NODE
    fingerprint="$(webapp_fingerprint adapters/webapp/package.json adapters/webapp/package-lock.json)" || return 1
    stamp="adapters/webapp/node_modules/.ulab-dependencies"
    if [[ ! -f "$stamp" || "$(cat "$stamp")" != "$fingerprint" || ! -f adapters/webapp/node_modules/playwright/cli.js || ! -f adapters/webapp/node_modules/typescript/bin/tsc ]]; then
        echo 'WEBAPP setup: installing adapter dependencies'
        npm --prefix adapters/webapp ci --include=dev --no-audit --no-fund || return 1
        printf '%s\n' "$fingerprint" >"$stamp" || return 1
        rm -f adapters/webapp/dist/.ulab-build
    fi
    fingerprint="$(webapp_fingerprint adapters/webapp/src adapters/webapp/tsconfig.json adapters/webapp/package-lock.json)" || return 1
    stamp="adapters/webapp/dist/.ulab-build"
    if [[ ! -f "$stamp" || "$(cat "$stamp")" != "$fingerprint" || ! -f adapters/webapp/dist/cli.js ]]; then
        echo 'WEBAPP setup: building browser adapter'
        npm --prefix adapters/webapp run build || return 1
        printf '%s\n' "$fingerprint" >"$stamp" || return 1
    fi
    if ! webapp_browser_check "$browser"; then
        echo "WEBAPP setup: installing $browser"
        "$WEBAPP_NODE" adapters/webapp/node_modules/playwright/cli.js install "$browser" || return 1
        if ! webapp_browser_check "$browser"; then
            echo 'WEBAPP setup: installing browser system libraries'
            if ((EUID == 0)); then
                "$WEBAPP_NODE" adapters/webapp/node_modules/playwright/cli.js install-deps "$browser" || return 1
            elif command -v sudo >/dev/null 2>&1 && sudo -n true 2>/dev/null; then
                sudo -n -- "$WEBAPP_NODE" "$(pwd)/adapters/webapp/node_modules/playwright/cli.js" install-deps "$browser" || return 1
            else
                echo 'error: browser libraries are missing; run sudo -v, then rerun the batch to install Playwright dependencies' >&2
                return 1
            fi
            webapp_browser_check "$browser" || { echo 'error: browser cannot start after installation' >&2; return 1; }
        fi
    fi
    flock -u "$WEBAPP_SETUP_LOCK"
    exec {WEBAPP_SETUP_LOCK}>&-
    echo "WEBAPP setup: $browser ready"
}
