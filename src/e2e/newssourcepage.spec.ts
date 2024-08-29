import { test, expect } from '@playwright/test';

test('test', async ({ page }) => {
    await page.goto('/');
    
    await test.step('Open first feed', async () => {
        const firstRow = page.locator('table tr').first();
        
        await firstRow.locator('a').click();
    });

    await test.step('has heading', async () => {
        await expect(page.getByRole('heading', { level: 1})).toContainText('Feed');
    });
});