import { test as setup } from '@playwright/test';
import path from 'path';

const authFile = path.join(__dirname, './.auth/user.json');

setup('authenticate', async ({ page }) => {
    await page.goto('/login');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password');
    await page.click('button[type="submit"]');

    // Wait for the final URL to ensure that the cookies are actually set.
    await page.waitForURL('/');

    await page.context().storageState({ path: authFile });
})
