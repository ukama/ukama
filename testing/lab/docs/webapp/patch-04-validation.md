# Patch 4 validation record

- Complete 55-source C CLI compiled and linked with `-Wall -Wextra -Werror`,
  project C11/POSIX flags and the existing per-workload warning/optimization
  exception. Host libraries: Jansson 2.14 headers with installed libjansson,
  libyaml 0.2.5 headers with installed libyaml, curl, pthread and math.
- 28 Python/C tests passed: 13 contract, 7 resource/runtime lifecycle and
  8 full CLI lifecycle tests. This includes legacy P0/smoke contract checks,
  invalid dynamic expectation references/properties, and separate coverage credit.
- Strict TypeScript build and 25 existing worker tests passed.
- 3 existing C CLI-to-Chromium lifecycle tests passed, including five iterations
  of active cancellation with failure evidence.
- 8 new C CLI-to-Chromium provisioning tests passed: two-site automatic recovery;
  all three node cards/details; wrong-tower rejection before submission;
  network and site UI timeouts after successful mutations; partial runtime
  setup failure; false teardown success; and unidentified submitted creation.
- The successful recovery test issued exactly one network creation and two site
  creations, observed 2/2 → 1/2 → 2/2, fetched the document only once, and deleted
  all nine owned resources while preserving the pre-existing network.
- The node journey passed 25 visible field assertions, using actual runtime IDs
  rather than internal lab node IDs. It verified every trio member in list/detail.
- All three executable webapp examples passed offline lint. The modified shell
  script passed `sh -n`; both repository diffs passed whitespace checks.
- All four console companion TSX files transpiled with TypeScript 5.9.3 without
  syntax diagnostics. Full console dependency-aware typecheck/build was not run;
  console dependencies and the live backend/auth stack were unavailable.

Runtime used Node.js 24.19.0, Playwright 1.62.1 and local Chromium 153.0.8010.0
via `ULAB_WEBAPP_EXECUTABLE_PATH`. The browser is an explicit override, not a
qualification of every bundled Playwright browser. Firefox/WebKit were not run.

All browser results above used controlled source-derived DOM/runtime/BFF fixtures.
No live console/auth/BFF/Factory/Podman acceptance run or coverage percentage is
claimed. Factory reservation/start scripts require host qualification; the tests
substitute runtime scripts and exercise the real C orchestration and teardown.
The inventory records five implemented requirements with verification `not_run`.

The deliverable packaging also applies each unified diff to its clean baseline
and compares every changed file's content and executable mode with the archive.
