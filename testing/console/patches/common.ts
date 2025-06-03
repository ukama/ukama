import * as fs from 'fs';
import path from 'path';

const defaultImportStatement =
  "import { TEST_USER_EMAIL, TEST_USER_PASSWORD, LOGIN_URL, CONSOLE_ROOT_URL } from '@/constants';\nimport { faker } from '@faker-js/faker';\nimport { selectRandomOption } from '@/utils';\ntest.setTimeout(180000);";

export const applyPatch = async (
  fileName: string,
  version: string,
  testDir: string,
  customReplacements: { regex: RegExp; replacement: string }[] = [],
) => {
  const generatedTestPath = path.join(
    __dirname,
    `../autogen/${version}/${fileName}.autogen.ts`,
  );
  const patchedTestPath = path.join(
    __dirname,
    `../tests/patched/${version}/${testDir}/${fileName}.spec.ts`,
  );

  const originalContent = fs.readFileSync(generatedTestPath, 'utf8');
  const lines = originalContent.split('\n');
  const insertIndex = lines.findIndex((line) => line.includes('import')) + 1;
  lines.splice(insertIndex, 0, defaultImportStatement);
  const contentWithImport = lines.join('\n');

  let patchedContent = contentWithImport
    .replace(/'http:\/\/localhost:4455\/auth\/login'/g, 'LOGIN_URL')
    .replace(/'admin@ukama\.com'/g, 'TEST_USER_EMAIL')
    .replace(/'@Pass2025\.'/g, 'TEST_USER_PASSWORD');

  customReplacements.forEach(({ regex, replacement }) => {
    patchedContent = patchedContent.replace(regex, replacement);
  });

  const finalContent = patchedContent.split('\n').join('\n');

  const patchedDir = path.dirname(patchedTestPath);
  if (!fs.existsSync(patchedDir)) {
    fs.mkdirSync(patchedDir, { recursive: true });
  }

  fs.writeFileSync(patchedTestPath, finalContent);

  return { patchedContent: finalContent };
};
