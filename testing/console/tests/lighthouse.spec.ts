import { LOGIN_URL, TEST_USER_EMAIL, TEST_USER_PASSWORD } from '@/constants';
import { expect, test } from '@playwright/test';
import chromeLauncher from 'chrome-launcher';
import CDP from 'chrome-remote-interface';
import fs from 'fs';
import lighthouse from 'lighthouse';
import path from 'path';
import { URL } from 'url';
// URLs to audit
const authenticatedUrls = [
  'http://localhost:3000/console/home',
  'http://localhost:3000/console/sites',
  'http://localhost:3000/console/subscribers',
  'http://localhost:3000/console/nodes',
  'http://localhost:3000/manage/billing',
  'http://localhost:3000/manage/data-plans',
  'http://localhost:3000/manage/sims',
  'http://localhost:3000/manage/nodes',
  'http://localhost:3000/manage/members',
];

// Define a type for the audit result for better type safety
interface LighthouseAuditResult {
  lhr: any; // Replace 'any' with a more specific Lighthouse report type if you have one
  report: string; // The HTML report content
}
// Function to run a single Lighthouse audit
async function runLighthouseAudit(
  url: string,
  authData: { cookies: any[]; localStorage: any[] }, // Accept authentication data
): Promise<LighthouseAuditResult | null> {
  const chrome = await chromeLauncher.launch({
    chromeFlags: [
      '--headless',
      '--disable-gpu',
      '--no-sandbox',
      '--disable-dev-shm-usage',
    ],
    logLevel: 'info',
  });
  // Connect to the launched Chrome instance
  const protocol = new URL(url).protocol;
  const hostname = new URL(url).hostname;
  const cdpConnection = await CDP({ port: chrome.port });
  const { Network: network, Emulation: emulation, Page: page } = cdpConnection;
  try {
    // Set cookies and local storage
    await network.enable();
    await page.enable();
    await emulation.setDeviceMetricsOverride({
      width: 1440,
      height: 1230,
      deviceScaleFactor: 1,
      mobile: false,
    });

    // Apply cookies
    for (const cookie of authData.cookies) {
      await network.setCookie({
        name: cookie.name,
        value: cookie.value,
        url: cookie.url || url, // Use the cookie's URL or the target URL
        domain: cookie.domain || hostname,
        path: cookie.path || '/',
        secure: cookie.secure || protocol === 'https:',
        httpOnly: cookie.httpOnly || false,
        sameSite: cookie.sameSite || 'None',
        expires: cookie.expires,
      });
    }

    // Apply local storage (requires evaluating JavaScript in the page context)
    // This is a bit more complex as Lighthouse controls the page lifecycle.
    // A common approach is to use a Lighthouse custom gatherer or a script
    // that runs before the audit. For simplicity in this example, we'll
    // focus on cookies which are usually sufficient for session management.
    // If you heavily rely on local storage for auth, you'll need a more
    // advanced approach.

    const runnerResult = await lighthouse(url, {
      port: chrome.port,
      output: ['html', 'json'],
      onlyCategories: [
        'performance',
        'accessibility',
        'best-practices',
        'seo',
        'pwa',
      ],
      screenEmulation: {
        mobile: false,
        width: 1440,
        height: 1230,
        deviceScaleFactor: 1,
        disabled: false,
      },
      formFactor: 'desktop',
      screenWidth: 1440,
      screenHeight: 1230,
    } as any);

    if (!runnerResult) {
      console.error(`Lighthouse audit failed for ${url}`);
      return null;
    }

    return {
      lhr: runnerResult.lhr,
      report: Array.isArray(runnerResult.report)
        ? runnerResult.report[0]
        : runnerResult.report,
    };
  } catch (error) {
    console.error(`Error running Lighthouse audit for ${url}:`, error);
    return null;
  } finally {
    await cdpConnection.close(); // Close CDP connection
    await chrome.kill();
  }
}
// Save report to file
function saveReport(url: string, result: LighthouseAuditResult) {
  const reportDir = path.join(
    `${process.cwd()}/testing/console`,
    'lighthouse-reports',
    url.replace(/[^a-zA-Z0-9]/g, '_'),
  );
  if (!fs.existsSync(reportDir)) {
    fs.mkdirSync(reportDir, { recursive: true });
  }
  const htmlFileName = 'report.html';
  const jsonFileName = 'report.json';
  const htmlPath = path.join(reportDir, htmlFileName);
  const jsonPath = path.join(reportDir, jsonFileName);
  try {
    fs.writeFileSync(htmlPath, result.report);
    console.log(`HTML report for ${url} saved to ${htmlPath}`);
    fs.writeFileSync(jsonPath, JSON.stringify(result.lhr, null, 2));
    console.log(`JSON report for ${url} saved to ${jsonPath}`);
  } catch (error) {
    console.error(`Failed to save report for ${url}:`, error);
  }
}
test.describe('Lighthouse Audits', () => {
  test.setTimeout(120000); // Set timeout to 2 minutes
  let authData: { cookies: any[]; localStorage: any[] }; // Declare variable to store auth data
  test.beforeEach(async ({ page }) => {
    // Login before each test
    console.log('Navigating to login page:', LOGIN_URL);
    await page.goto(LOGIN_URL);

    // Add a wait for the page to be ready
    await page.waitForLoadState('networkidle');

    const emailInput = page.getByRole('textbox', { name: 'EMAIL' });
    await emailInput.waitFor({ state: 'visible', timeout: 10000 });

    await emailInput.fill(TEST_USER_EMAIL);
    await emailInput.press('Tab');

    await page
      .getByRole('textbox', { name: 'PASSWORD' })
      .fill(TEST_USER_PASSWORD);

    await page.getByRole('button', { name: 'LOG IN' }).click();

    await page.waitForURL('**/console/home', { timeout: 30000 });

    // **Capture authentication state**
    authData = {
      cookies: await page.context().cookies(),
      localStorage: [], // Placeholder
    };
  });
  for (const url of authenticatedUrls) {
    test(`Lighthouse audit for ${url}`, async ({ page }) => {
      // Pass authentication data to the Lighthouse function
      const result = await runLighthouseAudit(url, authData);
      // If the audit failed, fail the test
      expect(result).not.toBeNull();
      if (!result) {
        return;
      }

      // Save report
      saveReport(url, result);

      // Log scores for debugging
      console.log('Lighthouse Scores:', {
        performance: result.lhr.categories.performance?.score,
        accessibility: result.lhr.categories.accessibility?.score,
        bestPractices: result.lhr.categories['best-practices']?.score,
        seo: result.lhr.categories.seo?.score,
      });

      // Assert minimum scores with null checks
      const scoreThreshold = 0.2;
      const performanceScore = result.lhr.categories.performance?.score ?? 0;
      const bestPracticesScore =
        result.lhr.categories['best-practices']?.score ?? 0;
      const seoScore = result.lhr.categories.seo?.score ?? 0;

      expect(performanceScore).toBeGreaterThanOrEqual(scoreThreshold);
      expect(bestPracticesScore).toBeGreaterThanOrEqual(scoreThreshold);
      expect(seoScore).toBeGreaterThanOrEqual(scoreThreshold);

      console.log(
        `Lighthouse report saved to for artifact upload: lighthouse-reports/${url.replace(/[^a-zA-Z0-9]/g, '_')}/report.html`,
      );
    });
  }
});
