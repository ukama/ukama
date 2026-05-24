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
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include "usys_log.h"

#include "config.h"
#include "exec.h"
#include "ovs.h"
#include "status.h"

#define OVS_MAX_STR 1024
#define SHELL_CMD_TIMEOUT_EXTRA 5

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
        usys_log_error("%s: %s", reason, cmd);
        return false;
    }

    return true;
}

static bool iface_exists(Config *config, const char *iface) {

    if (config == NULL || iface == NULL || iface[0] == '\0') return false;

    return exec_cmd(config->cmdTimeoutSec,
                    "ip", "link", "show", iface, NULL) == 0;
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

static bool check_tools(Config *config, AppStatus *status) {

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

    if (config->tunEnable && !exec_tool_exists("openvpn")) {
        status_fail(status, "required tool missing: openvpn");
        return false;
    }

    mark_bool(status, &status->toolsOk, true);
    return true;
}

static bool ovs_is_running(Config *config) {

    return exec_cmd(config->cmdTimeoutSec,
                    "ovs-vsctl", "--timeout=2", "show", NULL) == 0;
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

static bool setup_tun(Config *config, AppStatus *status) {

    int i;

    if (!config->tunEnable) {
        mark_bool(status, &status->tunReady, true);
        return true;
    }

    status_set(status, InitStateSetupTun, "setting up tun interface");

    exec_cmd(config->cmdTimeoutSec, "ip", "link", "delete", config->tunIf,
             NULL);

    if (exec_cmd(config->cmdTimeoutSec,
                 "openvpn", "--mktun", "--dev", config->tunIf,
                 NULL) != 0) {
        status_fail(status, "failed to create tun interface");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ip", "link", "set", config->tunIf, "up", NULL) != 0) {
        status_fail(status, "failed to bring tun interface up");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ip", "addr", "replace", config->tunPrimaryCidr,
                 "dev", config->tunIf, NULL) != 0) {
        status_fail(status, "failed to assign primary tun address");
        return false;
    }

    for (i = 0; i < config->tunExtraCount; i++) {
        exec_cmd(config->cmdTimeoutSec,
                 "ip", "addr", "replace", config->tunExtraCidrs[i],
                 "dev", config->tunIf, NULL);
    }

    mark_bool(status, &status->tunReady, true);
    return true;
}

static bool setup_bridge(Config *config, AppStatus *status) {

    char controller[OVS_MAX_STR];

    status_set(status, InitStateSetupBridge, "setting up OVS bridge");

    snprintf(controller, sizeof(controller), "punix:%s/%s.mgmt",
             config->mgmtDir, config->bridge);

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vsctl", "--may-exist", "add-br", config->bridge,
                 NULL) != 0) {
        status_fail(status, "failed to create OVS bridge");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vsctl", "set", "bridge", config->bridge,
                 "protocols=OpenFlow15", NULL) != 0) {
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
    }

    if (config->enableNat) {
        snprintf(rule, sizeof(rule), "-s %s -o %s -j MASQUERADE",
                 config->bridgeSubnet, config->externalIf);
        if (!ensure_iptables_rule(config, "nat", "POSTROUTING", rule)) {
            status_fail(status, "failed to add bridge NAT rule");
            return false;
        }

        snprintf(rule, sizeof(rule),
                 "-i %s -o %s -m state --state RELATED,ESTABLISHED -j ACCEPT",
                 config->externalIf, config->bridge);
        if (!ensure_iptables_rule(config, NULL, "FORWARD", rule)) {
            status_fail(status, "failed to add inbound FORWARD rule");
            return false;
        }

        snprintf(rule, sizeof(rule), "-i %s -o %s -j ACCEPT",
                 config->bridge, config->externalIf);
        if (!ensure_iptables_rule(config, NULL, "FORWARD", rule)) {
            status_fail(status, "failed to add outbound FORWARD rule");
            return false;
        }
    }

    mark_bool(status, &status->forwardingReady, true);
    return true;
}

static bool setup_gateway(Config *config, AppStatus *status) {

    if (!config->gatewayEnable) {
        mark_bool(status, &status->gatewayReady, true);
        return true;
    }

    status_set(status, InitStateSetupGateway, "setting up gateway namespace");

    if (strcmp(config->gatewayMode, "netns") != 0) {
        status_fail(status, "unsupported gateway mode");
        return false;
    }

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

    if (!shell_ok(config, "failed to recreate gateway veth",
                  "ip link show %s >/dev/null 2>&1 && ip link del %s || true; "
                  "ip netns exec %s ip link show %s >/dev/null 2>&1 && "
                  "ip netns exec %s ip link del %s || true; "
                  "ip link add %s type veth peer name %s",
                  config->gatewayBridgeIf, config->gatewayBridgeIf,
                  config->gatewayName, config->gatewayNamespaceIf,
                  config->gatewayName, config->gatewayNamespaceIf,
                  config->gatewayBridgeIf, config->gatewayNamespaceIf)) {
        status_fail(status, "failed to create gateway veth");
        return false;
    }

    if (exec_cmd(config->cmdTimeoutSec,
                 "ovs-vsctl", "--may-exist", "add-port", config->bridge,
                 config->gatewayBridgeIf, NULL) != 0) {
        status_fail(status, "failed to add gateway port to bridge");
        return false;
    }

    if (!shell_ok(config, "failed to configure gateway namespace",
                  "ip link set %s up && "
                  "ip link set %s netns %s && "
                  "ip netns exec %s ip link set lo up && "
                  "ip netns exec %s ip link set %s up && "
                  "ip netns exec %s ip addr replace %s dev %s && "
                  "ip netns exec %s ip route replace default via %s dev %s && "
                  "ip netns exec %s sysctl -w net.ipv4.ip_forward=1",
                  config->gatewayBridgeIf,
                  config->gatewayNamespaceIf, config->gatewayName,
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

    if (!shell_ok(config, "failed to configure gateway NAT",
                  "ip netns exec %s iptables -t nat -C POSTROUTING "
                  "-s %s -o %s -j MASQUERADE 2>/dev/null || "
                  "ip netns exec %s iptables -t nat -A POSTROUTING "
                  "-s %s -o %s -j MASQUERADE",
                  config->gatewayName, config->ueCidr,
                  config->gatewayNamespaceIf,
                  config->gatewayName, config->ueCidr,
                  config->gatewayNamespaceIf)) {
        status_fail(status, "failed to configure gateway NAT");
        return false;
    }

    mark_bool(status, &status->gatewayReady, true);
    return true;
}

static bool add_flow(Config *config, const char *flow) {

    return exec_cmd(config->cmdTimeoutSec,
                    "ovs-ofctl", "-O", config->openflow, "add-flow",
                    config->bridge, flow, NULL) == 0;
}

static bool setup_flows(Config *config, AppStatus *status) {

    char srcDrop[OVS_MAX_STR];
    char dstDrop[OVS_MAX_STR];

    status_set(status, InitStateSetupFlows, "setting up default OVS flows");

    exec_cmd(config->cmdTimeoutSec,
             "ovs-ofctl", "-O", config->openflow, "del-flows",
             config->bridge, "priority=0", NULL);

    if (!add_flow(config, "priority=0,actions=NORMAL")) {
        status_fail(status, "failed to add default NORMAL flow");
        return false;
    }

    if (config->defaultDrop) {
        snprintf(srcDrop, sizeof(srcDrop),
                 "priority=10,ip,nw_src=%s,actions=drop", config->ueCidr);
        snprintf(dstDrop, sizeof(dstDrop),
                 "priority=10,ip,nw_dst=%s,actions=drop", config->ueCidr);

        exec_cmd(config->cmdTimeoutSec,
                 "ovs-ofctl", "-O", config->openflow, "del-flows",
                 config->bridge, srcDrop, NULL);
        exec_cmd(config->cmdTimeoutSec,
                 "ovs-ofctl", "-O", config->openflow, "del-flows",
                 config->bridge, dstDrop, NULL);

        if (!add_flow(config, srcDrop)) {
            status_fail(status, "failed to add UE source drop flow");
            return false;
        }

        if (!add_flow(config, dstDrop)) {
            status_fail(status, "failed to add UE destination drop flow");
            return false;
        }
    }

    mark_bool(status, &status->flowsReady, true);
    return true;
}

static bool setup_policy_routing(Config *config, AppStatus *status) {

    status_set(status, InitStateSetupPolicyRouting,
               "setting up tun/ovs policy routing");

    if (!config->enablePolicyRouting) {
        mark_bool(status, &status->policyRoutingReady, true);
        return true;
    }

    if (!iface_exists(config, config->tunIf)) {
        status_fail(status, "tun interface not found for policy routing");
        mark_bool(status, &status->policyRoutingReady, false);
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

    if (!shell_ok(config, "failed to add tun policy rule",
                  "ip rule del iif %s table %d 2>/dev/null || true; "
                  "ip rule add iif %s table %d",
                  config->tunIf, config->tunTable,
                  config->tunIf, config->tunTable)) {
        status_fail(status, "failed to add tun policy rule");
        mark_bool(status, &status->policyRoutingReady, false);
        return false;
    }

    if (!shell_ok(config, "failed to add bridge policy rule",
                  "ip rule del iif %s table %d 2>/dev/null || true; "
                  "ip rule add iif %s table %d",
                  config->bridge, config->bridgeTable,
                  config->bridge, config->bridgeTable)) {
        status_fail(status, "failed to add bridge policy rule");
        mark_bool(status, &status->policyRoutingReady, false);
        return false;
    }

    mark_bool(status, &status->policyRoutingReady, true);
    return true;
}

bool ovs_reconcile(Config *config, AppStatus *status) {

    bool ok;

    if (config == NULL || status == NULL) return false;

    ok = true;

    if (!setup_forwarding(config, status)) ok = false;
    if (ok && !setup_gateway(config, status)) ok = false;
    if (ok && !setup_flows(config, status)) ok = false;
    if (ok && !setup_policy_routing(config, status)) ok = false;

    if (ok) {
        status_set(status, InitStateReady, "ready");
    }

    return ok;
}

bool ovs_setup(Config *config, AppStatus *status) {

    if (config == NULL || status == NULL) return false;

    if (!check_tools(config, status))       return false;
    if (!start_ovs(config, status))         return false;
    if (!setup_epc_aliases(config, status)) return false;
    if (!setup_tun(config, status))         return false;
    if (!setup_bridge(config, status))      return false;
    if (!setup_forwarding(config, status))  return false;
    if (!setup_gateway(config, status))     return false;
    if (!setup_flows(config, status))       return false;

    if (config->enablePolicyRouting && iface_exists(config, config->tunIf)) {
        if (!setup_policy_routing(config, status)) return false;
    }

    status_set(status, InitStateReady, "ready");
    return true;
}
