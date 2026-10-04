#!/usr/bin/env python3
import ipaddress
import os
from pathlib import Path
import re
import shlex
import shutil
from common import RunnerError, read_json

HERE = Path(__file__).resolve().parent
ROOT = Path("/opt/ukama-runner")


def load_config(filename):
    path = Path(filename).expanduser().resolve()
    if not path.is_file():
        raise RunnerError(f"AWS config missing: {path}; copy utils/runner/aws.example.json to aws.json and set VPN_CONFIG_FILE")
    defaults = read_json(HERE / "aws.example.json")
    values = read_json(path)
    if not isinstance(values, dict):
        raise RunnerError("AWS configuration must be a JSON object")
    # Optional explicit paths; otherwise use the caller's working kubectl setup.
    defaults.update({"KUBECTL_FILE": "", "KUBECONFIG_FILE": ""})
    if set(values) - set(defaults):
        raise RunnerError(f"unknown AWS configuration keys: {sorted(set(values) - set(defaults))}")
    cfg = {**defaults, **values}
    if not re.fullmatch(r"[a-z]{2}(?:-[a-z]+)+-\d", cfg["REGION"]):
        raise RunnerError("invalid REGION")
    if not re.fullmatch(r"[A-Za-z][A-Za-z0-9-]{0,90}", cfg["STACK_NAME"]):
        raise RunnerError("invalid STACK_NAME")
    for key, minimum, maximum in [("DISK_GB", 40, 16000), ("BOOT_TIMEOUT_MINUTES", 5, 180),
                                  ("MAX_WORKER_HOURS", 1, 36), ("UPLOAD_RETRY_MINUTES", 1, 120),
                                  ("MIN_FREE_DISK_GB", 1, 1000), ("RESULT_RETENTION_DAYS", 1, 3650)]:
        if type(cfg[key]) is not int or not minimum <= cfg[key] <= maximum:
            raise RunnerError(f"{key} must be an integer between {minimum} and {maximum}")
    networks = {key: ipaddress.IPv4Network(cfg[key]) for key in ("VPC_CIDR", "SUBNET_CIDR", "PODMAN_CIDR")}
    if not networks["SUBNET_CIDR"].subnet_of(networks["VPC_CIDR"]):
        raise RunnerError("SUBNET_CIDR must be within VPC_CIDR")
    if networks["VPC_CIDR"].overlaps(networks["PODMAN_CIDR"]) or networks["PODMAN_CIDR"].prefixlen > 23:
        raise RunnerError("PODMAN_CIDR must not overlap VPC_CIDR and must contain multiple /24 networks")
    if cfg["NETWORK_MODE"] not in ("vpn", "direct"):
        raise RunnerError("NETWORK_MODE must be vpn or direct")
    for key in ("VPN_CONFIG_FILES", "VPN_DOMAINS", "VPN_DNS_SERVERS", "VPN_ROUTES", "EXTRA_APT_PACKAGES", "EXTRA_PROBE_URLS", "REPO_PATHS"):
        if not isinstance(cfg[key], list) or not all(isinstance(v, str) for v in cfg[key]):
            raise RunnerError(f"{key} must be an array of strings")
    for domain in cfg["VPN_DOMAINS"]:
        if not re.fullmatch(r"[A-Za-z0-9](?:[A-Za-z0-9.-]*[A-Za-z0-9])?", domain):
            raise RunnerError("invalid VPN_DOMAINS entry")
    if cfg["NETWORK_MODE"] == "vpn" and not cfg["VPN_DOMAINS"]:
        raise RunnerError("VPN_DOMAINS must identify backend DNS domains")
    for address in cfg["VPN_DNS_SERVERS"]:
        ipaddress.IPv4Address(address)
    for cidr in cfg["VPN_ROUTES"]:
        route = ipaddress.IPv4Network(cidr)
        if route.prefixlen == 0 or any(route.overlaps(n) for n in networks.values()):
            raise RunnerError(f"VPN_ROUTES overlaps worker networking or contains a default: {cidr}")
    for package in cfg["EXTRA_APT_PACKAGES"]:
        if not re.fullmatch(r"[a-z0-9][a-z0-9+.-]*(?:=[A-Za-z0-9.+:~\-]+)?", package):
            raise RunnerError("invalid EXTRA_APT_PACKAGES entry")
    for value in cfg["REPO_PATHS"]:
        if Path(value).is_absolute() or ".." in Path(value).parts:
            raise RunnerError("REPO_PATHS must be relative paths inside UKAMA_REPO")
    for required in ("testing/node", "testing/ue"):
        if not any(Path(required).is_relative_to(Path(p)) for p in cfg["REPO_PATHS"]):
            raise RunnerError(f"REPO_PATHS must include {required}")
    if not isinstance(cfg["ENV"], dict) or not all(isinstance(k, str) and isinstance(v, (str, int, float)) for k, v in cfg["ENV"].items()):
        raise RunnerError("ENV must map variable names to strings or numbers")
    for key in ("VPN_CONFIG_FILE", "VPN_AUTH_FILE", "KUBECTL_FILE", "KUBECONFIG_FILE"):
        if cfg[key]:
            local = Path(os.path.expandvars(cfg[key])).expanduser()
            cfg[key] = str((path.parent / local).resolve())
    cfg["VPN_CONFIG_FILES"] = [str((path.parent / Path(os.path.expandvars(p)).expanduser()).resolve()) for p in cfg["VPN_CONFIG_FILES"]]
    return cfg


def environment(cfg):
    allowed = {"UKAMA_IDENTIFIER", "UKAMA_PASSWORD", "PAUTH_URL", "BFF_BASE_URL", "UKAMA_APP_VERSION",
               "NODE_RUNTIME", "BASE_IMAGE_REPO", "IMAGE_REPO", "FACTORY_ORG"}
    def accept(k):
        return (k in allowed or k.startswith(("ULAB_", "UKAMA_LAB_"))) and not k.startswith("ULAB_RUNNER_")
    if any(not accept(k) for k in cfg["ENV"]):
        raise RunnerError("ENV supports scenario ULAB_*/UKAMA_LAB_* settings and Ukama authentication/image settings only")
    result = {k: v for k, v in os.environ.items() if accept(k)}
    result.update({k: str(v) for k, v in cfg["ENV"].items()})
    for key in ("ULAB_KUBECTL", "ULAB_FACTORY_NODE_COUNT", "ULAB_RESILIENCE_NAME_SUFFIX"):
        result.pop(key, None)
    # payload() replaces the local kubectl path with the packaged worker path.
    # Each worker creates its own suffix. Never propagate host paths for tools.
    result["UKAMA_REPO"] = str(ROOT / "repo")
    result["ULAB_FACTORY_NODE_COUNT"] = "0"
    result["BASE_IMAGE_REPO"] = result.get("BASE_IMAGE_REPO", "localhost/testing/virtualnode-base")
    if not result.get("UKAMA_IDENTIFIER") or not result.get("UKAMA_PASSWORD"):
        raise RunnerError("export UKAMA_IDENTIFIER and UKAMA_PASSWORD before AWS execution")
    return result


def executable(value):
    result = shutil.which(value) or (str(Path(value).expanduser().resolve()) if value else "")
    if not result or not os.access(result, os.X_OK):
        raise RunnerError(f"executable not found: {value}")
    return Path(result)


def portable_vpn(filename, auth_file=""):
    """Inline certificates and replace laptop hooks; retain the profile's routing."""
    path = Path(filename)
    if not path.is_file():
        raise RunnerError(f"VPN profile missing: {path}")
    output, inline, inline_tags = [], None, set()
    stripped = {"up", "down", "route-up", "route-pre-down", "ipchange", "script-security",
                "daemon", "log", "log-append", "writepid", "status", "management", "user", "group",
                "dev", "dev-type", "persist-tun", "block-outside-dns"}
    file_options = {"ca", "cert", "key", "tls-auth", "tls-crypt", "tls-crypt-v2", "extra-certs", "auth-user-pass"}
    for raw in path.read_text().splitlines():
        line = raw.strip()
        if inline:
            output.append(raw)
            if line == f"</{inline}>":
                inline = None
            continue
        if line.startswith("<") and line.endswith(">") and not line.startswith("</"):
            inline = line[1:-1]
            if inline == "connection":
                raise RunnerError("VPN <connection> blocks need a flat worker profile")
            inline_tags.add(inline)
            output.append(raw)
            continue
        words = shlex.split(raw, comments=True)
        if not words:
            continue
        key = words[0].lstrip("-")
        if key in {"config", "plugin", "askpass", "pkcs12", "pkcs11-providers", "cd", "chroot", "route-noexec", "route-nopull"}:
            raise RunnerError(f"VPN option {key} needs a portable noninteractive worker profile")
        if key in stripped:
            continue
        if key in file_options:
            if len(words) > 1 and words[1] == "[inline]":
                continue
            source = Path(auth_file) if key == "auth-user-pass" and auth_file else (path.parent / words[1] if len(words) > 1 else None)
            if source is None:
                raise RunnerError("VPN needs VPN_AUTH_FILE for noninteractive auth-user-pass")
            if not source.is_file():
                raise RunnerError(f"VPN {key} file does not exist")
            output += [f"<{key}>", source.read_text().strip(), f"</{key}>"]
            inline_tags.add(key)
            if key == "tls-auth" and len(words) > 2:
                output.append("key-direction " + words[2])
        else:
            output.append(raw)
    if inline:
        raise RunnerError("unterminated inline block in VPN profile")
    if auth_file and "auth-user-pass" not in inline_tags:
        output += ["<auth-user-pass>", Path(auth_file).read_text().strip(), "</auth-user-pass>"]
    # Preserve IPv4 full-tunnel and split-tunnel behavior. Worker networking pins
    # S3, VPC and metadata routes to AWS before OpenVPN applies its redirect.
    return '\n'.join(['pull-filter ignore "block-outside-dns"',
                      'pull-filter ignore "route-ipv6"', 'pull-filter ignore "ifconfig-ipv6"',
                      *output, "dev ulab-vpn", "dev-type tun", "auth-retry nointeract", "verb 3", ""]) 
