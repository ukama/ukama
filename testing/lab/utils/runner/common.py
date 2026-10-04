#!/usr/bin/env python3
"""Small shared helpers. Only the AWS path imports this module."""
import hashlib
import json
import os
from pathlib import Path
import subprocess
import time


class RunnerError(RuntimeError):
    pass


def run(args, *, capture=True, check=True, timeout=180, env=None, cwd=None):
    result = subprocess.run([str(x) for x in args], text=True, capture_output=capture,
                            timeout=timeout, env=env, cwd=cwd)
    if check and result.returncode:
        # Do not print command arguments: they may contain credentials.
        detail = (result.stderr or "").strip() if capture else "see command output"
        raise RunnerError(f"{Path(str(args[0])).name} exited {result.returncode}: {detail}")
    return result


def atomic_json(path, value):
    path = Path(path)
    path.parent.mkdir(parents=True, exist_ok=True)
    temporary = path.with_name(path.name + ".tmp")
    temporary.write_text(json.dumps(value, indent=2) + "\n")
    temporary.chmod(0o600)
    temporary.replace(path)


def read_json(path):
    return json.loads(Path(path).read_text())


def sha256(path):
    digest = hashlib.sha256()
    with open(path, "rb") as stream:
        for data in iter(lambda: stream.read(8 * 1024 * 1024), b""):
            digest.update(data)
    return digest.hexdigest()


def now():
    return time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())


class AWS:
    def __init__(self, region):
        self.region = region
        self.env = {**os.environ, "AWS_PAGER": "", "AWS_DEFAULT_REGION": region,
                    "AWS_RETRY_MODE": "standard", "AWS_MAX_ATTEMPTS": "5"}

    def call(self, *args, timeout=180, check=True):
        return run(["aws", "--region", self.region, "--output", "json",
                    "--cli-connect-timeout", "15", "--cli-read-timeout", "60", *args],
                   timeout=timeout, check=check, env=self.env)

    def json(self, *args, **kwargs):
        output = self.call(*args, **kwargs).stdout
        return json.loads(output) if output.strip() else {}

    def copy(self, source, destination):
        # AWS CLI performs multipart transfers. A transfer error is never ignored.
        self.call("s3", "cp", str(source), str(destination), "--only-show-errors", timeout=7200)

    def sync(self, source, destination):
        self.call("s3", "sync", str(source), str(destination), "--only-show-errors", timeout=7200)

    def optional_json(self, bucket, key):
        # Distinguish a genuinely absent status object from IAM/network failures.
        p = self.call("s3", "cp", f"s3://{bucket}/{key}", "-", "--only-show-errors", check=False)
        if p.returncode:
            if "NoSuchKey" in p.stderr or "(404)" in p.stderr:
                return None
            raise RunnerError(f"cannot read S3 status: {p.stderr.strip()}")
        return json.loads(p.stdout)

    def head(self, bucket, key):
        p = self.call("s3api", "head-object", "--bucket", bucket, "--key", key, check=False)
        if not p.returncode:
            return json.loads(p.stdout)
        if "(404)" in p.stderr or "NoSuchKey" in p.stderr:
            return None
        raise RunnerError(f"cannot inspect S3 input: {p.stderr.strip()}")
