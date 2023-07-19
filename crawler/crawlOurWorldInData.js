const puppeteer = require("puppeteer");
const fs = require("fs");

(async () => {
  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();

  // Navigate to the page with the paginated table
  await page.goto("https://gitlab.pro.ai/");
  await page.type("#user_login", "PASSWORD");
  await page.type("#user_password", "PASSWORD");
  await page.click(".btn-confirm");

  await page.waitForSelector("ul.projects-list li.project-row a.project");

  // Loop through all pages of the table
  let hasNextPage = true;
  let totalResult = [];
  await fs.promises.unlink("./results.csv");

  while (hasNextPage) {
    // Get all the data in the current page of the table

    const tableItems = await page.$$eval("ul.projects-list li.project-row a.project", (lis) =>
      lis.map((li) => `https://gitlab.pro.ai${li.getAttribute("href")}.git`)
    );

    // Log the table items to the Node.js console
    totalResult.push(...tableItems)

    // Check if there is a next page button
    //await page.waitForSelector(".js-next-button a")
    const nextButton = await page.$('.js-next-button a');
    
    const IsActive = await page.$(".js-next-button.disabled") ? false : true;

    if (IsActive) {
      // Click the next page button
      await nextButton.click();
      // Wait for the next page to load
      await page.waitForSelector("ul.projects-list li.project-row a.project");
      await page.waitForSelector(".js-next-button a")
    } else {
      // Stop looping if there is no next page button
      hasNextPage = false;
    }
  }
  fs.appendFile('./results.csv', totalResult.join("\n") , (err) => {
    if(err) 
      throw err;
    console.log('The new_content was appended successfully');
  });
  console.log(totalResult.length);
  await browser.close();
})();
