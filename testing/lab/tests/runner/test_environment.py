import os
from pathlib import Path
import sys
import tempfile
import unittest

sys.path.insert(0, str(Path(__file__).resolve().parents[2] / 'utils/runner'))
from environment import resolve


class Environment(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.root = Path(self.tmp.name)
        (self.root / 'scenario.yaml').write_text('version: ${ULAB_SOFTWARE_TARGET_VERSION}\n')
        self.auth = {'UKAMA_IDENTIFIER': 'test@example.test', 'UKAMA_PASSWORD': 'test-only'}
    def tearDown(self): self.tmp.cleanup()
    def test_exported_settings_override_secret_defaults_without_exporting_credentials(self):
        secret = {**self.auth, 'ULAB_SOFTWARE_TARGET_VERSION': 'old', 'ULAB_VPN_CONFIG': 'private', 'ULAB_API_TOKEN': 'private'}
        env = resolve({}, {'ULAB_SOFTWARE_TARGET_VERSION': 'new', 'UKAMA_LAB_DUMP_BFF_CURL': '1', 'ULAB_RUN_ID':'wrong'}, secret, ['scenario.yaml'], self.root)
        self.assertEqual(env, {'ULAB_SOFTWARE_TARGET_VERSION':'new', 'UKAMA_LAB_DUMP_BFF_CURL':'1'})
    def test_missing_or_empty_expansion_is_caught_before_launch(self):
        for value in ({}, {'ULAB_SOFTWARE_TARGET_VERSION':''}):
            with self.subTest(value=value), self.assertRaisesRegex(ValueError, 'scenario.yaml:1'):
                resolve({}, value, self.auth, ['scenario.yaml'], self.root)
    def test_authentication_is_required(self):
        with self.assertRaisesRegex(ValueError, 'UKAMA_IDENTIFIER'):
            resolve({}, {}, {}, [], self.root)
    def test_authentication_tokens_are_supported_without_password(self):
        self.assertEqual(resolve({}, {}, {'UKAMA_SESSION_TOKEN':'s','UKAMA_BFF_TOKEN':'b'}, [], self.root), {})
    def test_explicit_tool_bundle_rebases_executable(self):
        tools = self.root / 'tools'; tools.mkdir()
        tool = tools / 'kubectl'; tool.write_text('#!/bin/sh\n'); tool.chmod(0o700)
        result = resolve({}, {'ULAB_KUBECTL':str(tool), 'P0_WORKER_FILES_DIR':str(tools), 'P0_SOURCE_ROOT':str(self.root/'lab')}, self.auth, [], self.root)
        self.assertEqual(result['ULAB_KUBECTL'], '/opt/ukama-p0/tools/kubectl')
    def test_host_tool_path_is_not_silently_sent_to_workers(self):
        with self.assertRaisesRegex(ValueError, 'P0_WORKER_FILES_DIR'):
            resolve({}, {'ULAB_KUBECTL':'/outside/lab/kubectl'}, self.auth, [], self.root)
    def test_dry_run_does_not_need_secret(self):
        self.assertEqual(resolve({}, {}, {}, [], self.root, verify_credentials=False), {})
