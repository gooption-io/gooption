import { GooptionAppPage } from './app.po';

describe('gooption-app App', () => {
  let page: GooptionAppPage;

  beforeEach(() => {
    page = new GooptionAppPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
