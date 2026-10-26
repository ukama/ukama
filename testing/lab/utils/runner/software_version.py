#!/usr/bin/env python3
"""Resolve the starter build version from the images used by this batch.

SPDX-License-Identifier: MPL-2.0
Copyright (c) 2026-present, Ukama Inc.
"""
import argparse
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import uuid

from common import RunnerError, atomic_json, run


CURRENT = "ULAB_SOFTWARE_CURRENT_VERSION"
DEFAULT_APP = "starter.d"
KINDS = ("anode", "cnode", "tnode")
VERSION = r"v?\d+\.\d+\.\d+(?:[-+][A-Za-z0-9][A-Za-z0-9.+-]*)?"

# The virtual-node image includes Python. Override its normal entrypoint and
# only read starter metadata. starter.d ignores argv, including --version.
PROBE = r'''
import json
from pathlib import Path
import re
import shutil
import subprocess
import sys

app, binary, version_file = sys.argv[1:]
starter = app in ("starter", "starterd", "starter.d")

def report(output, source):
    print(json.dumps({"output": output, "source": source}))

try:
    if version_file:
        report(Path(version_file).read_text(), version_file)
    elif starter:
        header = Path("/conf/starterd/version.h")
        if not binary and header.is_file():
            report(header.read_text(), str(header))
        else:
            path = Path(binary or "/sbin/starter.d")
            data = path.read_bytes()
            if not data.startswith(b"\x7fELF"):
                raise ValueError("starter binary is not ELF: " + str(path))
            # VERSION is a NUL-terminated string used by /v1/version.
            # Match complete strings, not compiler/library version fragments.
            pattern = (rb"v?\d+\.\d+\.\d+"
                       rb"(?:[-+][A-Za-z0-9][A-Za-z0-9.+-]*)?")
            versions = {part.decode("ascii") for part in data.split(b"\0")
                        if re.fullmatch(pattern, part)}
            if len(versions) != 1:
                raise ValueError("starter binary has no unique embedded "
                                 "version; candidates=" +
                                 repr(sorted(versions)) + "; preserve "
                                 "/conf/starterd/version.h in the image")
            report(versions.pop(), str(path) + " (embedded VERSION)")
    else:
        # Retain explicit custom-app support. Never execute starter.d here.
        binary = binary or shutil.which(app) or "/sbin/" + app
        if Path(binary).name in ("starter", "starterd", "starter.d"):
            raise ValueError("select ULAB_SOFTWARE_VERSION_APP=starter.d "
                             "to read starter metadata without executing it")
        result = subprocess.run([binary, "--version"], check=True,
                                capture_output=True, text=True, timeout=15)
        report(result.stdout + "\n" + result.stderr, binary + " --version")
except (OSError, ValueError, subprocess.SubprocessError) as exc:
    print(str(exc), file=sys.stderr)
    sys.exit(1)
'''


def parse_version(output):
    """Accept VERSION headers, bare versions or version-labelled output."""
    versions = set()
    for raw in output.splitlines():
        line = re.sub(r"\x1b\[[0-9;]*m", "", raw).strip()
        define = re.fullmatch(r'#\s*define\s+VERSION\s+"(' + VERSION +
                              r')"\s*(?://.*|/\*.*\*/)?', line)
        if define:
            versions.add(define.group(1))
            continue
        if re.fullmatch(VERSION, line):
            versions.add(line)
            continue
        match = re.search(r"\bversion\s*(?:[:=]|is)?\s*(" + VERSION +
                          r")(?![A-Za-z0-9_.+-])", line, re.IGNORECASE)
        if match:
            versions.add(match.group(1))
    if len(versions) != 1:
        raise RunnerError("image did not provide one unambiguous version")
    return versions.pop()


def image_version(image_id, env):
    app = env.get("ULAB_SOFTWARE_VERSION_APP", DEFAULT_APP)
    if not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9_.-]*", app):
        raise RunnerError("invalid ULAB_SOFTWARE_VERSION_APP")
    binary = env.get("ULAB_SOFTWARE_VERSION_BINARY", "")
    version_file = env.get("ULAB_SOFTWARE_VERSION_FILE", "")
    name = "ulab-version-" + uuid.uuid4().hex
    command = ["podman", "run", "--rm", "--name", name, "--pull=never",
               "--network=none", "--read-only", "--cap-drop=ALL",
               "--security-opt=no-new-privileges", "--entrypoint",
               "/usr/bin/python3", image_id, "-c", PROBE, app, binary,
               version_file]
    try:
        result = run(command, timeout=30, check=False)
        if result.returncode:
            raise RunnerError("cannot read " + app + " version from image " +
                              image_id +
                              ": " + result.stderr.strip())
        record = json.loads(result.stdout)
        return {"version": parse_version(record["output"]),
                "source": record["source"]}
    except subprocess.TimeoutExpired as exc:
        raise RunnerError("image version probe timed out: " + image_id) from exc
    finally:
        run(["podman", "rm", "-f", name], check=False, timeout=30)


def required(scenarios):
    return any("${" + CURRENT + "}" in Path(path).read_text()
               for path in scenarios)


def resolve(lab, repo, env, *, expected_images=None):
    """Build missing bases, then probe immutable image IDs, not repo HEAD."""
    base = env.get("BASE_IMAGE_REPO", "testing/virtualnode-base")
    runtime = env.get("NODE_RUNTIME", "starter")
    tags = [f"{base}:{kind}-{runtime}" for kind in KINDS]
    nodes = None
    records = []
    for kind, tag in zip(KINDS, tags):
        if run(["podman", "image", "exists", tag], check=False).returncode:
            if expected_images is not None:
                raise RunnerError(f"exported starter image is missing: {tag}")
            if nodes is None:
                # Reuse the AWS runner's existing Factory lookup/build route.
                from artifacts import factory_ids
                url = env.get("UKAMA_LAB_FACTORY_URL",
                              "http://factory-ukama.udev.ukama.com")
                nodes = factory_ids(url)
            if kind not in nodes:
                raise RunnerError(f"no Factory {kind} ID to build {tag}")
            print(f"software-version: building missing {tag}", file=sys.stderr)
            build_env = {**os.environ, **env, "BASE_IMAGE_REPO": base}
            run([lab / "scripts/build-node.sh", repo, nodes[kind], runtime],
                capture=False, timeout=14400, env=build_env, cwd=lab)
        info = json.loads(run(["podman", "image", "inspect", tag]).stdout)[0]
        image_id = info["Id"]
        if expected_images is not None:
            expected = expected_images.get(tag, "").removeprefix("sha256:")
            if image_id.removeprefix("sha256:") != expected:
                raise RunnerError(f"starter tag changed since export: {tag}")
        detected = image_version(image_id, env)
        records.append({"image": tag, "id": image_id, **detected})
        print(f"software-version: {tag} -> {detected['version']} "
              f"from {detected['source']} ({image_id[:19]})",
              file=sys.stderr, flush=True)
    versions = {record["version"] for record in records}
    if len(versions) != 1:
        detail = ", ".join(r["image"] + "=" + r["version"] for r in records)
        raise RunnerError("starter images have different app versions: " +
                          detail + "; rebuild matching bases")
    value = versions.pop()
    previous = env.get(CURRENT)
    if previous and previous != value:
        print(f"software-version: replacing exported {CURRENT}={previous}",
              file=sys.stderr)
    env[CURRENT] = value
    print(f"software-version: {CURRENT}={value}", file=sys.stderr, flush=True)
    return {"app": env.get("ULAB_SOFTWARE_VERSION_APP", DEFAULT_APP),
            "version": value, "images": records}


def prepare(lab, repo, scenarios, env, output, *, expected_images=None):
    if not required(scenarios):
        return False
    evidence = resolve(lab, repo, env, expected_images=expected_images)
    atomic_json(output, evidence)
    return True


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--lab-dir", type=Path, required=True)
    parser.add_argument("--repo", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    parser.add_argument("scenarios", nargs="*")
    args = parser.parse_args()
    env = dict(os.environ)
    try:
        if args.scenarios and not required(args.scenarios):
            return 0
        evidence = resolve(args.lab_dir.resolve(), args.repo.resolve(), env)
        atomic_json(args.output, evidence)
    except (OSError, ValueError, KeyError, RunnerError,
            subprocess.TimeoutExpired) as exc:
        print(f"software-version: ERROR: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
