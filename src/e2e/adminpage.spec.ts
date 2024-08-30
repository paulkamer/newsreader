import { test, expect } from '@playwright/test';

test.beforeEach(async ({ page }) => {
  await page.goto('/admin');
});

test('has title', async ({ page }) => {
  await expect(page).toHaveTitle(/Newsreader/);
});

test('has a heading', async ({ page }) => {
  await expect(page.getByRole('heading', { level: 1})).toHaveText('Admin');
});

test('has a list with rows', async ({ page }) => {
  const rows = await page.locator('ul li').count();
  expect(rows).toBeGreaterThan(0);
});


test('has an Add link', async ({ page }) => {
  expect(page.getByRole('link', { name: 'Add'})).toBeDefined();
});