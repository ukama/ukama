"""Compile and exercise the actual C loader and validators without services.

SPDX-License-Identifier: MPL-2.0
Copyright (c) 2026-present, Ukama Inc.
"""
import json
import os
from pathlib import Path
import re
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
PROBE = r'''
#include "scenario.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
int main(int argc, char **argv) {
    scenario_t *s;
    ulab_error_t err;
    int rc;
    int i;
    int failed = 0;
    size_t p;
    size_t n;
    if (argc < 3) return 2;
    s = calloc(1, sizeof(*s));
    if (!s) return 6;
    for (i = 2; i < argc; i++) {
        memset(&err, 0, sizeof(err));
        rc = scenario_load(argv[i], s, &err);
        if (!rc) rc = scenario_validate(s, &err);
        if (!rc && !strcmp(argv[1], "execute"))
            rc = scenario_execution_supported(s, &err);
        printf("%d\t%s\t%s\n", rc, argv[i], err.msg);
        if (!rc && !strcmp(argv[1], "inspect")) {
            for (p = 0; p < s->phase_count; p++) {
                for (n = 0; n < s->phases[p].event_count; n++)
                    printf("event=%s timeout=%u\n",
                           scenario_event_name(s->phases[p].events[n].type),
                           s->phases[p].events[n].timeout_seconds);
                for (n = 0; n < s->phases[p].check_count; n++)
                    printf("check=%s timeout=%u immediate=%d expected=%s\n",
                           scenario_check_name(s->phases[p].checks[n].type),
                           s->phases[p].checks[n].timeout_seconds,
                           s->phases[p].checks[n].immediate,
                           s->phases[p].checks[n].expected);
            }
        }
        if (rc) failed = 1;
        fflush(stdout);
    }
    free(s);
    return failed;
}
'''


class WebappContract(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.tmp = tempfile.TemporaryDirectory(prefix="ukama-webapp-contract-")
        cls.directory = Path(cls.tmp.name)
        (cls.directory / "version.h").write_text('#define VERSION "contract-test"\n')
        (cls.directory / "probe.c").write_text(PROBE)
        cls.probe = cls.directory / "probe"
        sources = [ROOT / "src" / f for f in
                   ("scenario.c", "validate.c", "scenario_webapp.c", "util.c")]
        subprocess.run([
            os.environ.get("CC", "cc"), "-std=gnu11", "-D_POSIX_C_SOURCE=200809L",
            "-Wall", "-Wextra", "-Werror", "-Wdeclaration-after-statement", "-O1",
            "-I" + str(cls.directory), "-I" + str(ROOT / "inc"),
            str(cls.directory / "probe.c"), *map(str, sources), "-o", str(cls.probe)
        ], check=True, capture_output=True, text=True)
        cls.example = (ROOT / "scenarios/webapp/p0/network/wb-001-sites-online-recovery.yaml").read_text()
        cls.env = dict(os.environ,
                       ULAB_WEBAPP_AUTH_ORIGIN="http://auth.example.test",
                       ULAB_WEBAPP_EXPECTED_CYCLE_USAGE="64 MB of 1 GB used this cycle",
                       ULAB_WEBAPP_EXPECTED_TOTAL_USAGE="64 MB",
                       ULAB_SOFTWARE_CURRENT_VERSION="contract-current",
                       ULAB_SOFTWARE_TARGET_VERSION="contract-target",
                       ULAB_CONTROLLER_RESTART_REASON="Controller busy",
                       ULAB_CONTROLLER_RF_REASON="Controller busy",
                       ULAB_CONTROLLER_SERVICE_REASON="Controller busy")
        for path in (ROOT / "scenarios/webapp/p0/auth").glob("*.yaml"):
            for name in re.findall(r"\$\{([A-Z0-9_]+)\}", path.read_text()):
                cls.env.setdefault(name, ".auth/test.json" if name.endswith("_STATE") else "dashboard" if name.endswith("_SURFACE") else "absent" if name.endswith("_CONTROL") else "Independent expected value")


    @classmethod
    def tearDownClass(cls):
        cls.tmp.cleanup()

    def run_scenario(self, text, mode="validate", extra_env=None):
        scenario = self.directory / "case.yaml"
        scenario.write_text(text)
        env = {**self.env, **(extra_env or {})}
        return subprocess.run([str(self.probe), mode, str(scenario)], env=env,
                              capture_output=True, text=True, timeout=15)

    def test_legacy_p0_and_smoke_unchanged(self):
        paths = sorted((ROOT / "scenarios/p0").rglob("*.yaml"))
        paths += sorted((ROOT / "scenarios/smoke").rglob("*.yaml"))
        result = subprocess.run([str(self.probe), "validate", *map(str, paths)],
                                env=self.env, capture_output=True, text=True, timeout=180)
        self.assertTrue(paths)
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertEqual(len(result.stdout.splitlines()), len(paths))

    def test_examples_have_immediate_checks(self):
        for path in sorted((ROOT / "scenarios/webapp").rglob("*.yaml")):
            with self.subTest(path=path.name):
                result = self.run_scenario(path.read_text(), "inspect")
                self.assertEqual(result.returncode, 0, result.stdout)
                checks = [line for line in result.stdout.splitlines() if line.startswith("check=")]
                self.assertTrue(checks)
                self.assertTrue(all("immediate=1" in line for line in checks))

    def test_deadlines_inherit_and_overrides_survive(self):
        text = self.example.replace("action_timeout_seconds: 30", "action_timeout_seconds: 7")
        text = text.replace("check_timeout_seconds: 30", "check_timeout_seconds: 11")
        text = text.replace("        timeout_seconds: 900\n", "", 1)
        result = self.run_scenario(text, "inspect")
        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn("event=web_open timeout=7", result.stdout)
        self.assertIn("check=web_kpi_equals timeout=11 immediate=1", result.stdout)
        self.assertEqual(result.stdout.count("timeout=900"), 2)

    def test_execution_gates_incomplete_handlers_without_blocking_skips(self):
        for status in ("wip", "active", "skip", "xfail"):
            with self.subTest(status=status):
                result = self.run_scenario(self.example.replace("status: active", "status: " + status), "execute")
                if status in ("wip", "active", "skip"):
                    self.assertEqual(result.returncode, 0, result.stdout)
                else:
                    self.assertNotEqual(result.returncode, 0)
                    self.assertIn("xfail", result.stdout)
        smoke = ROOT / "scenarios/webapp/p0/session/wb-000-authenticated-members.yaml"
        result = self.run_scenario(smoke.read_text(), "execute")
        self.assertEqual(result.returncode, 0, result.stdout)

    def test_reject_invalid_contracts(self):
        cases = [
            ("v1 browser scenario", "version: 2", "version: 1"),
            ("unsupported version", "version: 2", "version: 99"),
            ("overflow version", "version: 2", "version: 4294967298"),
            ("missing base URL", "  base_url: http://localhost:3000\n", ""),
            ("missing auth", "  auth_state: .auth/owner.json\n", ""),
            ("unknown browser", "browser: chromium", "browser: magic"),
            ("invalid bool", "headless: false", "headless: maybe"),
            ("invalid deadline", "action_timeout_seconds: 30", "action_timeout_seconds: 0"),
            ("negative deadline", "action_timeout_seconds: 30", "action_timeout_seconds: -1"),
            ("overflow deadline", "action_timeout_seconds: 30", "action_timeout_seconds: 4294967297"),
            ("overall deadline", "scenario_timeout_seconds: 3600", "scenario_timeout_seconds: 10"),
            ("unknown web field", "  browser: chromium", "  browzer: chromium"),
            ("duplicate web field", "  browser: chromium", "  browser: chromium\n  browser: firefox"),
            ("wrong setup driver", "create_via_webapp: [networks, sites]", "create_via_bff: [networks, sites]"),
            ("unknown setup item", "[networks, sites]", "[networks, sties]"),
            ("duplicate setup item", "[networks, sites]", "[networks, sites, sites]"),
            ("missing setup item", "[networks, sites]", "[networks]"),
            ("invalid view", "view: network_home", "view: made_up_view"),
            ("invalid network ref", "networks: net-001", "networks: net-002"),
            ("multiple network selector", "networks: net-001", "networks: all"),
            ("invalid fault reference", "tower-site-002-001", "tower-site-999-001"),
            ("missing expected", "        expected: \"2/2\"\n", ""),
            ("missing requirement", "        requirement: WEB-NET-001\n", ""),
            ("invalid requirement", "WEB-NET-001", "BFF-NET-001"),
            ("unknown check", "type: web_kpi_equals", "type: web_magic_check"),
            ("deferred browser check", "        label: Sites online", "        label: Sites online\n        immediate: false"),
            ("BFF operator action", "type: disconnect_nodes", "type: restart_site"),
            ("credential URL", "http://localhost:3000", "http://user:secret@localhost:3000"),
            ("empty URL host", "http://localhost:3000", "http:///console"),
            ("world overflow", "  networks: 1", "  networks: 4294967297"),
            ("duplicate check value", "        expected: \"2/2\"", "        expected: \"2/2\"\n        expected: \"3/3\""),
            ("duplicate phase", "name: recovery", "name: baseline"),
        ]
        for title, old, new in cases:
            with self.subTest(case=title):
                self.assertIn(old, self.example)
                result = self.run_scenario(self.example.replace(old, new, 1))
                self.assertNotEqual(result.returncode, 0, result.stdout)

    def test_bff_check_cannot_count_as_browser_acceptance(self):
        text = self.example[:self.example.index("    checks:")]
        text += """    checks:
      - type: kpi_value
        key: SITES_ONLINE
        expected_value: 2
"""
        result = self.run_scenario(text)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("acceptance checks must use web_*", result.stdout)

    def test_event_only_phase_and_reload_are_preserved(self):
        text = self.example.replace("  - name: baseline", """  - name: select
    events:
      - type: web_select_network
        networks: net-001
      - type: web_reload
  - name: baseline""")
        result = self.run_scenario(text, "inspect")
        self.assertEqual(result.returncode, 0, result.stdout)
        self.assertIn("event=web_select_network", result.stdout)
        self.assertIn("event=web_reload", result.stdout)

    def test_explicit_empty_text_and_zero_count(self):
        text = self.example.replace('expected: "2/2"', 'expected: ""', 1)
        self.assertEqual(self.run_scenario(text).returncode, 0)
        text = self.example.replace("type: web_kpi_equals", "type: web_table_count_equals", 1)
        text = text.replace('expected: "2/2"', "expected_count: 0", 1)
        self.assertEqual(self.run_scenario(text).returncode, 0)

    def test_unavailable_action_is_a_valid_expected_state(self):
        text = self.example.replace("type: web_kpi_equals", "type: web_action_available", 1)
        text = text.replace('expected: "2/2"', "available: false", 1)
        self.assertEqual(self.run_scenario(text).returncode, 0)
        self.assertNotEqual(self.run_scenario(text.replace("available: false", "available: no")).returncode, 0)

    def test_environment_expansion_is_explicit(self):
        text = self.example.replace("http://localhost:3000", "${ULAB_CONTRACT_TEST_URL}")
        env = dict(self.env)
        env.pop("ULAB_CONTRACT_TEST_URL", None)
        path = self.directory / "env-case.yaml"
        path.write_text(text)
        result = subprocess.run([str(self.probe), "validate", str(path)], env=env,
                                capture_output=True, text=True, timeout=15)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("missing environment variable", result.stdout)
        self.assertEqual(self.run_scenario(text, extra_env={"ULAB_CONTRACT_TEST_URL": "http://localhost:3000"}).returncode, 0)

    def test_orphan_phase_sections_fail_without_crashing(self):
        text = "version: 2\nphases:\n    checks:\n      - type: web_kpi_equals\n"
        result = self.run_scenario(text)
        self.assertGreater(result.returncode, 0, result.stdout + result.stderr)

    def test_world_expectations_reject_invalid_refs_properties_and_literal_conflicts(self):
        text = (ROOT / "scenarios/webapp/p0/network/wb-002-node-list-detail.yaml").read_text()
        self.assertEqual(self.run_scenario(text).returncode, 0)
        for modified in (
            text.replace("expected_ref: tower-site-001-001", "expected_ref: tower-site-999-001"),
            text.replace("expected_property: id", "expected_property: unknown"),
            text.replace("expected_property: id", 'expected_property: id\n        expected: "same"'),
            text.replace("nodes: tower-site-001-001", "nodes: all"),
        ):
            self.assertNotEqual(self.run_scenario(modified).returncode, 0)

    def test_operation_contract_rejects_unsafe_or_unrelated_fields(self):
        text = (ROOT / "scenarios/webapp/p0/operations/wb-012-software-update.yaml").read_text()
        self.assertEqual(self.run_scenario(text).returncode, 0)
        for old, new in (
            ('action: "update_software"', 'action: "made_up_operation"'),
            ('tag: "${ULAB_SOFTWARE_TARGET_VERSION}"', 'unknown_tag: "v2"'),
            ('action: "open_software"', 'action: "set_radio"'),
            ('app: "example"', 'app: ""'),
            ('label: "Current version"', 'label: "Current version"\n        match: contains'),
            ('action: "open_software"', 'action: "open_software"\n        value: nope'),
            ('nodes: "controller-site-001-001"', 'nodes: "controller-site-999-001"'),
        ):
            with self.subTest(old=old, new=new):
                self.assertNotEqual(self.run_scenario(text.replace(old, new)).returncode, 0)

    def test_site_confirmation_source_and_tab_names_are_closed_contracts(self):
        text = (ROOT / "scenarios/webapp/p0/operations/wb-011-site-operations.yaml").read_text()
        self.assertEqual(self.run_scenario(text).returncode, 0)
        for replacement in ('value_from: "anything"', 'value_from: "site_name"\n        value: explicit'):
            self.assertNotEqual(self.run_scenario(text.replace('value_from: "site_name"', replacement)).returncode, 0)
        text = (ROOT / "scenarios/webapp/p0/operations/wb-015-stale-confirmation.yaml").read_text()
        self.assertEqual(self.run_scenario(text).returncode, 0)
        self.assertNotEqual(self.run_scenario(text.replace('tab: "secondary"', 'tab: "third"')).returncode, 0)

    def test_commerce_duration_units_references_and_fields_are_strict(self):
        text = (ROOT / "scenarios/webapp/p0/commerce/wb-020-plan-unit-roundtrip.yaml").read_text()
        for old, new in [
            ("duration_days: 1", "duration_days: 2"),
            ("duration_minutes: 10080", "duration_minutes: 10079"),
            ("duration_days: 1", "duration_days: 1\n    duration_minutes: 1440"),
            ("data_mb: 1024", "data_mb: 1000"),
            ("amount: 10", "amount: -1"),
            ("package: weekly__net-001", "package: missing__net-001"),
            ("unit: GB", "unit: GiB"),
            ("unit: GB", "unit: GB\n        unit: GB"),
            ("action: create_plan", "action: purchase_by_api"),
        ]:
            with self.subTest(new=new):
                self.assertNotEqual(self.run_scenario(text.replace(old,new,1)).returncode,0)

    def test_commerce_ue_scope_and_runtime_order_contract(self):
        text = (ROOT / "scenarios/webapp/p0/commerce/wb-021-customer-sim-allocation.yaml").read_text()
        for old,new in [("ues: ue-000001", "ues: ue-000099"),
                        ("sims_per_network: 1", "sims_per_network: 101"),
                        ("action: allocate_sim", "action: allocate_sim\n        unit: GB"),
                        ("expected_property: \"iccid\"", "expected_property: \"password\"")]:
            with self.subTest(new=new):
                self.assertNotEqual(self.run_scenario(text.replace(old,new,1)).returncode,0)
        text = (ROOT / "scenarios/webapp/p0/commerce/wb-024-ue-usage.yaml").read_text()
        for old,new in [("amount_mb: 64", "amount_mb: 0"),
                        ("timeout_seconds: 900", "timeout_seconds: 901"),
                        ("ues_per_site: 1", "ues_per_site: 0")]:
            with self.subTest(new=new):
                self.assertNotEqual(self.run_scenario(text.replace(old,new,1)).returncode,0)

    def test_ui_interactions_reject_unknown_actions_labels_keys_and_refs(self):
        text = (ROOT / "scenarios/webapp/p0/expanded/wb-037-keyboard-dialog.yaml").read_text()
        for old, new in [("action: open_form", "action: submit_arbitrary"),
                         ("value: Escape", "value: Enter"),
                         ("label: Dialog title", "label: Whatever"),
                         ("value: keyboard", "value: something"),
                         ("action: press", "action: press\n        selector: '#secret'"),
                         ("subject: Create plan", "subject: ''"),
                         ("label: Dialog title", "label: Dialog title\n        package: missing__net-001")]:
            with self.subTest(new=new):
                self.assertNotEqual(self.run_scenario(text.replace(old,new,1)).returncode,0)

    def test_session_profiles_actions_subjects_and_modes_are_closed(self):
        text = (ROOT / "scenarios/webapp/p0/auth/wb-050-account-settings.yaml").read_text()
        self.assertEqual(self.run_scenario(text, "execute").returncode, 0)
        for old,new in [("session_mode: auth_test", "session_mode: maybe"),
                        ("session_mode: auth_test", "session_mode: authenticated"),
                        ("auth_origin: ${ULAB_WEBAPP_AUTH_ORIGIN}", "auth_origin: http://localhost:3000"),
                        ("auth_origin: ${ULAB_WEBAPP_AUTH_ORIGIN}", "auth_origin: http://auth.test/path"),
                        ("auth_origin: ${ULAB_WEBAPP_AUTH_ORIGIN}", "auth_origin: http://name:secret@auth.test"),
                        ("networks: 0", "networks: 1"),
                        ("action: navigate", "action: arbitrary_script"),
                        ("value: /business/settings", "value: //other.test"),
                        ("value: /business/settings", "value: /api/auth/logout"),
                        ("subject: Full name", "subject: Password"),
                        ("label: Surface", "label: Cookie value"),
                        ("action: reload", "action: reload\n        value: /business"),
                        ("type: web_session_equals", "type: web_ui_equals"),
                        ("action: open_account", "action: open_account\n        selector: html")]:
            with self.subTest(new=new):
                self.assertNotEqual(self.run_scenario(text.replace(old,new,1)).returncode,0)

    def test_remaining_roadmap_accounts_for_the_frozen_gap(self):
        catalog = json.loads((ROOT / "docs/webapp/coverage.json").read_text())
        plan = json.loads((ROOT / "docs/webapp/remaining-patches.json").read_text())
        ids = [rid for patch in plan['patches'] for rid in patch['requirements']]
        self.assertEqual(len(ids),96)
        self.assertEqual(len(set(ids)),96)
        requirements = {r['id']:r for r in catalog['requirements']}
        self.assertEqual(len(requirements),132)
        self.assertEqual(sum(r['priority']=='p0' for r in requirements.values()),86)
        self.assertEqual({r['id'] for r in catalog['requirements'] if r['automation']!='implemented'}, set(ids)-{'WEB-AUTH-'+f'{n:03}' for n in range(2,8)}-{'WEB-TEAM-004','WEB-UI-012'})
        self.assertEqual(requirements['WEB-SHELL-010']['automation'],'planned')
        self.assertEqual(plan['qualification_patch'],17)

    def test_coverage_inventory_has_no_unearned_credit(self):
        catalog = json.loads((ROOT / "docs/webapp/coverage.json").read_text())
        identifiers = {r["id"] for r in catalog["requirements"]}
        self.assertEqual(len(identifiers), len(catalog["requirements"]))
        for r in catalog["requirements"]:
            self.assertEqual(r["automation"], "implemented" if r["id"] in {"WEB-AUTH-002", "WEB-AUTH-003", "WEB-AUTH-004", "WEB-AUTH-005", "WEB-AUTH-006", "WEB-AUTH-007", "WEB-TEAM-004", "WEB-UI-012", "WEB-AUTH-001", "WEB-SHELL-001", "WEB-SHELL-008", "WEB-SHELL-009", "WEB-PLAN-006", "WEB-PLAN-009", "WEB-PLAN-010", "WEB-BIZ-001", "WEB-BIZ-008", "WEB-NET-001", "WEB-NET-002", "WEB-NET-003", "WEB-NODE-001", "WEB-NODE-002", "WEB-NODE-003", "WEB-NODE-005", "WEB-NODE-007", "WEB-OPS-001", "WEB-OPS-002", "WEB-OPS-003", "WEB-OPS-005", "WEB-OPS-006", "WEB-OPS-007", "WEB-OPS-009", "WEB-OPS-010", "WEB-CUSTOMER-002", "WEB-PAY-001", "WEB-PAY-002", "WEB-PAY-008", "WEB-PLAN-001", "WEB-PLAN-002", "WEB-PLAN-003", "WEB-PLAN-004", "WEB-PLAN-005", "WEB-PLAN-012", "WEB-SIM-001"} else "planned")
            self.assertEqual(r["verification"], "not_run")
            for path in r["scenarios"]:
                self.assertTrue((ROOT / path).is_file())
        for path in (ROOT / "scenarios/webapp").rglob("*.yaml"):
            for identifier in re.findall(r'requirement:\s+"?(WEB-[A-Z0-9-]+)', path.read_text()):
                self.assertIn(identifier, identifiers)


if __name__ == "__main__":
    unittest.main()
