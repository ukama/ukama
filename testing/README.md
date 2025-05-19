# Test generation rules

1. Use faker to generate random data
2. Use playwright to create test cases
3. While creating test cases, use format `test-<test-name>` as input values.
i.e: `test-network`, `test-user`, `test-site`, `test-node`
4. While coping tests to generated_test.ts file, name the test accordingly. i.e: `test('test', async ({ page })` to `test('TEST_NAME Test', async ({ page })`
5. Make sure steps not to be repeated.
6. Record each test from scratch.