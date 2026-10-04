#!/usr/bin/env python3
"""Package existing code and images; do not invent a second node build flow."""
import gzip
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import tarfile
import tempfile
from common import RunnerError, atomic_json, read_json, run, sha256
from config import ROOT, executable

EXCLUDE_DIRS = {".git", ".aws", ".ssh", ".kube", ".venv", "node_modules", "__pycache__", ".pytest_cache", "runs", "logs"}
SECRET_SUFFIXES = {".pem", ".key", ".p12", ".pfx", ".ovpn"}
IMAGE_CACHE_VERSION = 2


def validate_image_tar(path, images):
    """Validate the Docker archive without loading images or extracting files.

    A successful gzip checksum says nothing about whether the payload contains
    images. Check every expected tag/config and each referenced layer's bytes.
    Shared layers are hashed only once.
    """
    try:
        with tarfile.open(path, "r:") as archive:
            def read_member(name):
                member = archive.getmember(name)
                if not (member.isfile() or member.islnk()) or member.size < 0:
                    raise RunnerError(f"invalid image archive member: {name}")
                stream = archive.extractfile(member)
                if stream is None:
                    raise RunnerError(f"unreadable image archive member: {name}")
                return stream

            with read_member("manifest.json") as stream:
                manifest = json.load(stream)
            if not isinstance(manifest, list) or not manifest:
                raise RunnerError("image archive has no images in manifest.json")
            found, layers = {}, {}
            for entry in manifest:
                with read_member(entry["Config"]) as stream:
                    data = stream.read()
                config = json.loads(data)
                image_id = hashlib.sha256(data).hexdigest()
                if config.get("architecture") != "amd64" or config.get("os") != "linux":
                    raise RunnerError("image archive contains a non-Linux/amd64 image")
                for tag in entry.get("RepoTags") or []:
                    if tag in found and found[tag] != image_id:
                        raise RunnerError(f"ambiguous image tag in archive: {tag}")
                    found[tag] = image_id
                references = entry["Layers"]
                diffs = config["rootfs"]["diff_ids"]
                if len(references) != len(diffs):
                    raise RunnerError("image archive layer count differs from its image configuration")
                for name, expected in zip(references, diffs):
                    if name not in layers:
                        with read_member(name) as stream:
                            digest = hashlib.sha256()
                            for block in iter(lambda: stream.read(8 * 1024 * 1024), b""):
                                digest.update(block)
                        layers[name] = "sha256:" + digest.hexdigest()
                    if layers[name] != expected:
                        raise RunnerError(f"image archive layer checksum mismatch: {name}")
            if not images:
                raise RunnerError("no expected image identities provided")
            for tag, expected in images.items():
                if found.get(tag) != expected.removeprefix("sha256:"):
                    raise RunnerError(f"image archive is missing or has changed image: {tag}")
    except (OSError, EOFError, tarfile.TarError, ValueError, KeyError, TypeError) as exc:
        raise RunnerError(f"invalid Podman image archive: {exc}") from exc


def check_compressed_image(path, info, destination=None):
    """Verify compression round-trip; optionally materialize the plain tar.

    Source and destination must differ. No image load is attempted until both
    compressed and uncompressed checksums match the validated controller tar.
    """
    if info.get("cache_version") != IMAGE_CACHE_VERSION:
        raise RunnerError("image cache predates archive validation; export it again")
    if path.stat().st_size != info["size"] or sha256(path) != info["sha256"]:
        raise RunnerError("starter-image compressed checksum/size mismatch")
    if info.get("tar_size", 0) < 1024:
        raise RunnerError("starter-image archive is empty or truncated")
    if destination is not None and path.resolve() == destination.resolve():
        raise RunnerError("image source and destination must be different files")
    digest, size = hashlib.sha256(), 0
    try:
        with gzip.open(path, "rb") as source:
            with (destination.open("wb") if destination is not None else open(os.devnull, "wb")) as output:
                for block in iter(lambda: source.read(8 * 1024 * 1024), b""):
                    size += len(block)
                    if size > info["tar_size"]:
                        raise RunnerError("starter-image archive exceeds expected expanded size")
                    digest.update(block)
                    if destination is not None:
                        output.write(block)
        if size != info["tar_size"] or digest.hexdigest() != info["tar_sha256"]:
            raise RunnerError("starter-image expanded checksum/size mismatch")
    except (OSError, EOFError) as exc:
        raise RunnerError(f"cannot decompress starter images: {exc}") from exc


def relative_scenarios(values, lab):
    result = []
    for value in values:
        path = Path(value).resolve()
        try:
            relative = path.relative_to(lab)
        except ValueError:
            raise RunnerError(f"AWS scenarios must be inside the lab directory: {path}") from None
        if not path.is_file() or path.suffix not in (".yaml", ".yml") or "\n" in str(relative):
            raise RunnerError(f"invalid scenario path: {relative}")
        if str(relative) not in result:
            result.append(str(relative))
    return result


def add_tree(archive, source, destination, excludes=(), includes=None):
    source = Path(source).resolve()
    excluded = {Path(p).resolve() for p in excludes}
    skipped = []
    def included(path):
        return includes is None or any(path.is_relative_to(Path(i)) or Path(i).is_relative_to(path) for i in includes)
    def visit(directory):
        for p in sorted(directory.iterdir()):
            relative = p.relative_to(source)
            if not included(relative):
                continue
            if p.name in EXCLUDE_DIRS or p.name in {"setup", "aws.json", "config.env"} or p.name.startswith(".env") or p.suffix in SECRET_SUFFIXES:
                continue
            if any(p.resolve() == x or x in p.resolve().parents for x in excluded):
                continue
            arcname = str(Path(destination) / relative)
            if p.is_symlink():
                target = p.resolve()
                if not target.exists() or not target.is_relative_to(source) or not included(target.relative_to(source)):
                    skipped.append(str(relative))
                    continue
                info = archive.gettarinfo(str(p), arcname)
                info.linkname = os.path.relpath(Path(destination) / target.relative_to(source), Path(arcname).parent)
                info.uid = info.gid = 0
                info.uname = info.gname = "root"
                archive.addfile(info)
            elif p.is_dir():
                visit(p)
            elif p.is_file():
                info = archive.gettarinfo(str(p), arcname)
                info.uid = info.gid = 0
                info.uname = info.gname = "root"
                with p.open("rb") as stream:
                    archive.addfile(info, stream)
    visit(source)
    return skipped


def bundle_elf(binary, destination):
    """Use the caller's working binary and its loader/libraries, not a new build."""
    binary = executable(str(binary))
    with binary.open("rb") as stream:
        header = stream.read(20)
    if header[:4] != b"\x7fELF" or int.from_bytes(header[18:20], "little") != 62:
        raise RunnerError(f"AWS needs an x86_64 Linux ELF executable: {binary}")
    destination.mkdir(parents=True, exist_ok=True)
    shutil.copy2(binary, destination / "program")
    dependencies = run(["ldd", binary], check=False)
    output = dependencies.stdout + dependencies.stderr
    if "not found" in output:
        raise RunnerError(f"working binary has unresolved shared libraries: {binary}")
    if dependencies.returncode and not any(v in output for v in ("not a dynamic executable", "statically linked")):
        raise RunnerError(f"cannot inspect binary dependencies: {binary}")
    loader = None
    for line in output.splitlines():
        match = re.search(r"(?:=>\s+)?(/\S+)\s+\(", line)
        if not match:
            continue
        dependency = Path(match.group(1))
        shutil.copy2(dependency, destination / dependency.name)
        if dependency.name.startswith(("ld-linux", "ld-musl")):
            loader = dependency.name
    script = '#!/bin/sh\nset -eu\nd="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"\n'
    if loader:
        script += f'exec "$d/{loader}" --library-path "$d" "$d/program" "$@"\n'
    else:
        script += 'exec "$d/program" "$@"\n'
    (destination / "run").write_text(script)
    (destination / "run").chmod(0o755)


def factory_ids(url):
    p = run(["curl", "--fail", "--silent", "--show-error", "--connect-timeout", "10", "--max-time", "30",
             url.rstrip("/") + "/v1/nodefactory/nodes?isProvisioned=false"])
    data = json.loads(p.stdout)
    def find(value):
        if isinstance(value, list):
            return value
        if isinstance(value, dict):
            for key in ("nodes", "Nodes", "data", "Data", "items", "Items", "results", "Results"):
                if key in value:
                    found = find(value[key])
                    if found is not None:
                        return found
        return None
    result = {}
    for item in find(data) or []:
        node = next((str(item[k]) for k in ("id", "Id", "nodeId", "NodeId", "node_id") if item.get(k)), "")
        for kind in ("anode", "cnode", "tnode"):
            if f"-{kind}-" in node:
                result.setdefault(kind, node)
    return result


def image_archive(lab, repo, cache, env):
    base = env.get("BASE_IMAGE_REPO", "localhost/testing/virtualnode-base")
    runtime = env.get("NODE_RUNTIME", "starter")
    build_env = {**os.environ, **env, "UKAMA_REPO": str(repo)}
    tags = [f"{base}:{kind}-{runtime}" for kind in ("anode", "cnode", "tnode")]
    nodes = None
    for kind, tag in zip(("anode", "cnode", "tnode"), tags):
        if run(["podman", "image", "exists", tag], check=False).returncode:
            if nodes is None:
                nodes = factory_ids(env.get("UKAMA_LAB_FACTORY_URL", "http://factory-ukama.udev.ukama.com"))
            if kind not in nodes:
                raise RunnerError(f"missing {tag}; Factory has no {kind} ID to pass to the existing build-node.sh")
            print(f"Building missing starter base using scripts/build-node.sh: {tag}", flush=True)
            run([lab / "scripts/build-node.sh", repo, nodes[kind], runtime], capture=False,
                timeout=14400, env=build_env, cwd=lab)
        else:
            print(f"Reusing starter base: {tag}", flush=True)
    probe = env.get("ULAB_NET_PROBE_IMAGE", "docker.io/library/alpine:3.20")
    if run(["podman", "image", "exists", probe], check=False).returncode:
        run(["podman", "pull", probe], capture=False, timeout=1200)
    tags.append(probe)
    env["ULAB_NET_PROBE_IMAGE"] = probe
    # Carry existing UE/media images when present. Their existing scripts reuse
    # these images unless a rebuild is explicitly requested.
    for optional in ("localhost/ukama/ue:dev", "localhost/ukama/media:dev"):
        if not run(["podman", "image", "exists", optional], check=False).returncode:
            tags.append(optional)
    images = {}
    for tag in tags:
        info = json.loads(run(["podman", "image", "inspect", tag]).stdout)[0]
        if info.get("Architecture") != "amd64" or info.get("Os", "linux") != "linux":
            raise RunnerError(f"worker images must be Linux/amd64: {tag}")
        images[tag] = info["Id"].removeprefix("sha256:")
    fingerprint = hashlib.sha256(json.dumps(images, sort_keys=True).encode()).hexdigest()
    cache.mkdir(parents=True, exist_ok=True)
    target = cache / (fingerprint + ".tar.gz")
    metadata = cache / (fingerprint + ".json")
    if target.exists() and metadata.exists():
        try:
            saved = read_json(metadata)
            if saved.get("images") != images:
                raise RunnerError("cached image identities changed")
            check_compressed_image(target, saved)
            print("Reusing validated starter-image archive", flush=True)
            return target, saved
        except (RunnerError, OSError, ValueError, KeyError, TypeError) as exc:
            print(f"Discarding invalid image export cache: {exc}", flush=True)
            metadata.unlink(missing_ok=True)
            target.unlink(missing_ok=True)
    print("Exporting Podman images once (subsequent runs reuse this cache)...", flush=True)
    with tempfile.TemporaryDirectory(prefix="image-export-", dir=cache) as temporary:
        raw = Path(temporary) / "images.tar"
        compressed = Path(temporary) / "images.tar.gz"
        run(["podman", "save", "--format", "docker-archive", "--multi-image-archive",
             "--output", raw, *tags], capture=False, timeout=7200)
        print("Validating exported image tags, configurations and layers...", flush=True)
        validate_image_tar(raw, images)
        raw_info = {"cache_version": IMAGE_CACHE_VERSION, "tar_sha256": sha256(raw),
                    "tar_size": raw.stat().st_size, "images": images}
        with raw.open("rb") as source, compressed.open("wb") as target_stream:
            with gzip.GzipFile(filename="", mode="wb", fileobj=target_stream, compresslevel=1, mtime=0) as output:
                shutil.copyfileobj(source, output, 8 * 1024 * 1024)
        saved = {**raw_info, "sha256": sha256(compressed), "size": compressed.stat().st_size}
        check_compressed_image(compressed, saved)
        compressed.replace(target)
    atomic_json(metadata, saved)
    print(f"Image archive verified: {saved['size']:,} compressed bytes; {saved['tar_size']:,} tar bytes; {len(images)} images", flush=True)
    return target, saved


def bundle_kubectl(cfg, env, runtime):
    """Carry the caller's existing cleanup tool, without Kubernetes login/probes."""
    value = (cfg.get("KUBECTL_FILE") or cfg.get("ENV", {}).get("ULAB_KUBECTL")
             or os.environ.get("ULAB_KUBECTL"))
    if not value:
        home_tool = Path.home() / "kubectl"
        value = str(home_tool) if home_tool.is_file() else "kubectl"
    source = executable(os.path.expandvars(str(value)))
    destination = runtime / "kubectl"
    with source.open("rb") as stream:
        header = stream.read(4)
    if header == b"\x7fELF":
        bundle_elf(source, destination)
    elif header.startswith(b"#!"):
        destination.mkdir(parents=True, exist_ok=True)
        shutil.copy2(source, destination / "program")
        (destination / "run").write_text(
            '#!/bin/sh\nset -eu\nd="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"\n'
            'exec "$d/program" "$@"\n')
        (destination / "run").chmod(0o755)
    else:
        raise RunnerError("ULAB_KUBECTL must be a Linux executable or executable script")
    env["ULAB_KUBECTL"] = str(ROOT / "runtime/kubectl/run")
    env.pop("KUBECONFIG", None)

    # Kubectl may use an existing connection config even when VPN access needs
    # no login. Export only its selected context, embedding referenced files.
    # This command reads local config; it does not contact Kubernetes or execute
    # credential plugins. Do not impose authentication or connectivity gates.
    config_env = dict(os.environ)
    if cfg.get("KUBECONFIG_FILE"):
        config_env["KUBECONFIG"] = cfg["KUBECONFIG_FILE"]
    paths = [Path(p).expanduser() for p in config_env.get("KUBECONFIG", "").split(os.pathsep) if p]
    if not paths:
        paths = [Path.home() / ".kube/config"]
    if cfg.get("KUBECONFIG_FILE") and not paths[0].is_file():
        raise RunnerError(f"KUBECONFIG_FILE does not exist: {paths[0]}")
    if any(path.is_file() for path in paths):
        exported = run([source, "config", "view", "--raw", "--flatten", "--minify", "-o", "json"], env=config_env)
        try:
            contents = json.loads(exported.stdout)
            if not isinstance(contents, dict):
                raise ValueError("expected a configuration object")
        except ValueError as exc:
            raise RunnerError("kubectl could not export its existing connection configuration") from exc
        atomic_json(destination / "kubeconfig.json", contents)
        env["KUBECONFIG"] = str(ROOT / "runtime/kubectl/kubeconfig.json")
    print("Packaging existing kubectl for node-specific mesh cleanup" +
          (" (including current connection configuration)" if "KUBECONFIG" in env else ""), flush=True)
    return paths


def payload(lab, repo, lab_bin, cfg, env, work, excludes):
    with tempfile.TemporaryDirectory(prefix="runtime-", dir=work) as temporary:
        runtime = Path(temporary)
        bundle_elf(lab_bin, runtime / "lab")
        kube_paths = bundle_kubectl(cfg, env, runtime)
        excludes = [*excludes, *kube_paths]
        target = work / "source.tar.gz"
        print("Packaging current lab and Ukama repository...", flush=True)
        with tarfile.open(target, "w:gz", compresslevel=1) as archive:
            skipped = add_tree(archive, lab, "lab", excludes)
            # Avoid duplicating the lab when it lives inside the Ukama repository.
            skipped += add_tree(archive, repo, "repo", [*excludes, lab] if lab.is_relative_to(repo) else excludes, cfg["REPO_PATHS"])
            add_tree(archive, runtime, "runtime")
        atomic_json(work / "packaging.json", {"skipped_external_symlinks": skipped})
        return target


def extract(archive_path, destination):
    """Reject traversal, special files and escaping links before extracting."""
    destination = Path(destination).resolve()
    with tarfile.open(archive_path) as archive:
        members = archive.getmembers()
        links = {m.name for m in members if m.issym() or m.islnk()}
        for member in members:
            name = Path(member.name)
            if name.is_absolute() or ".." in name.parts or not (member.isfile() or member.isdir() or member.issym()):
                raise RunnerError("unsafe member in source archive")
            if any(str(parent) in links for parent in name.parents):
                raise RunnerError("source archive writes through a symlink")
            if member.issym():
                target = (destination / name.parent / member.linkname).resolve()
                if not target.is_relative_to(destination) or Path(member.linkname).is_absolute():
                    raise RunnerError("escaping symlink in source archive")
        # All names and symlinks validated, including on Python versions without filters.
        archive.extractall(destination, members=members)
