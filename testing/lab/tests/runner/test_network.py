"""Network boundary/DNS tests. No host configuration or cloud access."""
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
VPN = ROOT / 'utils/runner/vpn.sh'

class NetworkIsolation(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.root = Path(self.tmp.name)
        self.resolv = self.root / 'resolv.conf'
        self.resolv.write_text('nameserver 172.31.0.2\n')
        (self.root/'backend-hosts.txt').write_text('bff.example.test\n')
        self.env = {**os.environ, 'WORK_ROOT': str(self.root), 'ULAB_NETWORK_RESOLV_CONF': str(self.resolv),
                    'P0_LAB_NETNS':'ulab','P0_NETWORK_MODE':'vpn'}
    def tearDown(self): self.tmp.cleanup()
    def shell(self, body):
        return subprocess.run(['bash','-c','set -Eeuo pipefail\n. "$1"\n'+body,'_',str(VPN)],
                              env=self.env, capture_output=True,text=True,timeout=10)
    def test_dns_written_in_place_for_bound_processes(self):
        inode=self.resolv.stat().st_ino
        r=self.shell("export foreign_option_1='dhcp-option DNS 8.8.8.8'; script_type=up p0_vpn_hook")
        self.assertEqual(r.returncode,0,r.stderr)
        self.assertEqual(self.resolv.stat().st_ino,inode)
        self.assertIn('nameserver 8.8.8.8',self.resolv.read_text())
        self.assertIn('["8.8.8.8"]',(self.root/'containers.conf').read_text())
    def test_missing_dns_rejected(self):
        r=self.shell('script_type=up p0_vpn_hook')
        self.assertNotEqual(r.returncode,0)
    def test_loopback_dns_rejected(self):
        r=self.shell("export foreign_option_1='dhcp-option DNS 127.0.0.53'; script_type=up p0_vpn_hook")
        self.assertNotEqual(r.returncode,0)
    def test_hook_does_not_touch_host_resolver_or_routes(self):
        r=self.shell('''
ip() { echo 'unexpected host route command' >&2; return 99; }
resolvectl() { return 99; }
export foreign_option_1='dhcp-option DNS 8.8.8.8'
script_type=up p0_vpn_hook
script_type=route-up p0_vpn_hook
[[ -f "$WORK_ROOT/vpn.ready" ]]
script_type=down p0_vpn_hook
[[ ! -f "$WORK_ROOT/vpn.ready" ]]
''')
        self.assertEqual(r.returncode,0,r.stderr)
    def test_lab_exec_enters_namespace(self):
        r=self.shell('ip() { printf "%s\\n" "$@"; }; p0_lab_exec getent ahostsv4 bff.example.test')
        self.assertEqual(r.stdout.splitlines(),['netns','exec','ulab','getent','ahostsv4','bff.example.test'])
    def test_host_tunnel_is_rejected(self):
        r=self.shell('ip() { return 0; }; p0_network_probe')
        self.assertNotEqual(r.returncode,0)
    def test_dns_probe_occurs_in_lab_credentials_in_host(self):
        r=self.shell('''
p0_imds_role_check() { printf 'native-imds\\n' >>"$WORK_ROOT/calls"; }
ip() {
    [[ "$1" != link ]] || return 1
    printf '%s\\n' "$*" >>"$WORK_ROOT/calls"
}
export -f ip
# Avoid an external timeout subprocess for this command mock.
timeout() { shift; "$@"; }
touch "$WORK_ROOT/vpn.ready"
VPN_PID=$$
p0_network_probe
''')
        self.assertEqual(r.returncode,0,r.stderr)
        self.assertEqual((self.root/'calls').read_text().splitlines(),['native-imds','netns exec ulab getent ahostsv4 bff.example.test'])
    def test_transient_dns_failure_recovers_without_poisoning_worker(self):
        r=self.shell('''
n=0
p0_network_probe() { n=$((n+1)); if ((n==1)); then return 1; fi; touch "$WORK_ROOT/network.stop"; }
p0_network_snapshot() { :; }
sleep() { :; }
p0_network_monitor
[[ ! -e "$WORK_ROOT/network.failed" ]]
''')
        self.assertEqual(r.returncode,0,r.stderr)
