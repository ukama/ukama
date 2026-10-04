#!/usr/bin/env python3
"""AWS lifecycle only. Scenario execution remains in run-scenarios.sh."""
import argparse
import base64
import csv
import fcntl
import hashlib
import ipaddress
import json
import os
from pathlib import Path
import re
import shlex
import shutil
import subprocess
import sys
import tarfile
import time
import uuid
from common import AWS, RunnerError, atomic_json, now, read_json, run, sha256
from config import HERE, ROOT, environment, executable, load_config, portable_vpn
from artifacts import image_archive, payload, relative_scenarios


def lifecycle(work, name, message):
    """Print transitions once and retain a plain-text batch lifecycle log."""
    line = f"{now()} {name}: {message}"
    print("\n" + line, flush=True)
    with (work / "aws-lifecycle.log").open("a") as stream:
        stream.write(line + "\n")


def observe_instance(work, name, worker, instance):
    if not instance:
        return
    state = instance["State"]["Name"]
    if worker.get("observed_ec2_state") != state:
        suffix = "; decommissioned" if state == "terminated" else ""
        lifecycle(work, name, f"EC2 {instance['InstanceId']} {state}{suffix}")
        worker["observed_ec2_state"] = state


def observe_worker(work, name, worker, status):
    phase, stage = status.get("phase"), status.get("stage")
    key = (phase, stage)
    if phase and worker.get("observed_worker_stage") != list(key):
        descriptions = {
            "source-download": "downloading lab and repository",
            "vpn-and-dns": "setting up VPN and DNS",
            "image-download": "downloading shared starter images",
            "image-load": "loading starter images into local Podman",
            "runtime-check": "installing packaged runtime and kubectl cleanup tool",
            "connectivity": "checking existing worker connectivity",
            "ready": "worker setup complete",
            "scenario-execution": "worker setup complete; running assigned scenarios",
        }
        if phase in ("DONE", "STOPPED", "INFRA_ERROR", "BOOTSTRAP_ERROR"):
            message = f"worker {phase.lower()}"
        else:
            message = descriptions.get(stage, stage or phase.lower())
        lifecycle(work, name, message)
        worker["observed_worker_stage"] = list(key)
    current = status.get("progress", {}).get("current")
    if current and current != "-" and worker.get("observed_scenario") != current:
        lifecycle(work, name, f"scenario {current}")
        worker["observed_scenario"] = current


def request_termination(aws, work, name, worker, instance, reason):
    aws.call("ec2", "terminate-instances", "--instance-ids", instance["InstanceId"])
    if not worker.get("termination_requested_at"):
        worker["termination_requested_at"] = now()
        lifecycle(work, name, f"EC2 {instance['InstanceId']} termination requested ({reason})")


def confirm_termination(aws, work, state, timeout=60):
    """Observe teardown, without changing verdicts or retrying scenario work."""
    pending = {name for name, worker in state["workers"].items()
               if worker.get("instance_id") and
               (worker.get("termination_requested_at") or worker.get("observed_ec2_state") == "shutting-down")}
    if not pending:
        return
    deadline = time.monotonic() + timeout
    for attempt in range(max(1, int(timeout / 5) + 1)):
        try:
            current = instances(aws, state)
        except (RunnerError, OSError, ValueError, subprocess.TimeoutExpired) as exc:
            lifecycle(work, "AWS", f"could not confirm EC2 termination: {exc}; check batch instance IDs")
            break
        for name in sorted(pending.copy()):
            instance = current.get(name)
            if not instance:
                lifecycle(work, name, "instance no longer returned by EC2; termination not independently confirmed")
                pending.remove(name)
                continue
            observe_instance(work, name, state["workers"][name], instance)
            if instance["State"]["Name"] == "terminated":
                pending.remove(name)
        if not pending:
            break
        if time.monotonic() >= deadline or attempt >= int(timeout / 5):
            for name in sorted(pending):
                worker = state["workers"][name]
                lifecycle(work, name, f"EC2 {worker['instance_id']} termination not yet confirmed; last state={worker.get('observed_ec2_state', 'unknown')}")
            break
        time.sleep(min(5, max(0, deadline - time.monotonic())))
    atomic_json(work / "state.json", state)


def print_report(path):
    colors = {"PASS": "\033[1;32m", "FAIL": "\033[1;31m"} if sys.stdout.isatty() else {}
    for line in path.read_text().splitlines():
        label = line.split(" ", 1)[0]
        if label in colors:
            line = colors[label] + label + "\033[0m" + line[len(label):]
        print(line)


def ensure_stack(aws, cfg):
    print(f"Preparing AWS infrastructure: {cfg['STACK_NAME']} ({cfg['REGION']})", flush=True)
    parameters = [f"VpcCidr={cfg['VPC_CIDR']}", f"SubnetCidr={cfg['SUBNET_CIDR']}",
                  f"RetentionDays={cfg['RESULT_RETENTION_DAYS']}"]
    aws.call("cloudformation", "deploy", "--stack-name", cfg["STACK_NAME"],
             "--template-file", str(HERE / "infrastructure.json"), "--capabilities", "CAPABILITY_IAM",
             "--parameter-overrides", *parameters, "--no-fail-on-empty-changeset", timeout=1800)
    data = aws.json("cloudformation", "describe-stacks", "--stack-name", cfg["STACK_NAME"])["Stacks"][0]
    return {item["OutputKey"]: item["OutputValue"] for item in data["Outputs"]}


def find_ami(aws, cfg):
    ami = cfg["AMI_ID"]
    if not ami:
        ami = aws.json("ssm", "get-parameter", "--name",
                       "/aws/service/canonical/ubuntu/server/24.04/stable/current/amd64/hvm/ebs-gp3/ami-id")["Parameter"]["Value"]
    data = aws.json("ec2", "describe-images", "--image-ids", ami)["Images"]
    if len(data) != 1 or data[0]["Architecture"] != "x86_64" or data[0]["State"] != "available":
        raise RunnerError("AMI must be an available Ubuntu 24.04 x86_64 image")
    instance_type = aws.json("ec2", "describe-instance-types", "--instance-types", cfg["INSTANCE_TYPE"])["InstanceTypes"][0]
    if "x86_64" not in instance_type["ProcessorInfo"]["SupportedArchitectures"]:
        raise RunnerError("INSTANCE_TYPE must support x86_64")
    return ami, data[0]["RootDeviceName"]


def partition(scenarios, workers):
    return [scenarios[index::min(workers, len(scenarios))] for index in range(min(workers, len(scenarios)))]


def s3_routes(aws, region):
    """AWS's regional S3 gateway-endpoint prefixes, fetched once per batch."""
    name = f"com.amazonaws.{region}.s3"
    data = aws.json("ec2", "describe-prefix-lists", "--filters", f"Name=prefix-list-name,Values={name}")
    cidrs = sorted({str(ipaddress.IPv4Network(cidr))
                    for item in data.get("PrefixLists", []) if item.get("PrefixListName") == name
                    for cidr in item.get("Cidrs", [])})
    if not cidrs or any(ipaddress.IPv4Network(c).prefixlen < 2 for c in cidrs):
        raise RunnerError(f"cannot discover regional S3 routes for {name}; no workers launched")
    return cidrs


def user_data(state, worker):
    # No VPN profile, scenario password, kubeconfig or AWS key in EC2 user-data.
    values = {"ULAB_RUNNER_REGION": state["region"], "ULAB_RUNNER_BUCKET": state["bucket"],
              "ULAB_RUNNER_BATCH": state["batch"], "ULAB_RUNNER_WORKER": worker,
              "ULAB_RUNNER_MAX_HOURS": str(state["max_hours"])}
    prefix = f"s3://{state['bucket']}/inputs/{state['batch']}"
    assignments = "\n".join("export " + k + "=" + shlex.quote(v) for k, v in values.items())
    return f'''#!/bin/bash
set -Eeuo pipefail
umask 077
{assignments}
export HOME=/root AWS_DEFAULT_REGION="$ULAB_RUNNER_REGION" AWS_PAGER="" DEBIAN_FRONTEND=noninteractive
mkdir -p {ROOT}/bootstrap {ROOT}/results
exec > >(tee -a {ROOT}/results/bootstrap.log /dev/console) 2>&1
bootstrap_stage=packages
bootstrap_failed() {{
    rc=$?
    trap - ERR
    printf 'UKAMA_RUNNER_BOOTSTRAP_ERROR stage=%s exit=%s\\n' "$bootstrap_stage" "$rc"
    printf '{{"phase":"BOOTSTRAP_ERROR","error":"Bootstrap stage %s exited %s; see EC2 console","lab_exit":null}}\\n' "$bootstrap_stage" "$rc" >{ROOT}/results/status.json
    if command -v aws >/dev/null; then
        aws s3 cp {ROOT}/results/bootstrap.log {shlex.quote(f"s3://{state['bucket']}/results/{state['batch']}/{worker}/bootstrap.log")} --only-show-errors || true
        aws s3 cp {ROOT}/results/status.json {shlex.quote(f"s3://{state['bucket']}/results/{state['batch']}/{worker}/status.json")} --only-show-errors || true
    fi
    exit "$rc"
}}
trap bootstrap_failed ERR
# Independently bounds instance lifetime even if the controller disconnects.
shutdown -h +$((ULAB_RUNNER_MAX_HOURS * 60))
printf 'UKAMA_RUNNER_BOOTSTRAP stage=packages\\n'
apt-get -o DPkg::Lock::Timeout=300 -o Acquire::Retries=3 update
apt-get -o DPkg::Lock::Timeout=300 -o Acquire::Retries=3 install -y ca-certificates curl unzip gnupg python3
bootstrap_stage=aws-cli
printf 'UKAMA_RUNNER_BOOTSTRAP stage=aws-cli\\n'
# Ubuntu's configured repositories may not provide awscli. The official
# installer verifies the downloaded CLI signature using the installed GnuPG.
if ! aws --version 2>/dev/null | grep -q '^aws-cli/2\\.'; then
    curl --fail --silent --show-error --location --proto '=https' --tlsv1.2 \\
        --retry 3 --connect-timeout 15 --max-time 120 \\
        https://awscli.amazonaws.com/v2/install.sh -o {ROOT}/install-aws-cli.sh
    timeout 900 bash {ROOT}/install-aws-cli.sh --system --quiet
fi
hash -r
aws --version
bootstrap_stage=bootstrap-download
printf 'UKAMA_RUNNER_BOOTSTRAP stage=bootstrap-download\\n'
aws s3 cp {shlex.quote(prefix + '/bootstrap.tar.gz')} {ROOT}/bootstrap.tar.gz --only-show-errors
printf '%s  %s\\n' {shlex.quote(state['bootstrap_sha256'])} {ROOT}/bootstrap.tar.gz | sha256sum --check
tar -xzf {ROOT}/bootstrap.tar.gz -C {ROOT}/bootstrap
bootstrap_stage=worker-entry
printf 'UKAMA_RUNNER_BOOTSTRAP stage=worker-entry\\n'
bash {ROOT}/bootstrap/worker-entry.sh
'''


def prepare_factory(args, scenarios, cfg, work):
    if args.factory_nodes == "0" and not args.prepare_only:
        return
    selected = work / "factory-scenarios.txt"
    selected.write_text("\n".join(str(args.lab_dir / p) for p in scenarios) + "\n")
    env = {**os.environ, **{k: str(v) for k, v in cfg["ENV"].items()},
           "UKAMA_REPO": str(args.repo), "LAB_BIN": str(args.lab_bin),
           "SCENARIO_ROOT": str(args.scenario_root), "P0_RUNS_DIR": str(work / "factory")}
    run([args.lab_dir / "utils/run-scenarios.sh", args.suite, "--mode", "local", "--scenario-list", selected,
         "--factory-nodes", args.factory_nodes, "--prepare-only"], env=env, cwd=args.lab_dir,
        capture=False, timeout=7200)


def launch(args):
    cfg = load_config(args.config)
    args.lab_dir = Path(args.lab_dir).resolve()
    if not args.repo:
        raise RunnerError("export UKAMA_REPO to your working Ukama repository")
    args.repo = Path(args.repo).expanduser().resolve()
    args.lab_bin = executable(args.lab_bin)
    args.scenario_root = Path(args.scenario_root).resolve()
    if not args.scenario_root.is_relative_to(args.lab_dir):
        raise RunnerError("AWS SCENARIO_ROOT must be within the lab directory")
    if not os.access(args.repo / "testing/node/mk_local_vnode.sh", os.X_OK):
        raise RunnerError("UKAMA_REPO must contain executable testing/node/mk_local_vnode.sh")
    for file in ("testing/ue/ue/Containerfile", "testing/ue/media/Containerfile"):
        if not (args.repo / file).is_file():
            raise RunnerError(f"UKAMA_REPO is missing {file}")
    scenarios = relative_scenarios([v for v in sys.stdin.read().split("\0") if v], args.lab_dir)
    if not scenarios:
        raise RunnerError("no scenarios selected")
    if any(not (args.lab_dir / p).is_relative_to(args.scenario_root) for p in scenarios):
        raise RunnerError("AWS scenario-list entries must be within SCENARIO_ROOT; set SCENARIO_ROOT=scenarios for a mixed-suite list")
    if args.batch_id and not re.fullmatch(r"[A-Za-z0-9._-]{1,64}", args.batch_id):
        raise RunnerError("AWS --batch-id must be 1–64 letters, digits, dots, underscores or hyphens")
    batch = args.batch_id or (time.strftime("%Y%m%dt%H%M%Sz", time.gmtime()) + "-" + uuid.uuid4().hex[:8])
    work = Path(args.out).expanduser().resolve() / batch
    work.mkdir(parents=True, exist_ok=False)
    work.chmod(0o700)
    batch_lock = (work / ".controller.lock").open("w")
    fcntl.flock(batch_lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
    env = environment(cfg)
    print(f"AWS batch directory: {work}", flush=True)
    prepare_factory(args, scenarios, cfg, work)
    if args.prepare_only:
        return 0
    for command in ("aws", "podman", "ldd", "curl"):
        executable(command)
    # Validate all local material before creating infrastructure or launching machines.
    assignments = partition(scenarios, args.workers)
    profiles = cfg["VPN_CONFIG_FILES"] or [cfg["VPN_CONFIG_FILE"]]
    if cfg["NETWORK_MODE"] == "vpn" and len(profiles) != 1 and len(profiles) < len(assignments):
        raise RunnerError("VPN_CONFIG_FILES must have at least one profile per worker, or use VPN_CONFIG_FILE")
    vpns = [portable_vpn(p, cfg["VPN_AUTH_FILE"]) for p in profiles] if cfg["NETWORK_MODE"] == "vpn" else []
    cache = Path(args.out).expanduser().resolve() / "image-cache"
    # Prevent two controllers from modifying the same local export cache.
    cache.mkdir(parents=True, exist_ok=True)
    with (cache / ".export.lock").open("w") as cache_lock:
        fcntl.flock(cache_lock, fcntl.LOCK_EX)
        images_path, image_info = image_archive(args.lab_dir, args.repo, cache, env)
    input_dir = work / "input"
    input_dir.mkdir(mode=0o700)
    private_paths = [Path(cfg[key]) for key in ("VPN_AUTH_FILE",) if cfg[key]]
    source_path = payload(args.lab_dir, args.repo, args.lab_bin, cfg, env, input_dir,
                          [work.parent, Path(args.config).resolve(), *profiles, *private_paths])
    source_digest = sha256(source_path)
    bootstrap = input_dir / "bootstrap.tar.gz"
    with tarfile.open(bootstrap, "w:gz") as archive:
        for path in sorted(HERE.iterdir()):
            if path.suffix == ".py" or path.name == "worker-entry.sh":
                archive.add(path, arcname=path.name)
    aws = AWS(cfg["REGION"])
    identity = aws.json("sts", "get-caller-identity")
    infra = ensure_stack(aws, cfg)
    ami, root_device = find_ami(aws, cfg)
    aws_s3_cidrs = s3_routes(aws, cfg["REGION"]) if cfg["NETWORK_MODE"] == "vpn" else []
    image_key = "images/" + image_info["sha256"] + ".tar.gz"
    head = aws.head(infra["Bucket"], image_key)
    if head is None:
        print("Uploading shared starter-image archive...", flush=True)
        aws.copy(images_path, f"s3://{infra['Bucket']}/{image_key}")
    elif head["ContentLength"] != image_info["size"]:
        raise RunnerError("cached S3 image archive has an unexpected size")
    state = {"version": 1, "batch": batch, "suite": args.suite, "region": cfg["REGION"],
             "account": identity["Account"], "bucket": infra["Bucket"], "stack": cfg["STACK_NAME"],
             "ami": ami, "created_at": now(), "started_epoch": time.time(),
             "max_hours": cfg["MAX_WORKER_HOURS"], "boot_minutes": cfg["BOOT_TIMEOUT_MINUTES"],
             "fail_fast": args.fail_fast, "scenarios": scenarios, "workers": {},
             "bootstrap_sha256": sha256(bootstrap)}
    safe_cfg = {k: v for k, v in cfg.items() if k not in {"ENV", "VPN_CONFIG_FILE", "VPN_CONFIG_FILES", "VPN_AUTH_FILE", "KUBECTL_FILE", "KUBECONFIG_FILE"}}
    safe_cfg["AWS_S3_CIDRS"] = aws_s3_cidrs
    for index, assigned in enumerate(assignments):
        name = f"worker-{index + 1:03d}"
        batch_id = f"{batch}-w{index + 1:03d}"
        item = {"name": name, "batch_id": batch_id, "scenarios": assigned, "instance_id": None}
        state["workers"][name] = item
        worker_config = {"config": safe_cfg, "environment": env, "batch": batch, "suite": args.suite,
                         "worker": name, "batch_id": batch_id, "scenarios": assigned, "fail_fast": args.fail_fast,
                         "scenario_root": str(args.scenario_root.relative_to(args.lab_dir)), "bucket": infra["Bucket"],
                         "source_sha256": source_digest, "image_key": image_key, "image_info": image_info}
        atomic_json(input_dir / (name + ".json"), worker_config)
        if vpns:
            profile_path = input_dir / (name + ".ovpn")
            profile_path.write_text(vpns[0] if len(vpns) == 1 else vpns[index])
            profile_path.chmod(0o600)
    # Only new batches may use this prefix. Do not overwrite another controller's inputs.
    if aws.head(infra["Bucket"], f"inputs/{batch}/bootstrap.tar.gz"):
        raise RunnerError("this AWS batch ID already exists; choose another --batch-id or use --resume")
    aws.sync(input_dir, f"s3://{infra['Bucket']}/inputs/{batch}")
    atomic_json(work / "state.json", state)  # Saved before the first EC2 launch.
    lifecycle(work, "AWS", f"region={cfg['REGION']} stack={cfg['STACK_NAME']} instance_type={cfg['INSTANCE_TYPE']} disk={cfg['DISK_GB']} GiB AMI={ami}")
    lifecycle(work, "AWS", f"workers requested={args.workers}, launching={len(assignments)} for {len(scenarios)} selected scenarios; no SSH required")
    if len(assignments) > 1:
        print("Starting worker-001 first; remaining workers launch after its infrastructure is ready", flush=True)
    for name, item in state["workers"].items():
        script = user_data(state, name)
        if len(script.encode()) > 16 * 1024:
            raise RunnerError("EC2 user-data exceeds 16 KiB")
        tags = [{"Key": "Name", "Value": f"ukama-lab-{batch}-{name}"},
                {"Key": "UkamaLabBatch", "Value": batch}, {"Key": "UkamaLabWorker", "Value": name},
                {"Key": "UkamaLabStack", "Value": cfg["STACK_NAME"]}]
        request = {"ImageId": ami, "InstanceType": cfg["INSTANCE_TYPE"], "MinCount": 1, "MaxCount": 1,
                   "ClientToken": hashlib.sha256((state["account"] + batch + name).encode()).hexdigest(),
                   "IamInstanceProfile": {"Name": infra["InstanceProfile"]},
                   "NetworkInterfaces": [{"DeviceIndex": 0, "SubnetId": infra["Subnet"],
                                          "Groups": [infra["SecurityGroup"]], "AssociatePublicIpAddress": True}],
                   "BlockDeviceMappings": [{"DeviceName": root_device, "Ebs": {"VolumeType": "gp3", "VolumeSize": cfg["DISK_GB"], "Encrypted": True, "DeleteOnTermination": True}}],
                   "MetadataOptions": {"HttpTokens": "required", "HttpEndpoint": "enabled", "HttpPutResponseHopLimit": 1},
                   "InstanceInitiatedShutdownBehavior": "terminate",
                   "TagSpecifications": [{"ResourceType": "instance", "Tags": tags}, {"ResourceType": "volume", "Tags": tags}],
                   "UserData": base64.b64encode(script.encode()).decode()}
        request_file = work / (name + "-launch.json")
        atomic_json(request_file, request)
        # Retry the SAME client token/request for IAM propagation or transient launch errors.
        result = None
        for attempt in range(12):
            p = aws.call("ec2", "run-instances", "--cli-input-json", "file://" + str(request_file), check=False)
            if p.returncode == 0:
                result = json.loads(p.stdout)
                break
            if not any(code in p.stderr for code in ("Invalid IAM Instance Profile", "InvalidParameterValue", "RequestLimitExceeded", "ServiceUnavailable")):
                raise RunnerError(f"EC2 launch failed: {p.stderr.strip()}; use --resume or --cleanup {work}")
            time.sleep(5)
        if result is None:
            raise RunnerError(f"EC2 launch did not succeed; state retained at {work}")
        item["instance_id"] = result["Instances"][0]["InstanceId"]
        item["launched_epoch"] = time.time()
        atomic_json(work / "state.json", state)
        lifecycle(work, name, f"EC2 {item['instance_id']} created; {len(item['scenarios'])} scenarios assigned")
        if "State" in result["Instances"][0]:
            observe_instance(work, name, item, result["Instances"][0])
        if name == "worker-001" and len(assignments) > 1:
            if not wait_for_infrastructure(aws, work, state, name):
                for pending in state["workers"].values():
                    if not pending["instance_id"]:
                        pending["not_started_reason"] = "not launched: first worker infrastructure did not become ready"
                atomic_json(work / "state.json", state)
                print("Remaining workers were not launched; collecting the first worker's diagnostics", flush=True)
                break
    # Runtime credentials need not remain in local staging after every worker has an input.
    shutil.rmtree(input_dir)
    try:
        return monitor(work, state)
    finally:
        batch_lock.close()


def instances(aws, state):
    data = aws.json("ec2", "describe-instances", "--filters",
                    f"Name=tag:UkamaLabBatch,Values={state['batch']}",
                    f"Name=tag:UkamaLabStack,Values={state['stack']}")
    result = {}
    for reservation in data["Reservations"]:
        for instance in reservation["Instances"]:
            tags = {t["Key"]: t["Value"] for t in instance.get("Tags", [])}
            name = tags.get("UkamaLabWorker")
            if name in state["workers"]:
                result[name] = instance
    return result


def verify_account(aws, state):
    if aws.json("sts", "get-caller-identity")["Account"] != state["account"]:
        raise RunnerError("AWS account differs from this batch's account")


def collect(aws, work, state):
    target = work / "workers"
    target.mkdir(exist_ok=True)
    lifecycle(work, "AWS", "downloading worker logs and reports from S3")
    aws.sync(f"s3://{state['bucket']}/results/{state['batch']}", target)
    lifecycle(work, "AWS", f"worker logs and reports collected in {target}")
    return aggregate(work, state)


def wait_for_infrastructure(aws, work, state, name):
    """Admit the rest of the pool after real EC2 setup, never scenario PASS."""
    item = state["workers"][name]
    prefix = f"results/{state['batch']}/{name}"
    last_console, last_message = 0.0, ""
    while True:
        marker = aws.optional_json(state["bucket"], prefix + "/complete.json")
        status = aws.optional_json(state["bucket"], prefix + "/status.json") or {}
        ready = aws.optional_json(state["bucket"], prefix + "/ready.json")
        failure = next((record for record in (marker or {}, status)
                        if record.get("phase") in ("BOOTSTRAP_ERROR", "INFRA_ERROR")), None)
        if failure:
            print(f"{name}: {failure.get('error') or failure['phase']}", flush=True)
            return False
        if ready and ready.get("worker") == name and ready.get("batch") == state["batch"]:
            item["infrastructure_ready_at"] = ready["ready_at"]
            atomic_json(work / "state.json", state)
            lifecycle(work, name, "infrastructure ready; launching remaining workers")
            return True
        if marker:
            item["bootstrap_error"] = "worker completed without an infrastructure readiness marker"
            return False
        current = instances(aws, state).get(name)
        observe_instance(work, name, item, current)
        observe_worker(work, name, item, status)
        if current and current["State"]["Name"] in ("terminated", "shutting-down", "stopped"):
            item["bootstrap_error"] = "first worker stopped before infrastructure readiness"
            return False
        age = time.time() - item["launched_epoch"]
        if age >= state["boot_minutes"] * 60:
            item["bootstrap_error"] = "first worker exceeded infrastructure startup deadline"
            return False
        message = status.get("stage") or status.get("phase") or "bootstrapping"
        if not status and age >= 180 and time.time() - last_console >= 120:
            diagnostic = console_diagnostics(aws, work, name, item["instance_id"])
            last_console = time.time()
            message = diagnostic["stage"] or message
            if diagnostic["error"]:
                item["bootstrap_error"] = diagnostic["error"]
                return False
        if message != last_message:
            if not status:
                lifecycle(work, name, f"{message}; waiting for infrastructure readiness")
            last_message = message
        time.sleep(20)


def aggregate(work, state):
    results, workers, accounted = [], {}, set()
    for name, worker in state["workers"].items():
        directory = work / "workers" / name
        marker_path = directory / "complete.json"
        marker = read_json(marker_path) if marker_path.exists() else {}
        fallback_phase = "NOT_STARTED" if worker.get("not_started_reason") else ("BOOTSTRAP_ERROR" if worker.get("bootstrap_error") else "INCOMPLETE")
        workers[name] = {"instance_id": worker.get("instance_id"), "phase": marker.get("phase", fallback_phase),
                         "lab_exit": marker.get("lab_exit"), "error": marker.get("error", worker.get("not_started_reason") or worker.get("bootstrap_error", ""))}
        path = directory / "batch" / worker["batch_id"] / "batch-report.json"
        if path.exists():
            report = read_json(path)
            source_results = report.get("results", [])
        else:
            # On interruption the loop may have uploaded scenarios.tsv before the
            # final batch report existed. Read the same original report/exit-code rules.
            tsv = directory / "batch" / worker["batch_id"] / "scenarios.tsv"
            source_results = []
            if tsv.exists():
                with tsv.open() as stream:
                    for row in csv.DictReader(stream, delimiter="\t"):
                        source_results.append({**row, "exit_code": int(row["exit_code"])})
        for entry in source_results:
            row = dict(entry)
            for key in ("report", "log"):
                remote = Path(row.get(key, ""))
                try:
                    relative = remote.relative_to(ROOT / "results")
                except ValueError:
                    raise RunnerError(f"invalid {key} path in {name} report") from None
                row[key] = str((directory / relative).resolve())
            if "outcome" not in row:
                report_file = Path(row["report"])
                original = read_json(report_file) if report_file.is_file() else {}
                row["outcome"] = "SKIP" if row["exit_code"] == 0 and original.get("status") in ("skip", "wip") else (
                    "PASS" if row["exit_code"] == 0 and original.get("passed") is True else "FAIL")
            selected = next((p for p in worker["scenarios"] if p.endswith("/" + row["scenario"]) or p == row["scenario"]), None)
            if selected is None or selected in accounted:
                raise RunnerError(f"unexpected or duplicate scenario report from {name}")
            accounted.add(selected)
            results.append({**row, "worker": name, "selected_path": selected})
    pending = [p for p in state["scenarios"] if p not in accounted]
    counts = {field: sum(r["outcome"] == key for r in results)
              for field, key in (("passed", "PASS"), ("failed", "FAIL"), ("skipped", "SKIP"))}
    incomplete_workers = [k for k, v in workers.items() if v["phase"] not in ("DONE", "STOPPED")]
    summary = {"batch": state["batch"], "selected": len(state["scenarios"]), "completed": len(results), "total": len(results),
               **counts, "infrastructure_errors": incomplete_workers, "unfinished": pending,
               "workers": workers, "results": results}
    atomic_json(work / "batch-report.json", summary)
    lines = [f"Ukama {state['suite']} AWS batch {state['batch']}",
             f"selected={summary['selected']} complete={len(results)} pass={counts['passed']} fail={counts['failed']} skip={counts['skipped']} unfinished={len(pending)}",
             f"workers with infrastructure errors/incomplete collection={len(incomplete_workers)}", ""]
    infrastructure_reasons = {}
    for name in incomplete_workers:
        reason = workers[name]["error"] or workers[name]["phase"]
        infrastructure_reasons.setdefault(reason, []).append(name)
    for reason, names in infrastructure_reasons.items():
        lines.append(f"Infrastructure error ({len(names)} workers): {reason}")
    if infrastructure_reasons:
        lines.append("")
    lines.extend(f"{r['outcome']:<4} {r['worker']} {r['scenario']}" for r in results)
    if pending:
        lines += ["", "Unfinished assignments:", *pending]
    (work / "batch-report.txt").write_text("\n".join(lines) + "\n")
    return summary


def console_diagnostics(aws, work, name, iid):
    """Collect startup evidence even when the worker cannot install/use AWS CLI."""
    path = work / (name + "-console.json")
    try:
        result = aws.call("ec2", "get-console-output", "--instance-id", iid,
                          "--latest", check=False, timeout=90)
        if result.returncode:
            (work / (name + "-console-error.txt")).write_text(result.stderr or result.stdout)
            return {"error": "", "stage": "", "available": False}
        data = json.loads(result.stdout)
        output = data.get("Output") or ""
        if not output:
            return {"error": "", "stage": "", "available": False}
        # The AWS CLI has already decoded the API's base64 Output field.
        path.write_text(result.stdout)
        (work / (name + "-console.log")).write_text(output)
    except (RunnerError, ValueError, OSError, subprocess.TimeoutExpired) as exc:
        print(f"\n{name}: console collection unavailable: {exc}", flush=True)
        return {"error": "", "stage": "", "available": False}
    stages = re.findall(r"UKAMA_RUNNER_BOOTSTRAP stage=([a-z-]+)", output)
    failures = re.findall(r"UKAMA_RUNNER_BOOTSTRAP_ERROR stage=([a-z-]+) exit=(\d+)", output)
    error = ""
    if failures:
        stage, rc = failures[-1]
        error = f"bootstrap stage {stage} exited {rc}"
    elif "Failed to run module scripts_user" in output:
        # Also diagnose batches launched by the original bootstrap.
        reasons = [line for line in output.splitlines() if line.startswith("E: ")]
        error = reasons[-1] if reasons else "cloud-init failed to run the worker startup script"
    return {"error": error, "stage": stages[-1] if stages else "",
            "available": True}


def monitor(work, state):
    aws = AWS(state["region"])
    verify_account(aws, state)
    print(f"Monitoring batch. Ctrl-C detaches; reattach with --mode aws --resume {work}", flush=True)
    done = set()
    warned = set()
    console_checked = {}
    boot_messages = {}
    previous_line = ""
    last_print = 0.0
    while len(done) < len(state["workers"]):
        current = instances(aws, state)
        overall = {"completed": 0, "passed": 0, "failed": 0, "skipped": 0}
        for name, worker in state["workers"].items():
            instance = current.get(name)
            if instance:
                worker["instance_id"] = instance["InstanceId"]
                observe_instance(work, name, worker, instance)
            elif worker.get("not_started_reason"):
                done.add(name)
                continue
            prefix = f"results/{state['batch']}/{name}"
            marker = aws.optional_json(state["bucket"], prefix + "/complete.json")
            status = aws.optional_json(state["bucket"], prefix + "/status.json") or {}
            observe_worker(work, name, worker, status)
            for key in overall:
                overall[key] += int(status.get("progress", {}).get(key, 0))
            if marker:
                if name not in done:
                    lifecycle(work, name, f"worker finished ({marker.get('phase', 'complete')}); available results uploaded to S3")
                done.add(name)
                if instance and instance["State"]["Name"] in ("pending", "running", "stopping", "stopped"):
                    request_termination(aws, work, name, worker, instance, "results uploaded")
            elif name in done:
                continue
            elif status.get("phase") == "BOOTSTRAP_ERROR":
                diagnostic = {"available": False}
                if instance:
                    diagnostic = console_diagnostics(aws, work, name, instance["InstanceId"])
                worker["bootstrap_error"] = status.get("error") or "EC2 bootstrap failed"
                atomic_json(work / "state.json", state)
                detail = f"; console saved to {work / (name + '-console.log')}" if diagnostic["available"] else "; see worker bootstrap.log in S3"
                print(f"\n{name}: {worker['bootstrap_error']}{detail}", flush=True)
                if instance and instance["State"]["Name"] in ("pending", "running", "stopping", "stopped"):
                    request_termination(aws, work, name, worker, instance, "bootstrap failed")
                done.add(name)
            elif ((not instance and time.time() - worker.get("launched_epoch", state["started_epoch"]) > 180)
                  or (instance and instance["State"]["Name"] in ("terminated", "shutting-down", "stopped"))):
                # Missing workers are unfinished; never manufacture scenario FAIL rows.
                done.add(name)
            elif time.time() - float(status.get("heartbeat_epoch", worker.get("launched_epoch", state["started_epoch"]))) > 180:
                initial = not status.get("heartbeat_epoch")
                if instance and time.time() - console_checked.get(name, 0) >= 120:
                    diagnostic = console_diagnostics(aws, work, name, instance["InstanceId"])
                    console_checked[name] = time.time()
                    if initial and diagnostic["error"]:
                        worker["bootstrap_error"] = diagnostic["error"]
                        atomic_json(work / "state.json", state)
                        print(f"\n{name}: {diagnostic['error']}; console saved to {work / (name + '-console.log')}; stopping failed worker", flush=True)
                        request_termination(aws, work, name, worker, instance, "bootstrap failed")
                        done.add(name)
                    elif initial:
                        stage = diagnostic["stage"] or "startup"
                        message = f"{stage}; waiting for first worker heartbeat"
                        if boot_messages.get(name) != message:
                            print(f"\n{name}: {message}" + (f"; console saved to {work / (name + '-console.log')}" if diagnostic["available"] else "; console not yet available"), flush=True)
                            boot_messages[name] = message
                if not initial and name not in warned:
                    print(f"\n{name}: heartbeat stale in {status.get('phase', 'unknown')}; retaining instance", flush=True)
                    warned.add(name)
            else:
                warned.discard(name)
            if instance and not marker and name not in done:
                limit = state["boot_minutes"] * 60 if status.get("phase", "BOOTSTRAP") in ("BOOTSTRAP", "BOOTSTRAP_ERROR", "PREPARING") else state["max_hours"] * 3600
                if time.time() - worker.get("launched_epoch", state["started_epoch"]) > limit:
                    print(f"\n{name}: worker deadline reached; stopping instance", flush=True)
                    request_termination(aws, work, name, worker, instance, "worker deadline reached")
                    done.add(name)
        line = f"Overall: {overall['completed']}/{len(state['scenarios'])} complete | pass={overall['passed']} fail={overall['failed']} skip={overall['skipped']} | workers finished={len(done)}/{len(state['workers'])}"
        if sys.stdout.isatty():
            print("\r" + line.ljust(120), end="", flush=True)
        elif line != previous_line or time.monotonic() - last_print >= 120:
            print(line, flush=True)
            last_print = time.monotonic()
        previous_line = line
        atomic_json(work / "state.json", state)
        if len(done) < len(state["workers"]):
            time.sleep(20)
    if sys.stdout.isatty():
        print()
    summary = collect(aws, work, state)
    # Preserve EC2 console output for workers that never reached S3 publication.
    for name in summary["infrastructure_errors"]:
        iid = state["workers"][name].get("instance_id")
        if iid:
            console_diagnostics(aws, work, name, iid)
    confirm_termination(aws, work, state)
    print_report(work / "batch-report.txt")
    print(f"Reports: {work / 'batch-report.json'}", flush=True)
    return 2 if summary["infrastructure_errors"] or (summary["unfinished"] and not state["fail_fast"]) else (1 if summary["failed"] or summary["unfinished"] else 0)


def cleanup(work, state):
    aws = AWS(state["region"])
    verify_account(aws, state)
    # Preserve everything already uploaded before explicitly terminating active workers.
    current = instances(aws, state)
    for name, instance in current.items():
        diagnostic = console_diagnostics(aws, work, name, instance["InstanceId"])
        if diagnostic["error"]:
            state["workers"][name]["bootstrap_error"] = diagnostic["error"]
    atomic_json(work / "state.json", state)
    collect(aws, work, state)
    active = [v["InstanceId"] for v in current.values() if v["State"]["Name"] != "terminated"]
    if active:
        aws.call("ec2", "terminate-instances", "--instance-ids", *active)
        for name, instance in current.items():
            if instance["InstanceId"] in active:
                state["workers"][name]["termination_requested_at"] = now()
                lifecycle(work, name, f"EC2 {instance['InstanceId']} termination requested (explicit cleanup)")
        confirm_termination(aws, work, state)
    print(f"Requested termination of {len(active)} batch instances. Available results: {work}")
    return 0


def main():
    parser = argparse.ArgumentParser()
    sub = parser.add_subparsers(dest="action", required=True)
    start = sub.add_parser("launch")
    for name in ("suite", "config", "lab-dir", "repo", "lab-bin", "scenario-root", "out", "factory-nodes"):
        start.add_argument("--" + name, required=True)
    start.add_argument("--batch-id", default="")
    start.add_argument("--workers", type=int, default=1)
    start.add_argument("--fail-fast", action="store_true")
    start.add_argument("--prepare-only", action="store_true")
    for action in ("resume", "cleanup"):
        sub.add_parser(action).add_argument("--batch-dir", required=True)
    args = parser.parse_args()
    if args.action == "launch":
        return launch(args)
    work = Path(args.batch_dir).expanduser().resolve()
    state = read_json(work / "state.json")
    with (work / ".controller.lock").open("w") as lock:
        try:
            fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
        except BlockingIOError:
            raise RunnerError("another controller is already attached to this batch") from None
        return monitor(work, state) if args.action == "resume" else cleanup(work, state)


if __name__ == "__main__":
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        print("\nDetached. Workers continue within their configured lifetime; use --resume or --cleanup with the printed batch directory.", file=sys.stderr)
        sys.exit(130)
    except (RunnerError, OSError, ValueError) as exc:
        print(f"runner: {exc}", file=sys.stderr)
        sys.exit(2)
