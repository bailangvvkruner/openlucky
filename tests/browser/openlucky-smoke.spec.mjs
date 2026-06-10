import { test, expect } from '@playwright/test';

const baseURL = process.env.OPENLUCKY_E2E_URL || 'http://127.0.0.1:16601/lucky/';
const username = process.env.OPENLUCKY_E2E_USER || 'openlucky';
const password = process.env.OPENLUCKY_E2E_PASSWORD || 'openlucky-dev';

test('login and render core routes', async ({ page }) => {
  await page.goto(baseURL);
  await expect(page.getByRole('heading', { name: /openlucky/i })).toBeVisible();

  await page.getByLabel(/username/i).fill(username);
  await page.getByLabel(/password/i).fill(password);
  await page.getByRole('button', { name: /sign in/i }).click();

  await expect(page.getByText(/system overview/i)).toBeVisible();
  await page.getByRole('link', { name: /logs/i }).click();
  await expect(page.getByText(/recent events/i)).toBeVisible();
  await page.getByRole('link', { name: /modules/i }).click();
  await expect(page.getByText(/module registry/i)).toBeVisible();
});
