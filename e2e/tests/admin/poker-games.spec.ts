import { expect, test } from '../../fixtures/user-sessions';
import { AdminGamesPage } from '../../fixtures/admin/battles-page';

test.describe('The Admin Poker Games Page', () => {
  test.describe('Unauthenticated user', () => {
    test('redirects to login', async ({ page }) => {
      const adminPage = new AdminGamesPage(page);

      await adminPage.goto();

      const title = adminPage.page.locator('[data-formtitle="login"]');
      await expect(title).toHaveText('Login');
    });
  });

  test.describe('Guest user', () => {
    test('redirects to landing', async ({ guestPage }) => {
      const adminPage = new AdminGamesPage(guestPage.page);

      await adminPage.goto();

      const title = adminPage.page.locator('h1');
      await expect(title).toHaveText(
        'Thunderdome is an Agile Planning Poker app with a fun theme',
      );
    });
  });

  test.describe('Non Admin Registered User', () => {
    test('redirects to landing', async ({ registeredPage }) => {
      const adminPage = new AdminGamesPage(registeredPage.page);

      await adminPage.goto();

      const title = adminPage.page.locator('h1');
      await expect(title).toHaveText(
        'Thunderdome is an Agile Planning Poker app with a fun theme',
      );
    });
  });

  test.describe('Admin User', () => {
    test('loads Games page', async ({ adminPage }) => {
      const ap = new AdminGamesPage(adminPage.page);

      await ap.goto();

      const title = ap.page.locator('h1');
      await expect(title).toHaveText('Battles');
    });
  });
});
