# Serial browser matrix and requirement report

`utils/webapp/matrix.py` runs selected v2 scenarios once per browser and aggregates
complete scenario evidence against the fixed `docs/webapp/coverage.json` inventory.
It supports Chromium, Firefox and WebKit. Runs are serial, with unique run IDs,
private output folders, per-attempt logs and the normal C runner cleanup. It does
not deploy a stack or create accounts, retry mutations, or skip a browser whose
executable is missing.

## Prepare

Build the lab CLI and worker as in `runner.md` / `commerce.md`, then:

```sh
python3 -m pip install -r utils/webapp/requirements.txt
npm --prefix adapters/webapp ci
npm --prefix adapters/webapp run build
cd adapters/webapp
npx playwright install chromium firefox webkit
cd ../..
```

Set each selected YAML's `webapp.base_url`, `auth_state` and `headless` for the
host. Keep changes in your checkout: source/profile hashes identify the cohort.
Capture fresh authentication with the existing auth command. Review the selected
scenario's prerequisites, including runtime endpoints and independently specified
usage/software expectations. Export any `${VARIABLE}` used by a scenario, both
for execution and report generation. Only a hash of these inputs enters the
matrix manifest. Credentials and raw environment values are not copied there.

Choose explicit scenario paths for a small run. Omitting `--scenario` selects all
scenarios referenced by the inventory, including partial requirements. Omitting
`--browser` selects all three browsers. Non-active scenarios remain skipped and
receive no credit. Missing browsers, auth, services or prerequisites fail the
attempt and do not silently reduce the requested matrix.

```sh
python3 utils/webapp/matrix.py run \
  --binary ./bin/ukama-lab \
  --out ./out/matrix-001 \
  --browser chromium \
  --scenario scenarios/webapp/p0/expanded/wb-032-network-plan-isolation.yaml \
  --evidence live --app-build CONSOLE_COMMIT --backend-build BACKEND_BUILD \
  -- --repo /path/to/ukama --scripts ./scripts \
  --warehouse-url http://localhost:8080 \
  --factory-url http://localhost:8081 --asr-url http://localhost:8082 \
  --bff http://localhost:8000/graphql
```

Replace endpoints with those of your local stack. Only the six endpoint/repo/
scripts flags shown may be forwarded after `--`; the matrix owns run ID, worker
and output flags. `--evidence` defaults to `unverified`. Use `fixture` for
controlled responses. `live` requires both build identifiers and represents an
operator declaration that the run uses the real target stack, without intercepted
or fabricated product responses. The tool cannot independently attest to that
declaration or a remote build. Build hashes should identify the deployed app and
backend, and the lab binary must be freshly built from the checkout.

A run exits nonzero if any process fails or a selected scenario is skipped. On
Ctrl+C or a matrix deadline it asks the C runner to terminate and clean up, then
kills the owned process group if the grace period expires. Inspect failed cleanup
journals before another run. No output directory is overwritten. Every planned
attempt is recorded before execution, so cancellation leaves pending attempts
visible. Browser reports now include the declared browser profile.

## Generate the report

```sh
python3 utils/webapp/matrix.py report \
  --manifest ./out/matrix-001/matrix.json \
  --out ./out/coverage-001 \
  --browser chromium \
  --app-build CONSOLE_COMMIT --backend-build BACKEND_BUILD \
  --gate
```

Repeat `--manifest` to include later runs and `--browser` for every required
browser. The command writes `coverage.json` and a standalone, escaped
`coverage.html` without external scripts or services. Without `--gate`, incomplete
coverage is reported with a normal zero command exit. With `--gate`, success
requires **100% of the 86 P0 requirements and at least 90% of all 132 requirements**
to be live verified on every requested browser. These thresholds are fixed; the
report does not shrink its denominator to the selected scenarios.

An inventory-only snapshot needs no manifest:

```sh
python3 utils/webapp/matrix.py report --out ./out/coverage-not-run --gate
```

This intentionally exits 1 and records zero live verification.

## Evidence rules

A requirement earns credit only when marked `implemented` and all its mapped
scenarios pass on every requested browser for the same app/backend build pair.
A mapped scenario must contain an assertion for that requirement. Every planned
check and event must appear with matching phase/type/requirement/label counts;
a bare PASS flag, empty report, missing assertion or duplicate result cannot pass.
The process result, report outcome, counters, browser, run identity and cleanup
must all agree. Expected and actual values must be present.

Original and generated scenario hashes, the inventory, lab source cohort,
referenced environment inputs and report file hashes must match. A matrix-generated
scenario may differ from its source only in the browser line. A stale, modified,
missing, skipped, partial, failed or unverified artifact receives no live credit.
The manifest also records the launched binary hash; this is provenance, not proof
that the binary matches source. These local artifacts are not cryptographically
signed attestations.

The latest attempt is authoritative for a scenario/browser within the selected
build pair. Passing and failing live attempts in that same current source cohort
are labeled flaky and withhold credit even if the last run passes. Include all
attempt manifests when reviewing a release; omitting a failed manifest hides that
history from any offline aggregator. Existing direct C reports without a matrix
manifest remain useful diagnostics but do not qualify for verified coverage.

The HTML contains requirement text, evidence classifications and gap notes. Raw
traces, screenshots, auth state and observed customer/payment values remain in
private per-run artifacts; the report does not embed them.

## Onboarding injected faults

Patch 9 scenarios containing `web_onboard` / `arm_fault` are controlled UI
evidence. The report derives this from the unchanged scenario, overrides a
`live` declaration with `controlled_ui`, and withholds live coverage credit.
Real service-side fault qualification remains required; see `onboarding.md`.

Patch 11 applies the same rule to every `web_action` / `status_fault` scenario,
including fault removal. A passing injected operation read-error or optimistic
fallback journey is controlled UI evidence and cannot become live credit by
changing the manifest classification. Backend lease-timeout qualification remains
separate; see `operation-recovery.md`.
