import { browser } from 'k6/experimental/browser';
import { check } from 'k6';

export const options = {
  scenarios: {
    ui: {
      executor: 'shared-iterations',
      options: {
        browser: {
          type: 'chromium',
        },
      },
    },
  },
  thresholds: {
    checks: ["rate==1.0"]
  }
}

export default async function () {
  const page = browser.newPage({
    ignoreHTTPSErrors: true,
  });

  try {
    await page.goto('https://wall.local', {
        
    });

    page.locator('input[id=":r1:"]').type('admin');
    page.locator('textarea[id=":r2:"').type('My Hello World !');

    const submitButton = page.locator('button[type="submit"]');

    await Promise.all([submitButton.click()]);

    // check(page, {
    //   'header': p => p.locator('h2').textContent() == 'Welcome, admin!',
    // });
  } finally {
    page.close();
  }
}
