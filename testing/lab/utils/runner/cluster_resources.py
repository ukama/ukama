# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2026-present, Ukama Inc.
"""Container limits, observed counter changes, and node pressure."""
import math


LABELS = {
    "cpu_mcpu": "CPU usage (mCPU, 5m rate)",
    "cpu_request_mcpu": "CPU request (mCPU)",
    "cpu_limit_mcpu": "CPU limit (mCPU)",
    "cpu_request_pct": "CPU / request (%)",
    "cpu_limit_pct": "CPU / limit (%)",
    "cpu_throttle_pct": "CPU throttled periods (%)",
    "memory_working_mib": "Memory working set (MiB)",
    "memory_usage_mib": "Memory total usage (MiB)",
    "memory_rss_mib": "Memory RSS (MiB)",
    "memory_request_mib": "Memory request (MiB)",
    "memory_limit_mib": "Memory limit (MiB)",
    "memory_request_pct": "Memory total / request (%)",
    "memory_limit_pct": "Memory total / limit (%)",
    "memory_working_pct": "Memory working set / limit (%)",
    "memory_headroom_mib": "Memory limit minus total usage (MiB)",
    "cpu_psi_pct": "CPU PSI some (%)",
    "memory_psi_pct": "Memory PSI some (%)",
    "memory_full_psi_pct": "Memory PSI full (%)",
}
NODE_LABELS = {
    "cpu_busy_pct": "CPU busy (%)",
    "memory_used_pct": "Memory unavailable (%)",
    "memory_available_mib": "Memory available (MiB)",
    "cpu_psi_pct": "CPU PSI some (%)",
    "memory_psi_pct": "Memory PSI some (%)",
    "memory_full_psi_pct": "Memory PSI full (%)",
    "memory_pressure": "Kubernetes MemoryPressure (0/1)",
}
TOTAL_LABELS = {
    "throttled_seconds": "Observed throttled seconds >=",
    "oom_events": "Observed OOM events >=",
    "oom_kills": "Observed OOMKilled terminations >=",
}
NOTES = [
    "Resource percentages are computed per container at matching timestamps. "
    "Service resource percentages/PSI show the worst container at each "
    "sample; headroom shows the smallest container headroom. Other service "
    "resource gauges sum observed containers. Missing coverage can hide load.",
    "Missing or zero requests/limits give N/A percentages, not zero usage. "
    "Absent limits are not proof of unlimited resources. Pod/ancestor limits "
    "can also apply and are not inferred from container resource settings.",
    "CPU throttled periods (%) counts periods with any quota throttling; "
    "it is not percent wall time lost. Throttled seconds may accumulate "
    "across CPUs and are not converted to a wall-time percentage.",
    "Memory working set and RSS are diagnostic views, not exact OOM "
    "thresholds. Total memory includes reclaimable cache; low headroom "
    "alone does not establish an OOM. PSI measures stalls, not limit hits.",
    "Counter totals include only observed adjacent samples inside each "
    "window, handle resets, and never count the initial counter value. "
    "OOMKilled counts deduplicate termination timestamps with matching "
    "reasons. Old reasons are not new kills. These counts are lower bounds.",
    "The 90% marker is a diagnostic threshold, not a scenario verdict. "
    "Sampled seconds require adjacent observations; gaps are not filled. "
    "Configuration ranges show min..max values observed during execution.",
    "Node pressure covers node-exporter targets and Kubernetes nodes at "
    "the selected Prometheus endpoint, including unrelated workloads. "
    "Exporter targets are identified by job/instance; nodenames use explicit "
    "node_uname_info labels. Node metrics are not AWS test-worker metrics.",
]


def queries(namespace):
    labels = namespace + ',pod!="",container!="",container!="POD"'
    identity = "namespace, pod, container"

    def gauge(metric, extra="", scale=""):
        return f"max by ({identity}) ({metric}{{{labels}{extra}}}){scale}"

    def rate(metric, scale=""):
        return (f"max by ({identity}) "
                f"(rate({metric}{{{labels}}}[5m])){scale}")

    values = {
        "cpu_mcpu": rate("container_cpu_usage_seconds_total", " * 1000"),
        "cpu_request_mcpu": gauge("kube_pod_container_resource_requests",
                                  ',resource="cpu"', " * 1000"),
        "cpu_limit_mcpu": gauge("kube_pod_container_resource_limits",
                                ',resource="cpu"', " * 1000"),
        "memory_request_mib": gauge("kube_pod_container_resource_requests",
                                    ',resource="memory"', " / 1048576"),
        "memory_limit_mib": gauge("kube_pod_container_resource_limits",
                                  ',resource="memory"', " / 1048576"),
        "memory_working_mib": gauge("container_memory_working_set_bytes",
                                    scale=" / 1048576"),
        "memory_usage_mib": gauge("container_memory_usage_bytes",
                                  scale=" / 1048576"),
        "memory_rss_mib": gauge("container_memory_rss",
                                scale=" / 1048576"),
        "cpu_periods_rate": rate("container_cpu_cfs_periods_total"),
        "cpu_throttled_periods_rate": rate(
            "container_cpu_cfs_throttled_periods_total"),
    }
    for name, metric in (
            ("cpu_psi_pct", "container_pressure_cpu_waiting_seconds_total"),
            ("memory_psi_pct",
             "container_pressure_memory_waiting_seconds_total"),
            ("memory_full_psi_pct",
             "container_pressure_memory_stalled_seconds_total")):
        values[name] = rate(metric, " * 100")
    for name, metric in (
            ("throttled_seconds_counter",
             "container_cpu_cfs_throttled_seconds_total"),
            ("oom_events_counter", "container_oom_events_total")):
        values[name] = (f"max by ({identity}, id) "
                        f"({metric}{{{labels}}})")
    values["terminated_timestamp"] = (
        f"max by ({identity}, uid) "
        "(kube_pod_container_status_last_terminated_timestamp{"
        + labels + "})")
    result = {"container_" + name: ("range", query)
              for name, query in values.items()}
    node = "job, instance"
    available = f"max by ({node}) (node_memory_MemAvailable_bytes)"
    total = f"max by ({node}) (node_memory_MemTotal_bytes)"
    values = {
        "cpu_busy_pct": "100 * (1 - avg by (job, instance) "
                        '(rate(node_cpu_seconds_total{mode="idle"}[5m])))',
        "memory_available_mib": available + " / 1048576",
        "memory_used_pct": f"100 * (1 - {available} / ({total} > 0))",
        "memory_pressure": "max by (node) (kube_node_status_condition"
                           '{condition="MemoryPressure",status="true"})',
        "identity": "max by (job, instance, nodename) (node_uname_info)",
    }
    for name, metric in (
            ("cpu_psi_pct", "node_pressure_cpu_waiting_seconds_total"),
            ("memory_psi_pct", "node_pressure_memory_waiting_seconds_total"),
            ("memory_full_psi_pct",
             "node_pressure_memory_stalled_seconds_total")):
        values[name] = f"100 * max by ({node}) (rate({metric}[5m]))"
    result.update({"node_" + name: ("range", query)
                   for name, query in values.items()})
    result["resource_pod_info"] = ("range", "max by (namespace, pod, node) "
                                   f"(kube_pod_info{{{namespace}}})")
    return result


def points(series):
    values = series.get("values", [])
    if "value" in series:
        values = [series["value"]]
    result = []
    for timestamp, value in values:
        try:
            timestamp, value = float(timestamp), float(value)
            if math.isfinite(timestamp) and math.isfinite(value):
                result.append((timestamp, value))
        except (ValueError, TypeError):
            continue
    return sorted(result)


def windows(raw):
    return {"before": (raw["start"] - raw["config"]["before"], raw["start"]),
            "during": (raw["start"], raw["end"]),
            "after": (raw["end"], raw["observed_end"])}


def in_window(t, name, bounds):
    lo, hi = bounds
    return lo <= t < hi if name == "before" else (
        lo < t <= hi if name == "after" else lo <= t <= hi)


def summarize(series, raw):
    result = {}
    for name, bounds in windows(raw).items():
        samples = [(t, v) for t, v in series
                   if in_window(t, name, bounds)]
        numbers = [v for _, v in samples]
        covered, high = 0, 0
        for (t, v), (next_t, _) in zip(samples, samples[1:]):
            if 0 < next_t - t <= raw["config"]["step"] * 1.5:
                covered += next_t - t
                if v >= 90:
                    high += next_t - t
        result[name] = {
            "mean": sum(numbers) / len(numbers) if numbers else None,
            "peak": max(numbers) if numbers else None,
            "min": min(numbers) if numbers else None,
            "last": numbers[-1] if numbers else None,
            "samples": len(numbers), "covered_seconds": covered,
            "above_90_seconds": high if covered else None,
        }
    return result


def record(identity, **metadata):
    return {"identity": identity, **metadata, "series": {}, "totals": {}}


def merge_values(target, name, values):
    merged = target["series"].setdefault(name, {})
    for t, v in values:
        merged[t] = max(v, merged.get(t, v))


def add_total(target, name, phase, delta):
    totals = target["totals"].setdefault(name, {})
    totals[phase] = totals.get(phase, 0) + delta


def counter_changes(series, raw):
    values = points(series)
    for name, bounds in windows(raw).items():
        for (t0, v0), (t1, v1) in zip(values, values[1:]):
            if (v0 < 0 or v1 < 0 or t1 <= t0 or t0 < bounds[0]
                    or t1 > bounds[1]):
                continue
            yield name, t1, v1 - v0 if v1 >= v0 else v1


def derived_values(container):
    series = container["series"]
    for name, numerator, denominator in (
            ("cpu_request_pct", "cpu_mcpu", "cpu_request_mcpu"),
            ("cpu_limit_pct", "cpu_mcpu", "cpu_limit_mcpu"),
            ("cpu_throttle_pct", "cpu_throttled_periods_rate",
             "cpu_periods_rate"),
            ("memory_request_pct", "memory_usage_mib", "memory_request_mib"),
            ("memory_limit_pct", "memory_usage_mib", "memory_limit_mib"),
            ("memory_working_pct", "memory_working_mib", "memory_limit_mib")):
        den = series.get(denominator, {})
        series[name] = {t: 100 * v / den[t]
                        for t, v in series.get(numerator, {}).items()
                        if den.get(t, 0) > 0}
    limit = series.get("memory_limit_mib", {})
    series["memory_headroom_mib"] = {
        t: limit[t] - v for t, v in series.get("memory_usage_mib", {}).items()
        if limit.get(t, 0) > 0}


def finish(target, raw):
    target["series"] = {name: sorted(values.items())
                        for name, values in target["series"].items()}
    target["statistics"] = {name: summarize(values, raw)
                            for name, values in target["series"].items()}
    for name in TOTAL_LABELS:
        values = target["totals"].setdefault(name, {})
        for phase in windows(raw):
            values.setdefault(phase, None)


def build(raw, responses, service_for):
    containers, events = {}, []
    pod_nodes = {}
    for item in responses.get("resource_pod_info", []):
        labels = item["metric"]
        if labels.get("node") and any(v == 1 for _, v in points(item)):
            pod_nodes.setdefault((labels.get("namespace"), labels.get("pod")),
                                 set()).add(labels["node"])

    def get_container(labels):
        if not labels.get("pod") or not labels.get("container"):
            return None
        service = service_for(labels)
        if service is None:
            return None
        key = tuple(labels.get(k, "unknown")
                    for k in ("namespace", "pod", "container"))
        identity = "/".join(key)
        if identity not in containers:
            containers[identity] = record(
                identity, service=service["service"], namespace=key[0],
                pod=key[1], container=key[2],
                nodes=sorted(pod_nodes.get(key[:2], [])))
        return containers[identity]

    # Keep UID and cgroup identity when deduplicating counter series.
    reasons = {}
    for item in responses.get("termination_reason", []):
        labels = item["metric"]
        key = tuple(labels.get(k, "")
                    for k in ("namespace", "pod", "container", "uid"))
        for t, v in points(item):
            if v == 1:
                reasons.setdefault((key, t), set()).add(labels.get("reason"))
    seen_kills = set()
    for query, rows in responses.items():
        if not query.startswith("container_"):
            continue
        metric = query.removeprefix("container_")
        # PromQL removes scrape replicas; also tolerate duplicate raw input.
        unique = {}
        for item in rows:
            key = tuple(item["metric"].get(k, "") for k in
                        ("namespace", "pod", "container", "id", "uid"))
            merged = unique.setdefault(key, {"metric": item["metric"],
                                             "values": {}})
            for t, v in points(item):
                merged["values"][t] = max(v, merged["values"].get(t, v))
        for item in unique.values():
            item["values"] = sorted(item["values"].items())
            container = get_container(item["metric"])
            if container is None:
                continue
            if metric.endswith("_counter"):
                name = metric.removesuffix("_counter")
                for phase, t, delta in counter_changes(item, raw):
                    add_total(container, name, phase, delta)
                    if name == "oom_events" and phase == "during" and delta:
                        events.append({"time": t,
                                       "service": container["service"],
                                       "event": "OOM events increased "
                                       f"{delta:g}"
                                       f" ({container['pod']}/"
                                       f"{container['container']})"})
            elif metric == "terminated_timestamp":
                labels = item["metric"]
                key = tuple(labels.get(k, "") for k in
                            ("namespace", "pod", "container", "uid"))
                for t, terminated in points(item):
                    reason = reasons.get((key, t), set())
                    if len(reason) != 1:
                        continue
                    for phase, bounds in windows(raw).items():
                        if in_window(t, phase, bounds):
                            add_total(container, "oom_kills", phase, 0)
                        claim = key, terminated
                        if (reason == {"OOMKilled"} and 0 < terminated <= t
                                and in_window(terminated, phase, bounds)
                                and claim not in seen_kills):
                            seen_kills.add(claim)
                            add_total(container, "oom_kills", phase, 1)
                            if phase == "during":
                                events.append({"time": terminated,
                                    "service": container["service"],
                                    "event": "OOMKilled "
                                    + container["pod"] + "/"
                                    + container["container"]})
            else:
                merge_values(container, metric, points(item))

    for container in containers.values():
        derived_values(container)
        service = service_for(container)
        resource = service.setdefault("resources", record(service["service"]))
        resource.setdefault("containers", []).append(container["identity"])
        worst = resource.setdefault("worst", {})
        for name, values in container["series"].items():
            if name not in LABELS:
                continue
            target = resource["series"].setdefault(name, {})
            for t, v in values.items():
                if name.endswith("_pct"):
                    target[t] = max(target.get(t, v), v)
                elif name == "memory_headroom_mib":
                    target[t] = min(target.get(t, v), v)
                else:
                    target[t] = target.get(t, 0) + v
                if (name.endswith("_pct") and raw["start"] <= t <= raw["end"]
                        and (name not in worst or v > worst[name]["value"])):
                    worst[name] = {"container": container["identity"],
                                   "value": v, "time": t}
        for name, totals in container["totals"].items():
            for phase, delta in totals.items():
                add_total(resource, name, phase, delta)
        finish(container, raw)
    services = {c["service"]: service_for(c) for c in containers.values()}
    for service in services.values():
        finish(service["resources"], raw)
    return ([containers[k] for k in sorted(containers)],
            build_nodes(raw, responses), events)


def build_nodes(raw, responses):
    names, nodes = {}, {}
    for item in responses.get("node_identity", []):
        labels = item["metric"]
        if labels.get("nodename") and any(v == 1 for _, v in points(item)):
            key = labels.get("job", ""), labels.get("instance", "")
            names.setdefault(key, set()).add(labels["nodename"])
    for query, rows in responses.items():
        if not query.startswith("node_") or query == "node_identity":
            continue
        for item in rows:
            labels = item["metric"]
            if query == "node_memory_pressure":
                identity = "kubernetes/" + labels.get("node", "unknown")
                name = labels.get("node", "unknown")
            else:
                key = labels.get("job", ""), labels.get("instance", "")
                identity = "exporter/" + "/".join(key)
                found = names.get(key, set())
                name = next(iter(found)) if len(found) == 1 else key[1]
            target = nodes.setdefault(identity, record(identity, node=name))
            merge_values(target, query.removeprefix("node_"), points(item))
    for target in nodes.values():
        finish(target, raw)
    return [nodes[k] for k in sorted(nodes)]
