import { test, expect } from '@playwright/test';

test('opens build info and see populated table', async ({ page }) => {
  await page.goto('/');

  // Expect a title "to contain" a substring.
  await page.getByLabel('build-info').click()
  await expect(page.getByLabel('status-data-table')).toBeVisible()
});
