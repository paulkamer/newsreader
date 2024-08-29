import { test, expect } from '@playwright/test';


test.beforeEach(async ({ page }) => {
  await page.goto('/');
});

test('has title', async ({ page }) => {
  await expect(page).toHaveTitle(/Newsreader/);
});

test('has a heading', async ({ page }) => {
  await expect(page.getByRole('heading')).toHaveText('Home');
});

test('has a table with rows', async ({ page }) => {
  const rows = await page.locator('table tr').count();
  expect(rows).toBeGreaterThan(0);
});
