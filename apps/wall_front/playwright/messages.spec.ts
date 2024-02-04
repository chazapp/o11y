import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('/');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/The Wall/);
});

test('sends a message and sees it in stream', async ({ page }) => {
  await page.goto('/');

  // Fill the MessageForm and send
  await page.getByLabel('Nick').fill("Playwright");
  await page.getByLabel('Message').fill("This is automated via Playwright");
  await page.getByRole('button', {name: 'Send'}).click();
  // Expect MessageForm has been emptied
  await expect(page.getByLabel('Nick')).toBeEmpty();
  await expect(page.getByLabel('Message')).toBeEmpty();
  
  // Expect MessageStream to be populated
  // Expects page to have a heading with the name of Installation.
  await expect(page.getByText("This is automated via Playwright")).toBeVisible();

});
