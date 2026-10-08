"""Scenario-owned catalog counts through the actual C loader and BFF client.

Requires a C compiler, libcurl and Jansson development files. CPPFLAGS,
LDFLAGS and LDLIBS can point to the lab's vendor dependencies.
SPDX-License-Identifier: MPL-2.0
"""
import json
import os
from pathlib import Path
import shlex
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
AFFECTED = {1: 1, 2: 1, 3: 1, 4: 1, 80: 16, 121: 3, 122: 3,
            123: 3, 124: 3, 125: 3, 316: 1, 317: 1, 318: 1, 319: 1}


def catalog(ids, foreign=12):
    return {'data': {'getPackages': {'packages': [
        *({'uuid': value, 'name': 'renamed plan', 'active': False}
          for value in ids),
        *({'uuid': f'foreign-{i}', 'name': 'renamed plan', 'networkId': ''}
          for i in range(foreign)),
    ]}}}


class CatalogCount(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.build = tempfile.TemporaryDirectory(prefix='ulab-catalog-build-')
        cls.addClassCleanup(cls.build.cleanup)
        build = Path(cls.build.name)
        (build / 'version.h').write_text('#define VERSION "catalog-test"\n')
        cls.binary = build / 'probe'
        subprocess.run([
            *shlex.split(os.environ.get('CC', 'cc')),
            '-std=gnu11', '-D_POSIX_C_SOURCE=200809L', '-O0',
            '-Wall', '-Wextra', '-Werror', '-Wdeclaration-after-statement',
            '-ffunction-sections', '-fdata-sections',
            *shlex.split(os.environ.get('CPPFLAGS', '')),
            '-I' + str(ROOT / 'inc'), '-I' + str(build),
            str(ROOT / 'tests/runner/catalog_count_probe.c'),
            *(str(ROOT / 'src' / f) for f in (
                'bff.c', 'bff_queries.c', 'scenario.c', 'validate.c',
                'scenario_webapp.c', 'world.c', 'selector.c', 'util.c', 'log.c')),
            '-Wl,--gc-sections', *shlex.split(os.environ.get('LDFLAGS', '')),
            *shlex.split(os.environ.get('LDLIBS', '-lcurl -ljansson -lm')),
            '-o', str(cls.binary),
        ], check=True)

    def setUp(self):
        tmp = tempfile.TemporaryDirectory(prefix='ulab-catalog-')
        self.addCleanup(tmp.cleanup)
        self.directory = Path(tmp.name)

    def run_probe(self, scenario, response, mode='catalog'):
        fixture = self.directory / 'response.json'
        fixture.write_text(json.dumps(response))
        result = subprocess.run(
            [str(self.binary), str(scenario), str(fixture), mode],
            env={**os.environ, 'ULAB_RESILIENCE_NAME_SUFFIX': 'catalog-test'},
            capture_output=True, text=True, timeout=15)
        self.assertEqual(result.stderr, '', result.stderr)
        return result

    def scenario(self, count=1, expected=None, target='packages', networks=1,
                 scope='network'):
        path = self.directory / 'scenario.yaml'
        packages = ''.join(f'''
  - ref: p{i}
    name: Original Plan {i}
    scope: {scope}
    data_mb: 100
    duration_days: 1
    amount: 1.01
    currency: USD
    country: USA
    assign_percent: {100 if i == 0 else 0}
''' for i in range(count))
        path.write_text(f'''version: 1
name: catalog-probe
seed: 1774
provider:
  type: virtual
world:
  networks: {networks}
  sites_per_network: 0
  ues_per_site: 0
setup:
  create_via_bff: [networks, packages]
packages:{packages if packages else ' []'}
phases:
  - name: count
    checks:
      - type: list_count_equals
        target: {target}
        networks: net-001
        expected_count: {count if expected is None else expected}
''')
        return path

    def assert_passes(self, result, expected):
        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn(f'expected={expected} actual={expected}', result.stdout)
        self.assertIn('REQUEST getPackages {"networkId":"network-0"}', result.stdout)

    def test_all_14_reported_scenarios_ignore_12_foreign_plans(self):
        scenarios = sorted(p for p in (ROOT / 'scenarios/resilience').rglob('r*.yaml')
                           if int(p.name[1:4]) in AFFECTED)
        self.assertEqual(len(scenarios), 14)
        for scenario in scenarios:
            expected = AFFECTED[int(scenario.name[1:4])]
            with self.subTest(scenario=scenario.name):
                result = self.run_probe(scenario, catalog(
                    [f'owned-{i}' for i in range(expected)]))
                self.assert_passes(result, expected)

    def test_both_aliases_count_ids_despite_rename_and_deactivation(self):
        for target in ('packages', 'plans'):
            with self.subTest(target=target):
                result = self.run_probe(self.scenario(target=target),
                                        catalog(['owned-0']))
                self.assert_passes(result, 1)

    def test_missing_owned_plan_is_not_replaced_by_foreign_plans(self):
        for count in (1, 3, 16):
            with self.subTest(count=count):
                result = self.run_probe(self.scenario(count=count), catalog(
                    [f'owned-{i}' for i in range(count - 1)]))
                self.assertNotEqual(result.returncode, 0)
                self.assertIn(f'expected={count} actual={count - 1}', result.stdout)

    def test_duplicate_row_cannot_replace_missing_owned_id(self):
        result = self.run_probe(self.scenario(count=2),
                                catalog(['owned-0', 'owned-0']))
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('expected=2 actual=1', result.stdout)

    def test_empty_scenario_catalog_ignores_foreign_plans(self):
        self.assert_passes(self.run_probe(self.scenario(count=0), catalog([])), 0)

    def test_organization_scoped_owned_plan_is_counted(self):
        self.assert_passes(self.run_probe(
            self.scenario(scope='organization', networks=2), catalog(['owned-0'])), 1)

    def test_own_packages_in_other_networks_are_not_assumed_visible(self):
        self.assert_passes(self.run_probe(
            self.scenario(networks=2), catalog(['owned-0'])), 1)

    def test_leaked_own_package_from_other_network_still_fails(self):
        result = self.run_probe(self.scenario(networks=2),
                                catalog(['owned-0', 'owned-1']))
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('expected=1 actual=2', result.stdout)

    def test_unrecorded_package_id_is_an_execution_error(self):
        result = self.run_probe(self.scenario(), catalog(['owned-0']), 'uncreated')
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('package list count requires id for package', result.stdout)
        self.assertNotIn('REQUEST', result.stdout)

    def test_missing_list_and_transport_errors_are_not_empty_catalogs(self):
        for response, message in (
            ({'data': {'getPackages': {}}}, 'getPackages missing packages list'),
            ({'data': {'getPackages': {'packages': None}}},
             'getPackages missing packages list'),
            ({'transport_error': True}, 'fixture transport failed'),
        ):
            with self.subTest(response=response):
                result = self.run_probe(self.scenario(count=0), response)
                self.assertNotEqual(result.returncode, 0)
                self.assertIn(message, result.stdout)

    def test_non_catalog_list_counts_are_unchanged(self):
        response = {'data': {'getSubscribersByNetwork': {'subscribers': [
            {'id': 'customer-0'}, {'id': 'customer-1'}]}}}
        result = self.run_probe(self.scenario(count=0, target='customers', expected=2),
                                response, 'all')
        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn('target=customers expected=2 actual=2', result.stdout)
        self.assertIn('REQUEST getSubscribersByNetwork', result.stdout)


if __name__ == '__main__':
    unittest.main()
