#!/usr/bin/env bash
# Local controller for disposable EC2 P0 workers.
# Protocol: tar archives + text files in S3.

set -Eeuo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
LAB_ROOT="$(CDPATH= cd -- "$SCRIPT_DIR/../.." && pwd)"
# shellcheck source=lib.sh
. "$SCRIPT_DIR/lib.sh"
p0_load_config "$SCRIPT_DIR"

usage() {
    cat <<EOF_USAGE
usage:
  $0 [--workers N] [--dry-run] [--batch-id ID] [--scenario PATH ...]
  $0 [--workers N] [--dry-run] [--batch-id ID] --scenario-list FILE
  $0 [--workers N] [--dry-run] [--batch-id ID] [category ...]
  $0 --status BATCH_ID
  $0 --resume BATCH_ID
  $0 --collect BATCH_ID
  $0 --cleanup BATCH_ID

Normal execution packages the current ukama-lab tree and \$UKAMA_REPO,
launches disposable EC2 workers, displays live status, downloads all results,
prints the combined report, and terminates any workers still alive.
EOF_USAGE
}

WORKERS="${DEFAULT_WORKERS:-10}"
DRY_RUN=0
MODE=run
BATCH_ID=""
MODE_BATCH_ID=""
categories=()
explicit_scenarios=()
controller_scenario_list=""

while (($#)); do
    case "$1" in
        --workers)
            (($# >= 2)) || p0_die '--workers requires a number'
            WORKERS="$2"
            shift 2
            ;;
        --workers=*)
            WORKERS="${1#*=}"
            shift
            ;;
        --dry-run)
            DRY_RUN=1
            shift
            ;;
        --batch-id)
            (($# >= 2)) || p0_die '--batch-id requires a value'
            BATCH_ID="$2"
            shift 2
            ;;
        --batch-id=*)
            BATCH_ID="${1#*=}"
            shift
            ;;
        --scenario)
            (($# >= 2)) || p0_die '--scenario requires a path'
            explicit_scenarios+=("$2")
            shift 2
            ;;
        --scenario=*)
            explicit_scenarios+=("${1#*=}")
            shift
            ;;
        --scenario-list)
            (($# >= 2)) || p0_die '--scenario-list requires a file'
            controller_scenario_list="$2"
            shift 2
            ;;
        --scenario-list=*)
            controller_scenario_list="${1#*=}"
            shift
            ;;
        --status|--resume|--collect|--cleanup)
            (($# >= 2)) || p0_die "$1 requires a batch ID"
            MODE="${1#--}"
            MODE_BATCH_ID="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        --)
            shift
            categories+=("$@")
            break
            ;;
        -*)
            p0_die "unknown option: $1"
            ;;
        *)
            categories+=("$1")
            shift
            ;;
    esac
done

p0_require_cmd aws
p0_require_cmd jq
p0_require_cmd tar
p0_require_cmd sha256sum
p0_require_config AWS_REGION S3_BUCKET S3_PREFIX AMI_ID INSTANCE_TYPE \
    INSTANCE_PROFILE_NAME SUBNET_ID SECURITY_GROUP_ID SECRET_ID
p0_validate_simple_id S3_BUCKET "$S3_BUCKET"
p0_validate_simple_id S3_PREFIX "$S3_PREFIX"
p0_validate_simple_id AMI_ID "$AMI_ID"

export AWS_DEFAULT_REGION="$AWS_REGION"
p0_aws sts get-caller-identity >/dev/null

status_value() {
    local file="$1"
    local key="$2"
    awk -F '\t' -v wanted="$key" '$1 == wanted {print substr($0, index($0, $2)); exit}' "$file" 2>/dev/null
}

load_remote_manifest() {
    local batch_id="$1"
    local local_dir="$2"
    local uri

    uri="$(p0_s3_batch_root "$batch_id")"
    mkdir -p "$local_dir"
    if [[ ! -r "$local_dir/manifest.env" ]]; then
        p0_aws s3 cp "$uri/input/manifest.env" \
            "$local_dir/manifest.env" --only-show-errors
    fi
    # shellcheck disable=SC1090
    . "$local_dir/manifest.env"
}

sync_status() {
    local batch_id="$1"
    local local_dir="$2"
    local uri

    uri="$(p0_s3_batch_root "$batch_id")"
    mkdir -p "$local_dir/status"
    p0_aws s3 sync "$uri/output/" "$local_dir/status/" \
        --exclude '*' \
        --include '*.status' \
        --include '*.done' \
        --include '*.failed' \
        --only-show-errors >/dev/null || true
}

display_status() {
    local batch_id="$1"
    local local_dir="$2"
    local workers="$3"
    local i worker file state completed total passed failed skipped current
    local done_count=0
    local failed_worker_count=0
    local overall_completed=0
    local overall_total=0
    local overall_passed=0
    local overall_failed=0
    local overall_skipped=0

    if [[ -t 1 && "${P0_NO_CLEAR:-0}" != "1" ]]; then
        printf '\033[2J\033[H'
    fi

    printf 'Ukama distributed P0\n'
    printf 'Batch: %s\n\n' "$batch_id"
    printf '%-10s %-11s %9s %6s %6s %6s  %s\n' \
        WORKER STATE DONE PASS FAIL SKIP CURRENT

    for ((i = 1; i <= workers; i++)); do
        printf -v worker 'worker-%02d' "$i"
        file="$local_dir/status/$worker.status"
        state=LAUNCHING
        completed=0
        total=0
        passed=0
        failed=0
        skipped=0
        current=-

        if [[ -r "$file" ]]; then
            state="$(status_value "$file" state)"
            completed="$(status_value "$file" completed)"
            total="$(status_value "$file" total)"
            passed="$(status_value "$file" passed)"
            failed="$(status_value "$file" failed)"
            skipped="$(status_value "$file" skipped)"
            current="$(status_value "$file" current)"
            state="${state:-UNKNOWN}"
            completed="${completed:-0}"
            total="${total:-0}"
            passed="${passed:-0}"
            failed="${failed:-0}"
            skipped="${skipped:-0}"
            current="${current:--}"
        fi

        if [[ -f "$local_dir/status/$worker.done" ]]; then
            state=DONE
            done_count=$((done_count + 1))
        elif [[ -f "$local_dir/status/$worker.failed" ]]; then
            state=FAILED
            failed_worker_count=$((failed_worker_count + 1))
        fi

        overall_completed=$((overall_completed + completed))
        overall_total=$((overall_total + total))
        overall_passed=$((overall_passed + passed))
        overall_failed=$((overall_failed + failed))
        overall_skipped=$((overall_skipped + skipped))

        if ((${#current} > 72)); then
            current="...${current: -69}"
        fi
        printf '%-10s %-11s %4s/%-4s %6s %6s %6s  %s\n' \
            "$worker" "$state" "$completed" "$total" \
            "$passed" "$failed" "$skipped" "$current"
    done

    printf '\nOverall: %s/%s complete | pass=%s fail=%s skip=%s | workers done=%s failed=%s\n' \
        "$overall_completed" "$overall_total" "$overall_passed" \
        "$overall_failed" "$overall_skipped" "$done_count" \
        "$failed_worker_count"

    STATUS_COMPLETE_WORKERS=$((done_count + failed_worker_count))
    STATUS_FAILED_WORKERS=$failed_worker_count
}

watch_batch() {
    local batch_id="$1"
    local local_dir="$2"
    local workers="$3"
    local started now deadline
    local controller_minutes="${CONTROLLER_TIMEOUT_MINUTES:-$(( ${WORKER_TIMEOUT_MINUTES:-180} + 30 ))}"

    started="$(date +%s)"
    deadline=$((started + controller_minutes * 60))

    while :; do
        sync_status "$batch_id" "$local_dir"
        display_status "$batch_id" "$local_dir" "$workers"
        if ((STATUS_COMPLETE_WORKERS >= workers)); then
            return 0
        fi
        now="$(date +%s)"
        if ((now >= deadline)); then
            printf 'controller timeout after %s minutes\n' "$controller_minutes" >&2
            return 1
        fi
        sleep "${POLL_INTERVAL_SECONDS:-10}"
    done
}

collect_batch() {
    local batch_id="$1"
    local local_dir="$2"
    local workers="$3"
    local uri i worker worker_dir archive

    uri="$(p0_s3_batch_root "$batch_id")"
    mkdir -p "$local_dir/raw-output" "$local_dir/workers"
    p0_aws s3 sync "$uri/output/" "$local_dir/raw-output/" \
        --only-show-errors

    for ((i = 1; i <= workers; i++)); do
        printf -v worker 'worker-%02d' "$i"
        worker_dir="$local_dir/workers/$worker"
        archive="$local_dir/raw-output/$worker.tar.gz"
        mkdir -p "$worker_dir"

        for suffix in status log done failed; do
            if [[ -f "$local_dir/raw-output/$worker.$suffix" ]]; then
                cp "$local_dir/raw-output/$worker.$suffix" "$worker_dir/"
            fi
        done
        if [[ -f "$archive" ]]; then
            tar -xzf "$archive" -C "$worker_dir"
        fi
    done

    "$SCRIPT_DIR/report.sh" "$local_dir"
}

if [[ "$MODE" != run ]]; then
    BATCH_ID="$MODE_BATCH_ID"
    p0_validate_simple_id BATCH_ID "$BATCH_ID"
    LOCAL_BATCH_DIR="$(p0_local_batch_root "$LAB_ROOT" "$BATCH_ID")"

    case "$MODE" in
        cleanup)
            p0_cleanup_batch "$BATCH_ID"
            exit 0
            ;;
        status)
            load_remote_manifest "$BATCH_ID" "$LOCAL_BATCH_DIR"
            sync_status "$BATCH_ID" "$LOCAL_BATCH_DIR"
            display_status "$BATCH_ID" "$LOCAL_BATCH_DIR" "$WORKER_COUNT"
            exit 0
            ;;
        collect)
            load_remote_manifest "$BATCH_ID" "$LOCAL_BATCH_DIR"
            collect_batch "$BATCH_ID" "$LOCAL_BATCH_DIR" "$WORKER_COUNT"
            exit $?
            ;;
        resume)
            load_remote_manifest "$BATCH_ID" "$LOCAL_BATCH_DIR"
            watch_rc=0
            watch_batch "$BATCH_ID" "$LOCAL_BATCH_DIR" "$WORKER_COUNT" || watch_rc=$?
            report_rc=0
            collect_batch "$BATCH_ID" "$LOCAL_BATCH_DIR" "$WORKER_COUNT" || report_rc=$?
            p0_cleanup_batch "$BATCH_ID" || true
            ((watch_rc == 0 && report_rc == 0))
            exit $?
            ;;
    esac
fi

[[ "$WORKERS" =~ ^[1-9][0-9]*$ ]] || p0_die '--workers must be a positive integer'
[[ "${MAX_WORKERS:-20}" =~ ^[1-9][0-9]*$ ]] || p0_die 'MAX_WORKERS must be positive'
((WORKERS <= MAX_WORKERS)) || p0_die "workers=$WORKERS exceeds MAX_WORKERS=$MAX_WORKERS"

: "${UKAMA_REPO:?set UKAMA_REPO to the local Ukama source tree}"
[[ -d "$UKAMA_REPO" ]] || p0_die "UKAMA_REPO does not exist: $UKAMA_REPO"
[[ -d "$LAB_ROOT/scenarios/p0" ]] || p0_die "P0 scenarios not found under $LAB_ROOT"
[[ -x "$LAB_ROOT/bin/ukama-lab" ]] ||
    p0_die "build ukama-lab first; missing executable $LAB_ROOT/bin/ukama-lab"

# The udev control plane is private. Route workers through the known-good EC2
# host that already has the working backend tunnel/route.
if [[ -n "${BACKEND_GATEWAY_INSTANCE_ID:-}" ]]; then
    if ((DRY_RUN)); then
        BACKEND_GATEWAY_PRIVATE_IP="$(p0_aws ec2 describe-instances \
            --instance-ids "$BACKEND_GATEWAY_INSTANCE_ID" \
            --query 'Reservations[0].Instances[0].PrivateIpAddress' \
            --output text)"
        export BACKEND_GATEWAY_PRIVATE_IP
    elif [[ "${BACKEND_GATEWAY_AUTO_SETUP:-true}" == "true" ]]; then
        "$SCRIPT_DIR/setup-backend-gateway.sh" --quiet
        p0_load_config "$SCRIPT_DIR"
    fi
    p0_require_config BACKEND_GATEWAY_PRIVATE_IP BACKEND_ROUTE_CIDR \
        BACKEND_TEST_IP
fi

if [[ -z "$BATCH_ID" ]]; then
    BATCH_ID="$(date -u +%Y%m%dt%H%M%Sz)"
fi
p0_validate_simple_id BATCH_ID "$BATCH_ID"
LOCAL_BATCH_DIR="$(p0_local_batch_root "$LAB_ROOT" "$BATCH_ID")"
[[ ! -e "$LOCAL_BATCH_DIR" ]] || p0_die "local batch already exists: $LOCAL_BATCH_DIR"
mkdir -p "$LOCAL_BATCH_DIR/input/shards"

SCENARIO_LIST="$LOCAL_BATCH_DIR/input/scenarios.txt"
: >"$SCENARIO_LIST"

if [[ -n "$controller_scenario_list" && ${#explicit_scenarios[@]} -gt 0 ]]; then
    p0_die '--scenario and --scenario-list cannot be combined'
fi
if [[ -n "$controller_scenario_list" || ${#explicit_scenarios[@]} -gt 0 ]]; then
    ((${#categories[@]} == 0)) ||
        p0_die 'categories cannot be combined with --scenario or --scenario-list'
fi

add_scenario_path() {
    local supplied="$1"
    local absolute
    local relative

    supplied="${supplied%%#*}"
    supplied="$(sed 's/^[[:space:]]*//;s/[[:space:]]*$//' <<<"$supplied")"
    [[ -n "$supplied" ]] || return 0

    if [[ "$supplied" == /* ]]; then
        absolute="$supplied"
    else
        absolute="$LAB_ROOT/$supplied"
    fi
    absolute="$(readlink -f -- "$absolute")" ||
        p0_die "scenario does not exist: $supplied"
    [[ -f "$absolute" ]] || p0_die "scenario is not a file: $supplied"
    case "$absolute" in
        "$LAB_ROOT/scenarios/p0/"*.yaml|"$LAB_ROOT/scenarios/p0/"*.yml)
            ;;
        *)
            p0_die "scenario must be below scenarios/p0: $supplied"
            ;;
    esac
    relative="${absolute#"$LAB_ROOT/"}"
    printf '%s\n' "$relative" >>"$SCENARIO_LIST"
}

if [[ -n "$controller_scenario_list" ]]; then
    [[ -r "$controller_scenario_list" ]] ||
        p0_die "scenario list is not readable: $controller_scenario_list"
    while IFS= read -r scenario || [[ -n "$scenario" ]]; do
        add_scenario_path "$scenario"
    done <"$controller_scenario_list"
elif ((${#explicit_scenarios[@]} > 0)); then
    for scenario in "${explicit_scenarios[@]}"; do
        add_scenario_path "$scenario"
    done
elif ((${#categories[@]} == 0)); then
    find "$LAB_ROOT/scenarios/p0" -type f \
        \( -name '*.yaml' -o -name '*.yml' \) -print |
        sort >"$SCENARIO_LIST.absolute"
else
    for category in "${categories[@]}"; do
        [[ "$category" =~ ^[a-z0-9-]+$ ]] || p0_die "invalid category: $category"
        [[ -d "$LAB_ROOT/scenarios/p0/$category" ]] || p0_die "unknown P0 category: $category"
        find "$LAB_ROOT/scenarios/p0/$category" -type f \
            \( -name '*.yaml' -o -name '*.yml' \) -print
    done | sort -u >"$SCENARIO_LIST.absolute"
fi

if [[ -f "$SCENARIO_LIST.absolute" ]]; then
    while IFS= read -r path; do
        printf '%s\n' "${path#"$LAB_ROOT/"}" >>"$SCENARIO_LIST"
    done <"$SCENARIO_LIST.absolute"
    rm -f "$SCENARIO_LIST.absolute"
else
    sort -u -o "$SCENARIO_LIST" "$SCENARIO_LIST"
fi

SCENARIO_COUNT="$(wc -l <"$SCENARIO_LIST" | tr -d ' ')"
((SCENARIO_COUNT > 0)) || p0_die 'no scenarios selected'
if ((WORKERS > SCENARIO_COUNT)); then
    WORKERS=$SCENARIO_COUNT
fi

EXCLUSIVE_FILE="${P0_EXCLUSIVE_FILE:-$SCRIPT_DIR/exclusive.txt}"
NORMAL_LIST="$LOCAL_BATCH_DIR/input/normal.txt"
EXCLUSIVE_LIST="$LOCAL_BATCH_DIR/input/exclusive.txt"
: >"$NORMAL_LIST"
: >"$EXCLUSIVE_LIST"

exclusive_patterns=()
if [[ -r "$EXCLUSIVE_FILE" ]]; then
    while IFS= read -r pattern; do
        pattern="${pattern%%#*}"
        pattern="$(sed 's/^[[:space:]]*//;s/[[:space:]]*$//' <<<"$pattern")"
        [[ -n "$pattern" ]] || continue
        exclusive_patterns+=("$pattern")
    done <"$EXCLUSIVE_FILE"
fi

is_exclusive() {
    local candidate="$1"
    local pattern

    for pattern in "${exclusive_patterns[@]}"; do
        # Intentional shell-pattern matching from the trusted exclusive file.
        if [[ "$candidate" == $pattern ]]; then
            return 0
        fi
    done
    return 1
}

while IFS= read -r path; do
    if is_exclusive "$path"; then
        printf '%s\n' "$path" >>"$EXCLUSIVE_LIST"
    else
        printf '%s\n' "$path" >>"$NORMAL_LIST"
    fi
done <"$SCENARIO_LIST"

normal_count="$(wc -l <"$NORMAL_LIST" | tr -d ' ')"
exclusive_count="$(wc -l <"$EXCLUSIVE_LIST" | tr -d ' ')"
exclusive_worker=0
if ((exclusive_count > 0 && WORKERS > 1)); then
    exclusive_worker=1
fi

normal_workers=$WORKERS
if ((exclusive_worker)); then
    normal_workers=$((WORKERS - 1))
fi
if ((normal_count == 0)); then
    normal_workers=0
elif ((normal_workers > normal_count)); then
    normal_workers=$normal_count
fi

worker_index=0
if ((normal_workers > 0)); then
    i=0
    while IFS= read -r scenario; do
        worker_index=$((i % normal_workers + 1))
        printf -v worker_id 'worker-%02d' "$worker_index"
        printf '%s\n' "$scenario" \
            >>"$LOCAL_BATCH_DIR/input/shards/$worker_id.txt"
        i=$((i + 1))
    done <"$NORMAL_LIST"
fi

WORKER_COUNT=$normal_workers
if ((exclusive_count > 0)); then
    WORKER_COUNT=$((WORKER_COUNT + 1))
    printf -v worker_id 'worker-%02d' "$WORKER_COUNT"
    cat "$EXCLUSIVE_LIST" >"$LOCAL_BATCH_DIR/input/shards/$worker_id.txt"
fi
((WORKER_COUNT > 0)) || p0_die 'no worker shards were created'

cat >"$LOCAL_BATCH_DIR/input/manifest.env" <<EOF_MANIFEST
BATCH_ID=$BATCH_ID
WORKER_COUNT=$WORKER_COUNT
SCENARIO_COUNT=$SCENARIO_COUNT
CREATED_AT=$(date -u +%Y-%m-%dT%H:%M:%SZ)
EOF_MANIFEST

{
    printf 'AWS_REGION=%q\n' "$AWS_REGION"
    printf 'SECRET_ID=%q\n' "$SECRET_ID"
    printf 'WORKER_TIMEOUT_MINUTES=%q\n' "${WORKER_TIMEOUT_MINUTES:-180}"
    printf 'STATUS_INTERVAL_SECONDS=%q\n' "${STATUS_INTERVAL_SECONDS:-10}"
    printf 'UKAMA_LAB_SIM_TYPE=%q\n' "${UKAMA_LAB_SIM_TYPE:-ukama_data}"
    printf 'UKAMA_LAB_WAREHOUSE_URL=%q\n' "${UKAMA_LAB_WAREHOUSE_URL:-http://warehouse-ukama.udev.ukama.com}"
    printf 'UKAMA_LAB_FACTORY_URL=%q\n' "${UKAMA_LAB_FACTORY_URL:-http://factory-ukama.udev.ukama.com}"
    printf 'ULAB_FACTORY_SEED_URL=%q\n' "${ULAB_FACTORY_SEED_URL:-https://factory-ukama.udev.ukama.com}"
    printf 'ULAB_UKAMA_AGENT_NODE_GW_URL=%q\n' "${ULAB_UKAMA_AGENT_NODE_GW_URL:-https://ukamaagent-nodegateway-ukama.udev.ukama.com:8080/}"
    printf 'ULAB_UKAMA_AGENT_API_GW_URL=%q\n' "${ULAB_UKAMA_AGENT_API_GW_URL:-https://ukamaagent-ukama.udev.ukama.com:8080/}"
    printf 'BACKEND_GATEWAY_PRIVATE_IP=%q\n' "${BACKEND_GATEWAY_PRIVATE_IP:-}"
    printf 'BACKEND_ROUTE_CIDR=%q\n' "${BACKEND_ROUTE_CIDR:-10.0.0.0/8}"
    printf 'BACKEND_TEST_IP=%q\n' "${BACKEND_TEST_IP:-10.0.101.241}"
    printf 'P0_CONNECT_TIMEOUT_SECONDS=%q\n' "${P0_CONNECT_TIMEOUT_SECONDS:-15}"
    printf 'P0_HTTP_TIMEOUT_SECONDS=%q\n' "${P0_HTTP_TIMEOUT_SECONDS:-25}"
    printf 'UKAMA_LAB_DUMP_BFF_CURL=%q\n' "${UKAMA_LAB_DUMP_BFF_CURL:-1}"
    printf 'ULAB_CDR_WAIT_SEC=%q\n' "${ULAB_CDR_WAIT_SEC:-5}"
    printf 'ULAB_CDR_DIAG_STRICT=%q\n' "${ULAB_CDR_DIAG_STRICT:-0}"
    printf 'ULAB_OVS_ALLOW_UNMETERED_FALLBACK=%q\n' "${ULAB_OVS_ALLOW_UNMETERED_FALLBACK:-1}"
    printf 'ULAB_CDR_DIAG_DISABLE=%q\n' "${ULAB_CDR_DIAG_DISABLE:-1}"
    printf 'ULAB_SOFTWARE_CURRENT_VERSION=%q\n' "${ULAB_SOFTWARE_CURRENT_VERSION:-}"
    printf 'ULAB_SOFTWARE_TARGET_VERSION=%q\n' "${ULAB_SOFTWARE_TARGET_VERSION:-}"
    printf 'ULAB_SOFTWARE_UNAVAILABLE_VERSION=%q\n' "${ULAB_SOFTWARE_UNAVAILABLE_VERSION:-}"
    printf 'ULAB_SOFTWARE_NEXT_VERSION=%q\n' "${ULAB_SOFTWARE_NEXT_VERSION:-}"
} >"$LOCAL_BATCH_DIR/input/worker.env"

printf '\nBatch: %s\n' "$BATCH_ID"
printf 'Scenarios: %s\n' "$SCENARIO_COUNT"
printf 'Workers: %s\n' "$WORKER_COUNT"
for shard in "$LOCAL_BATCH_DIR"/input/shards/*.txt; do
    printf '  %-12s %s scenario(s)\n' \
        "$(basename "$shard" .txt)" "$(wc -l <"$shard" | tr -d ' ')"
done

if [[ "${FACTORY_NODE_TARGET:-0}" != "0" ]]; then
    CREDENTIALS_FILE="${P0_AWS_CREDENTIALS:-$SCRIPT_DIR/credentials.env}"
    [[ -r "$CREDENTIALS_FILE" ]] ||
        p0_die "factory preparation requires $CREDENTIALS_FILE"
    set -a
    # shellcheck disable=SC1090
    . "$CREDENTIALS_FILE"
    set +a
    printf '\nPreparing factory nodes locally...\n'
    (
        cd "$LAB_ROOT"
        ULAB_FACTORY_HEADROOM_PERCENT="${FACTORY_HEADROOM_PERCENT:-10}" \
        ./utils/run-p0-scenarios.sh \
            --scenario-list "$SCENARIO_LIST" \
            --factory-nodes "$FACTORY_NODE_TARGET" \
            --prepare-only \
            --batch-id "factory-$BATCH_ID"
    )
fi

printf '\nPackaging current local source trees...\n'
tar --exclude-from="$SCRIPT_DIR/lab-excludes.txt" \
    -C "$LAB_ROOT" -czf "$LOCAL_BATCH_DIR/input/ukama-lab.tar.gz" .

ukama_tar_args=(--exclude-from="$SCRIPT_DIR/ukama-excludes.txt")
case "$LAB_ROOT/" in
    "$UKAMA_REPO"/*)
        lab_relative="${LAB_ROOT#"$UKAMA_REPO"/}"
        ukama_tar_args+=(
            --exclude="$lab_relative"
            --exclude="$lab_relative/**"
        )
        ;;
esac

tar "${ukama_tar_args[@]}" \
    -C "$UKAMA_REPO" -czf "$LOCAL_BATCH_DIR/input/ukama.tar.gz" .
cp "$SCRIPT_DIR/worker.sh" "$LOCAL_BATCH_DIR/input/worker.sh"

(
    cd "$LOCAL_BATCH_DIR/input"
    sha256sum ukama-lab.tar.gz ukama.tar.gz >checksums.sha256
)

if ((DRY_RUN)); then
    printf '\nDry run complete. Nothing was uploaded or launched.\n'
    printf 'plan: %s\n' "$LOCAL_BATCH_DIR"
    exit 0
fi

BATCH_URI="$(p0_s3_batch_root "$BATCH_ID")"
printf 'Uploading input to %s/input/\n' "$BATCH_URI"
p0_aws s3 sync "$LOCAL_BATCH_DIR/input/" "$BATCH_URI/input/" \
    --only-show-errors

p0_aws ec2 describe-images --image-ids "$AMI_ID" \
    --query 'Images[0].ImageId' --output text | grep -q "^$AMI_ID$" ||
    p0_die "AMI is unavailable in $AWS_REGION: $AMI_ID"

ROOT_DEVICE_NAME="$(p0_aws ec2 describe-images --image-ids "$AMI_ID" \
    --query 'Images[0].RootDeviceName' --output text)"
block_file="$LOCAL_BATCH_DIR/block-device.json"
network_file="$LOCAL_BATCH_DIR/network-interface.json"
jq -n \
    --arg device "$ROOT_DEVICE_NAME" \
    --arg type "${ROOT_VOLUME_TYPE:-gp3}" \
    --argjson size "${ROOT_VOLUME_GB:-100}" \
    '[{DeviceName:$device,Ebs:{VolumeSize:$size,VolumeType:$type,DeleteOnTermination:true,Encrypted:true}}]' \
    >"$block_file"
jq -n \
    --arg subnet "$SUBNET_ID" \
    --arg group "$SECURITY_GROUP_ID" \
    --argjson public_ip "${ASSOCIATE_PUBLIC_IP:-true}" \
    '[{DeviceIndex:0,SubnetId:$subnet,Groups:[$group],AssociatePublicIpAddress:$public_ip,DeleteOnTermination:true}]' \
    >"$network_file"

printf 'worker\tinstance_id\n' >"$LOCAL_BATCH_DIR/instances.tsv"
launched=1
cleanup_on_exit() {
    local rc=$?
    if ((launched)); then
        p0_cleanup_batch "$BATCH_ID" || true
    fi
    exit "$rc"
}
trap cleanup_on_exit EXIT INT TERM

printf '\nLaunching %s EC2 workers...\n' "$WORKER_COUNT"
for ((i = 1; i <= WORKER_COUNT; i++)); do
    printf -v worker_id 'worker-%02d' "$i"
    user_data="$LOCAL_BATCH_DIR/$worker_id-user-data.sh"
    cat >"$user_data" <<EOF_USER_DATA
#!/usr/bin/env bash
set -Eeuo pipefail
exec > >(tee -a /var/log/ukama-p0-bootstrap.log) 2>&1
trap '/usr/sbin/shutdown -P now || true' EXIT
export AWS_DEFAULT_REGION=$(printf '%q' "$AWS_REGION")
if ! command -v aws >/dev/null 2>&1; then
    if command -v dnf >/dev/null 2>&1; then
        dnf install -y awscli
    elif command -v apt-get >/dev/null 2>&1; then
        apt-get update
        DEBIAN_FRONTEND=noninteractive apt-get install -y awscli
    fi
fi
aws s3 cp $(printf '%q' "$BATCH_URI/input/worker.sh") /tmp/ukama-p0-worker.sh --only-show-errors
chmod 700 /tmp/ukama-p0-worker.sh
BATCH_ID=$(printf '%q' "$BATCH_ID") \\
WORKER_ID=$(printf '%q' "$worker_id") \\
S3_BUCKET=$(printf '%q' "$S3_BUCKET") \\
S3_PREFIX=$(printf '%q' "$S3_PREFIX") \\
AWS_REGION=$(printf '%q' "$AWS_REGION") \\
/tmp/ukama-p0-worker.sh
EOF_USER_DATA

    instance_id="$(p0_aws ec2 run-instances \
        --image-id "$AMI_ID" \
        --instance-type "$INSTANCE_TYPE" \
        --iam-instance-profile "Name=$INSTANCE_PROFILE_NAME" \
        --network-interfaces "file://$network_file" \
        --block-device-mappings "file://$block_file" \
        --metadata-options 'HttpEndpoint=enabled,HttpTokens=required,HttpPutResponseHopLimit=2' \
        --instance-initiated-shutdown-behavior terminate \
        --user-data "file://$user_data" \
        --tag-specifications \
            "ResourceType=instance,Tags=[{Key=Name,Value=ukama-p0-$BATCH_ID-$worker_id},{Key=UkamaP0Batch,Value=$BATCH_ID},{Key=UkamaP0Worker,Value=$worker_id},{Key=Ephemeral,Value=true}]" \
        --query 'Instances[0].InstanceId' --output text)"
    printf '%s\t%s\n' "$worker_id" "$instance_id" \
        >>"$LOCAL_BATCH_DIR/instances.tsv"
    printf '  %-10s %s\n' "$worker_id" "$instance_id"
done

watch_rc=0
watch_batch "$BATCH_ID" "$LOCAL_BATCH_DIR" "$WORKER_COUNT" || watch_rc=$?
report_rc=0
collect_batch "$BATCH_ID" "$LOCAL_BATCH_DIR" "$WORKER_COUNT" || report_rc=$?

p0_cleanup_batch "$BATCH_ID" || true
launched=0
trap - EXIT INT TERM

if [[ "${DELETE_S3_AFTER_COLLECT:-false}" == "true" ]]; then
    p0_aws s3 rm "$BATCH_URI/" --recursive --only-show-errors
fi

((watch_rc == 0 && report_rc == 0))
