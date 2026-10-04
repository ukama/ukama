#!/usr/bin/env python3
"""Worker-only networking. The laptop's VPN, DNS and Podman are never modified."""
import ipaddress
import json
import os
from pathlib import Path
import shlex
import shutil
import socket
import sys
import time
from urllib.parse import urlsplit
from common import RunnerError, atomic_json, read_json, run as command_run
from config import ROOT

READY = ROOT / "vpn.ready"
VPN_ERROR = ROOT / "vpn-error.json"
DNS_CONF = Path("/etc/ukama-runner-dns.conf")
SYSTEM_PATH = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"


def run(args, **kwargs):
    # OpenVPN hooks do not inherit the login shell's PATH. Resolve every worker
    # network tool against explicit system directories, including sbin.
    command = shutil.which(str(args[0]), path=SYSTEM_PATH)
    if command is None:
        raise RunnerError(f"required worker network command is missing: {args[0]}")
    env = {**os.environ, **kwargs.pop("env", {}), "PATH": SYSTEM_PATH}
    return command_run([command, *args[1:]], env=env, **kwargs)


def aws_routes(cfg, host):
    if not host.get("gateway"):
        raise RunnerError("worker has no AWS gateway")
    cidrs = cfg.get("AWS_S3_CIDRS", [])
    if cfg["NETWORK_MODE"] == "vpn" and not cidrs:
        raise RunnerError("worker input is missing AWS S3 routes; launch a new batch with the updated controller")
    for cidr in dict.fromkeys([cfg["VPC_CIDR"], "169.254.169.254/32", "169.254.169.253/32", *cidrs]):
        net = ipaddress.IPv4Network(cidr)
        if net.prefixlen < 2:
            raise RunnerError(f"invalid AWS bypass route: {cidr}")
        run(["ip", "route", "replace", str(net), "via", host["gateway"], "dev", host["device"]])


def host_network():
    route = json.loads(run(["ip", "-j", "route", "get", "1.1.1.1"]).stdout)[0]
    return {"device": route["dev"], "gateway": route.get("gateway"), "address": route["prefsrc"]}


def service(name, text):
    Path(f"/etc/systemd/system/{name}.service").write_text(text)
    run(["systemctl", "daemon-reload"])
    run(["systemctl", "restart", name])


def dns_config(cfg, host, servers):
    lines = ["no-resolv", "bind-interfaces", f"listen-address=127.0.0.1,{host['address']}",
             "server=169.254.169.253", "cache-size=0", "domain-needed"]
    if cfg["NETWORK_MODE"] == "vpn":
        for domain in cfg["VPN_DOMAINS"]:
            if servers:
                lines.extend(f"server=/{domain}/{address}" for address in servers)
            else:
                lines.append(f"local=/{domain}/")
    return "\n".join(lines) + "\n"


def firewall(cfg, host):
    def ensure(table, chain, *rule):
        base = ["iptables", "-w", "10", "-t", table]
        if run([*base, "-C", chain, *rule], check=False).returncode:
            run([*base, "-I", chain, "1", *rule])
    pool = cfg["PODMAN_CIDR"]
    for proto in ("tcp", "udp"):
        ensure("filter", "INPUT", "-s", pool, "-d", host["address"], "-p", proto, "--dport", "53", "-j", "ACCEPT")
    if cfg["NETWORK_MODE"] == "vpn":
        ensure("nat", "POSTROUTING", "-s", pool, "-o", "ulab-vpn", "-j", "MASQUERADE")
        ensure("filter", "FORWARD", "-s", pool, "-o", "ulab-vpn", "-j", "ACCEPT")
        ensure("filter", "FORWARD", "-d", pool, "-i", "ulab-vpn", "-m", "conntrack", "--ctstate", "ESTABLISHED,RELATED", "-j", "ACCEPT")
        ensure("mangle", "FORWARD", "-o", "ulab-vpn", "-p", "tcp", "--tcp-flags", "SYN,RST", "SYN", "-j", "TCPMSS", "--clamp-mss-to-pmtu")


def configure(cfg):
    host = host_network()
    atomic_json(ROOT / "host-network.json", host)
    aws_routes(cfg, host)
    run(["sysctl", "-w", "net.ipv4.ip_forward=1", "net.ipv4.conf.all.rp_filter=2", "net.ipv4.conf.default.rp_filter=2"])
    conf = Path("/etc/containers/containers.conf")
    conf.parent.mkdir(parents=True, exist_ok=True)
    pool = ipaddress.IPv4Network(cfg["PODMAN_CIDR"])
    default_subnet = next(pool.subnets(new_prefix=24))
    conf.write_text('[containers]\ndns_servers = ' + json.dumps([host["address"]]) + '\n'
                    '[engine]\ncgroup_manager = "systemd"\n'
                    '[network]\nnetwork_backend = "netavark"\n'
                    f'default_subnet = "{default_subnet}"\n'
                    f'default_subnet_pools = [{{base = "{pool}", size = 24}}]\n')
    DNS_CONF.write_text(dns_config(cfg, host, []))
    service("ukama-runner-dns", "[Unit]\nAfter=network-online.target\n[Service]\n"
            "ExecStart=/usr/sbin/dnsmasq --keep-in-foreground --conf-file=/etc/ukama-runner-dns.conf\n"
            "Restart=on-failure\nRestartSec=3\n")
    resolver = Path("/etc/resolv.conf")
    if resolver.is_symlink():
        resolver.unlink()
    resolver.write_text("nameserver 127.0.0.1\noptions timeout:2 attempts:2\n")
    firewall(cfg, host)
    if cfg["NETWORK_MODE"] == "vpn":
        READY.unlink(missing_ok=True)
        VPN_ERROR.unlink(missing_ok=True)
        script = Path(__file__).resolve()
        up = ROOT / "vpn-up.sh"
        down = ROOT / "vpn-down.sh"
        up.write_text(f"#!/bin/sh\nexport PATH={SYSTEM_PATH}\nexec /usr/bin/python3 {shlex.quote(str(script))} up\n")
        down.write_text(f"#!/bin/sh\nexport PATH={SYSTEM_PATH}\nexec /usr/bin/python3 {shlex.quote(str(script))} down\n")
        up.chmod(0o700)
        down.chmod(0o700)
        service("ukama-runner-vpn", "[Unit]\nAfter=network-online.target ukama-runner-dns.service\n"
                "[Service]\nType=simple\n"
                f"ExecStart=/usr/sbin/openvpn --config {ROOT}/client.ovpn --script-security 2 "
                f"--route-up {up} --down {down} --down-pre --up-restart\n"
                "Restart=always\nRestartSec=10\n")
        deadline = time.monotonic() + 180
        while not READY.exists():
            if VPN_ERROR.exists():
                raise RunnerError("VPN route-up failed: " + read_json(VPN_ERROR)["error"])
            if time.monotonic() >= deadline:
                run(["journalctl", "-u", "ukama-runner-vpn", "-n", "50", "--no-pager"], capture=False, check=False)
                raise RunnerError("VPN did not become ready; see VPN journal (routes/DNS/authentication)")
            time.sleep(2)


def up(cfg):
    READY.unlink(missing_ok=True)
    VPN_ERROR.unlink(missing_ok=True)
    host = read_json(ROOT / "host-network.json")
    # Reapply after pushed routes and on every reconnect. These routes are more
    # specific than OpenVPN's full-tunnel default and remain usable without VPN.
    aws_routes(cfg, host)
    servers = list(cfg["VPN_DNS_SERVERS"])
    if not servers:
        for key, value in sorted(os.environ.items()):
            words = value.split()
            if key.startswith("foreign_option_") and len(words) == 3 and words[:2] == ["dhcp-option", "DNS"]:
                servers.append(str(ipaddress.IPv4Address(words[2])))
            elif key.startswith("dns_server_") and "_address_" in key:
                try:
                    servers.append(str(ipaddress.IPv4Address(value)))
                except ValueError:
                    pass
    servers = list(dict.fromkeys(servers))
    if not servers:
        raise RunnerError("VPN supplied no IPv4 DNS servers; set VPN_DNS_SERVERS in aws.json")
    pool = ipaddress.IPv4Network(cfg["PODMAN_CIDR"])
    vpc = ipaddress.IPv4Network(cfg["VPC_CIDR"])
    for address in servers:
        ip = ipaddress.IPv4Address(address)
        if ip in pool or ip in vpc or ip.is_loopback or ip.is_link_local:
            raise RunnerError("VPN DNS overlaps worker networking; choose nonoverlapping VPC/PODMAN CIDRs")
        run(["ip", "route", "replace", address + "/32", "dev", "ulab-vpn"])
    for cidr in cfg["VPN_ROUTES"]:
        run(["ip", "route", "replace", cidr, "dev", "ulab-vpn"])
    for route in json.loads(run(["ip", "-j", "route", "show", "dev", "ulab-vpn"]).stdout):
        destination = route.get("dst", "default")
        net = ipaddress.IPv4Network("0.0.0.0/0" if destination == "default" else destination, strict=False)
        if net.prefixlen <= 1:
            continue  # Full-tunnel default or the two redirect-gateway def1 routes.
        if net.overlaps(pool) or net.overlaps(vpc):
            raise RunnerError(f"VPN route {destination} overlaps worker networking; change VPC_CIDR/PODMAN_CIDR")
    DNS_CONF.write_text(dns_config(cfg, host, servers))
    run(["systemctl", "restart", "ukama-runner-dns"])
    firewall(cfg, host)
    atomic_json(READY, {"dns": servers, "connected": True})
    print("VPN routes, DNS and firewall: ready", flush=True)


def probe_urls(cfg, env):
    # Agent gateway URLs belong to the lab's optional CDR diagnostics. Their
    # presence in the environment does not make them worker prerequisites.
    # Explicit operator-selected probes are supported via EXTRA_PROBE_URLS.
    values = [env.get("PAUTH_URL", "https://pauth.udev.ukama.com"),
              env.get("UKAMA_LAB_BFF", env.get("BFF_BASE_URL", "https://bff.udev.ukama.com") + "/gateway/graphql"),
              env.get("UKAMA_LAB_FACTORY_URL", "http://factory-ukama.udev.ukama.com"),
              env.get("UKAMA_LAB_WAREHOUSE_URL", "http://warehouse-ukama.udev.ukama.com"),
              *cfg["EXTRA_PROBE_URLS"]]
    result = []
    for value in values:
        if not value:
            continue
        parsed = urlsplit(value)
        if parsed.scheme not in ("http", "https") or not parsed.hostname or parsed.username or parsed.password:
            raise RunnerError("probe URLs must be HTTP(S) URLs without embedded credentials")
        # Connectivity only. No GraphQL operation, auth token, mutation or KPI check.
        origin = parsed.scheme + "://" + parsed.netloc + "/"
        if origin not in result:
            result.append(origin)
    return result


def check(cfg, env):
    if cfg["NETWORK_MODE"] == "vpn" and not READY.exists():
        raise RunnerError("VPN is disconnected")
    host = read_json(ROOT / "host-network.json")
    firewall(cfg, host)
    name = "ulab-infrastructure-probe"
    run(["podman", "network", "rm", name], check=False)
    run(["podman", "network", "create", "--disable-dns", name])
    base = env.get("BASE_IMAGE_REPO", "localhost/testing/virtualnode-base") + ":anode-" + env.get("NODE_RUNTIME", "starter")
    try:
        for url in probe_urls(cfg, env):
            hostname = urlsplit(url).hostname
            socket.getaddrinfo(hostname, None, family=socket.AF_INET)
            args = ["--ipv4", "--silent", "--show-error", "--connect-timeout", "10", "--max-time", "20", "--output", "/dev/null", url]
            run(["curl", *args], timeout=30)
            run(["podman", "run", "--rm", "--network", name, "--entrypoint", "/usr/bin/curl", base, *args], timeout=45)
        print("Host/container DNS, TCP and TLS connectivity: OK", flush=True)
    finally:
        run(["podman", "network", "rm", name], check=False)


def snapshot(destination):
    destination = Path(destination)
    destination.mkdir(parents=True, exist_ok=True)
    commands = {
        "routes": ["ip", "route", "show", "table", "all"],
        "addresses": ["ip", "address"],
        "dns": ["cat", "/etc/resolv.conf", "/etc/ukama-runner-dns.conf"],
        "vpn": ["journalctl", "-u", "ukama-runner-vpn", "--no-pager", "-n", "300"],
        "dns-journal": ["journalctl", "-u", "ukama-runner-dns", "--no-pager", "-n", "80"],
        "podman": ["podman", "info"], "containers": ["podman", "ps", "-a"],
        "disk": ["df", "-h"], "firewall": ["iptables-save"]}
    for name, command in commands.items():
        try:
            p = run(command, check=False, timeout=30)
            (destination / (name + ".txt")).write_text(p.stdout + p.stderr)
        except Exception as exc:
            (destination / (name + ".txt")).write_text(str(exc))


def main():
    action = sys.argv[1] if len(sys.argv) > 1 else ""
    try:
        cfg = read_json(ROOT / "worker.json")["config"]
        if action == "up":
            up(cfg)
        elif action == "down":
            READY.unlink(missing_ok=True)
        else:
            raise RunnerError("unknown network action")
    except Exception as exc:
        if action == "up":
            atomic_json(VPN_ERROR, {"error": str(exc)})
        print(f"worker network: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
