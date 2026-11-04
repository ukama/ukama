"""P0 KPI and scenario regression coverage through the production C checks.

Requires libcurl and Jansson development files; accepts CPPFLAGS/LDFLAGS/LDLIBS.
SPDX-License-Identifier: MPL-2.0
"""
import os
from pathlib import Path
import shlex
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]


class P0Checks(unittest.TestCase):
    def test_checks_and_scenarios(self):
        with tempfile.TemporaryDirectory(prefix='ulab-p0-checks-') as directory:
            build = Path(directory)
            (build / 'version.h').write_text('#define VERSION "p0-checks-test"\n')
            binary = build / 'probe'
            subprocess.run([
                *shlex.split(os.environ.get('CC', 'cc')),
                '-std=gnu11', '-D_POSIX_C_SOURCE=200809L', '-O0',
                '-Wall', '-Wextra', '-Werror', '-Wdeclaration-after-statement',
                '-ffunction-sections', '-fdata-sections',
                *shlex.split(os.environ.get('CPPFLAGS', '')),
                '-I' + str(ROOT / 'inc'), '-I' + str(build),
                str(ROOT / 'tests/runner/p0_checks_probe.c'),
                *(str(ROOT / 'src' / name) for name in (
                    'bff.c', 'bff_queries.c', 'scenario.c', 'validate.c',
                    'scenario_webapp.c', 'world.c', 'selector.c', 'model.c',
                    'util.c', 'log.c')),
                '-Wl,--gc-sections', *shlex.split(os.environ.get('LDFLAGS', '')),
                *shlex.split(os.environ.get('LDLIBS', '-lcurl -ljansson -lm')),
                '-o', str(binary),
            ], check=True)
            scenarios = [ROOT / 'scenarios/p0' / name for name in (
                'console-kpi/network/tc-185-site-uptime.yaml',
                'software-update/tc-056-disconnect-during-update-terminal.yaml',
                'ue/pdn/service-cycle-preserves-pdn.yaml',
            )]
            subprocess.run([str(binary), *map(str, scenarios)], check=True,
                           env={**os.environ, 'ULAB_SOFTWARE_TARGET_VERSION': 'test-version'},
                           timeout=30)


if __name__ == '__main__':
    unittest.main()
