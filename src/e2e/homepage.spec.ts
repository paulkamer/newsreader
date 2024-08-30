import { test, expect } from '@playwright/test';


test.beforeEach(async ({ page }) => {
  await page.goto('/');
});

test('has title', async ({ page }) => {
  await expect(page).toHaveTitle(/Newsreader/);
});

test('has a heading', async ({ page }) => {
  await expect(page.getByRole('heading', { level: 1})).toHaveText('All feeds');
});

test('has a list with items', async ({ page }) => {
  const rows = await page.locator('ul li').count();
  expect(rows).toBeGreaterThan(0);
});
