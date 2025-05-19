# Test generation rules

1. Use faker library for data generation
2. Use playwright codegen to create test cases `npx playwright codegen http://localhost:3000/`
3. While recording test cases, use input values format as `test-<test-name>`
i.e: `test-network`, `test-user`, `test-site`, `test-node`
4. While coping tests to generated_test.ts file, name the test accordingly. i.e: `test('test', async ({ page })` to `test('TEST_NAME Test', async ({ page })`
5. DO NOT edit the generated_test.ts and recorded tests manually.
6. Make sure steps not to be repeated.
7. Record each test from scratch.
8. Updating patch file with required changes.

## How to generate tests

After updating patch file with required changes, run the `pnpm test` in terminal.
