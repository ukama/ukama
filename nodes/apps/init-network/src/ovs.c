/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include "usys_log.h"

#include "config.h"
#include "exec.h"
#include "ovs.h"
#include "status.h"

#define OVS_MAX_STR             1024
#define SHELL_CMD_TIMEOUT_EXTRA 5
#define BASE_FLOW_COOKIE        "0x554b414d00000000"
#define BASE_FLOW_MASK          "0xffffffff00000000"
#define TUN_RULE_PRIORITY       20000
#define BRIDGE_RULE_PRIORITY    20001

static pthread_mutex_t networkMutex = PTHREAD_MUTEX_INITIALIZER;

static bool path_exists(const char *path) {

    struct stat st;

    if (path == NULL || path[0] == '\0') return false;

    return stat(path, &st) == 0;
}

static bool mkdir_p(const char *path) {

    char tmp[OVS_MAX_STR];
    char *p;
    size_t len;

    if (path == NULL || path[0] == '\0') return false;

    snprintf(tmp, sizeof(tmp), "%s", path);

    len = strlen(tmp);
    if (len == 0) return false;

    if (tmp[len - 1] == '/') tmp[len - 1] = '\0';

    for (p = tmp + 1; *p != '\0'; p++) {
        if (*p == '/') {
            *p = '\0';
            if (mkdir(tmp, 0755) != 0 && errno != EEXIST) return false;
            *p = '/';
        }
    }

    if (mkdir(tmp, 0755) != 0 && errno != EEXIST) return false;

    return true;
}

static void mark_bool(AppStatus *status, bool *field, bool value) {

    pthread_mutex_lock(&status->mutex);
    *field = value;
    pthread_mutex_unlock(&status->mutex);
}

static bool shell_ok(Config *config, const char *reason, const char *fmt, ...) {

    va_list ap;
    char cmd[OVS_MAX_STR * 2];
    int rc;

    if (config == NULL || fmt == NULL) return false;

    va_start(ap, fmt);
    rc = vsnprintf(cmd, sizeof(cmd), fmt, ap);
    va_end(ap);

    if (rc < 0 || (size_t)rc >= sizeof(cmd)) {
        usys_log_error("shell command too long: %s", reason);
        return false;
    }

    rc = exec_cmd(config->cmdTimeoutSec + SHELL_CMD_TIMEOUT_EXTRA,
                  "sh", "-c", cmd, NULL);
    if (rc != 0) {
        usys_log_error("%s: %s",
                       reason ? reason : "command failed",
                       cmd);
        return false;
    }

    return true;
}

static bool iface_exists(Config *config, const char *iface) {

    char cmd[OVS_MAX_STR];

    if (config == NULL || iface == NULL || iface[0] == '\0') return false;

    snprintf(cmd, sizeof(cmd), "ip link show %s >/dev/null 2>&1", iface);

    return shell_ok(config, NULL, "%s", cmd);
}

static bool ensure_mgmt_dir(Config *config, AppStatus *status) {

    char tmp[OVS_MAX_STR];

    if (strcmp(config->mgmtDir, config->runDir) == 0) {
        if (!mkdir_p(config->mgmtDir)) {
            status_fail(status, "failed to create OVS management dir");
            return false;
        }
        return true;
    }

    if (!mkdir_p(config->runDir)) {
        status_fail(status, "failed to create OVS run dir");
        return false;
    }

    if (path_exists(config->mgmtDir)) return true;

    snprintf(tmp, sizeof(tmp), "%s.tmp.%d", config->mgmtDir, getpid());
    unlink(tmp);

    if (symlink(config->runDir, tmp) != 0) {
        status_fail(status, "failed to create OVS management symlink");
        return false;
    }

    if (rename(tmp, config->mgmtDir) != 0) {
        unlink(tmp);
        status_fail(status, "failed to activate OVS management symlink");
        return false;
    }

    return true;
}

static bool check_tools(AppStatus *status) {

    const char *baseTools[] = {
        "ovs-vsctl",
        "ovs-ofctl",
        "ovsdb-server",
        "ovs-vswitchd",
        "ovsdb-tool",
        "ip",
        "ifconfig",
        "sysctl",
        "iptables",
        "sh",
        NULL
    };
    char reason[STATUS_REASON_LEN];
    int i;

    status_set(status, InitStateCheckTools, "checking required tools");

    for (i = 0; baseTools[i] != NULL; i++) {
        if (!exec_tool_exists(baseTools[i])) {
            snprintf(reason, sizeof(reason), "required tool missing: %s",
                     baseTools[i]);
            status_fail(status, reason);
            return false;
        }
    }

    mark_bool(status, &status->toolsOk, true);
    return true;
}

static bool ovs_is_running(Config *config) {

    if (config == NULL) return false;

    return shell_ok(config, NULL,
                    "ovs-vsctl show >/dev/null 2>&1");
}

static bool start_ovs(Config *config, AppStatus *status) {

    char dbPath[OVS_MAX_STR];
    char dbSock[OVS_MAX_STR];
    char dbRemote[OVS_MAX_STR];
    char ovsdbPid[OVS_MAX_STR];
    char vswitchdPid[OVS_MAX_STR];
    char vswitchdDb[OVS_MAX_STR];
    char ovsdbPidOpt[OVS_MAX_STR + 32];
    char vswitchdPidOpt[OVS_MAX_STR + 32];
    int rc;

    status_set(status, InitStateStartOvs, "starting ovs");

    if (!mkdir_p(config->runDir)) {
        status_fail(status, "failed to create OVS run directory");
        return false;
    }

    if (!mkdir_p(config->dbDir)) {
        status_fail(status, "failed to create OVS db directory");
        return false;
    }

    if (!ensure_mgmt_dir(config, status)) return false;

    if (ovs_is_running(config)) {
        usys_log_info("OVS is already running");
        mark_bool(status, &status->ovsdbRunning, true);
        mark_bool(status, &status->vswitchdRunning, true);
        return true;
    }

    snprintf(dbPath, sizeof(dbPath), "%s/conf.db", config->dbDir);
    snprintf(dbSock, sizeof(dbSock), "punix:%s/db.sock", config->runDir);
    snprintf(dbRemote, sizeof(dbRemote),
             "db:Open_vSwitch,Open_vSwitch,manager_options");
    snprintf(ovsdbPid, sizeof(ovsdbPid), "%s/ovsdb-server.pid",
             config->runDir);
    snprintf(vswitchdPid, sizeof(vswitchdPid), "%s/ovs-vswitchd.pid",
             config->runDir);
    snprintf(vswitchdDb, sizeof(vswitchdDb), "unix:%s/db.sock",
             config->runDir);

    rc = snprintf(ovsdbPidOpt, sizeof(ovsdbPidOpt), "--pidfile=%s",
                  ovsdbPid);
    if (rc < 0 || (size_t)rc >= sizeof(ovsdbPidOpt)) {
        status_fail(status, "ovsdb pidfile option too long");
        return false;
    }

    rc = snprintf(vswitchdPidOpt, sizeof(vswitchdPidOpt), "--pidfile=%s",
                  vswitchdPid);
    if (rc < 0 || (size_t)rc >= sizeof(vswitchdPidOpt)) {
        status_fail(status, "ovs-vswitchd pidfile option too long");
        return false;
    }

    if (!path_exists(dbPath)) {
        if (!path_exists(config->schema)) {
            status_fail(status, "OVS schema not found");
            return false;
        }

        if (exec_cmd(config->cmdTimeoutSec,
                     "ovsdb-tool", "create", dbPath, config->schema,
                     NULL) != 0) {
            status_fail(status, "failed to create OVS database");
            return false;
        }
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovsdb-server", dbPath, "--remote", dbSock,
                 "--remote", dbRemote, ovsdbPidOpt, "--detach",
                 NULL) != 0) {
        status_fail(status, "failed to start ovsdb-server");
        return false;
    }

    mark_bool(status, &status->ovsdbRunning, true);

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vsctl", "--no-wait", "init", NULL) != 0) {
        status_fail(status, "failed to initialize OVS database");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vswitchd", vswitchdDb, vswitchdPidOpt, "--detach",
                 NULL) != 0) {
        status_fail(status, "failed to start ovs-vswitchd");
        return false;
    }

    mark_bool(status, &status->vswitchdRunning, true);

    if (!ovs_is_running(config)) {
        status_fail(status, "OVS did not become ready");
        return false;
    }

    return true;
}

static bool setup_epc_aliases(Config *config, AppStatus *status) {

    char epcSctp[OVS_MAX_STR];
    char epcGtpu[OVS_MAX_STR];

    if (!config->epcEnable) {
        mark_bool(status, &status->epcIfReady, true);
        return true;
    }

    status_set(status, InitStateSetupEpcIf, "setting up EPC host aliases");

    snprintf(epcSctp, sizeof(epcSctp), "%s:0", config->epcIf);
    snprintf(epcGtpu, sizeof(epcGtpu), "%s:1", config->epcIf);

    if (exec_cmd(config->cmdTimeoutSec,
                 "ifconfig", epcSctp, config->epcSctpAddr, NULL) != 0) {
        status_fail(status, "failed to assign EPC SCTP address");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ifconfig", epcGtpu, config->epcGtpuAddr, "up",
                 NULL) != 0) {
        status_fail(status, "failed to assign EPC GTPU address");
        return false;
    }

    mark_bool(status, &status->epcIfReady, true);
    return true;
}

static bool tun_exists(Config *config) {

    return shell_ok(config, "persistent TUN is unavailable",
                    "ip tuntap show dev %s | grep -F '%s: tun' | "
                    "grep -qw persist",
                    config->tunIf, config->tunIf);
}

static bool setup_tun(Config *config, AppStatus *status) {

    int i;

    mark_bool(status, &status->tunReady, false);
    if (!config->tunEnable) {
        if (config->enablePolicyRouting && !tun_exists(config)) {
            status_fail(status, "policy routing requires a persistent TUN");
            return false;
        }
        mark_bool(status, &status->tunReady,
                  iface_exists(config, config->tunIf));
        return true;
    }

    status_set(status, InitStateSetupTun, "setting up tun interface");

    if (!iface_exists(config, config->tunIf) &&
        exec_cmd(config->cmdTimeoutSec,
                 "ip", "tuntap", "add", "dev", config->tunIf,
                 "mode", "tun", NULL) != 0) {
        status_fail(status, "failed to create tun interface");
        return false;
    }

    if (!tun_exists(config)) {
        status_fail(status, "configured interface is not a persistent TUN");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ip", "addr", "replace", config->tunPrimaryCidr,
                 "dev", config->tunIf, NULL) != 0) {
        status_fail(status, "failed to assign primary tun address");
        return false;
    }

    for (i = 0; i < config->tunExtraCount; i++) {
        if (exec_cmd(config->cmdTimeoutSec,
                     "ip", "addr", "replace", config->tunExtraCidrs[i],
                     "dev", config->tunIf, NULL) != 0) {
            status_fail(status, "failed to assign extra tun address");
            return false;
        }
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ip", "link", "set", config->tunIf, "up", NULL) != 0) {
        status_fail(status, "failed to bring tun interface up");
        return false;
    }

    mark_bool(status, &status->tunReady, true);
    return true;
}

static bool setup_bridge(Config *config, AppStatus *status) {

    char controller[OVS_MAX_STR];
    char protocols[INIT_NETWORK_MAX_STR + 16];

    status_set(status, InitStateSetupBridge, "setting up OVS bridge");

    snprintf(controller, sizeof(controller), "punix:%s/%s.mgmt",
             config->mgmtDir, config->bridge);

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vsctl", "--may-exist", "add-br", config->bridge,
                 NULL) != 0) {
        status_fail(status, "failed to create OVS bridge");
        return false;
    }

    snprintf(protocols, sizeof(protocols), "protocols=%s", config->openflow);
    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vsctl", "set", "bridge", config->bridge,
                 protocols, NULL) != 0) {
        status_fail(status, "failed to set bridge OpenFlow version");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vsctl", "set-controller", config->bridge,
                 controller, NULL) != 0) {
        status_fail(status, "failed to set bridge controller socket");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vsctl", "set-fail-mode", config->bridge,
                 "standalone", NULL) != 0) {
        status_fail(status, "failed to set bridge fail-mode");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ip", "link", "set", config->bridge, "up", NULL) != 0) {
        status_fail(status, "failed to bring bridge up");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ip", "addr", "replace", config->bridgeCidr,
                 "dev", config->bridge, NULL) != 0) {
        status_fail(status, "failed to configure bridge address");
        return false;
    }

    mark_bool(status, &status->bridgeReady, true);
    return true;
}

static bool ensure_iptables_rule(Config *config, const char *table,
                                 const char *chain, const char *rule) {

    if (table != NULL && table[0] != '\0') {
        return shell_ok(config, "failed to ensure iptables rule",
                        "iptables -t %s -C %s %s 2>/dev/null || "
                        "iptables -t %s -A %s %s",
                        table, chain, rule, table, chain, rule);
    }

    return shell_ok(config, "failed to ensure iptables rule",
                    "iptables -C %s %s 2>/dev/null || iptables -A %s %s",
                    chain, rule, chain, rule);
}

static bool setup_forwarding(Config *config, AppStatus *status) {

    char rule[OVS_MAX_STR];

    status_set(status, InitStateSetupForwarding, "setting up forwarding");

    if (config->enableIpForward) {
        if (exec_cmd(config->cmdTimeoutSec,
                     "sysctl", "-w", "net.ipv4.ip_forward=1",
                     NULL) != 0) {
            status_fail(status, "failed to enable ip_forward");
            return false;
        }

        shell_ok(config, "failed to disable rp_filter",
                 "sysctl -w net.ipv4.conf.all.rp_filter=0 >/dev/null 2>&1 "
                 "|| true; "
                 "sysctl -w net.ipv4.conf.default.rp_filter=0 "
                 ">/dev/null 2>&1 || true; "
                 "sysctl -w net.ipv4.conf.%s.rp_filter=0 "
                 ">/dev/null 2>&1 || true; "
                 "sysctl -w net.ipv4.conf.%s.rp_filter=0 "
                 ">/dev/null 2>&1 || true",
                 config->bridge,
                 config->externalIf);

        /* The TUN is now created before default.rp_filter is disabled. */
        if ((config->tunEnable || config->enablePolicyRouting) &&
            !shell_ok(config, "failed to disable tun rp_filter",
                      "sysctl -w net.ipv4.conf.%s.rp_filter=0",
                      config->tunIf)) {
            status_fail(status, "failed to disable tun rp_filter");
            return false;
        }
    }

    if (config->enableNat) {
        snprintf(rule, sizeof(rule), "-s %s -o %s -j MASQUERADE",
                 config->bridgeSubnet, config->externalIf);
        if (!ensure_iptables_rule(config, "nat", "POSTROUTING", rule)) {
            status_fail(status, "failed to add bridge NAT rule");
            return false;
        }

        snprintf(rule, sizeof(rule), "-s %s -o %s -j MASQUERADE",
                 config->ueCidr, config->externalIf);
        if (!ensure_iptables_rule(config, "nat", "POSTROUTING", rule)) {
            status_fail(status, "failed to add UE NAT rule");
            return false;
        }

        snprintf(rule, sizeof(rule),
                 "-i %s -o %s -m state --state RELATED,ESTABLISHED -j ACCEPT",
                 config->externalIf, config->bridge);
        if (!ensure_iptables_rule(config, NULL, "FORWARD", rule)) {
            status_fail(status, "failed to add bridge inbound FORWARD rule");
            return false;
        }

        snprintf(rule, sizeof(rule), "-i %s -o %s -j ACCEPT",
                 config->bridge, config->externalIf);
        if (!ensure_iptables_rule(config, NULL, "FORWARD", rule)) {
            status_fail(status, "failed to add bridge outbound FORWARD rule");
            return false;
        }
    }

    mark_bool(status, &status->forwardingReady, true);
    return true;
}

static bool ensure_gateway_link(Config *config, AppStatus *status,
                                 bool netns) {

    bool bridgeExists;
    bool peerExists;

    bridgeExists = iface_exists(config, config->gatewayBridgeIf);
    if (netns) {
        peerExists = shell_ok(config, NULL,
                              "ip -n %s link show dev %s >/dev/null 2>&1",
                              config->gatewayName,
                              config->gatewayNamespaceIf);
    } else {
        peerExists = iface_exists(config, config->gatewayNamespaceIf);
    }

    if (bridgeExists && peerExists) {
        /* Refuse unrelated interfaces instead of replacing live links. */
        if (!shell_ok(config, "gateway bridge port is not a veth",
                      "ip -d link show dev %s | grep -qw veth",
                      config->gatewayBridgeIf)) {
            status_fail(status, "gateway bridge port is not a veth");
            return false;
        }
        if (netns) {
            peerExists = shell_ok(config, "gateway peer is not a veth",
                                  "ip -n %s -d link show dev %s | "
                                  "grep -qw veth",
                                  config->gatewayName,
                                  config->gatewayNamespaceIf);
        } else {
            peerExists = shell_ok(config, "gateway veth peers differ",
                                  "test \"$(cat /sys/class/net/%s/iflink)\" "
                                  "= \"$(cat /sys/class/net/%s/ifindex)\" && "
                                  "test \"$(cat /sys/class/net/%s/iflink)\" "
                                  "= \"$(cat /sys/class/net/%s/ifindex)\"",
                                  config->gatewayBridgeIf,
                                  config->gatewayNamespaceIf,
                                  config->gatewayNamespaceIf,
                                  config->gatewayBridgeIf);
        }
        if (!peerExists) status_fail(status, "gateway veth peers differ");
        return peerExists;
    }

    if (bridgeExists || peerExists) {
        status_fail(status, "gateway veth pair is incomplete");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ip", "link", "add", config->gatewayBridgeIf,
                 "type", "veth", "peer", "name",
                 config->gatewayNamespaceIf, NULL) != 0) {
        status_fail(status, "failed to create gateway veth");
        return false;
    }

    if (netns && exec_cmd(config->cmdTimeoutSec,
                         "ip", "link", "set", config->gatewayNamespaceIf,
                         "netns", config->gatewayName, NULL) != 0) {
        status_fail(status, "failed to move gateway peer to namespace");
        return false;
    }

    return true;
}

static bool setup_gateway_root(Config *config, AppStatus *status) {

    char rule[OVS_MAX_STR];

    status_set(status, InitStateSetupGateway, "setting up root gateway");

    if (!ensure_gateway_link(config, status, false)) return false;

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vsctl", "--may-exist", "add-port", config->bridge,
                 config->gatewayBridgeIf, NULL) != 0) {
        status_fail(status, "failed to add gateway port to bridge");
        return false;
    }

    if (!shell_ok(config, "failed to configure root gateway interface",
                  "ip link set %s up && "
                  "ip link set %s up && "
                  "ip addr replace %s/32 dev %s && "
                  "ip route replace %s/32 dev %s && "
                  "ip route replace %s via %s dev %s && "
                  "(sysctl -w net.ipv4.conf.%s.rp_filter=0 "
                  ">/dev/null 2>&1 || true)",
                  config->gatewayBridgeIf,
                  config->gatewayNamespaceIf,
                  config->gatewayIp, config->gatewayNamespaceIf,
                  config->bridgeAddr, config->gatewayNamespaceIf,
                  config->ueCidr, config->bridgeAddr,
                  config->gatewayNamespaceIf,
                  config->gatewayNamespaceIf)) {
        status_fail(status, "failed to configure root gateway interface");
        return false;
    }

    if (config->enableNat) {
        snprintf(rule, sizeof(rule), "-s %s -o %s -j MASQUERADE",
                 config->ueCidr, config->externalIf);
        if (!ensure_iptables_rule(config, "nat", "POSTROUTING", rule)) {
            status_fail(status, "failed to add gateway UE NAT rule");
            return false;
        }

        snprintf(rule, sizeof(rule),
                 "-i %s -o %s -m state --state RELATED,ESTABLISHED -j ACCEPT",
                 config->externalIf, config->gatewayNamespaceIf);
        if (!ensure_iptables_rule(config, NULL, "FORWARD", rule)) {
            status_fail(status, "failed to add gateway inbound FORWARD rule");
            return false;
        }

        snprintf(rule, sizeof(rule), "-i %s -o %s -j ACCEPT",
                 config->gatewayNamespaceIf, config->externalIf);
        if (!ensure_iptables_rule(config, NULL, "FORWARD", rule)) {
            status_fail(status, "failed to add gateway outbound FORWARD rule");
            return false;
        }
    }

    return true;
}

static bool setup_gateway_netns(Config *config, AppStatus *status) {

    status_set(status, InitStateSetupGateway, "setting up gateway namespace");

    if (!shell_ok(config, "failed to create netns directory",
                  "mkdir -p /var/run/netns")) {
        status_fail(status, "failed to create netns directory");
        return false;
    }

    if (!shell_ok(config, "failed to create gateway namespace",
                  "ip netns list | awk '{print $1}' | grep -qx %s || "
                  "ip netns add %s",
                  config->gatewayName, config->gatewayName)) {
        status_fail(status, "failed to create gateway namespace");
        return false;
    }

    if (!ensure_gateway_link(config, status, true)) return false;

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vsctl", "--may-exist", "add-port", config->bridge,
                 config->gatewayBridgeIf, NULL) != 0) {
        status_fail(status, "failed to add gateway port to bridge");
        return false;
    }

    if (!shell_ok(config, "failed to configure gateway namespace",
                  "ip link set %s up && "
                  "ip netns exec %s ip link set lo up && "
                  "ip netns exec %s ip link set %s up && "
                  "ip netns exec %s ip addr replace %s dev %s && "
                  "ip netns exec %s ip route replace default via %s dev %s && "
                  "ip netns exec %s sysctl -w net.ipv4.ip_forward=1",
                  config->gatewayBridgeIf,
                  config->gatewayName,
                  config->gatewayName, config->gatewayNamespaceIf,
                  config->gatewayName, config->gatewayAddr,
                  config->gatewayNamespaceIf,
                  config->gatewayName, config->bridgeAddr,
                  config->gatewayNamespaceIf,
                  config->gatewayName)) {
        status_fail(status, "failed to configure gateway namespace");
        return false;
    }

    return true;
}

static bool setup_gateway(Config *config, AppStatus *status) {

    if (!config->gatewayEnable) {
        mark_bool(status, &status->gatewayReady, true);
        return true;
    }

    if (strcmp(config->gatewayMode, "root") == 0) {
        if (!setup_gateway_root(config, status)) return false;
    } else if (strcmp(config->gatewayMode, "netns") == 0) {
        if (!setup_gateway_netns(config, status)) return false;
    } else {
        status_fail(status, "unsupported gateway mode");
        return false;
    }

    mark_bool(status, &status->gatewayReady, true);
    return true;
}

static bool setup_flows(Config *config, AppStatus *status) {

    char path[] = "/tmp/init-network-flows.XXXXXX";
    FILE *file;
    int fd;
    bool ok;

    status_set(status, InitStateSetupFlows, "setting up default OVS flows");
    mark_bool(status, &status->flowsReady, false);

    fd = mkstemp(path);
    if (fd < 0) {
        status_fail(status, "failed to create baseline flow file");
        return false;
    }
    file = fdopen(fd, "w");
    if (file == NULL) {
        close(fd);
        unlink(path);
        status_fail(status, "failed to open baseline flow file");
        return false;
    }

    /* Replace only our baseline, atomically. PCRF uses 32-bit cookies. */
    fprintf(file, "delete cookie=%s/%s,table=0\n",
            BASE_FLOW_COOKIE, BASE_FLOW_MASK);
    if (config->defaultDrop) {
        fprintf(file, "add cookie=0x554b414d00000001,table=0,"
                "priority=10,ip,nw_src=%s,actions=drop\n", config->ueCidr);
        fprintf(file, "add cookie=0x554b414d00000002,table=0,"
                "priority=10,ip,nw_dst=%s,actions=drop\n", config->ueCidr);
    }
    fprintf(file, "add cookie=%s,table=0,priority=0,actions=NORMAL\n",
            BASE_FLOW_COOKIE);

    ok = !ferror(file);
    if (fclose(file) != 0) ok = false;
    if (ok) {
        ok = exec_cmd(config->cmdTimeoutSec,
                      "ovs-ofctl", "-O", config->openflow, "--bundle",
                      "add-flows", config->bridge, path, NULL) == 0;
    }
    unlink(path);
    if (!ok) {
        status_fail(status, "failed to install baseline OVS flows");
        return false;
    }

    mark_bool(status, &status->flowsReady, true);
    return true;
}

static bool ensure_policy_rule(Config *config, const char *iface,
                                int table, int priority) {

    /* A different owner at our priority is an error, not a deletion. */
    return shell_ok(config, "failed to ensure policy rule",
                    "rules=\"$(ip -N -4 rule show priority %d)\" && "
                    "if [ -z \"$rules\" ]; then "
                    "ip -4 rule add priority %d iif %s table %d; "
                    "else printf '%%s\n' \"$rules\" | "
                    "awk 'NF {n++; if ($1 != \"%d:\" || "
                    "$2 != \"from\" || $3 != \"all\" || "
                    "$4 != \"iif\" || $5 != \"%s\" || "
                    "$6 != \"lookup\" || $7 != \"%d\" || NF != 7) "
                    "bad=1} END {exit n != 1 || bad}'; fi",
                    priority, priority, iface, table,
                    priority, iface, table);
}

static bool setup_policy_routing(Config *config, AppStatus *status) {

    status_set(status, InitStateSetupPolicyRouting,
               "setting up tun/ovs policy routing");

    if (!config->enablePolicyRouting) {
        mark_bool(status, &status->policyRoutingReady, true);
        return true;
    }

    if (!iface_exists(config, config->tunIf)) {
        mark_bool(status, &status->policyRoutingReady, false);
        status_fail(status, "tun interface is not ready for policy routing");
        return false;
    }

    if (!iface_exists(config, config->bridge)) {
        status_fail(status, "bridge not found for policy routing");
        mark_bool(status, &status->policyRoutingReady, false);
        return false;
    }

    if (!shell_ok(config, "failed to add tun table default route",
                  "ip route replace default via %s dev %s table %d",
                  config->gatewayIp, config->bridge, config->tunTable)) {
        status_fail(status, "failed to add tun table default route");
        mark_bool(status, &status->policyRoutingReady, false);
        return false;
    }

    if (!shell_ok(config, "failed to add bridge table UE route",
                  "ip route replace %s dev %s table %d",
                  config->ueCidr, config->tunIf, config->bridgeTable)) {
        status_fail(status, "failed to add bridge table UE route");
        mark_bool(status, &status->policyRoutingReady, false);
        return false;
    }

    if (!ensure_policy_rule(config, config->tunIf, config->tunTable,
                             TUN_RULE_PRIORITY) ||
        !ensure_policy_rule(config, config->bridge, config->bridgeTable,
                             BRIDGE_RULE_PRIORITY)) {
        status_fail(status, "failed to configure policy rules");
        mark_bool(status, &status->policyRoutingReady, false);
        return false;
    }

    mark_bool(status, &status->policyRoutingReady, true);
    return true;
}

static bool interface_up(Config *config, const char *iface) {

    return shell_ok(config, "interface is not up",
                    "ip -o link show dev %s | "
                    "grep -Eq '<([^>]*,)?UP[,>]'", iface);
}

static bool check_policy_rule(Config *config, const char *iface,
                              int table, int priority) {

    return shell_ok(config, "policy rule is unavailable",
                    "ip -N -4 rule show priority %d | "
                    "awk 'NF {n++; if ($1 != \"%d:\" || "
                    "$2 != \"from\" || $3 != \"all\" || "
                    "$4 != \"iif\" || $5 != \"%s\" || "
                    "$6 != \"lookup\" || $7 != \"%d\" || NF != 7) "
                    "bad=1} END {exit n != 1 || bad}'",
                    priority, priority, iface, table);
}

static bool check_policy_routing(Config *config) {

    if (!config->enablePolicyRouting) return true;

    if (!check_policy_rule(config, config->tunIf, config->tunTable,
                            TUN_RULE_PRIORITY) ||
        !check_policy_rule(config, config->bridge, config->bridgeTable,
                            BRIDGE_RULE_PRIORITY)) return false;

    return shell_ok(config, "policy routes are unavailable",
                    "ip -4 route show table %d exact default | "
                    "awk '$1 == \"default\" && $2 == \"via\" && "
                    "$3 == \"%s\" && $4 == \"dev\" && $5 == \"%s\" "
                    "{found=1} END {exit !found}' && "
                    "ip -4 route show table %d exact %s | "
                    "awk '$1 == \"%s\" && $2 == \"dev\" && $3 == \"%s\" "
                    "{found=1} END {exit !found}'",
                    config->tunTable, config->gatewayIp, config->bridge,
                    config->bridgeTable, config->ueCidr,
                    config->ueCidr, config->tunIf);
}

static bool check_baseline(Config *config) {

    if (!shell_ok(config, "baseline NORMAL rule is unavailable",
                  "ovs-ofctl -O %s dump-flows %s "
                  "'table=0,cookie=%s/-1' | grep -q 'actions=NORMAL'",
                  config->openflow, config->bridge, BASE_FLOW_COOKIE)) {
        return false;
    }
    if (!config->defaultDrop) return true;

    return shell_ok(config, "baseline UE drop rules are unavailable",
                    "ovs-ofctl -O %s dump-flows %s "
                    "'table=0,cookie=0x554b414d00000001/-1,ip,nw_src=%s' "
                    "| grep -q 'actions=drop' && "
                    "ovs-ofctl -O %s dump-flows %s "
                    "'table=0,cookie=0x554b414d00000002/-1,ip,nw_dst=%s' "
                    "| grep -q 'actions=drop'",
                    config->openflow, config->bridge, config->ueCidr,
                    config->openflow, config->bridge, config->ueCidr);
}

static bool check_gateway(Config *config) {

    if (!config->gatewayEnable) return true;
    if (!interface_up(config, config->gatewayBridgeIf)) return false;

    if (!shell_ok(config, "gateway port is not on the configured bridge",
                  "test \"$(ovs-vsctl port-to-br %s)\" = \"%s\"",
                  config->gatewayBridgeIf, config->bridge)) return false;

    if (strcmp(config->gatewayMode, "netns") == 0) {
        return shell_ok(config, "gateway namespace interface is down",
                        "ip -n %s -o link show dev %s | "
                        "grep -Eq '<([^>]*,)?UP[,>]'",
                        config->gatewayName, config->gatewayNamespaceIf);
    }

    if (!interface_up(config, config->gatewayNamespaceIf)) return false;
    return shell_ok(config, "gateway return route is unavailable",
                    "ip -4 route show table main exact %s | "
                    "awk '$1 == \"%s\" && $2 == \"via\" && $3 == \"%s\" && "
                    "$4 == \"dev\" && $5 == \"%s\" "
                    "{found=1} END {exit !found}'",
                    config->ueCidr, config->ueCidr, config->bridgeAddr,
                    config->gatewayNamespaceIf);
}

static bool check_network(Config *config, char *reason, size_t reasonSize) {

    const char *error = NULL;

    if (!ovs_is_running(config)) {
        error = "ovs is not running";
    } else if (!interface_up(config, config->bridge)) {
        error = "ovs bridge is unavailable or down";
    } else if ((config->tunEnable || config->enablePolicyRouting) &&
               (!tun_exists(config) ||
                !interface_up(config, config->tunIf))) {
        error = "persistent tun interface is unavailable or down";
    } else if (!check_baseline(config)) {
        error = "baseline OVS flows are unavailable";
    } else if (!check_gateway(config)) {
        error = "gateway path is unavailable";
    } else if (!check_policy_routing(config)) {
        error = "tun/bridge policy routing is unavailable";
    } else if (config->enableIpForward &&
               !shell_ok(config, "ip forwarding is disabled",
                         "test \"$(sysctl -n net.ipv4.ip_forward)\" = 1")) {
        error = "ip forwarding is disabled";
    } else if (config->enableNat &&
               !shell_ok(config, "UE NAT rule is unavailable",
                         "iptables -t nat -C POSTROUTING "
                         "-s %s -o %s -j MASQUERADE",
                         config->ueCidr, config->externalIf)) {
        error = "UE NAT rule is unavailable";
    }

    snprintf(reason, reasonSize, "%s", error ? error : "ready");
    return error == NULL;
}

bool ovs_check(Config *config, char *reason, size_t reasonSize) {

    bool ok;

    if (config == NULL || reason == NULL || reasonSize == 0) return false;

    pthread_mutex_lock(&networkMutex);
    ok = check_network(config, reason, reasonSize);
    pthread_mutex_unlock(&networkMutex);
    return ok;
}

static bool configure_network(Config *config, AppStatus *status) {

    char reason[STATUS_REASON_LEN];

    if (!setup_epc_aliases(config, status)) return false;
    if (!setup_tun(config, status))         return false;
    if (!setup_bridge(config, status))      return false;
    if (!setup_flows(config, status))       return false;
    if (!setup_forwarding(config, status))  return false;
    if (!setup_gateway(config, status))     return false;
    if (!setup_policy_routing(config, status)) return false;

    if (!check_network(config, reason, sizeof(reason))) {
        status_fail(status, reason);
        return false;
    }

    status_set(status, InitStateReady, "ready");
    return true;
}

bool ovs_reconcile(Config *config, AppStatus *status) {

    bool ok;

    if (config == NULL || status == NULL) return false;

    pthread_mutex_lock(&networkMutex);
    ok = configure_network(config, status);
    pthread_mutex_unlock(&networkMutex);
    return ok;
}

bool ovs_setup(Config *config, AppStatus *status) {

    bool ok;

    if (config == NULL || status == NULL) return false;

    pthread_mutex_lock(&networkMutex);
    ok = check_tools(status) && start_ovs(config, status);
    if (ok) ok = configure_network(config, status);
    pthread_mutex_unlock(&networkMutex);
    return ok;
}
