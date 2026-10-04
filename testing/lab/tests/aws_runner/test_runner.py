#!/usr/bin/env python3
"""Offline contract tests; no AWS account, backend, VPN or Podman daemon needed."""
import contextlib
import argparse
import copy
import gzip
import hashlib
import io
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tarfile
import tempfile
import unittest
from unittest.mock import patch

LAB = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(LAB / "utils/runner"))
import artifacts
import common
import config
import controller
import network
import worker


IMAGE_TAGS = [f"localhost/testing/virtualnode-base:{kind}-starter" for kind in ("anode", "cnode", "tnode")] + ["docker.io/library/alpine:3.20"]


def image_fixture(path, *, damaged_layer=False, missing_layer=False, empty_manifest=False):
    """Real Docker-archive bytes with shared layer tar, configs and four tags."""
    layer = io.BytesIO()
    with tarfile.open(fileobj=layer, mode="w") as archive:
        data = b"image fixture: contents survive save/compress/download/decompress/load\n"
        member = tarfile.TarInfo("etc/fixture")
        member.size = len(data)
        archive.addfile(member, io.BytesIO(data))
    content = layer.getvalue()
    diff = "sha256:" + hashlib.sha256(content).hexdigest()
    images, manifest, members = {}, [], {}
    for tag in IMAGE_TAGS:
        cfg = json.dumps({"architecture": "amd64", "os": "linux", "config": {"Env": ["TAG=" + tag]},
                          "rootfs": {"type": "layers", "diff_ids": [diff]}}).encode()
        digest = hashlib.sha256(cfg).hexdigest()
        images[tag] = digest
        members[digest + ".json"] = cfg
        manifest.append({"Config": digest + ".json", "RepoTags": [tag], "Layers": ["shared/layer.tar"]})
    if not missing_layer:
        members["shared/layer.tar"] = (b"broken" + content[6:]) if damaged_layer else content
    members["manifest.json"] = json.dumps([] if empty_manifest else manifest).encode()
    with tarfile.open(path, "w") as archive:
        for name, data in members.items():
            member = tarfile.TarInfo(name)
            member.size = len(data)
            archive.addfile(member, io.BytesIO(data))
    return images


def compressed_fixture(raw, images):
    target = raw.with_suffix(".tar.gz")
    with raw.open("rb") as source, gzip.open(target, "wb") as output:
        shutil.copyfileobj(source, output)
    return target, {"cache_version": artifacts.IMAGE_CACHE_VERSION, "images": images,
                    "sha256": common.sha256(target), "size": target.stat().st_size,
                    "tar_sha256": common.sha256(raw), "tar_size": raw.stat().st_size}


class Fixture(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.root = Path(self.tmp.name)
        self.cfg = json.loads((LAB / "utils/runner/aws.example.json").read_text())

    def tearDown(self):
        self.tmp.cleanup()

    def write(self, path, text, mode=0o644):
        file = self.root / path
        file.parent.mkdir(parents=True, exist_ok=True)
        file.write_text(text)
        file.chmod(mode)
        return file


class CLITests(Fixture):
    def setUp(self):
        super().setUp()
        for name in ("a-pass", "b-skip", "c-fail"):
            self.write(f"scenarios/p0/demo/{name}.yaml", f"version: 1\nname: {name}\nstatus: active\n")
        self.binary = self.write("bin/fake-lab", '''#!/usr/bin/env python3
import json, sys, os
from pathlib import Path
args=sys.argv[1:]
scenario=Path(args[1]).stem
run_id=args[args.index('--run-id')+1]
out=Path(args[args.index('--out')+1])/run_id
out.mkdir(parents=True,exist_ok=True)
skip='skip' in scenario
fail='fail' in scenario
report={'status':'skip' if skip else 'active','passed':not fail,'duration_sec':2,'cleanup':'ok'}
(out/'report.json').write_text(json.dumps(report))
print('scenario',scenario)
print('kubectl='+os.environ.get('ULAB_KUBECTL',''))
sys.exit(1 if fail else 0)
''', 0o755)
        (self.root / "repo").mkdir()
        self.env = {**os.environ, "UKAMA_IDENTIFIER": "test", "UKAMA_PASSWORD": "test-password",
                    "UKAMA_REPO": str(self.root / "repo"), "LAB_BIN": str(self.binary)}
        for key in list(self.env):
            if key.startswith(("ULAB_RUNNER_", "P0_")):
                del self.env[key]

    def cli(self, *args, extra=None):
        return subprocess.run(["bash", str(LAB / "utils/run-scenarios.sh"), "p0", *args],
                              cwd=self.root, env={**self.env, **(extra or {})}, text=True, capture_output=True)

    def test_default_and_explicit_local_have_same_results(self):
        a = self.cli("--batch-id", "default")
        b = self.cli("--mode", "local", "--batch-id", "explicit")
        self.assertEqual((a.returncode, b.returncode), (1, 1), (a.stderr, b.stderr))
        for name in ("default", "explicit"):
            result = common.read_json(self.root / "runs/p0-batches" / name / "batch-report.json")
            self.assertEqual([r["outcome"] for r in result["results"]], ["PASS", "SKIP", "FAIL"])
            self.assertEqual((result["passed"], result["skipped"], result["failed"]), (1, 1, 1))

    def test_aws_list_requires_no_aws_config_credentials_or_binary(self):
        result = self.cli("--mode", "aws", "--workers", "10", "--list",
                          extra={"UKAMA_IDENTIFIER": "", "UKAMA_PASSWORD": "", "LAB_BIN": "/missing"})
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("total=3", result.stdout)
        self.assertFalse((self.root / "runs").exists())

    def test_local_preserves_optional_kubectl_environment(self):
        result = self.cli("--batch-id", "local-kubectl", extra={"ULAB_KUBECTL": "/local/kubectl"})
        self.assertEqual(result.returncode, 1)
        self.assertIn("kubectl=/local/kubectl", result.stdout)

    def test_local_worker_option_is_rejected(self):
        self.assertEqual(self.cli("--workers", "3").returncode, 2)

    def test_invalid_modes_and_counts(self):
        for args in [("--mode", "invalid"), ("--mode", "aws", "--workers", "0"),
                     ("--mode", "aws", "--workers", "101"), ("--mode=aws", "--workers=NaN")]:
            self.assertEqual(self.cli(*args).returncode, 2, args)

    def test_exact_list_remains_sequential(self):
        selected = self.write("selected.txt", "scenarios/p0/demo/c-fail.yaml\nscenarios/p0/demo/a-pass.yaml\n")
        p = self.cli("--scenario-list", str(selected), "--batch-id", "ordered")
        self.assertEqual(p.returncode, 1, p.stderr)
        data = common.read_json(self.root / "runs/p0-batches/ordered/batch-report.json")
        self.assertEqual([r["outcome"] for r in data["results"]], ["FAIL", "PASS"])

    def test_local_fail_fast_preserved(self):
        selected = self.write("selected.txt", "scenarios/p0/demo/c-fail.yaml\nscenarios/p0/demo/a-pass.yaml\n")
        p = self.cli("--scenario-list", str(selected), "--batch-id", "fast", "--fail-fast")
        self.assertEqual(p.returncode, 1)
        data = common.read_json(self.root / "runs/p0-batches/fast/batch-report.json")
        self.assertEqual(data["total"], 1)

    def test_worker_stop_does_not_create_fake_failures(self):
        hook = self.write("hook", "#!/bin/sh\nexit 10\n", 0o755)
        p = self.cli("--batch-id", "stopped", extra={"ULAB_RUNNER_HOOK": str(hook)})
        self.assertEqual(p.returncode, 0, p.stderr)
        data = common.read_json(self.root / "runs/p0-batches/stopped/batch-report.json")
        self.assertEqual((data["total"], data["failed"]), (0, 0))

    def test_worker_infrastructure_hook_failure_has_separate_exit(self):
        hook = self.write("hook", "#!/bin/sh\nexit 70\n", 0o755)
        p = self.cli("--batch-id", "infra", extra={"ULAB_RUNNER_HOOK": str(hook)})
        self.assertEqual(p.returncode, 70)
        data = common.read_json(self.root / "runs/p0-batches/infra/batch-report.json")
        self.assertEqual(data["failed"], 0)

    def test_upload_failure_preserves_completed_scenario(self):
        hook = self.write("hook", '#!/bin/sh\n[ "$1" = before ] && exit 0\nexit 70\n', 0o755)
        p = self.cli("--batch-id", "upload", extra={"ULAB_RUNNER_HOOK": str(hook)})
        self.assertEqual(p.returncode, 70)
        data = common.read_json(self.root / "runs/p0-batches/upload/batch-report.json")
        self.assertEqual([r["outcome"] for r in data["results"]], ["PASS"])


class ConfigurationTests(Fixture):
    def load(self, **overrides):
        self.cfg.update(overrides)
        return config.load_config(self.write("aws.json", json.dumps(self.cfg)))

    def test_defaults_valid(self):
        self.assertEqual(self.load()["INSTANCE_TYPE"], "c7i.xlarge")

    def test_overlapping_networks_rejected(self):
        with self.assertRaises(common.RunnerError):
            self.load(PODMAN_CIDR="10.240.0.0/16")

    def test_route_overlaps_rejected(self):
        with self.assertRaises(common.RunnerError):
            self.load(VPN_ROUTES=["0.0.0.0/0"])

    def test_unknown_key_rejected(self):
        with self.assertRaises(common.RunnerError):
            self.load(WORKER_COUT=10)

    def test_kubectl_config_paths_are_resolved_relative_to_aws_json(self):
        result = self.load(KUBECTL_FILE="tools/kubectl", KUBECONFIG_FILE="cluster/config")
        self.assertEqual(result["KUBECTL_FILE"], str(self.root / "tools/kubectl"))
        self.assertEqual(result["KUBECONFIG_FILE"], str(self.root / "cluster/config"))

    def test_configured_kubectl_is_not_forwarded_to_aws(self):
        self.cfg["ENV"] = {"UKAMA_IDENTIFIER": "test", "UKAMA_PASSWORD": "test", "ULAB_KUBECTL": "/missing/kubectl"}
        self.assertNotIn("ULAB_KUBECTL", config.environment(self.cfg))

    def test_environment_values_are_data_and_aws_keys_not_forwarded(self):
        self.cfg["ENV"] = {"UKAMA_PASSWORD": "$(do-not-run) `literal` 'secret'", "ULAB_CDR_WAIT_SEC": 60}
        with patch.dict(os.environ, {"UKAMA_IDENTIFIER": "test", "UKAMA_PASSWORD": "default", "AWS_SECRET_ACCESS_KEY": "never-copy", "ULAB_KUBECTL": "/local/tool", "ULAB_RESILIENCE_NAME_SUFFIX": "controller"}, clear=True):
            result = config.environment(self.cfg)
        self.assertEqual(result["UKAMA_PASSWORD"], self.cfg["ENV"]["UKAMA_PASSWORD"])
        self.assertNotIn("AWS_SECRET_ACCESS_KEY", result)
        self.assertNotIn("ULAB_KUBECTL", result)
        self.assertNotIn("ULAB_RESILIENCE_NAME_SUFFIX", result)

    def test_vpn_inlines_files_and_removes_laptop_hooks(self):
        self.write("ca.crt", "TEST CA")
        self.write("userpass", "user\npassword")
        p = self.write("client.ovpn", 'client\nremote vpn.example.test 1194\nca ca.crt\nauth-user-pass userpass\nup /laptop/hook\nredirect-gateway def1\n')
        text = config.portable_vpn(p)
        self.assertIn("<ca>\nTEST CA\n</ca>", text)
        self.assertIn("<auth-user-pass>", text)
        self.assertNotIn("/laptop/hook", text)
        self.assertIn("redirect-gateway def1", text)
        self.assertNotIn('pull-filter ignore "redirect-gateway"', text)
        self.assertIn('pull-filter ignore "block-outside-dns"', text)
        self.assertIn("dev ulab-vpn", text)

    def test_vpn_preserves_explicit_full_tunnel_routes(self):
        p = self.write("client.ovpn", "client\nroute 0.0.0.0 128.0.0.0\nroute 128.0.0.0 128.0.0.0\nblock-outside-dns\n")
        text = config.portable_vpn(p)
        self.assertIn("route 0.0.0.0 128.0.0.0", text)
        self.assertIn("route 128.0.0.0 128.0.0.0", text)
        self.assertNotIn("\nblock-outside-dns\n", text)

    def test_vpn_interactive_auth_rejected(self):
        p = self.write("client.ovpn", "client\nauth-user-pass\n")
        with self.assertRaises(common.RunnerError):
            config.portable_vpn(p)

    def test_vpn_included_configs_rejected(self):
        p = self.write("client.ovpn", "config /laptop/other.conf\n")
        with self.assertRaises(common.RunnerError):
            config.portable_vpn(p)

    def test_repository_paths_cannot_escape(self):
        with self.assertRaises(common.RunnerError):
            self.load(REPO_PATHS=["../private"])


class PackagingTests(Fixture):
    def test_working_native_binary_runs_with_packaged_loader(self):
        destination = self.root / "runtime"
        artifacts.bundle_elf(Path(shutil.which("true")), destination)
        result = subprocess.run([destination / "run"])
        self.assertEqual(result.returncode, 0)

    def test_payload_packages_kubectl_and_rewrites_only_the_worker_path(self):
        self.write("lab/scripts/test.sh", "echo test")
        self.write("repo/testing/node/mk_local_vnode.sh", "echo test")
        work = self.root / "work"
        work.mkdir()
        tool = self.write("tools/kubectl", "#!/bin/sh\nprintf 'fixture kubectl\\n'\n", 0o755)
        env = {}
        with patch.dict(os.environ, {"ULAB_KUBECTL": str(tool), "KUBECONFIG": "/missing/kubeconfig"}):
            with contextlib.redirect_stdout(io.StringIO()):
                result = artifacts.payload(self.root / "lab", self.root / "repo", Path(shutil.which("true")), self.cfg, env, work, [])
        with tarfile.open(result) as archive:
            names = archive.getnames()
        self.assertIn("runtime/lab/program", names)
        self.assertIn("runtime/kubectl/program", names)
        self.assertEqual(env["ULAB_KUBECTL"], "/opt/ukama-runner/runtime/kubectl/run")
        self.assertNotIn("KUBECONFIG", env)
        artifacts.extract(result, self.root / "worker")
        process = subprocess.run([self.root / "worker/runtime/kubectl/run"], text=True, capture_output=True)
        self.assertEqual((process.returncode, process.stdout), (0, "fixture kubectl\n"))

    def test_source_excludes_secrets_results_and_external_symlinks(self):
        self.write("source/scripts/test.sh", "echo ok")
        self.write("source/client.pem", "SECRET")
        self.write("source/.aws/credentials", "SECRET")
        self.write("source/runs/report.txt", "old results")
        self.write("source/setup", "SECRET")
        (self.root / "source/external").symlink_to("/etc/passwd")
        target = self.root / "source.tar.gz"
        with tarfile.open(target, "w:gz") as archive:
            artifacts.add_tree(archive, self.root / "source", "lab")
        with tarfile.open(target) as archive:
            self.assertEqual(archive.getnames(), ["lab/scripts/test.sh"])

    def test_repo_contains_only_needed_build_contexts(self):
        self.write("repo/testing/node/mk_local_vnode.sh", "echo ok")
        self.write("repo/testing/ue/ue/Containerfile", "FROM alpine")
        self.write("repo/nodes/huge-build.bin", "not needed")
        path = self.root / "repo.tar.gz"
        with tarfile.open(path, "w:gz") as archive:
            artifacts.add_tree(archive, self.root / "repo", "repo", includes=self.cfg["REPO_PATHS"])
        with tarfile.open(path) as archive:
            self.assertEqual(len(archive.getnames()), 2)

    def test_archive_path_traversal_rejected(self):
        path = self.root / "bad.tar"
        with tarfile.open(path, "w") as archive:
            info = tarfile.TarInfo("../outside")
            info.size = 1
            archive.addfile(info, io.BytesIO(b"x"))
        with self.assertRaises(common.RunnerError):
            artifacts.extract(path, self.root / "output")

    def test_archive_symlink_escape_rejected(self):
        path = self.root / "bad.tar"
        with tarfile.open(path, "w") as archive:
            info = tarfile.TarInfo("lab/link")
            info.type = tarfile.SYMTYPE
            info.linkname = "../../outside"
            archive.addfile(info)
        with self.assertRaises(common.RunnerError):
            artifacts.extract(path, self.root / "output")

    def test_selection_rejects_host_paths_outside_lab(self):
        file = self.write("elsewhere/a.yaml", "name: test")
        with self.assertRaises(common.RunnerError):
            artifacts.relative_scenarios([file], self.root / "lab")

    def test_partition_has_no_duplicate_or_missing_assignments(self):
        scenarios = [f"scenario-{i}" for i in range(146)]
        parts = controller.partition(scenarios, 10)
        self.assertEqual(sorted(sum(parts, [])), sorted(scenarios))
        self.assertEqual(max(map(len, parts)) - min(map(len, parts)), 1)
        self.assertEqual(controller.partition(["one"], 10), [["one"]])


class ImagePipelineTests(Fixture):
    def setUp(self):
        super().setUp()
        self.raw = self.root / "fixture.tar"
        self.images = image_fixture(self.raw)
        self.cache = self.root / "cache"
        self.exports = 0

    def command(self, args, **kwargs):
        args = list(map(str, args))
        if args[1:3] == ["image", "exists"]:
            return subprocess.CompletedProcess(args, int(args[-1] not in self.images), "", "")
        if args[1:3] == ["image", "inspect"]:
            return subprocess.CompletedProcess(args, 0, json.dumps([{"Architecture": "amd64", "Os": "linux", "Id": self.images[args[-1]]}]), "")
        self.assertEqual(args[1], "save")
        self.exports += 1
        output = Path(args[args.index("--output") + 1])
        shutil.copyfile(self.raw, output)
        return subprocess.CompletedProcess(args, 0, "", "")

    def export(self):
        with patch.object(artifacts, "run", self.command), contextlib.redirect_stdout(io.StringIO()):
            return artifacts.image_archive(self.root / "lab", self.root / "repo", self.cache, {})

    def test_export_roundtrip_keeps_all_tags_configs_layers_and_reuses_cache(self):
        target, info = self.export()
        self.assertEqual(gzip.decompress(target.read_bytes()), self.raw.read_bytes())
        self.assertGreater(info["tar_size"], 1024)
        expanded = self.root / "expanded.tar"
        artifacts.check_compressed_image(target, info, expanded)
        artifacts.validate_image_tar(expanded, self.images)
        self.assertEqual(expanded.read_bytes(), self.raw.read_bytes())
        self.assertEqual(self.export(), (target, info))
        self.assertEqual(self.exports, 1)

    def test_exact_old_temp_path_collision_cache_is_replaced_even_with_matching_sha(self):
        self.cache.mkdir()
        fingerprint = hashlib.sha256(json.dumps(self.images, sort_keys=True).encode()).hexdigest()
        target = self.cache / (fingerprint + ".tar.gz")
        raw = self.cache / (fingerprint + ".tar.tmp")
        compressed = target.with_suffix(".tmp")
        self.assertEqual(raw, compressed)  # Reproduce the delivered defect.
        raw.write_bytes(self.raw.read_bytes())
        with raw.open("rb") as source, gzip.open(compressed, "wb") as output:
            shutil.copyfileobj(source, output)
        compressed.replace(target)
        self.assertEqual(gzip.decompress(target.read_bytes()), b"")
        common.atomic_json(target.with_name(fingerprint + ".json"),
                           {"sha256": common.sha256(target), "size": target.stat().st_size, "images": self.images})
        repaired, info = self.export()
        self.assertEqual(gzip.decompress(repaired.read_bytes()), self.raw.read_bytes())
        self.assertEqual(info["cache_version"], artifacts.IMAGE_CACHE_VERSION)
        self.assertEqual(self.exports, 1)

    def test_damaged_cached_gzip_is_reexported(self):
        target, _ = self.export()
        target.write_bytes(target.read_bytes()[:30])
        self.export()
        self.assertEqual(self.exports, 2)

    def test_empty_podman_export_never_creates_cache_metadata(self):
        self.raw.write_bytes(b"")
        with self.assertRaisesRegex(common.RunnerError, "invalid Podman image archive"):
            self.export()
        self.assertFalse(list(self.cache.glob("*.json")))
        self.assertFalse(list(self.cache.glob("*.tar.gz")))
        self.assertFalse(list(self.cache.glob("image-export-*")))

    def test_rejects_empty_manifest_missing_layer_and_changed_layer(self):
        for flag in ("empty_manifest", "missing_layer", "damaged_layer"):
            with self.subTest(flag=flag):
                images = image_fixture(self.raw, **{flag: True})
                with self.assertRaises(common.RunnerError):
                    artifacts.validate_image_tar(self.raw, images)

    def test_rejects_missing_tag_and_changed_config(self):
        for images in ({"missing:tag": "0" * 64}, {IMAGE_TAGS[0]: "0" * 64}):
            with self.subTest(images=images), self.assertRaises(common.RunnerError):
                artifacts.validate_image_tar(self.raw, images)

    def test_empty_gzip_rejected_even_when_its_compressed_checksum_matches(self):
        target, info = compressed_fixture(self.raw, self.images)
        target.write_bytes(gzip.compress(b""))
        info.update(sha256=common.sha256(target), size=target.stat().st_size)
        with self.assertRaisesRegex(common.RunnerError, "expanded checksum/size"):
            artifacts.check_compressed_image(target, info)

    def test_truncated_gzip_rejected_even_when_compressed_checksum_matches(self):
        target, info = compressed_fixture(self.raw, self.images)
        target.write_bytes(target.read_bytes()[:-5])
        info.update(sha256=common.sha256(target), size=target.stat().st_size)
        with self.assertRaisesRegex(common.RunnerError, "cannot decompress"):
            artifacts.check_compressed_image(target, info)

    def test_same_input_output_is_rejected_without_truncating(self):
        target, info = compressed_fixture(self.raw, self.images)
        before = target.read_bytes()
        with self.assertRaisesRegex(common.RunnerError, "different files"):
            artifacts.check_compressed_image(target, info, target)
        self.assertEqual(target.read_bytes(), before)

    def test_worker_loads_verified_plain_tar_and_checks_each_image_identity(self):
        target, info = self.export()
        downloaded = self.root / "images.tar.gz"
        shutil.copyfile(target, downloaded)
        calls = []
        (self.root / "results").mkdir()
        def command(args, **kwargs):
            args = list(map(str, args))
            calls.append(args)
            if args[1] == "load":
                archive = Path(args[-1])
                self.assertEqual(archive.suffix, ".tar")
                self.assertEqual(archive.read_bytes(), self.raw.read_bytes())
                artifacts.validate_image_tar(archive, self.images)
                return subprocess.CompletedProcess(args, 0, "Loaded images\n", "")
            return subprocess.CompletedProcess(args, 0, json.dumps([{"Id": self.images[args[-1]]}]), "")
        with patch.object(worker, "ROOT", self.root), patch.object(worker, "RESULTS", self.root / "results"), patch.object(worker, "run", command), contextlib.redirect_stdout(io.StringIO()):
            worker.load_images(downloaded, info, 0)
        self.assertEqual(len(calls), 5)
        self.assertFalse(downloaded.exists())
        self.assertFalse((self.root / "images.tar").exists())

    def test_worker_does_not_call_podman_for_corrupt_archive(self):
        target, info = compressed_fixture(self.raw, self.images)
        target.write_bytes(gzip.compress(b""))
        info.update(sha256=common.sha256(target), size=target.stat().st_size)
        with patch.object(worker, "ROOT", self.root), patch.object(worker, "run") as command:
            with self.assertRaises(common.RunnerError):
                worker.load_images(target, info, 0)
            command.assert_not_called()

    def test_worker_surfaces_actual_podman_error(self):
        target, info = compressed_fixture(self.raw, self.images)
        (self.root / "results").mkdir()
        failure = subprocess.CompletedProcess([], 125, "", "Error: layer storage failed")
        with patch.object(worker, "ROOT", self.root), patch.object(worker, "RESULTS", self.root / "results"), patch.object(worker, "run", return_value=failure), contextlib.redirect_stdout(io.StringIO()):
            with self.assertRaisesRegex(common.RunnerError, "Error: layer storage failed"):
                worker.load_images(target, info, 0)
        self.assertIn("layer storage failed", (self.root / "results/image-load.log").read_text())

    def test_disk_reserve_checked_before_decompression_or_podman(self):
        target, info = compressed_fixture(self.raw, self.images)
        disk = shutil._ntuple_diskusage(100, 99, 1)
        with patch.object(worker.shutil, "disk_usage", return_value=disk), patch.object(worker, "check_compressed_image") as decompress, patch.object(worker, "run") as command:
            with self.assertRaisesRegex(common.RunnerError, "insufficient worker disk"):
                worker.load_images(target, info, 0)
            decompress.assert_not_called()
            command.assert_not_called()


class InfrastructureTests(Fixture):
    def test_cloudformation_has_no_ingress_and_storage_is_private(self):
        resources = common.read_json(LAB / "utils/runner/infrastructure.json")["Resources"]
        self.assertNotIn("SecurityGroupIngress", resources["WorkerSecurityGroup"]["Properties"])
        self.assertTrue(all(resources["Bucket"]["Properties"]["PublicAccessBlockConfiguration"].values()))
        self.assertEqual(resources["S3Endpoint"]["Properties"]["VpcEndpointType"], "Gateway")
        actions = resources["WorkerRole"]["Properties"]["Policies"][0]["PolicyDocument"]["Statement"]
        self.assertFalse(any("ec2:*" in item["Action"] for item in actions))

    def test_user_data_is_bounded_valid_bash_and_has_no_credentials(self):
        state = {"region": "us-east-1", "bucket": "bucket", "batch": "test", "max_hours": 24, "bootstrap_sha256": "a"*64}
        text = controller.user_data(state, "worker-001")
        self.assertLess(len(text.encode()), 16384)
        self.assertNotIn("UKAMA_PASSWORD", text)
        self.assertNotIn("auth-user-pass", text)
        self.assertIn("shutdown -h", text)
        self.assertNotIn("ssh ", text)
        file = self.write("user-data.sh", text)
        self.assertEqual(subprocess.run(["bash", "-n", file]).returncode, 0)

    def test_dns_is_split_and_container_address_is_not_loopback(self):
        text = network.dns_config(self.cfg, {"address": "10.240.1.10"}, ["172.20.0.10"])
        self.assertIn("server=/udev.ukama.com/172.20.0.10", text)
        self.assertIn("server=169.254.169.253", text)
        self.assertIn("listen-address=127.0.0.1,10.240.1.10", text)

    def test_backend_names_are_not_sent_to_public_dns_before_vpn(self):
        text = network.dns_config(self.cfg, {"address": "10.240.1.10"}, [])
        self.assertIn("local=/udev.ukama.com/", text)

    def test_probes_have_no_mutating_operations_or_credentials(self):
        env = {"BFF_BASE_URL": "https://bff.udev.ukama.com", "UKAMA_PASSWORD": "secret"}
        urls = network.probe_urls(self.cfg, env)
        self.assertIn("https://bff.udev.ukama.com/", urls)
        self.assertTrue(all(url.endswith("/") for url in urls))
        self.assertNotIn("secret", " ".join(urls))

    def test_s3_access_denied_is_not_mistaken_for_absent_status(self):
        aws = common.AWS("us-east-1")
        with patch.object(aws, "call", return_value=subprocess.CompletedProcess([], 1, "", "AccessDenied")):
            with self.assertRaises(common.RunnerError):
                aws.optional_json("bucket", "key")

    def test_absent_status_is_allowed(self):
        aws = common.AWS("us-east-1")
        with patch.object(aws, "call", return_value=subprocess.CompletedProcess([], 1, "", "(404) HeadObject")):
            self.assertIsNone(aws.optional_json("bucket", "key"))


class ReportsTests(Fixture):
    def state(self):
        return {"batch": "test", "suite": "p0", "scenarios": ["scenarios/p0/demo/a.yaml", "scenarios/p0/demo/b.yaml"],
                "workers": {"worker-001": {"batch_id": "test-w001", "instance_id": "i-test", "scenarios": ["scenarios/p0/demo/a.yaml", "scenarios/p0/demo/b.yaml"]}}}

    def create_report(self, outcome="PASS", marker="DONE"):
        root = self.root / "workers/worker-001"
        original = str(config.ROOT / "results/batch/test-w001/runs/run1/report.json")
        log = str(config.ROOT / "results/batch/test-w001/logs/a.log")
        common.atomic_json(root / "batch/test-w001/batch-report.json", {"results": [
            {"scenario": "demo/a.yaml", "outcome": outcome, "exit_code": int(outcome == "FAIL"), "report": original, "log": log}]})
        common.atomic_json(root / "complete.json", {"phase": marker, "lab_exit": int(outcome == "FAIL")})

    def test_missing_scenarios_are_unfinished_not_failed(self):
        self.create_report(marker="INFRA_ERROR")
        result = controller.aggregate(self.root, self.state())
        self.assertEqual((result["completed"], result["passed"], result["failed"]), (1, 1, 0))
        self.assertEqual(result["unfinished"], ["scenarios/p0/demo/b.yaml"])
        self.assertEqual(result["infrastructure_errors"], ["worker-001"])

    def test_scenario_failure_does_not_mean_worker_failed(self):
        self.create_report(outcome="FAIL")
        result = controller.aggregate(self.root, self.state())
        self.assertEqual(result["failed"], 1)
        self.assertEqual(result["infrastructure_errors"], [])

    def test_report_paths_are_rebased_without_changing_original(self):
        self.create_report()
        source = self.root / "workers/worker-001/batch/test-w001/batch-report.json"
        before = source.read_bytes()
        result = controller.aggregate(self.root, self.state())
        self.assertTrue(result["results"][0]["report"].startswith(str(self.root)))
        self.assertEqual(source.read_bytes(), before)

    def test_partial_tsv_uses_original_scenario_verdict(self):
        directory = self.root / "workers/worker-001/batch/test-w001"
        report = directory / "runs/run1/report.json"
        common.atomic_json(report, {"status": "active", "passed": True})
        directory.joinpath("scenarios.tsv").write_text("category\tscenario\trun_id\texit_code\treport\tlog\n" +
            "demo\tdemo/a.yaml\trun1\t0\t/opt/ukama-runner/results/batch/test-w001/runs/run1/report.json\t/opt/ukama-runner/results/batch/test-w001/logs/a.log\n")
        result = controller.aggregate(self.root, self.state())
        self.assertEqual(result["passed"], 1)
        self.assertEqual(len(result["unfinished"]), 1)

    def test_fail_fast_stop_is_not_worker_failure(self):
        self.create_report(marker="STOPPED")
        self.assertEqual(controller.aggregate(self.root, self.state())["infrastructure_errors"], [])


class VPNHookTests(Fixture):
    def exercise(self, *, fail_firewall=False, missing_firewall=False, routes=None):
        """Execute the actual hook entry point with OpenVPN's restricted PATH.

        Kernel commands are harmless executable fixtures; Python command lookup,
        hook control flow, DNS output, and readiness/error files are real.
        """
        cfg = {**self.cfg, "AWS_S3_CIDRS": ["52.216.0.0/15", "54.231.0.0/16"]}
        common.atomic_json(self.root / "worker.json", {"config": cfg})
        common.atomic_json(self.root / "host-network.json", {"device": "ens5", "gateway": "10.240.1.1", "address": "10.240.1.20"})
        tool = '''#!{python}
import json, os, sys
from pathlib import Path
name=Path(sys.argv[0]).name
with open(os.environ['FIXTURE_LOG'],'a') as out:
    out.write(json.dumps([name,*sys.argv[1:]])+'\\n')
if name=='ip' and sys.argv[1:]==['-j','route','show','dev','ulab-vpn']:
    print(os.environ['FIXTURE_ROUTES'])
if name=='iptables':
    if os.environ.get('FIXTURE_FAIL_FIREWALL')=='1':
        print('fixture firewall failure',file=sys.stderr)
        sys.exit(2)
    if '-C' in sys.argv: sys.exit(1)
'''.format(python=sys.executable)
        self.write("bin/ip", tool, 0o755)
        self.write("bin/systemctl", tool, 0o755)
        if not missing_firewall:
            self.write("sbin/iptables", tool, 0o755)
        driver = self.write("hook.py", f'''import sys
from pathlib import Path
sys.path.insert(0, {str(LAB / 'utils/runner')!r})
import network
network.ROOT = Path({str(self.root)!r})
network.READY = network.ROOT / 'vpn.ready'
network.VPN_ERROR = network.ROOT / 'vpn-error.json'
network.DNS_CONF = network.ROOT / 'dns.conf'
network.SYSTEM_PATH = {str(self.root / 'bin') + ':' + str(self.root / 'sbin')!r}
sys.exit(network.main())
''')
        log = self.root / "commands.jsonl"
        env = {"PATH": str(self.root / "bin"), "FIXTURE_LOG": str(log),
               "FIXTURE_ROUTES": json.dumps(routes if routes is not None else [
                   {"dst": "0.0.0.0/1"}, {"dst": "128.0.0.0/1"}, {"dst": "10.1.133.0/27"}]),
               "FIXTURE_FAIL_FIREWALL": str(int(fail_firewall)),
               "foreign_option_1": "dhcp-option DNS 8.8.8.8",
               "foreign_option_2": "dhcp-option DNS 8.8.4.4"}
        result = subprocess.run([sys.executable, driver, "up"], env=env, text=True, capture_output=True, timeout=15)
        calls = [json.loads(line) for line in log.read_text().splitlines()]
        return result, calls

    def test_route_up_finds_sbin_tools_and_marks_ready_with_full_tunnel(self):
        result, calls = self.exercise()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertTrue((self.root / "vpn.ready").exists())
        self.assertFalse((self.root / "vpn-error.json").exists())
        self.assertTrue(any(c[0] == "iptables" and "-I" in c for c in calls))
        dns = (self.root / "dns.conf").read_text()
        self.assertIn("server=/udev.ukama.com/8.8.8.8", dns)
        self.assertIn("server=/udev.ukama.com/8.8.4.4", dns)

    def test_s3_and_metadata_use_aws_while_dns_uses_vpn(self):
        result, calls = self.exercise()
        self.assertEqual(result.returncode, 0, result.stderr)
        for cidr in ["52.216.0.0/15", "54.231.0.0/16", "169.254.169.254/32", "169.254.169.253/32", "10.240.0.0/16"]:
            self.assertIn(["ip", "route", "replace", cidr, "via", "10.240.1.1", "dev", "ens5"], calls)
        self.assertIn(["ip", "route", "replace", "8.8.8.8/32", "dev", "ulab-vpn"], calls)

    def test_single_default_full_tunnel_is_accepted(self):
        result, _ = self.exercise(routes=[{"dst": "default"}, {"dst": "10.1.133.0/27"}])
        self.assertEqual(result.returncode, 0, result.stderr)

    def test_split_tunnel_is_accepted(self):
        result, _ = self.exercise(routes=[{"dst": "10.1.133.0/27"}, {"dst": "172.20.0.0/16"}])
        self.assertEqual(result.returncode, 0, result.stderr)

    def test_overlapping_specific_vpn_route_is_rejected(self):
        result, _ = self.exercise(routes=[{"dst": "10.241.1.0/24"}])
        self.assertEqual(result.returncode, 1)
        self.assertFalse((self.root / "vpn.ready").exists())
        self.assertIn("overlaps worker networking", common.read_json(self.root / "vpn-error.json")["error"])

    def test_hook_failure_records_actual_error_and_does_not_mark_ready(self):
        result, _ = self.exercise(fail_firewall=True)
        self.assertEqual(result.returncode, 1)
        self.assertFalse((self.root / "vpn.ready").exists())
        self.assertIn("fixture firewall failure", common.read_json(self.root / "vpn-error.json")["error"])

    def test_missing_command_is_explicit(self):
        result, _ = self.exercise(missing_firewall=True)
        self.assertEqual(result.returncode, 1)
        self.assertIn("required worker network command is missing: iptables", result.stderr)

    def test_controller_discovers_only_regional_s3_routes(self):
        aws = common.AWS("us-east-1")
        payload = {"PrefixLists": [
            {"PrefixListName": "com.amazonaws.us-east-1.s3", "Cidrs": ["52.216.0.0/15", "54.231.0.0/16"]},
            {"PrefixListName": "com.amazonaws.us-east-1.dynamodb", "Cidrs": ["3.0.0.0/8"]}]}
        with patch.object(aws, "json", return_value=payload) as call:
            routes = controller.s3_routes(aws, "us-east-1")
        self.assertEqual(routes, ["52.216.0.0/15", "54.231.0.0/16"])
        self.assertIn("Name=prefix-list-name,Values=com.amazonaws.us-east-1.s3", call.call_args.args)

    def test_missing_s3_routes_stop_before_launch(self):
        aws = common.AWS("us-east-1")
        with patch.object(aws, "json", return_value={"PrefixLists": []}):
            with self.assertRaisesRegex(common.RunnerError, "no workers launched"):
                controller.s3_routes(aws, "us-east-1")


class BootstrapTests(Fixture):
    def exercise(self, *, apt_failure=False, install_failure=False, existing=False):
        """Run generated user-data in Bash; isolate host/cloud commands in PATH."""
        fake_bin = self.root / "fake-bin"
        fake_bin.mkdir()
        events = self.root / "events"
        aws = self.write("fake-aws", '''#!/bin/bash
set -eu
if [[ "$1" == --version ]]; then
    echo 'aws-cli/2.0.0 test'
elif [[ "$1 $2" == 's3 cp' && "$3" == s3:* ]]; then
    cp "$FIXTURE_ARCHIVE" "$4"
else
    printf 'upload %s\n' "$3" >> "$FIXTURE_EVENTS"
fi
''', 0o755)
        installer = self.write("fake-installer.sh", '''#!/bin/bash
set -eu
[[ "$*" == '--system --quiet' ]]
printf 'install\n' >> "$FIXTURE_EVENTS"
[[ "$FIXTURE_INSTALL_FAIL" != 1 ]] || exit 42
cp "$FIXTURE_AWS" "$FIXTURE_BIN/aws"
chmod +x "$FIXTURE_BIN/aws"
''')
        self.write("fake-bin/curl", '''#!/bin/bash
set -eu
while [[ "$1" != -o ]]; do shift; done
cp "$FIXTURE_INSTALLER" "$2"
''', 0o755)
        self.write("fake-bin/apt-get", '''#!/bin/bash
printf 'apt %s\n' "$*" >> "$FIXTURE_EVENTS"
for arg in "$@"; do
    if [[ "$arg" == awscli ]]; then
        echo "E: Package 'awscli' has no installation candidate"
        exit 100
    fi
done
[[ "$FIXTURE_APT_FAIL" != 1 ]] || exit 100
exit 0
''', 0o755)
        self.write("fake-bin/shutdown", '#!/bin/bash\nexit 0\n', 0o755)
        entry = self.write("worker-entry.sh", '#!/bin/bash\nprintf "worker-started\\n" >> "$FIXTURE_EVENTS"\n')
        archive = self.root / "bootstrap-fixture.tar.gz"
        with tarfile.open(archive, "w:gz") as target:
            target.add(entry, arcname="worker-entry.sh")
        if existing:
            shutil.copy2(aws, fake_bin / "aws")
        state = {"region": "us-east-1", "bucket": "bucket", "batch": "batch", "max_hours": 24,
                 "bootstrap_sha256": common.sha256(archive)}
        root = self.root / "worker"
        with patch.object(controller, "ROOT", root):
            script = controller.user_data(state, "worker-001")
        script = script.replace("/dev/console", str(self.root / "console.log"))
        path = self.write("user-data.sh", script)
        env = {**os.environ, "PATH": str(fake_bin) + ":" + os.environ["PATH"],
               "FIXTURE_EVENTS": str(events), "FIXTURE_ARCHIVE": str(archive),
               "FIXTURE_AWS": str(aws), "FIXTURE_BIN": str(fake_bin),
               "FIXTURE_INSTALLER": str(installer), "FIXTURE_APT_FAIL": str(int(apt_failure)),
               "FIXTURE_INSTALL_FAIL": str(int(install_failure))}
        result = subprocess.run(["bash", path], env=env, text=True, capture_output=True, timeout=15)
        return result, events.read_text(), root

    def test_bootstrap_reaches_worker_without_apt_awscli(self):
        result, events, root = self.exercise()
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn("install\n", events)
        self.assertIn("worker-started\n", events)
        self.assertNotIn("install -y awscli", events)
        self.assertFalse((root / "results/status.json").exists())

    def test_bootstrap_reuses_working_cli_v2(self):
        result, events, _ = self.exercise(existing=True)
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertNotIn("install\n", events)
        self.assertIn("worker-started\n", events)

    def test_apt_failure_is_visible_without_aws_cli(self):
        result, events, root = self.exercise(apt_failure=True)
        self.assertEqual(result.returncode, 100)
        self.assertIn("UKAMA_RUNNER_BOOTSTRAP_ERROR stage=packages exit=100", result.stdout)
        self.assertEqual(common.read_json(root / "results/status.json")["phase"], "BOOTSTRAP_ERROR")
        self.assertNotIn("worker-started", events)

    def test_installer_failure_stops_before_worker(self):
        result, events, _ = self.exercise(install_failure=True)
        self.assertEqual(result.returncode, 42)
        self.assertIn("UKAMA_RUNNER_BOOTSTRAP_ERROR stage=aws-cli exit=42", result.stdout)
        self.assertNotIn("worker-started", events)


class StartupMonitorTests(Fixture):
    def exercise(self, output, *, status=None, console_error=False, stops=True, cleanup=False):
        state = {"region": "us-east-1", "bucket": "bucket", "batch": "test", "suite": "p0",
                 "started_epoch": 700, "boot_minutes": 45, "max_hours": 24, "fail_fast": False,
                 "scenarios": ["scenarios/p0/demo/a.yaml"],
                 "workers": {"worker-001": {"batch_id": "test-w001", "instance_id": "i-test",
                             "scenarios": ["scenarios/p0/demo/a.yaml"], "launched_epoch": 700}}}
        calls = []
        instance = {"InstanceId": "i-test", "State": {"Name": "running"}}
        class Cloud:
            def optional_json(self, bucket, key):
                return status if key.endswith("status.json") else None
            def sync(self, *args): pass
            def call(self, *args, **kwargs):
                calls.append(args)
                if args[:2] == ("ec2", "get-console-output"):
                    return subprocess.CompletedProcess([], int(console_error), json.dumps({"Output": output}), "AccessDenied" if console_error else "")
                assert args[:2] == ("ec2", "terminate-instances"), args
                instance["State"]["Name"] = "terminated"
                return subprocess.CompletedProcess([], 0, "{}", "")
        cloud = Cloud()
        captured = io.StringIO()
        with contextlib.ExitStack() as stack:
            stack.enter_context(patch.object(controller, "AWS", lambda region: cloud))
            stack.enter_context(patch.object(controller, "verify_account"))
            stack.enter_context(patch.object(controller, "instances", return_value={"worker-001": instance}))
            stack.enter_context(patch.object(controller.time, "time", return_value=1000))
            stack.enter_context(patch.object(controller.time, "sleep", side_effect=RuntimeError("next-poll")))
            stack.enter_context(contextlib.redirect_stdout(captured))
            if cleanup:
                self.assertEqual(controller.cleanup(self.root, state), 0)
            elif stops:
                self.assertEqual(controller.monitor(self.root, state), 2)
            else:
                with self.assertRaisesRegex(RuntimeError, "next-poll"):
                    controller.monitor(self.root, state)
        return calls, captured.getvalue(), state

    def test_original_apt_failure_is_collected_and_reported_without_scenario_fail(self):
        output = "E: Package 'awscli' has no installation candidate\nFailed to run module scripts_user\n"
        calls, text, state = self.exercise(output)
        self.assertIn("no installation candidate", text)
        self.assertTrue(any(c[1] == "terminate-instances" for c in calls))
        self.assertEqual((self.root / "worker-001-console.log").read_text(), output)
        report = common.read_json(self.root / "batch-report.json")
        self.assertEqual((report["failed"], report["completed"]), (0, 0))
        self.assertEqual(len(report["unfinished"]), 1)
        self.assertEqual(report["workers"]["worker-001"]["phase"], "BOOTSTRAP_ERROR")
        self.assertIn("no installation candidate", state["workers"]["worker-001"]["bootstrap_error"])

    def test_new_failure_marker_identifies_exact_stage(self):
        _, text, _ = self.exercise("UKAMA_RUNNER_BOOTSTRAP_ERROR stage=aws-cli exit=42\n")
        self.assertIn("bootstrap stage aws-cli exited 42", text)

    def test_s3_bootstrap_error_stops_even_without_console(self):
        calls, text, _ = self.exercise(None, status={"phase": "BOOTSTRAP_ERROR", "error": "download failed"})
        self.assertIn("download failed", text)
        self.assertTrue(any(c[1] == "terminate-instances" for c in calls))

    def test_slow_installation_is_retained_and_stage_is_shown(self):
        calls, text, _ = self.exercise("UKAMA_RUNNER_BOOTSTRAP stage=packages\n", stops=False)
        self.assertIn("packages; waiting for first worker heartbeat", text)
        self.assertFalse(any(c[1] == "terminate-instances" for c in calls))

    def test_ssm_warning_alone_does_not_stop_worker(self):
        calls, _, _ = self.exercise("SSM Agent unable to acquire credentials: AccessDeniedException\n", stops=False)
        self.assertFalse(any(c[1] == "terminate-instances" for c in calls))

    def test_unavailable_console_does_not_stop_worker(self):
        calls, _, _ = self.exercise(None, console_error=True, stops=False)
        self.assertFalse(any(c[1] == "terminate-instances" for c in calls))

    def test_stale_running_worker_is_retained(self):
        calls, text, _ = self.exercise("", status={"phase": "RUNNING", "heartbeat_epoch": 710}, stops=False)
        self.assertIn("heartbeat stale in RUNNING", text)
        self.assertFalse(any(c[1] == "terminate-instances" for c in calls))

    def test_cleanup_saves_console_before_termination(self):
        output = "E: Package 'awscli' has no installation candidate\nFailed to run module scripts_user\n"
        calls, _, _ = self.exercise(output, cleanup=True)
        self.assertEqual(calls[0][1], "get-console-output")
        self.assertEqual(calls[-1][1], "terminate-instances")
        self.assertEqual((self.root / "worker-001-console.log").read_text(), output)

    def test_console_timeout_is_nonfatal(self):
        aws = common.AWS("us-east-1")
        with patch.object(aws, "call", side_effect=subprocess.TimeoutExpired("aws", 90)), contextlib.redirect_stdout(io.StringIO()):
            result = controller.console_diagnostics(aws, self.root, "worker-001", "i-test")
        self.assertFalse(result["available"])
        self.assertEqual(result["error"], "")


class InfrastructureAdmissionTests(Fixture):
    def exercise(self, records, *, running=True, age=10, console=None):
        state = {"bucket": "bucket", "batch": "batch", "boot_minutes": 5,
                 "workers": {"worker-001": {"instance_id": "i-test", "launched_epoch": 1000 - age}}}
        class Cloud:
            def optional_json(self, bucket, key):
                return records.get(Path(key).name)
        with contextlib.ExitStack() as stack:
            stack.enter_context(patch.object(controller.time, "time", return_value=1000))
            stack.enter_context(patch.object(controller.time, "sleep", side_effect=RuntimeError("still-waiting")))
            stack.enter_context(patch.object(controller, "instances", return_value={"worker-001": {"InstanceId": "i-test", "State": {"Name": "running" if running else "terminated"}}}))
            stack.enter_context(patch.object(controller, "console_diagnostics", return_value=console or {"error": "", "stage": "packages", "available": True}))
            stack.enter_context(contextlib.redirect_stdout(io.StringIO()))
            return controller.wait_for_infrastructure(Cloud(), self.root, state, "worker-001"), state

    def ready(self):
        return {"ready.json": {"worker": "worker-001", "batch": "batch", "ready_at": "2026-10-03T00:00:00Z"}}

    def test_readiness_survives_transition_to_running(self):
        accepted, state = self.exercise({**self.ready(), "status.json": {"phase": "RUNNING"}})
        self.assertTrue(accepted)
        self.assertIn("infrastructure_ready_at", state["workers"]["worker-001"])

    def test_scenario_failure_does_not_block_launches_after_infrastructure_ready(self):
        accepted, _ = self.exercise({**self.ready(), "complete.json": {"phase": "DONE", "lab_exit": 1}})
        self.assertTrue(accepted)

    def test_infrastructure_error_prevents_admission_even_before_completion_upload(self):
        accepted, _ = self.exercise({"status.json": {"phase": "INFRA_ERROR", "error": "image-load: empty tar"}})
        self.assertFalse(accepted)

    def test_known_infrastructure_error_takes_precedence_over_old_ready_marker(self):
        accepted, _ = self.exercise({**self.ready(), "complete.json": {"phase": "INFRA_ERROR"}})
        self.assertFalse(accepted)

    def test_slow_start_does_not_admit_or_reject_early(self):
        with self.assertRaisesRegex(RuntimeError, "still-waiting"):
            self.exercise({}, age=200)

    def test_marker_from_other_batch_is_not_accepted(self):
        records = self.ready()
        records["ready.json"]["batch"] = "different"
        with self.assertRaisesRegex(RuntimeError, "still-waiting"):
            self.exercise(records)

    def test_startup_deadline_stops_admission(self):
        accepted, state = self.exercise({}, age=301)
        self.assertFalse(accepted)
        self.assertIn("deadline", state["workers"]["worker-001"]["bootstrap_error"])

    def test_early_console_bootstrap_failure_stops_admission(self):
        accepted, state = self.exercise({}, age=200, console={"error": "bootstrap stage packages exited 100", "stage": "packages", "available": True})
        self.assertFalse(accepted)
        self.assertIn("packages", state["workers"]["worker-001"]["bootstrap_error"])

    def test_terminated_worker_without_marker_stops_admission(self):
        accepted, _ = self.exercise({}, running=False)
        self.assertFalse(accepted)


class WorkerContractTests(Fixture):
    def exercise(self, *, fail_upload=False, fail_network=False, fail_images=False):
        raw = self.root / "fixture.tar"
        images = image_fixture(raw)
        image_file, image_info = compressed_fixture(raw, images)
        if fail_images:
            image_file.write_bytes(gzip.compress(b""))
            image_info.update(sha256=common.sha256(image_file), size=image_file.stat().st_size)
        kubectl = self.write("runtime/kubectl/run", "#!/bin/sh\nexit 0\n", 0o755)
        kubeconfig = self.write("runtime/kubectl/kubeconfig.json", "{}", 0o600)
        cfg = {"config": self.cfg, "bucket": "bucket", "batch": "batch", "worker": "worker-001",
               "batch_id": "batch-w001", "suite": "p0", "scenario_root": "scenarios/p0",
               "scenarios": ["scenarios/p0/demo/fail.yaml"], "environment": {"ULAB_KUBECTL": str(kubectl), "KUBECONFIG": str(kubeconfig)}, "fail_fast": False,
               "source_sha256": hashlib.sha256(b"archive").hexdigest(), "image_key": "images/digest.tar.gz",
               "image_info": image_info}
        cfg["config"]["MIN_FREE_DISK_GB"] = 0
        calls, uploaded = [], {}
        root = self.root
        class Cloud:
            def copy(self, source, destination):
                calls.append(("copy", str(destination)))
                if str(destination).startswith("s3:"):
                    uploaded[str(destination)] = Path(source).read_bytes()
                else:
                    Path(destination).parent.mkdir(parents=True, exist_ok=True)
                    if str(source).endswith("images/digest.tar.gz"):
                        shutil.copyfile(image_file, destination)
                    else:
                        Path(destination).write_text("archive")
            def sync(self, source, destination):
                calls.append(("sync", str(destination)))
                if fail_upload:
                    raise common.RunnerError("upload unavailable")
        class IdleThread:
            def __init__(self, **kwargs): pass
            def start(self): pass
            def join(self, **kwargs): pass
        class Process:
            def __init__(self, args, **kwargs):
                self.args = args
                self.env = kwargs["env"]
                assert self.env["ULAB_KUBECTL"] == str(kubectl)
                assert self.env["KUBECONFIG"] == str(kubeconfig)
                assert any(key.endswith("/ready.json") for key in uploaded)
                calls.append(("scenario-start", None))
                Path(self.env["P0_STATUS_FILE"]).write_text("completed\t1\nfailed\t1\n")
                common.atomic_json(Path(self.env["P0_RUNS_DIR"]) / "batch-w001/batch-report.json", {"results": [{"outcome": "FAIL"}]})
            def wait(self): return 1
        def command(args, **kwargs):
            values = list(map(str, args))
            calls.append(("command", values))
            if values[:2] == ["podman", "load"]:
                assert Path(values[-1]).read_bytes() == raw.read_bytes()
            text = json.dumps([{"Id": images[values[-1]]}]) if "inspect" in values else ""
            return subprocess.CompletedProcess(values, 0, text, "")
        with contextlib.ExitStack() as stack:
            for attr, value in [("ROOT", root), ("RESULTS", root / "results"), ("AWS", lambda region: Cloud()),
                                ("run", command), ("extract", lambda *args: (root / "lab").mkdir()),
                                ("bounded_retry", lambda callback, *args: callback())]:
                stack.enter_context(patch.object(worker, attr, value))
            stack.enter_context(patch.object(worker.threading, "Thread", IdleThread))
            stack.enter_context(patch.object(worker.subprocess, "Popen", Process))
            stack.enter_context(patch.object(worker.network, "configure", side_effect=common.RunnerError("VPN failed") if fail_network else None))
            stack.enter_context(patch.object(worker.network, "check"))
            stack.enter_context(patch.object(worker.network, "snapshot"))
            if fail_upload:
                with self.assertRaises(common.RunnerError):
                    worker.worker(cfg)
            else:
                self.assertEqual(worker.worker(cfg), 0)
        return calls, uploaded

    def test_scenario_fail_is_done_and_marker_follows_all_uploads(self):
        calls, uploaded = self.exercise()
        marker_key = "s3://bucket/results/batch/worker-001/complete.json"
        marker = json.loads(uploaded[marker_key])
        self.assertEqual((marker["phase"], marker["lab_exit"]), ("DONE", 1))
        sync_index = next(i for i, c in enumerate(calls) if c[0] == "sync")
        marker_index = calls.index(("copy", marker_key))
        shutdown_index = next(i for i, c in enumerate(calls) if c[0] == "command" and c[1][0] == "shutdown")
        self.assertLess(sync_index, marker_index)
        self.assertLess(marker_index, shutdown_index)

    def test_upload_failure_never_publishes_completion_or_shuts_down(self):
        calls, uploaded = self.exercise(fail_upload=True)
        self.assertFalse(any(key.endswith("complete.json") for key in uploaded))
        self.assertFalse(any(c[0] == "command" and c[1][0] == "shutdown" for c in calls))

    def test_vpn_failure_is_infrastructure_error_without_scenario_exit(self):
        with contextlib.redirect_stdout(io.StringIO()):
            calls, uploaded = self.exercise(fail_network=True)
        marker = json.loads(uploaded["s3://bucket/results/batch/worker-001/complete.json"])
        self.assertEqual(marker["phase"], "INFRA_ERROR")
        self.assertIsNone(marker["lab_exit"])
        self.assertFalse(any(key.endswith("ready.json") for key in uploaded))

    def test_invalid_images_publish_error_without_ready_or_scenario_execution(self):
        with contextlib.redirect_stdout(io.StringIO()):
            calls, uploaded = self.exercise(fail_images=True)
        marker = json.loads(uploaded["s3://bucket/results/batch/worker-001/complete.json"])
        self.assertEqual(marker["phase"], "INFRA_ERROR")
        self.assertIn("image-load: starter-image expanded checksum/size mismatch", marker["error"])
        self.assertIsNone(marker["lab_exit"])
        self.assertFalse(any(key.endswith("ready.json") for key in uploaded))
        self.assertFalse(any(kind == "scenario-start" for kind, _ in calls))


class LaunchContractTests(Fixture):
    def exercise(self, *, ready=True):
        self.write("lab/scenarios/p0/demo/a.yaml", "name: a")
        self.write("lab/scenarios/p0/demo/b.yaml", "name: b")
        self.write("repo/testing/node/mk_local_vnode.sh", "#!/bin/sh\n", 0o755)
        self.write("repo/testing/ue/ue/Containerfile", "FROM alpine")
        self.write("repo/testing/ue/media/Containerfile", "FROM alpine")
        self.cfg["VPN_CONFIG_FILE"] = str(self.write("client.ovpn", "client\nremote vpn.example.test 1194\n<ca>\ncertificate\n</ca>\n"))
        config_file = self.write("aws.json", json.dumps(self.cfg))
        image = self.write("images.tar.gz", "cached images")
        captured, launches = {}, []
        class Cloud:
            def json(self, *args):
                if args[:2] == ("sts", "get-caller-identity"):
                    return {"Account": "account"}
                if args[:2] == ("ec2", "describe-prefix-lists"):
                    return {"PrefixLists": [{"PrefixListName": "com.amazonaws.us-east-1.s3", "Cidrs": ["52.216.0.0/15"]}]}
                raise AssertionError(args)
            def head(self, *args): return None
            def copy(self, *args): pass
            def sync(self, source, destination):
                captured.update({p.name: p.read_bytes() for p in Path(source).iterdir() if p.is_file()})
            def call(self, *args, **kwargs):
                assert args[:2] == ("ec2", "run-instances")
                request = common.read_json(args[3].removeprefix("file://"))
                launches.append(request)
                return subprocess.CompletedProcess([], 0, json.dumps({"Instances": [{"InstanceId": "i-" + str(len(launches))}]}), "")
        args = argparse.Namespace(config=str(config_file), lab_dir=str(self.root / "lab"), repo=str(self.root / "repo"),
                                  lab_bin=shutil.which("true"), scenario_root=str(self.root / "lab/scenarios/p0"),
                                  suite="p0", batch_id="test", out=str(self.root / "runs"), workers=2,
                                  factory_nodes="0", prepare_only=False, fail_fast=False)
        def payload(*values):
            target = values[-2] / "source.tar.gz"
            target.write_text("sources")
            return target
        def admission(*values):
            self.assertEqual(len(launches), 1)
            return ready
        with contextlib.ExitStack() as stack:
            stack.enter_context(patch.dict(os.environ, {"UKAMA_IDENTIFIER": "test", "UKAMA_PASSWORD": "test"}))
            stack.enter_context(patch.object(sys, "stdin", io.StringIO(str(self.root / "lab/scenarios/p0/demo/a.yaml") + "\0" + str(self.root / "lab/scenarios/p0/demo/b.yaml") + "\0")))
            stack.enter_context(patch.object(controller, "AWS", lambda region: Cloud()))
            stack.enter_context(patch.object(controller, "executable", lambda value: Path(value)))
            stack.enter_context(patch.object(controller, "image_archive", return_value=(image, {"sha256": "image", "size": 13, "images": {}})))
            stack.enter_context(patch.object(controller, "payload", payload))
            stack.enter_context(patch.object(controller, "ensure_stack", return_value={"Bucket": "bucket", "Subnet": "subnet", "SecurityGroup": "sg", "InstanceProfile": "profile"}))
            stack.enter_context(patch.object(controller, "find_ami", return_value=("ami-test", "/dev/sda1")))
            stack.enter_context(patch.object(controller, "monitor", return_value=0))
            stack.enter_context(patch.object(controller, "wait_for_infrastructure", side_effect=admission))
            stack.enter_context(contextlib.redirect_stdout(io.StringIO()))
            self.assertEqual(controller.launch(args), 0)
        return launches, captured

    def test_launch_prepares_exact_worker_inputs_and_uses_no_ssh(self):
        launches, captured = self.exercise()
        self.assertEqual(len(launches), 2)
        self.assertEqual(len({r["ClientToken"] for r in launches}), 2)
        for request in launches:
            self.assertNotIn("KeyName", request)
            self.assertEqual(request["MetadataOptions"]["HttpTokens"], "required")
            self.assertTrue(request["BlockDeviceMappings"][0]["Ebs"]["Encrypted"])
        first = json.loads(captured["worker-001.json"])
        second = json.loads(captured["worker-002.json"])
        self.assertEqual(first["scenarios"], ["scenarios/p0/demo/a.yaml"])
        self.assertEqual(second["scenarios"], ["scenarios/p0/demo/b.yaml"])
        self.assertEqual(first["environment"]["UKAMA_REPO"], "/opt/ukama-runner/repo")
        self.assertEqual(first["config"]["AWS_S3_CIDRS"], ["52.216.0.0/15"])
        state = common.read_json(self.root / "runs/test/state.json")
        self.assertEqual(state["workers"]["worker-002"]["instance_id"], "i-2")
        self.assertNotIn("UKAMA_PASSWORD", json.dumps(state))

    def test_first_worker_failure_prevents_remaining_launches(self):
        launches, _ = self.exercise(ready=False)
        self.assertEqual(len(launches), 1)
        state = common.read_json(self.root / "runs/test/state.json")
        self.assertIsNone(state["workers"]["worker-002"]["instance_id"])
        self.assertIn("not launched", state["workers"]["worker-002"]["not_started_reason"])



class AWSCleanupAndLifecycleTests(Fixture):
    def test_existing_kubeconfig_is_exported_locally_without_auth_checks(self):
        original = self.write("lab/connection.yaml", "private original config")
        self.write("repo/testing/node/mk_local_vnode.sh", "echo test")
        export = {"apiVersion": "v1", "kind": "Config", "current-context": "vpn",
                  "clusters": [{"name": "vpn", "cluster": {"server": "https://cluster.internal"}}],
                  "users": [{"name": "existing", "user": {"exec": {"command": "existing-plugin"}}}]}
        tool = self.write("tools/kubectl", "#!/bin/sh\nexit 0\n", 0o755)
        env, work = {}, self.root / "work"
        work.mkdir()
        self.cfg.update(KUBECTL_FILE=str(tool), KUBECONFIG_FILE=str(original))
        actual_run = artifacts.run
        commands = []
        def command(args, **kwargs):
            if str(args[0]) == str(tool):
                commands.append(list(map(str, args)))
                self.assertEqual(kwargs["env"]["KUBECONFIG"], str(original))
                return subprocess.CompletedProcess(args, 0, json.dumps(export), "")
            return actual_run(args, **kwargs)
        with patch.object(artifacts, "run", side_effect=command):
            archive_path = artifacts.payload(self.root / "lab", self.root / "repo", Path(shutil.which("true")), self.cfg, env, work, [])
        self.assertEqual(commands, [[str(tool), "config", "view", "--raw", "--flatten", "--minify", "-o", "json"]])
        with tarfile.open(archive_path) as archive:
            self.assertNotIn("lab/connection.yaml", archive.getnames())
            member = archive.getmember("runtime/kubectl/kubeconfig.json")
            self.assertEqual(member.mode & 0o777, 0o600)
            self.assertEqual(json.load(archive.extractfile(member)), export)
        self.assertEqual(env["KUBECONFIG"], "/opt/ukama-runner/runtime/kubectl/kubeconfig.json")

    def test_native_kubectl_is_bundled_without_a_config_requirement(self):
        env = {}
        self.cfg["KUBECTL_FILE"] = shutil.which("true")
        with patch.dict(os.environ, {"HOME": str(self.root), "KUBECONFIG": ""}):
            artifacts.bundle_kubectl(self.cfg, env, self.root / "runtime")
        result = subprocess.run([self.root / "runtime/kubectl/run"])
        self.assertEqual(result.returncode, 0)
        self.assertNotIn("KUBECONFIG", env)

    def cleanup_with_packaged_tool(self, delete_rc):
        tool = self.write("tools/kubectl", f"""#!{sys.executable}
import json, os, sys
with open(os.environ['CLEANUP_CALLS'], 'a') as stream:
    print(json.dumps(sys.argv[1:]), file=stream)
if sys.argv[1:3] == ['get', 'pods']:
    print('ukama-mesh-node-other-abc')
    print('ukama-mesh-node-owned-xyz')
    print('unrelated-pod')
    sys.exit(0)
if sys.argv[1:3] == ['delete', 'pod']:
    sys.exit({delete_rc})
sys.exit(99)
""", 0o755)
        self.cfg["KUBECTL_FILE"] = str(tool)
        packaged_env = {}
        with patch.dict(os.environ, {"HOME": str(self.root), "KUBECONFIG": ""}):
            artifacts.bundle_kubectl(self.cfg, packaged_env, self.root / "runtime")
        script = self.root / "scripts/stop-node.sh"
        script.parent.mkdir()
        shutil.copy2(LAB / "scripts/stop-node.sh", script)
        self.write("run/runtime-nodes/node-1.env", "FACTORY_NODE_ID=owned\n")
        self.write("bin/podman", "#!/bin/sh\nexit 1\n", 0o755)
        commands = self.root / "calls.jsonl"
        env = {**os.environ, "ULAB_KUBECTL": str(self.root / "runtime/kubectl/run"),
               "PATH": str(self.root / "bin") + os.pathsep + os.environ['PATH'],
               "CLEANUP_CALLS": str(commands)}
        result = subprocess.run([script, "node-1", self.root / "run"], env=env, capture_output=True, text=True)
        return result, [json.loads(line) for line in commands.read_text().splitlines()]

    def test_unchanged_stop_node_deletes_only_its_own_mesh_pod(self):
        result, calls = self.cleanup_with_packaged_tool(0)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(calls, [["get", "pods", "-n", "ukama-messaging", "-o", "custom-columns=NAME:.metadata.name", "--no-headers"],
                                 ["delete", "pod", "ukama-mesh-node-owned-xyz", "-n", "ukama-messaging", "--ignore-not-found=true"]])

    def test_existing_cleanup_failure_remains_a_failure(self):
        result, _ = self.cleanup_with_packaged_tool(1)
        self.assertEqual(result.returncode, 1)
        self.assertIn("failed to delete mesh pod", result.stderr)

    def test_lifecycle_logs_only_transitions_and_confirms_real_termination(self):
        item = {"instance_id": "i-test", "termination_requested_at": "now"}
        state = {"workers": {"worker-001": item}}
        running = {"InstanceId": "i-test", "State": {"Name": "running"}}
        stopped = {"InstanceId": "i-test", "State": {"Name": "terminated"}}
        output = io.StringIO()
        with contextlib.redirect_stdout(output):
            controller.observe_instance(self.root, "worker-001", item, running)
            controller.observe_instance(self.root, "worker-001", item, running)
            with patch.object(controller, "instances", side_effect=[{"worker-001": running}, {"worker-001": stopped}]), patch.object(controller.time, "sleep"):
                controller.confirm_termination(object(), self.root, state)
        self.assertEqual(output.getvalue().count("EC2 i-test running"), 1)
        self.assertIn("EC2 i-test terminated; decommissioned", output.getvalue())
        self.assertIn("decommissioned", (self.root / "aws-lifecycle.log").read_text())

    def test_unconfirmed_termination_does_not_claim_decommissioned(self):
        state = {"workers": {"worker-001": {"instance_id": "i-test", "termination_requested_at": "now"}}}
        instance = {"InstanceId": "i-test", "State": {"Name": "shutting-down"}}
        output = io.StringIO()
        with patch.object(controller, "instances", return_value={"worker-001": instance}), contextlib.redirect_stdout(output):
            controller.confirm_termination(object(), self.root, state, timeout=0)
        self.assertIn("termination not yet confirmed", output.getvalue())
        self.assertNotIn("decommissioned", output.getvalue())

    def test_summary_colors_only_terminal_output_and_preserves_report(self):
        report = self.write("batch-report.txt", "FAIL worker-001 billing/a.yaml\nPASS worker-002 node/b.yaml\n")
        original = report.read_bytes()
        output = io.StringIO()
        with contextlib.redirect_stdout(output), patch.object(output, "isatty", return_value=True):
            controller.print_report(report)
        self.assertIn("\033[1;31mFAIL\033[0m", output.getvalue())
        self.assertIn("\033[1;32mPASS\033[0m", output.getvalue())
        self.assertEqual(report.read_bytes(), original)
        output = io.StringIO()
        with contextlib.redirect_stdout(output):
            controller.print_report(report)
        self.assertEqual(output.getvalue(), original.decode())


if __name__ == "__main__":
    unittest.main()
