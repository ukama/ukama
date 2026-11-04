"""Offline image-version and local/AWS preparation regressions."""
import contextlib
import io
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch


LAB = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(LAB / "utils/runner"))
import software_version as software
from common import RunnerError


class SoftwareVersion(unittest.TestCase):
    def setUp(self):
        temporary = tempfile.TemporaryDirectory()
        self.addCleanup(temporary.cleanup)
        self.root = Path(temporary.name)
        self.images = self.root / "images.json"
        self.calls = self.root / "calls.jsonl"
        self.tags = ["testing/virtualnode-base:" + kind + "-starter"
                     for kind in software.KINDS]
        self.state = {tag: {"Id": "sha256:" + str(i) * 64,
                            "output": "v0.1.0-698199a0c\n"}
                      for i, tag in enumerate(self.tags, 1)}
        self.images.write_text(json.dumps(self.state))
        bindir = self.root / "bin"
        bindir.mkdir()
        podman = bindir / "podman"
        podman.write_text("#!" + sys.executable + '''
import json, os, sys
from pathlib import Path
args = sys.argv[1:]
with open(os.environ['VERSION_CALLS'], 'a') as stream:
    print(json.dumps(args), file=stream)
images = json.loads(Path(os.environ['VERSION_IMAGES']).read_text())
if args[:2] == ['image', 'exists']:
    sys.exit(0 if args[2] in images else 1)
if args[:2] == ['image', 'inspect']:
    print(json.dumps([images[args[2]]]))
elif args[0] == 'run':
    image_id = args[args.index('/usr/bin/python3') + 1]
    image = next(value for value in images.values() if value['Id'] == image_id)
    print(json.dumps({'output': image['output'],
                      'source': '/sbin/starter.d (embedded VERSION)'}))
elif args[:2] != ['rm', '-f']:
    sys.exit(99)
''')
        podman.chmod(0o755)
        env = {"PATH": str(bindir) + ":" + os.environ["PATH"],
               "VERSION_IMAGES": str(self.images),
               "VERSION_CALLS": str(self.calls)}
        self.environment = patch.dict(os.environ, env)
        self.environment.start()
        self.addCleanup(self.environment.stop)
        self.repo = self.root / "repo"
        self.repo.mkdir()

    def resolve(self, env=None, **kwargs):
        with contextlib.redirect_stderr(io.StringIO()):
            return software.resolve(LAB, self.repo, env or {}, **kwargs)

    def test_image_version_overrides_stale_export_and_preserves_targets(self):
        env = {software.CURRENT: "v0.1.0-wrong-repo-head",
               "ULAB_SOFTWARE_TARGET_VERSION": "1.2.3-target",
               "ULAB_SOFTWARE_NEXT_VERSION": "v1.2.4-next",
               "ULAB_SOFTWARE_UNAVAILABLE_VERSION": "99.99.99-missing"}
        evidence = self.resolve(env)
        self.assertEqual(env[software.CURRENT], "v0.1.0-698199a0c")
        self.assertEqual(env["ULAB_SOFTWARE_TARGET_VERSION"], "1.2.3-target")
        self.assertEqual(env["ULAB_SOFTWARE_NEXT_VERSION"], "v1.2.4-next")
        self.assertEqual(env["ULAB_SOFTWARE_UNAVAILABLE_VERSION"],
                         "99.99.99-missing")
        self.assertEqual(len(evidence["images"]), 3)
        self.assertEqual(evidence["app"], "starter.d")
        self.assertEqual(evidence["images"][0]["source"],
                         "/sbin/starter.d (embedded VERSION)")
        lines = self.calls.read_text().splitlines()
        calls = [json.loads(line) for line in lines]
        probes = [c for c in calls if c[0] == "run"]
        self.assertEqual(len(probes), 3)
        for probe in probes:
            self.assertIn("--network=none", probe)
            self.assertIn("--read-only", probe)
            self.assertEqual(probe[probe.index("-c") + 2], "starter.d")
            self.assertTrue(any(v.startswith("sha256:") for v in probe))
        self.assertEqual(len([c for c in calls if c[:2] == ["rm", "-f"]]), 3)

    def test_mixed_versions_fail_without_changing_current(self):
        self.state[self.tags[2]]["output"] = "Version: v0.1.0-different"
        self.images.write_text(json.dumps(self.state))
        env = {software.CURRENT: "old"}
        with self.assertRaisesRegex(RunnerError, "different app versions"):
            self.resolve(env)
        self.assertEqual(env[software.CURRENT], "old")

    def test_missing_or_ambiguous_binary_version_fails(self):
        for output in ("", "help: use --version", "hash=1234567",
                       "Version: v1.0.0-a\nVersion: v2.0.0-b"):
            with self.subTest(output=output):
                with self.assertRaises(RunnerError):
                    software.parse_version(output)

    def test_preserves_exact_leading_v(self):
        for value in ("v0.1.0-698199a0c", "1.2.3-accdefgh"):
            self.assertEqual(software.parse_version(value), value)
            self.assertEqual(software.parse_version("Version: " + value), value)

    def test_generated_version_header(self):
        header = ('#ifndef VERSION_H_\n#define VERSION_H_\n\n'
                  '#define VERSION "v0.1.0-698199a0c"\n\n#endif\n')
        self.assertEqual(software.parse_version(header), "v0.1.0-698199a0c")
        with self.assertRaises(RunnerError):
            software.parse_version('#define VERSION "latest"\n')

    def test_aws_checks_exported_image_ids(self):
        expected = {tag: value["Id"] for tag, value in self.state.items()}
        self.assertEqual(self.resolve(expected_images=expected)["version"],
                         "v0.1.0-698199a0c")
        expected[self.tags[0]] = "sha256:wrong"
        with self.assertRaisesRegex(RunnerError, "changed since export"):
            self.resolve(expected_images=expected)

    def test_prepare_exports_and_saves_aws_environment(self):
        scenario = self.root / "scenario.yaml"
        scenario.write_text("version: ${ULAB_SOFTWARE_CURRENT_VERSION}\n")
        output = self.root / "software-version.json"
        env = {}
        with contextlib.redirect_stderr(io.StringIO()):
            self.assertTrue(software.prepare(LAB, self.repo, [scenario], env,
                                             output))
        self.assertEqual(env[software.CURRENT],
                         json.loads(output.read_text())["version"])
        scenario.write_text("version: 1\n")
        with patch.object(software, "resolve") as resolve:
            self.assertFalse(software.prepare(LAB, self.repo, [scenario], {},
                                              output))
            resolve.assert_not_called()

    def test_missing_base_is_built_before_reading_version(self):
        self.state.pop(self.tags[0])
        self.images.write_text(json.dumps(self.state))
        lab = self.root / "lab"
        (lab / "scripts").mkdir(parents=True)
        build = lab / "scripts/build-node.sh"
        build.write_text("#!" + sys.executable + '''
import json, os
from pathlib import Path
p = Path(os.environ['VERSION_IMAGES'])
images = json.loads(p.read_text())
images['testing/virtualnode-base:anode-starter'] = {
    'Id': 'sha256:' + '1' * 64, 'output': 'Version: v0.1.0-698199a0c'}
p.write_text(json.dumps(images))
''')
        build.chmod(0o755)
        nodes = {"anode": "a-anode-1"}
        with patch("artifacts.factory_ids", return_value=nodes):
            with contextlib.redirect_stderr(io.StringIO()):
                evidence = software.resolve(lab, self.repo, {})
        self.assertEqual(evidence["version"], "v0.1.0-698199a0c")

    def test_probe_executes_binary_with_version_argument(self):
        binary = self.root / "example"
        binary.write_text('#!/bin/sh\n[ "$1" = --version ] || exit 2\n'
                          'echo "Version: v0.1.0-698199a0c"\n')
        binary.chmod(0o755)
        result = subprocess.run([sys.executable, "-c", software.PROBE,
                                 "example", str(binary), ""],
                                text=True, capture_output=True, check=True)
        self.assertEqual(software.parse_version(
            json.loads(result.stdout)["output"]),
                         "v0.1.0-698199a0c")

    def test_version_file_fallback(self):
        version_file = self.root / "version"
        version_file.write_text("v0.1.0-698199a0c\n")
        result = subprocess.run([sys.executable, "-c", software.PROBE,
                                 "example", "", str(version_file)],
                                text=True, capture_output=True, check=True)
        self.assertEqual(software.parse_version(
            json.loads(result.stdout)["output"]),
                         "v0.1.0-698199a0c")

    def starter_probe(self, binary="", version_file=""):
        # Map only image paths into an isolated fixture filesystem. Execute
        # the actual probe, rather than mocking its version output.
        probe = software.PROBE.replace(
            "/conf/starterd/version.h", str(self.root / "version.h"))
        probe = probe.replace("/sbin/starter.d", str(self.root / "starter.d"))
        return subprocess.run([sys.executable, "-c", probe, "starter.d",
                               str(binary), str(version_file)], text=True,
                              capture_output=True, timeout=10)

    def compile_starter(self, versions):
        compiler = shutil.which("cc")
        if not compiler:
            self.skipTest("C compiler required for real ELF probe fixture")
        source = self.root / "starter.c"
        marker = self.root / "daemon-was-started"
        source.write_text(
            '#include <stdio.h>\n'
            'const char *versions[] = {' +
            ','.join(json.dumps(v) for v in versions) + '};\n'
            'const char *noise = "libcurl/8.5.0 compiler 13.2.0";\n'
            'int main(void) {\n'
            '  FILE *f = fopen(' + json.dumps(str(marker)) + ', "w");\n'
            '  if (f) fclose(f);\n'
            '  puts(versions[0]);\n'
            '  return 0;\n}\n')
        subprocess.run([compiler, "-g", "-O0", str(source), "-o",
                        str(self.root / "starter.d")], check=True,
                       capture_output=True, text=True)
        return marker

    def test_default_reads_installed_header_without_binary(self):
        (self.root / "version.h").write_text(
            '#define VERSION "v0.1.0-698199a0c"\n')
        result = self.starter_probe()
        self.assertEqual(result.returncode, 0, result.stderr)
        record = json.loads(result.stdout)
        self.assertEqual(software.parse_version(record["output"]),
                         "v0.1.0-698199a0c")
        self.assertEqual(record["source"], str(self.root / "version.h"))

    def test_headerless_image_reads_real_elf_without_starting_daemon(self):
        marker = self.compile_starter(["v0.1.0-698199a0c"])
        result = self.starter_probe()
        self.assertEqual(result.returncode, 0, result.stderr)
        record = json.loads(result.stdout)
        self.assertEqual(record["output"], "v0.1.0-698199a0c")
        self.assertTrue(record["source"].endswith(" (embedded VERSION)"))
        self.assertFalse(marker.exists())

    def test_ambiguous_elf_versions_fail_without_starting_daemon(self):
        marker = self.compile_starter(["v1.2.3-aaa", "v1.2.3-bbb"])
        result = self.starter_probe()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("no unique embedded version", result.stderr)
        self.assertFalse(marker.exists())

    def test_missing_starter_metadata_does_not_use_example_or_repo(self):
        (self.root / "example-version").write_text("v1.2.3-unrelated\n")
        result = self.starter_probe()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("starter.d", result.stderr)

    def test_starter_script_is_never_executed(self):
        path = self.root / "starter.d"
        marker = self.root / "started"
        path.write_text('#!/bin/sh\ntouch "' + str(marker) + '"\n')
        path.chmod(0o755)
        result = self.starter_probe()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("not ELF", result.stderr)
        self.assertFalse(marker.exists())

    def test_explicit_starter_header_override(self):
        path = self.root / "custom-version.h"
        path.write_text('#define VERSION "1.2.3-custom"\n')
        result = self.starter_probe(version_file=path)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(software.parse_version(
            json.loads(result.stdout)["output"]), "1.2.3-custom")

    def test_batch_passes_detected_version_to_scenario_process(self):
        scenario_dir = self.root / "scenarios/p0/software-update"
        scenario_dir.mkdir(parents=True)
        (scenario_dir / "test.yaml").write_text(
            "version: 1\nname: version_test\nstatus: active\n"
            "expected: ${ULAB_SOFTWARE_CURRENT_VERSION}\n")
        binary = self.root / "fake-lab"
        binary.write_text("#!" + sys.executable + '''
import json, os, sys
from pathlib import Path
args = sys.argv
version = os.environ['ULAB_SOFTWARE_CURRENT_VERSION']
assert version == 'v0.1.0-698199a0c', version
assert os.environ['ULAB_SOFTWARE_TARGET_VERSION'] == 'manual-target'
out = Path(args[args.index('--out')+1]) / args[args.index('--run-id')+1]
out.mkdir(parents=True, exist_ok=True)
(out / 'report.json').write_text(json.dumps({
    'status': 'active', 'passed': True, 'duration_sec': 1, 'cleanup': 'ok'}))
print('detected-current=' + version)
''')
        binary.chmod(0o755)
        env = {**os.environ, "UKAMA_IDENTIFIER": "test",
               "UKAMA_PASSWORD": "test", "UKAMA_REPO": str(self.repo),
               "LAB_BIN": str(binary), software.CURRENT: "stale",
               "ULAB_SOFTWARE_TARGET_VERSION": "manual-target",
               "ULAB_FACTORY_NODE_COUNT": "0"}
        result = subprocess.run(
            ["bash", str(LAB / "utils/run-scenarios.sh"), "p0",
             "--batch-id", "version-test"], cwd=self.root, env=env,
            capture_output=True, text=True, timeout=15)
        self.assertEqual(result.returncode, 0, result.stderr + result.stdout)
        self.assertIn("detected-current=v0.1.0-698199a0c", result.stdout)
        self.assertTrue((self.root / "runs/p0-batches/version-test/"
                        "software-version.json").is_file())


if __name__ == "__main__":
    unittest.main()
