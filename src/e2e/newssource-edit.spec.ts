import { test, expect, Page } from '@playwright/test';

test.beforeEach(async ({ page }) => {
  await page.goto('/admin');
});

let orgTitle = '';

async function editNewsSource(page: Page, newTitle: string) {
  await test.step('Find first feed and go to edit page', async () => {
    const row = page.locator('ul li').first();
    const editLink = row.locator('a', { hasText: 'Edit' });
    
    const editLinkCount = await editLink.count();
    expect(editLinkCount).toBe(1);
  
    await editLink.click();
  });
  
  await test.step('Edit feed - update title and submit', async () => {
    await expect(page.getByRole('heading', { level: 1})).toContainText('Edit ');
  
    expect(page.locator('form')).toBeDefined();
  
    const titleInput = page.locator('input[name="title"]');

    orgTitle = await titleInput.inputValue();
    await titleInput.fill(newTitle);
  
    const submitButton = page.locator('form button[type="submit"]');
  
    const submitButtonCount = await submitButton.count();
    expect(submitButtonCount).toBe(1);
    
    await Promise.all([
      submitButton.click(),
      page.waitForURL('/admin')
    ]);
  });
  
  await test.step('Find renamed feed', async () => {
    const listItem = page.locator('ul li').filter({ has: page.locator('a', { hasText: newTitle }) });
  
    const listItemCount = await listItem.count();
    expect(listItemCount).toBe(1);
  });
}

test('Edit first news source', async ({ page }) => {
  const newTitle = `New Title ${Math.random()}`;

  await editNewsSource(page, newTitle);

  await editNewsSource(page, orgTitle); // Restore original title
});