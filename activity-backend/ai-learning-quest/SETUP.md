# AI Learning Quest setup

The public page is `/activity/ai-learning-quest/`. It is deliberately absent
from the site navigation and search-engine sitemap. It uses a Google Sheet plus
Apps Script to store anonymous responses and control whether intake is open.

## 1. Test the page safely

The committed `config.js` starts in demo mode. Demo submissions stay only in
the current browser's local storage and never leave the device. This makes it
possible to test the full four-level flow before connecting a real collector.

## 2. Create the private collector

1. Create a blank Google Sheet in the Google account that will own the data.
2. In the sheet, choose **Extensions → Apps Script**.
3. Replace the editor contents with `Code.gs` from this folder.
4. In **Project Settings**, enable display of the `appsscript.json` manifest,
   then replace it with the manifest in this folder.
5. Save the project.
6. In the function selector, choose `setupCollector`, select **Run**, and grant
   the requested permissions.
7. Reload the Google Sheet. A **Learning Quest** menu will appear.

## 3. Publish the collector

1. In Apps Script, choose **Deploy → New deployment**.
2. Select **Web app**.
3. Set **Execute as** to yourself.
4. Set access to **Anyone**. If a university administrator blocks this option,
   deploy from an approved account or use an institution-approved backend.
5. Deploy and copy the `/exec` web app URL.
6. In `activity/ai-learning-quest/config.js`, paste that URL into `endpoint`
   and change `demoMode` to `false`.
7. Commit and push the site changes. GitHub Pages will publish the activity.

The endpoint is visible to browsers by design. The class code and open/closed
state prevent ordinary unwanted submissions. Do not place secrets in
`config.js`.

## 4. Run a class session

There is no server to start or stop on your computer.

1. Open the private Google Sheet.
2. Choose **Learning Quest → Start a new session…**.
3. Enter the session title and a class code.
4. Choose **Learning Quest → Open intake**.
5. Show students the hidden page URL and class code.
6. At the end, choose **Learning Quest → Close intake**.
7. Choose **Learning Quest → Export current session as Markdown**.
8. Open the generated `.md` file in Google Drive and download or upload it to
   ChatGPT for group-level analysis.

Starting another session does not delete earlier responses. Each export uses
only the current session ID.

## Privacy notes

- The activity does not request names, email addresses, student IDs, or the
  full ChatGPT transcript.
- Responses are stored in the private Google Sheet owned by the lecturer.
- The page asks students to review and consent to the exact text submitted.
- A hidden URL is not authentication. Keep the class code short-lived and close
  intake promptly.
- Follow university requirements for notice, retention, and use of classroom
  response data.
