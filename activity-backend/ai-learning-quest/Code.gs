const QUEST = Object.freeze({
  settingsSheet: "Quest Settings",
  responsesSheet: "Quest Responses",
  menuName: "Learning Quest",
  maxTextLength: 1000,
});

function onOpen() {
  SpreadsheetApp.getUi()
    .createMenu(QUEST.menuName)
    .addItem("Set up collector", "setupCollector")
    .addSeparator()
    .addItem("Start a new session…", "startNewSession")
    .addItem("Open intake", "openIntake")
    .addItem("Close intake", "closeIntake")
    .addItem("Show status", "showStatus")
    .addSeparator()
    .addItem("Export current session as Markdown", "exportCurrentSessionAsMarkdown")
    .addToUi();
}

function setupCollector() {
  const spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
  PropertiesService.getScriptProperties().setProperty("spreadsheetId", spreadsheet.getId());
  let settings = spreadsheet.getSheetByName(QUEST.settingsSheet);
  let responses = spreadsheet.getSheetByName(QUEST.responsesSheet);

  if (!settings) {
    settings = spreadsheet.insertSheet(QUEST.settingsSheet);
    settings.getRange(1, 1, 1, 2).setValues([["Setting", "Value"]]);
    settings.getRange(2, 1, 5, 2).setValues([
      ["intakeOpen", "FALSE"],
      ["sessionId", "not-started"],
      ["sessionTitle", "Engineering Mathematics I"],
      ["sessionCode", "CHANGE-ME"],
      ["sessionStartedAt", ""],
    ]);
    settings.setFrozenRows(1);
    settings.autoResizeColumns(1, 2);
  }

  if (!responses) {
    responses = spreadsheet.insertSheet(QUEST.responsesSheet);
    responses.getRange(1, 1, 1, 9).setValues([[
      "Submitted at",
      "Response ID",
      "Session ID",
      "Session title",
      "Mission",
      "Provisional claim",
      "Human challenge",
      "Final contribution",
      "Consent",
    ]]);
    responses.setFrozenRows(1);
    responses.autoResizeColumns(1, 9);
  }

  SpreadsheetApp.getUi().alert(
    "Collector ready",
    "Use Learning Quest → Start a new session before opening intake.",
    SpreadsheetApp.getUi().ButtonSet.OK,
  );
}

function startNewSession() {
  ensureCollector_();
  const ui = SpreadsheetApp.getUi();
  const titlePrompt = ui.prompt(
    "New Learning Quest session",
    "Enter a title students will see, for example: Engineering Mathematics I · Week 1",
    ui.ButtonSet.OK_CANCEL,
  );
  if (titlePrompt.getSelectedButton() !== ui.Button.OK) return;

  const codePrompt = ui.prompt(
    "Class code",
    "Enter a short class code. Students must enter it before submitting.",
    ui.ButtonSet.OK_CANCEL,
  );
  if (codePrompt.getSelectedButton() !== ui.Button.OK) return;

  const title = cleanText_(titlePrompt.getResponseText(), 120);
  const code = normalizeCode_(codePrompt.getResponseText());
  if (!title || code.length < 3) {
    ui.alert("Please provide a title and a class code containing at least 3 characters.");
    return;
  }

  const sessionId = Utilities.formatDate(new Date(), Session.getScriptTimeZone(), "yyyyMMdd-HHmmss");
  writeSettings_({
    intakeOpen: "FALSE",
    sessionId,
    sessionTitle: title,
    sessionCode: code,
    sessionStartedAt: new Date().toISOString(),
  });
  ui.alert(
    "Session prepared",
    `Title: ${title}\nClass code: ${code}\n\nIntake is still closed. Open it when students are ready.`,
    ui.ButtonSet.OK,
  );
}

function openIntake() {
  ensureCollector_();
  const settings = readSettings_();
  if (!settings.sessionId || settings.sessionId === "not-started") {
    SpreadsheetApp.getUi().alert("Start a new session before opening intake.");
    return;
  }
  writeSettings_({ intakeOpen: "TRUE" });
  SpreadsheetApp.getUi().alert(`Intake is open for “${settings.sessionTitle}”.`);
}

function closeIntake() {
  ensureCollector_();
  writeSettings_({ intakeOpen: "FALSE" });
  SpreadsheetApp.getUi().alert("Intake is closed. New submissions will be rejected.");
}

function showStatus() {
  ensureCollector_();
  const settings = readSettings_();
  const count = getCurrentSessionRows_(settings.sessionId).length;
  SpreadsheetApp.getUi().alert(
    "Learning Quest status",
    [
      `Intake: ${asBoolean_(settings.intakeOpen) ? "OPEN" : "CLOSED"}`,
      `Session: ${settings.sessionTitle}`,
      `Class code: ${settings.sessionCode}`,
      `Responses: ${count}`,
    ].join("\n"),
    SpreadsheetApp.getUi().ButtonSet.OK,
  );
}

function doGet(event) {
  ensureCollector_();
  const action = String((event && event.parameter && event.parameter.action) || "");
  if (action !== "status") {
    return HtmlService.createHtmlOutput("Learning Quest collector is online.");
  }

  const callback = safeCallback_((event.parameter && event.parameter.callback) || "callback");
  const settings = readSettings_();
  const payload = {
    open: asBoolean_(settings.intakeOpen),
    sessionTitle: settings.sessionTitle || "Engineering Mathematics I",
  };
  return ContentService
    .createTextOutput(`${callback}(${JSON.stringify(payload)});`)
    .setMimeType(ContentService.MimeType.JAVASCRIPT);
}

function doPost(event) {
  ensureCollector_();
  const parameters = (event && event.parameter) || {};
  const requestId = cleanText_(parameters.requestId, 100);
  const settings = readSettings_();

  if (!asBoolean_(settings.intakeOpen)) {
    return postMessageResult_(requestId, false, "This class intake is closed.");
  }
  if (normalizeCode_(parameters.sessionCode) !== normalizeCode_(settings.sessionCode)) {
    return postMessageResult_(requestId, false, "The class code is not correct.");
  }
  if (String(parameters.consent || "") !== "yes") {
    return postMessageResult_(requestId, false, "Please approve the text before submitting.");
  }

  const allowedMissions = [
    "Why attend class?",
    "What will I learn?",
    "How can AI help?",
    "Where can mathematics take me?",
  ];
  const mission = cleanText_(parameters.mission, 120);
  const aiClaim = cleanText_(parameters.aiClaim, 700);
  const humanChallenge = cleanText_(parameters.humanChallenge, 700);
  const revisedClaim = cleanText_(parameters.revisedClaim, 900);

  if (!allowedMissions.includes(mission) || !aiClaim || !humanChallenge || !revisedClaim) {
    return postMessageResult_(requestId, false, "One or more required responses are missing.");
  }

  const lock = LockService.getScriptLock();
  try {
    lock.waitLock(10000);
    const freshSettings = readSettings_();
    if (!asBoolean_(freshSettings.intakeOpen) || freshSettings.sessionId !== settings.sessionId) {
      return postMessageResult_(requestId, false, "The session closed before this response arrived.");
    }
    const sheet = getSpreadsheet_().getSheetByName(QUEST.responsesSheet);
    sheet.appendRow([
      new Date(),
      Utilities.getUuid(),
      safeCell_(settings.sessionId),
      safeCell_(settings.sessionTitle),
      safeCell_(mission),
      safeCell_(aiClaim),
      safeCell_(humanChallenge),
      safeCell_(revisedClaim),
      "yes",
    ]);
  } catch (error) {
    console.error(error);
    return postMessageResult_(requestId, false, "The collector is busy. Please try again.");
  } finally {
    if (lock.hasLock()) lock.releaseLock();
  }

  return postMessageResult_(requestId, true, "Response collected.");
}

function exportCurrentSessionAsMarkdown() {
  ensureCollector_();
  const settings = readSettings_();
  const rows = getCurrentSessionRows_(settings.sessionId);
  if (!rows.length) {
    SpreadsheetApp.getUi().alert("There are no responses in the current session.");
    return;
  }

  const closedAt = Utilities.formatDate(new Date(), Session.getScriptTimeZone(), "yyyy-MM-dd HH:mm z");
  const lines = [
    `# ${escapeMarkdown_(settings.sessionTitle)}`,
    "",
    `- Session ID: \`${escapeMarkdown_(settings.sessionId)}\``,
    `- Exported: ${closedAt}`,
    `- Responses: ${rows.length}`,
    "",
    "> Anonymous classroom responses. Analyse themes and tensions at group level; do not infer identities or grade individuals.",
    "",
  ];

  rows.forEach((row, index) => {
    lines.push(
      `## Response ${index + 1}`,
      "",
      `**Mission:** ${escapeMarkdown_(row[4])}`,
      "",
      "**Provisional claim after AI conversation**",
      "",
      escapeMarkdown_(row[5]),
      "",
      "**Challenge from a human partner**",
      "",
      escapeMarkdown_(row[6]),
      "",
      "**Final contribution**",
      "",
      escapeMarkdown_(row[7]),
      "",
      "---",
      "",
    );
  });

  const filename = `learning-quest-${settings.sessionId}.md`;
  const file = createFileBesideSpreadsheet_(filename, lines.join("\n"));
  const html = HtmlService.createHtmlOutput(
    `<p><strong>${rows.length} responses exported.</strong></p><p><a href="${file.getUrl()}" target="_blank">Open ${filename}</a></p>`,
  ).setWidth(420).setHeight(150);
  SpreadsheetApp.getUi().showModalDialog(html, "Markdown export ready");
}

function ensureCollector_() {
  const spreadsheet = getSpreadsheet_();
  if (!spreadsheet.getSheetByName(QUEST.settingsSheet) || !spreadsheet.getSheetByName(QUEST.responsesSheet)) {
    throw new Error("Collector is not set up. Run setupCollector from the script editor or Learning Quest menu.");
  }
}

function readSettings_() {
  const sheet = getSpreadsheet_().getSheetByName(QUEST.settingsSheet);
  const lastRow = sheet.getLastRow();
  const values = lastRow > 1 ? sheet.getRange(2, 1, lastRow - 1, 2).getDisplayValues() : [];
  return values.reduce((result, row) => {
    result[String(row[0])] = String(row[1]);
    return result;
  }, {});
}

function writeSettings_(updates) {
  const sheet = getSpreadsheet_().getSheetByName(QUEST.settingsSheet);
  const values = sheet.getRange(2, 1, sheet.getLastRow() - 1, 2).getValues();
  Object.keys(updates).forEach((key) => {
    const index = values.findIndex((row) => String(row[0]) === key);
    if (index >= 0) {
      sheet.getRange(index + 2, 2).setValue(updates[key]);
    } else {
      sheet.appendRow([key, updates[key]]);
    }
  });
}

function getCurrentSessionRows_(sessionId) {
  const sheet = getSpreadsheet_().getSheetByName(QUEST.responsesSheet);
  if (sheet.getLastRow() < 2) return [];
  return sheet
    .getRange(2, 1, sheet.getLastRow() - 1, 9)
    .getDisplayValues()
    .filter((row) => String(row[2]) === String(sessionId));
}

function postMessageResult_(requestId, ok, message) {
  const payload = JSON.stringify({
    type: "learning-quest-submit",
    requestId: cleanText_(requestId, 100),
    ok: Boolean(ok),
    message: cleanText_(message, 240),
  }).replace(/</g, "\\u003c");
  return HtmlService
    .createHtmlOutput(`<!doctype html><meta charset="utf-8"><script>window.parent.postMessage(${payload}, "*");<\/script>`)
    .setXFrameOptionsMode(HtmlService.XFrameOptionsMode.ALLOWALL);
}

function createFileBesideSpreadsheet_(name, content) {
  const spreadsheetFile = DriveApp.getFileById(getSpreadsheet_().getId());
  const parents = spreadsheetFile.getParents();
  const folder = parents.hasNext() ? parents.next() : DriveApp.getRootFolder();
  return folder.createFile(name, content, MimeType.PLAIN_TEXT);
}

function getSpreadsheet_() {
  const properties = PropertiesService.getScriptProperties();
  const savedId = properties.getProperty("spreadsheetId");
  if (savedId) return SpreadsheetApp.openById(savedId);

  const active = SpreadsheetApp.getActiveSpreadsheet();
  if (!active) throw new Error("Spreadsheet not found. Run setupCollector from the bound spreadsheet first.");
  properties.setProperty("spreadsheetId", active.getId());
  return active;
}

function safeCallback_(value) {
  const callback = String(value || "callback");
  return /^[A-Za-z_$][0-9A-Za-z_$\.]*$/.test(callback) ? callback : "callback";
}

function normalizeCode_(value) {
  return String(value || "").trim().toUpperCase().replace(/\s+/g, "-").slice(0, 24);
}

function cleanText_(value, limit) {
  return String(value || "")
    .replace(/\u0000/g, "")
    .replace(/\r\n?/g, "\n")
    .trim()
    .slice(0, Math.min(limit || QUEST.maxTextLength, QUEST.maxTextLength));
}

function safeCell_(value) {
  const text = cleanText_(value, QUEST.maxTextLength);
  return /^[=+\-@]/.test(text) ? `'${text}` : text;
}

function escapeMarkdown_(value) {
  return String(value || "").replace(/\\/g, "\\\\").replace(/([*_`\[\]])/g, "\\$1");
}

function asBoolean_(value) {
  return String(value || "").toUpperCase() === "TRUE";
}
