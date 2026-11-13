"""traffic.sh exhaust mode, run against fake podman, timeout and iperf3.

SPDX-License-Identifier: MPL-2.0
"""
import os
from pathlib import Path
import subprocess
import tempfile
import textwrap
import unittest

ROOT = Path(__file__).resolve().parents[2]
SCRIPT = ROOT / 'scripts/traffic.sh'

PODMAN = """\
#!/bin/sh
echo "$*" >> "$FAKE_DIR/podman.log"
case "$1 $3 $4" in
    "inspect "*) echo true ;;
    "exec pkill -x") [ -f "$FAKE_DIR/iperf.pid" ] && kill "$(cat "$FAKE_DIR/iperf.pid")" 2>/dev/null ;;
    "exec cat /sys/class/net/eth0/statistics/rx_bytes")
        # First read is before the transfer, later reads are after it.
        if [ -f "$FAKE_DIR/rx.read" ]; then echo $((1000 + FAKE_DELIVERED)); else echo 1000; fi
        touch "$FAKE_DIR/rx.read" ;;
esac
exit 0
"""

# Drops timeout's options and duration; the fake iperf decides its own exit code.
TIMEOUT = """\
#!/bin/sh
while [ "$#" -gt 0 ]; do
    case "$1" in
        -k) shift 2 ;;
        *s) shift; break ;;
        *) break ;;
    esac
done
exec "$@"
"""

TRAFFIC_UE = """\
#!/bin/sh
echo $$ > "$FAKE_DIR/iperf.pid"
echo "[ ID] Interval           Transfer     Bitrate         Retr  Cwnd"
echo "[  5]   0.00-1.00   sec  $FAKE_SENT  8.38 Mbits/sec    0    124 KBytes"
case "$FAKE_IPERF" in
    done) echo "iperf Done."; exit 0 ;;
    timeout) exit 124 ;;
    hang) sleep 30 & wait ;;
esac
"""


class TrafficExhaustion(unittest.TestCase):
    def setUp(self):
        tmp = tempfile.TemporaryDirectory(prefix='ulab-traffic-exhaust-')
        self.addCleanup(tmp.cleanup)
        self.dir = Path(tmp.name)
        bin_dir = self.dir / 'bin'
        ue_scripts = self.dir / 'ue/scripts'
        run_ues = self.dir / 'run/runtime-ues'
        for d in (bin_dir, ue_scripts, run_ues):
            d.mkdir(parents=True)
        for path, body in ((bin_dir / 'podman', PODMAN),
                           (bin_dir / 'timeout', TIMEOUT),
                           (ue_scripts / 'traffic-ue.sh', TRAFFIC_UE)):
            path.write_text(body)
            path.chmod(0o755)
        (run_ues / 'ue-1.env').write_text(textwrap.dedent(f"""\
            UE_CONTAINER=ue-001
            MEDIA_CONTAINER=media-001
            MEDIA_IP=10.241.1.6
            IMSI=001010000000001
            UE_DIR={self.dir / 'ue'}
            """))
        self.env = {**os.environ, 'FAKE_DIR': str(self.dir),
                    'PATH': f"{bin_dir}:{os.environ['PATH']}",
                    'ULAB_TRAFFIC_EXHAUST_GRACE_SEC': '1'}

    # 1 MB of iperf3 payload reached media as 1097309 bytes in batch 150953z.
    def run_traffic(self, iperf, sent, *extra, mb='1', delivered=1097309):
        env = {**self.env, 'FAKE_IPERF': iperf, 'FAKE_SENT': sent,
               'FAKE_DELIVERED': str(delivered)}
        return subprocess.run([os.environ.get('ULAB_TEST_SH', 'sh'), str(SCRIPT), 'ue-1', mb,
                               str(self.dir / 'run'), *extra],
                              env=env, capture_output=True, text=True,
                              timeout=30)

    def podman_calls(self):
        log = self.dir / 'podman.log'
        return log.read_text() if log.exists() else ''

    def test_quota_cut_after_full_allowance_passes_and_resets_iperf(self):
        result = self.run_traffic('hang', '1.00 MBytes', 'exhaust')
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn('traffic-exhausted ue=ue-1', result.stdout)
        self.assertIn('delivered_bytes=1097309', result.stdout)
        self.assertIn('exec ue-001 pkill -x iperf3', self.podman_calls())
        self.assertIn('exec media-001 sh -c', self.podman_calls())
        self.assertFalse(list((self.dir / 'run/runtime-ues').glob('traffic-*')))

    def test_backend_cutting_user_early_fails_even_if_client_sent_all(self):
        result = self.run_traffic('hang', '1.00 MBytes', 'exhaust',
                                  delivered=524288)
        self.assertEqual(result.returncode, 1, result.stdout + result.stderr)
        self.assertIn('traffic cut before the allowance was used', result.stderr)
        self.assertIn('delivered_bytes=524288', result.stderr)
        self.assertNotIn('exec media-001 sh -c', self.podman_calls())

    def test_allowance_less_ack_share_passes(self):
        result = self.run_traffic('timeout', '2048 KBytes', 'exhaust', mb='2',
                                  delivered=2097152 * 96 // 100)
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn('sent_bytes=2097152', result.stdout)

    def test_partial_delivery_fails(self):
        result = self.run_traffic('timeout', '512 KBytes', 'exhaust',
                                  delivered=550000)
        self.assertEqual(result.returncode, 1, result.stdout + result.stderr)
        self.assertIn('delivered_bytes=550000', result.stderr)

    def test_completed_transfer_does_not_restart_server(self):
        result = self.run_traffic('done', '1.00 MBytes', 'exhaust')
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn('traffic-complete ue=ue-1', result.stdout)
        self.assertNotIn('exec media-001 sh -c', self.podman_calls())

    def test_default_mode_still_fails_on_timeout(self):
        result = self.run_traffic('timeout', '1.00 MBytes')
        self.assertEqual(result.returncode, 1, result.stdout + result.stderr)
        self.assertIn('timed out', result.stderr)

    def test_unknown_mode_is_rejected(self):
        result = self.run_traffic('done', '1.00 MBytes', 'drain')
        self.assertEqual(result.returncode, 2)


if __name__ == '__main__':
    unittest.main()
