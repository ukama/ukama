# Test generation rules

1. Use playwright codegen to record test cases `npx playwright codegen http://localhost:3000/`
   a. While recording test cases, use input values format as `test-<test-name>`
   i.e: `test-network`, `test-user`, `test-site`, `test-node`.
   b. DO NOT modify the recorded test steps manually - any changes should be made through patches to ensure consistency and maintainability.
   c. Make sure steps not to be repeated.
   d. Record each test from scratch.
2. Create file under `console/autogen` dir according to the test name.
3. Add patch for the test case in `console/patches` dir (if needed).
4. Call the patch apply method in tests/apply_patches.spec.ts file.

## Apply patch to generate tests

To apply patches on auto gen files, run the `pnpm patch-tests` in terminal. This will generate tests in `tests/patched` dir.

## Run tests

To run tests, run the `pnpm test` in terminal.
