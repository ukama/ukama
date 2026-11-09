#!/usr/bin/env python3
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2026-present, Ukama Inc.
"""Apply batch browser settings to a copy of a scenario."""
import json
import os
from pathlib import Path
import re
import sys
from urllib.parse import urlsplit


def prepare(source, environment):
    url = environment["ULAB_WEBAPP_BASE_URL"].rstrip("/")
    parsed = urlsplit(url)
    if (parsed.scheme not in ("http", "https") or not parsed.hostname or
            parsed.username or parsed.password or parsed.path or parsed.query or parsed.fragment):
        raise ValueError("ULAB_WEBAPP_BASE_URL must be an http(s) origin")
    browser = environment["ULAB_WEBAPP_BROWSER"]
    headless = environment["ULAB_WEBAPP_HEADLESS"]
    if browser not in ("chromium", "firefox", "webkit") or headless not in ("true", "false"):
        raise ValueError("invalid webapp browser or headless setting")
    changes = {"base_url": json.dumps(url), "browser": browser, "headless": headless}
    lines = source.splitlines(keepends=True)
    starts = [i for i, line in enumerate(lines) if re.fullmatch(r"webapp:\s*(?:#.*)?\n?", line)]
    if len(starts) != 1:
        raise ValueError("scenario must contain one webapp block")
    start = starts[0] + 1
    end = next((i for i in range(start, len(lines))
                if lines[i].strip() and not lines[i][0].isspace() and not lines[i].startswith("#")), len(lines))
    seen = set()
    for i in range(start, end):
        field = re.match(r"^  ([a-z_]+):\s*(.*?)\s*$", lines[i])
        if not field:
            continue
        key, value = field.groups()
        if key in seen:
            raise ValueError("duplicate webapp setting")
        seen.add(key)
        if key == "auth_state" and value.split("#", 1)[0].strip().strip("\"'") == ".auth/owner.json":
            changes[key] = json.dumps(environment["ULAB_WEBAPP_BATCH_AUTH_STATE"])
        if key in changes:
            lines[i] = f"  {key}: {changes[key]}\n"
    lines[end:end] = [f"  {key}: {value}\n" for key, value in changes.items() if key not in seen]
    return "".join(lines)


if __name__ == "__main__":
    os.umask(0o077)
    try:
        source, destination = map(Path, sys.argv[1:])
        result = prepare(source.read_text(), os.environ)
        with destination.open("x") as stream:
            stream.write(result)
    except (OSError, ValueError, KeyError) as exc:
        print(f"webapp profile: {exc}", file=sys.stderr)
        sys.exit(2)
