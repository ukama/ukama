#!/usr/bin/env python3
"""Hold spare EKS nodes for the length of an AWS batch.

kube-system/capacity-placeholder runs low-priority pause pods, one per node.
Scaling it up before workers start makes the autoscaler add nodes ahead of
the mesh-pod burst; real pods preempt the placeholders when they need room.
"""
import os
from pathlib import Path
import subprocess

from common import RunnerError, atomic_json, now, run
from config import executable

NAMESPACE = "kube-system"
DEPLOYMENT = "capacity-placeholder"


def kubectl_command(cfg):
    value = (cfg.get("KUBECTL_FILE") or cfg.get("ENV", {}).get("ULAB_KUBECTL")
             or os.environ.get("ULAB_KUBECTL"))
    if not value:
        home_tool = Path.home() / "kubectl"
        value = str(home_tool) if home_tool.is_file() else "kubectl"
    command = [str(executable(os.path.expandvars(str(value))))]
    if cfg.get("KUBECONFIG_FILE"):
        command += ["--kubeconfig", cfg["KUBECONFIG_FILE"]]
    return command


def scale(command, replicas):
    """Return an error message, or an empty string when the scale succeeded."""
    try:
        result = run([*command, "-n", NAMESPACE, "scale", "deployment/" + DEPLOYMENT,
                      f"--replicas={replicas}"], check=False, timeout=60)
    except (OSError, subprocess.TimeoutExpired) as exc:
        return str(exc)
    if result.returncode:
        return (result.stderr or result.stdout or "").strip() or f"kubectl exited {result.returncode}"
    return ""


def scale_up(work, state, cfg, lifecycle):
    replicas = cfg["CAPACITY_PLACEHOLDER_REPLICAS"]
    if not replicas:
        return
    try:
        command = kubectl_command(cfg)
    except RunnerError as exc:
        lifecycle(work, "EKS", f"capacity-placeholder not scaled up: {exc}; continuing without spare nodes")
        return
    error = scale(command, replicas)
    if error:
        lifecycle(work, "EKS", f"capacity-placeholder not scaled up: {error}; continuing without spare nodes")
        return
    state["capacity_placeholder"] = {"kubectl": command, "replicas": replicas,
                                     "idle_replicas": cfg["CAPACITY_PLACEHOLDER_IDLE_REPLICAS"],
                                     "scaled_up_at": now()}
    atomic_json(work / "state.json", state)
    lifecycle(work, "EKS", f"capacity-placeholder scaled to {replicas}; the autoscaler adds spare nodes while workers boot")


def scale_down(work, state, lifecycle):
    placeholder = state.get("capacity_placeholder")
    if not placeholder or placeholder.get("scaled_down_at"):
        return
    replicas = placeholder["idle_replicas"]
    error = scale(placeholder["kubectl"], replicas)
    if error:
        lifecycle(work, "EKS", f"capacity-placeholder still at {placeholder['replicas']}: {error}; "
                               f"run kubectl -n {NAMESPACE} scale deploy {DEPLOYMENT} --replicas={replicas}")
        return
    placeholder["scaled_down_at"] = now()
    atomic_json(work / "state.json", state)
    lifecycle(work, "EKS", f"capacity-placeholder scaled back to {replicas}")
