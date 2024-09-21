import { test, expect, Page } from '@playwright/test';

test.beforeEach(async ({ page }) => {
  await page.goto('/admin');
});

async function addNewsSource(page: Page, title: string) {
  await test.step('Go to add page', async () => {
    await page.locator('a', { hasText: 'Add' }).click();
    await expect(page.getByRole('heading', { level: 1})).toContainText('Add new news source');
  });

  await test.step('Fill in form and submit', async () => {
    await page.locator('input[name="title"]').fill(title);
    await page.locator('input[name="url"]').fill('https://www.nasa.gov/feed/');
    await page.locator('select[name="feedtype"]').selectOption({ label: 'RSS' });

    const submitButton = page.getByRole('button', { name: /submit/i })

    await Promise.all([
      submitButton.click(),
      page.waitForURL('/admin')
    ]);
  });

  await test.step('Find new feed', async () => {
    const listItem = page.locator('ul li').filter({ has: page.locator('a', { hasText: title }) });

    const listItemCount = await listItem.count();
    expect(listItemCount).toBe(1);
});
}

test('Add new news source and delete via Admin page', async ({ page }) => {
  const title = `Test new news Source ${new Date().toISOString()}`;

  await addNewsSource(page, title);

  await test.step('Delete new feed on Admin page', async () => {
    const listItem = page.locator('ul li').filter({ has: page.locator('a', { hasText: title }) });
    const deleteBtn = listItem.locator('button', { hasText: 'Delete' });

    const deleteBtnCount = await deleteBtn.count();
    expect(deleteBtnCount).toBe(1);

    await deleteBtn.click();
  });

  await test.step('Check if feed is deleted', async () => {
    await page.reload();

    const listItem = page.locator('ul li').filter({ has: page.locator('a', { hasText: title }) });

    const listItemCount = await listItem.count();
    expect(listItemCount).toBe(0);
  });
});

test('Add new news source and delete via Edit page', async ({ page }) => {
  const title = `Test new news Source ${new Date().toISOString()}`;

  addNewsSource(page, title);

  await test.step('Navigate to its Edit page', async () => {
    const listItem = page.locator('ul li').filter({ has: page.locator('a', { hasText: title }) });
    const editLink = listItem.locator('a', { hasText: 'Edit' });

    await editLink.click();
    await page.waitForLoadState('domcontentloaded');
  });

  await test.step('Delete feed on its Edit page', async () => {
    const deleteBtn = page.locator('button', { hasText: 'Delete' });

    const deleteBtnCount = await deleteBtn.count();
    expect(deleteBtnCount).toBe(1);

    await Promise.all([
      deleteBtn.click(),
      page.waitForURL('/admin')
    ]);
  });

  await test.step('Check if feed is deleted', async () => {
    const listItem = page.locator('ul li').filter({ has: page.locator('a', { hasText: title }) });

    const listItemCount = await listItem.count();
    expect(listItemCount).toBe(0);
  });
});
