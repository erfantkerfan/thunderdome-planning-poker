import { Page } from '@playwright/test';

export class AdminTeamsPage {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  async goto() {
    await this.page.goto(`/admin/teams`);
  }
}
