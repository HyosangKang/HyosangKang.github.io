(() => {
  "use strict";

  const DAYS = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"];
  const SHORT_DAYS = ["Mon", "Tue", "Wed", "Thu", "Fri"];
  const SLOTS_PER_DAY = 52;
  const ROW_HEIGHT = 26;
  const BLOCK_COLORS = [
    "#dceefa",
    "#dff1ea",
    "#eee8fb",
    "#faebdb",
    "#dceee9",
    "#f5e6ed",
    "#e7edf8",
    "#e9f0dc",
  ];

  const state = {
    engineReady: false,
    activeTab: "courses",
    sourceName: "Sample data.xlsx",
    courses: [],
    instructors: [],
    students: [],
    classrooms: [],
    summary: {},
    included: new Set(),
    assignments: [],
    assignmentByCourse: new Map(),
    selected: {
      courses: "",
      instructors: "",
      students: "",
      classrooms: "",
    },
    query: "",
    courseFilter: "all",
    timetableView: "all",
    timetableQuery: "",
    running: false,
  };

  const app = document.querySelector("#app-content");
  const fileInput = document.querySelector("#workbook-file");
  const engineBanner = document.querySelector("#engine-banner");
  const engineStatus = document.querySelector("#engine-status");

  document.querySelector(".main-nav").addEventListener("click", (event) => {
    const button = event.target.closest("[data-tab]");
    if (!button || state.running) return;
    state.activeTab = button.dataset.tab;
    state.query = "";
    render();
  });

  fileInput.addEventListener("change", async () => {
    const [file] = fileInput.files;
    if (!file) return;
    await importWorkbook(new Uint8Array(await file.arrayBuffer()), file.name);
    fileInput.value = "";
  });

  app.addEventListener("click", (event) => {
    const choose = event.target.closest(".choose-workbook");
    if (choose && state.engineReady) {
      fileInput.click();
      return;
    }

    const tab = event.target.closest("[data-switch-tab]");
    if (tab) {
      state.activeTab = tab.dataset.switchTab;
      render();
      return;
    }

    const filter = event.target.closest("[data-course-filter]");
    if (filter) {
      state.courseFilter = filter.dataset.courseFilter;
      renderCourses();
      return;
    }

    const includeAction = event.target.closest("[data-include-action]");
    if (includeAction) {
      if (includeAction.dataset.includeAction === "all") {
        state.courses.forEach((course) => state.included.add(course.id));
      } else {
        state.included.clear();
      }
      renderCourses();
      return;
    }

    const toggle = event.target.closest("[data-toggle-course]");
    if (toggle) {
      event.stopPropagation();
      const id = toggle.dataset.toggleCourse;
      if (state.included.has(id)) state.included.delete(id);
      else state.included.add(id);
      renderCourses();
      return;
    }

    const entity = event.target.closest("[data-entity-key]");
    if (entity) {
      const type = entity.dataset.entityType;
      state.selected[type] = entity.dataset.entityKey;
      renderEntityDetail(type);
      refreshSelectedRows(type);
      return;
    }

    const view = event.target.closest("[data-timetable-view]");
    if (view) {
      state.timetableView = view.dataset.timetableView;
      renderTimetableGrid();
      refreshActiveSegments();
      return;
    }

    const generate = event.target.closest("#generate-timetable");
    if (generate) {
      if (state.running) window.DataWeb.cancel();
      else generateTimetable();
      return;
    }

    const cancel = event.target.closest("#cancel-generation");
    if (cancel) {
      window.DataWeb.cancel();
      return;
    }

    const save = event.target.closest("#download-timetable");
    if (save) {
      downloadTimetable();
      return;
    }

    const block = event.target.closest("[data-course-block]");
    if (block) {
      state.selected.courses = block.dataset.courseBlock;
      state.activeTab = "courses";
      render();
    }
  });

  app.addEventListener("input", (event) => {
    if (event.target.matches("#entity-search")) {
      state.query = event.target.value;
      if (state.activeTab === "courses") renderCourses();
      else renderDirectory(state.activeTab);
    }
    if (event.target.matches("#timetable-search")) {
      state.timetableQuery = event.target.value;
      filterTimetableBlocks();
    }
  });

  async function initialize() {
    try {
      if (typeof Go !== "function") {
        throw new Error("WebAssembly support could not be loaded.");
      }
      const go = new Go();
      const response = await fetch("dataweb.wasm");
      if (!response.ok) {
        throw new Error(`The timetabling engine could not be downloaded (${response.status}).`);
      }

      let result;
      try {
        result = await WebAssembly.instantiateStreaming(response.clone(), go.importObject);
      } catch {
        result = await WebAssembly.instantiate(await response.arrayBuffer(), go.importObject);
      }
      go.run(result.instance);
      await waitForApi();
      await window.DataWeb.ready;
      state.engineReady = true;
      await loadSample();
    } catch (error) {
      setEngineStatus(readError(error), "error");
      renderError(readError(error));
    }
  }

  function waitForApi() {
    return new Promise((resolve, reject) => {
      let attempts = 0;
      const check = () => {
        if (window.DataWeb) return resolve();
        if (attempts++ > 240) {
          return reject(new Error("The timetabling engine did not start."));
        }
        window.setTimeout(check, 25);
      };
      check();
    });
  }

  async function loadSample() {
    setEngineStatus("Loading the protected sample workbook…");
    const response = await fetch("sample-data.xlsx");
    if (!response.ok) {
      throw new Error(`The sample workbook could not be loaded (${response.status}).`);
    }
    await importWorkbook(new Uint8Array(await response.arrayBuffer()), "Sample data.xlsx");
  }

  async function importWorkbook(bytes, sourceName) {
    if (!state.engineReady) return;
    setEngineStatus(`Reading ${sourceName}…`);
    await new Promise((resolve) => window.setTimeout(resolve, 20));
    const result = window.DataWeb.importWorkbook(bytes);
    if (!result.ok) {
      const message = (result.errors || ["The workbook could not be imported."]).join(" · ");
      setEngineStatus(message, "error");
      return;
    }

    state.sourceName = sourceName;
    state.courses = Array.from(result.courses || []).sort(courseSort);
    state.instructors = Array.from(result.instructors || []).sort((a, b) =>
      text(a.name).localeCompare(text(b.name)),
    );
    state.students = Array.from(result.students || []).sort((a, b) =>
      text(a.id).localeCompare(text(b.id), undefined, { numeric: true }),
    );
    state.classrooms = Array.from(result.classrooms || []).sort((a, b) =>
      text(a.name).localeCompare(text(b.name), undefined, { numeric: true }),
    );
    state.summary = result.summary || {};
    state.included = new Set(state.courses.map((course) => course.id));
    state.assignments = [];
    state.assignmentByCourse = new Map();
    state.selected.courses = state.courses[0]?.id || "";
    state.selected.instructors = state.instructors[0]?.name || "";
    state.selected.students = state.students[0]?.id || "";
    state.selected.classrooms = state.classrooms[0]?.name || "";
    state.query = "";
    state.courseFilter = "all";
    state.timetableQuery = "";
    setEngineStatus(
      `${sourceName} loaded — ${formatNumber(state.summary.courses)} courses, ${formatNumber(state.summary.instructors)} instructors, ${formatNumber(state.summary.students)} students, ${formatNumber(state.summary.classrooms)} classrooms.`,
      "ready",
    );
    render();
  }

  function setEngineStatus(message, type = "loading") {
    engineStatus.textContent = message;
    engineBanner.classList.toggle("is-ready", type === "ready");
    engineBanner.classList.toggle("is-error", type === "error");
  }

  function render() {
    document.querySelectorAll(".nav-item").forEach((button) => {
      button.classList.toggle("is-active", button.dataset.tab === state.activeTab);
    });
    if (!state.courses.length) {
      renderLoading();
      return;
    }
    if (state.activeTab === "courses") renderCoursesPage();
    else if (["instructors", "students", "classrooms"].includes(state.activeTab)) {
      renderDirectoryPage(state.activeTab);
    } else if (state.activeTab === "timetable") renderTimetablePage();
    else if (state.activeTab === "settings") renderSettings();
    else renderHelp();
  }

  function renderLoading() {
    app.innerHTML = `
      <div class="panel empty">
        <div>
          <p class="eyebrow">DATA browser simulation</p>
          <h1>Preparing the application</h1>
          <p>The sample workbook is being parsed locally by the WebAssembly engine.</p>
        </div>
      </div>`;
  }

  function renderError(message) {
    app.innerHTML = `
      <div class="panel empty">
        <div>
          <p class="eyebrow">Unable to start</p>
          <h1>The demo could not be loaded</h1>
          <p>${escapeHtml(message)}</p>
        </div>
      </div>`;
  }

  function pageHeading(title, copy, eyebrow = "2026 spring semester") {
    return `
      <div class="page-heading">
        <div>
          <p class="eyebrow">${escapeHtml(eyebrow)}</p>
          <h1>${escapeHtml(title)}</h1>
          <p class="heading-copy">${escapeHtml(copy)}</p>
        </div>
        <div class="summary-line">${summaryText()}</div>
      </div>`;
  }

  function fileToolbar() {
    return `
      <div class="file-toolbar">
        <button class="button button-primary choose-workbook" type="button">Open data file</button>
        <span class="loaded-file">${escapeHtml(state.sourceName)}</span>
        <span class="privacy-copy">Processed locally in this browser</span>
      </div>`;
  }

  function summaryText() {
    return `${formatNumber(state.summary.courses)} courses · ${formatNumber(state.summary.instructors)} instructors · ${formatNumber(state.summary.students)} students · ${formatNumber(state.summary.classrooms)} classrooms`;
  }

  function renderCoursesPage() {
    app.innerHTML = `
      ${pageHeading(
        "Courses",
        "Review imported courses, decide which classes to include, and inspect their scheduling constraints.",
      )}
      ${fileToolbar()}
      <section class="panel entity-layout">
        <div class="entity-browser">
          <div class="browser-toolbar">
            <label class="search">
              <span aria-hidden="true">⌕</span>
              <input id="entity-search" type="search" value="${escapeAttr(state.query)}" placeholder="Search course code, name, or instructor" aria-label="Search courses">
            </label>
            <div class="filter-row" aria-label="Course filters">
              ${courseFilterButton("all", "All")}
              ${courseFilterButton("include", "Include")}
              ${courseFilterButton("exclude", "Exclude")}
              ${courseFilterButton("allocated", "Allocation")}
              ${courseFilterButton("unallocated", "Unallocated")}
            </div>
            <div class="browser-actions">
              <small><span id="visible-count">0</span> visible · <span id="included-count">${formatNumber(state.included.size)}</span> included</small>
              <span>
                <button class="text-button" type="button" data-include-action="all">Include all</button>
                <button class="text-button" type="button" data-include-action="none">Clear all</button>
              </span>
            </div>
          </div>
          <div id="entity-list" class="entity-list"></div>
        </div>
        <aside id="entity-detail" class="detail-pane"></aside>
      </section>`;
    renderCourses();
  }

  function courseFilterButton(value, label) {
    return `<button class="filter-chip${state.courseFilter === value ? " is-active" : ""}" type="button" data-course-filter="${value}">${label}</button>`;
  }

  function renderCourses() {
    const list = document.querySelector("#entity-list");
    if (!list) return;
    const query = state.query.trim().toLowerCase();
    const visible = state.courses.filter((course) => {
      const included = state.included.has(course.id);
      const allocated = state.assignmentByCourse.get(course.id)?.assigned;
      const filterMatches =
        state.courseFilter === "all" ||
        (state.courseFilter === "include" && included) ||
        (state.courseFilter === "exclude" && !included) ||
        (state.courseFilter === "allocated" && allocated) ||
        (state.courseFilter === "unallocated" && !allocated);
      const haystack = [
        course.code,
        course.section,
        course.name,
        course.nameKorean,
        ...(course.instructors || []),
      ]
        .join(" ")
        .toLowerCase();
      return filterMatches && haystack.includes(query);
    });

    document.querySelector("#visible-count").textContent = formatNumber(visible.length);
    document.querySelector("#included-count").textContent = formatNumber(state.included.size);
    const groups = new Map();
    visible.forEach((course) => {
      const year = Number(course.year) || 0;
      const label = year > 0 ? `Year ${year}` : "Other courses";
      if (!groups.has(label)) groups.set(label, []);
      groups.get(label).push(course);
    });

    list.innerHTML = groups.size
      ? Array.from(groups, ([label, courses]) => `
          <section>
            <h2 class="year-heading">${escapeHtml(label)} · ${courses.length}</h2>
            ${courses.map(courseRow).join("")}
          </section>`).join("")
      : `<div class="empty"><p>No courses match this view.</p></div>`;

    if (!visible.some((course) => course.id === state.selected.courses)) {
      state.selected.courses = visible[0]?.id || state.courses[0]?.id || "";
    }
    renderEntityDetail("courses");
  }

  function courseRow(course) {
    const included = state.included.has(course.id);
    const assignment = state.assignmentByCourse.get(course.id);
    const allocated = assignment?.assigned;
    return `
      <button class="entity-row${state.selected.courses === course.id ? " is-selected" : ""}" type="button" data-entity-type="courses" data-entity-key="${escapeAttr(course.id)}">
        <span>
          <strong>${escapeHtml(courseLabel(course))} · ${escapeHtml(course.name || course.nameKorean || "Untitled course")}</strong>
          <small>${escapeHtml(joinOr(course.instructors, "Instructor not specified"))} · ${formatNumber(course.students || course.capacity || 0)} students</small>
        </span>
        <span class="row-tags">
          <span class="tag ${included ? "include" : "exclude"}" data-toggle-course="${escapeAttr(course.id)}">${included ? "Include" : "Exclude"}</span>
          <span class="tag ${allocated ? "assigned" : "unallocated"}">${allocated ? "Allocated" : "Unallocated"}</span>
          ${course.required ? `<span class="tag required">Required</span>` : ""}
        </span>
      </button>`;
  }

  function renderDirectoryPage(type) {
    const titles = {
      instructors: ["Instructors", "Review teaching assignments and instructor time preferences."],
      students: ["Students", "Review imported student plans and their requested courses."],
      classrooms: ["Classrooms", "Review room capacity, special facilities, and current allocations."],
    };
    app.innerHTML = `
      ${pageHeading(titles[type][0], titles[type][1])}
      ${fileToolbar()}
      <section class="panel entity-layout">
        <div class="entity-browser">
          <div class="browser-toolbar">
            <label class="search">
              <span aria-hidden="true">⌕</span>
              <input id="entity-search" type="search" value="${escapeAttr(state.query)}" placeholder="Search ${escapeAttr(type)}" aria-label="Search ${escapeAttr(type)}">
            </label>
            <div class="browser-actions">
              <small><span id="visible-count">0</span> records</small>
            </div>
          </div>
          <div id="entity-list" class="entity-list"></div>
        </div>
        <aside id="entity-detail" class="detail-pane"></aside>
      </section>`;
    renderDirectory(type);
  }

  function renderDirectory(type) {
    const list = document.querySelector("#entity-list");
    if (!list) return;
    const query = state.query.trim().toLowerCase();
    const records = state[type].filter((record) =>
      Object.values(record)
        .flat()
        .join(" ")
        .toLowerCase()
        .includes(query),
    );
    document.querySelector("#visible-count").textContent = formatNumber(records.length);
    list.innerHTML = records.length
      ? records.map((record) => directoryRow(type, record)).join("")
      : `<div class="empty"><p>No records match this search.</p></div>`;
    const keys = records.map((record) => entityKey(type, record));
    if (!keys.includes(state.selected[type])) state.selected[type] = keys[0] || "";
    renderEntityDetail(type);
  }

  function directoryRow(type, record) {
    const key = entityKey(type, record);
    let title;
    let subtitle;
    if (type === "instructors") {
      title = record.name;
      subtitle = `${record.courseIDs?.length || 0} courses · ${record.email || "No email"}`;
    } else if (type === "students") {
      title = `Student ${record.id}`;
      subtitle = `${record.courseIDs?.length || 0} requested courses`;
    } else {
      title = record.name;
      const assigned = state.assignments.filter((item) => item.classroom === record.name).length;
      subtitle = `Capacity ${formatNumber(record.capacity || 0)} · ${assigned} allocations`;
    }
    return `
      <button class="entity-row${state.selected[type] === key ? " is-selected" : ""}" type="button" data-entity-type="${type}" data-entity-key="${escapeAttr(key)}">
        <span><strong>${escapeHtml(title)}</strong><small>${escapeHtml(subtitle)}</small></span>
        <span class="row-tags"><span class="tag ${type === "classrooms" && record.special ? "required" : "include"}">${type === "classrooms" && record.special ? "Special" : "Imported"}</span></span>
      </button>`;
  }

  function refreshSelectedRows(type) {
    document.querySelectorAll(`[data-entity-type="${type}"]`).forEach((row) => {
      row.classList.toggle("is-selected", row.dataset.entityKey === state.selected[type]);
    });
  }

  function renderEntityDetail(type) {
    const pane = document.querySelector("#entity-detail");
    if (!pane) return;
    if (type === "courses") {
      const course = state.courses.find((item) => item.id === state.selected.courses);
      pane.innerHTML = course ? courseDetail(course) : emptyDetail("course");
      return;
    }
    const record = state[type].find((item) => entityKey(type, item) === state.selected[type]);
    pane.innerHTML = record ? directoryDetail(type, record) : emptyDetail(type);
  }

  function courseDetail(course) {
    const assignment = state.assignmentByCourse.get(course.id);
    const status = assignment?.assigned ? "Allocated" : "Not allocated";
    return `
      <p class="detail-overline">${escapeHtml(course.category?.join(" · ") || "Course record")}</p>
      <h2>${escapeHtml(courseLabel(course))}</h2>
      <p class="detail-subtitle">${escapeHtml(course.name || course.nameKorean || "Untitled course")}</p>
      <div class="detail-tabs"><span class="is-active">Timetable</span><span>Basic information</span></div>
      <div class="detail-grid">
        ${dataCard("Status", status)}
        ${dataCard("Assigned classroom", assignment?.classroom || course.classroom || "Not assigned")}
        ${dataCard("Assigned time", assignment?.classtime || course.classtime || "Not assigned")}
        ${dataCard("Credit / hours", `${course.credit ?? "—"} credits · ${course.lectureHours ?? 0} lecture · ${course.labHours ?? 0} lab`)}
        ${dataCard("Instructors", joinOr(course.instructors, "Not specified"), true)}
        ${dataCard("Capacity / student plans", `${formatNumber(course.capacity || 0)} capacity · ${formatNumber(course.students || 0)} plans`)}
        ${dataCard("Required course", course.required ? "Yes" : "No")}
        ${dataCard("Possible classrooms", joinOr(course.possibleClassrooms, "Any compatible room"), true)}
        ${dataCard("Preferred times", course.preferredTimes || "No preference", true)}
        ${dataCard("Avoided times", course.avoidedTimes || "No restriction", true)}
        ${dataCard("Possible times", course.possibleTimes || "Weekday operating hours", true)}
      </div>`;
  }

  function directoryDetail(type, record) {
    if (type === "instructors") {
      return `
        <p class="detail-overline">Instructor record</p>
        <h2>${escapeHtml(record.name)}</h2>
        <p class="detail-subtitle">${escapeHtml(record.email || "No email address")}</p>
        <div class="detail-tabs"><span class="is-active">Teaching constraints</span><span>Basic information</span></div>
        <div class="detail-grid">
          ${dataCard("Employee ID", record.employeeID || "Not specified")}
          ${dataCard("Teaching days", record.teachingDays || "No limit")}
          ${dataCard("Courses", courseNames(record.courseIDs), true)}
          ${dataCard("Preferred times", record.preferredTimes || "No preference", true)}
          ${dataCard("Avoided times", record.avoidedTimes || "No restriction", true)}
        </div>`;
    }
    if (type === "students") {
      return `
        <p class="detail-overline">Student plan</p>
        <h2>${escapeHtml(record.id)}</h2>
        <p class="detail-subtitle">${record.courseIDs?.length || 0} requested courses</p>
        <div class="detail-tabs"><span class="is-active">Course plan</span></div>
        <div class="detail-grid">
          ${dataCard("Student ID", record.id)}
          ${dataCard("Requested courses", courseNames(record.courseIDs), true)}
        </div>`;
    }
    const allocations = state.assignments.filter((item) => item.classroom === record.name);
    return `
      <p class="detail-overline">Classroom record</p>
      <h2>${escapeHtml(record.name)}</h2>
      <p class="detail-subtitle">${record.special ? "Special-purpose classroom" : "General teaching classroom"}</p>
      <div class="detail-tabs"><span class="is-active">Allocation</span><span>Basic information</span></div>
      <div class="detail-grid">
        ${dataCard("Capacity", formatNumber(record.capacity || 0))}
        ${dataCard("Special facility", record.special ? "Yes" : "No")}
        ${dataCard("Allocated courses", allocations.length ? allocations.map((item) => item.code).join(", ") : "No current allocations", true)}
      </div>`;
  }

  function dataCard(label, value, wide = false) {
    return `
      <div class="data-card${wide ? " is-wide" : ""}">
        <span class="data-label">${escapeHtml(label)}</span>
        <span class="data-value">${escapeHtml(text(value))}</span>
      </div>`;
  }

  function emptyDetail(noun) {
    return `<div class="empty"><p>Select a ${escapeHtml(noun)} to inspect its imported constraints.</p></div>`;
  }

  function renderTimetablePage() {
    app.innerHTML = `
      ${pageHeading(
        "Timetable",
        "Generate and inspect a weekly timetable using the same search engine as the desktop web application.",
      )}
      ${fileToolbar()}
      <section class="panel timetable-panel">
        <div class="timetable-toolbar">
          <button id="generate-timetable" class="button ${state.running ? "button-danger" : "button-primary"}" type="button" ${state.included.size ? "" : "disabled"}>
            ${state.running ? "Stop generating" : state.assignments.length ? "Generate again" : "Auto-generate"}
          </button>
          <button id="download-timetable" class="button button-secondary" type="button" ${state.assignments.length ? "" : "disabled"}>Download CSV</button>
          <div class="segmented" aria-label="Timetable view">
            ${DAYS.map((day, index) => segmentButton(String(index), SHORT_DAYS[index])).join("")}
            ${segmentButton("all", "All · compact")}
          </div>
          <label class="search">
            <span aria-hidden="true">⌕</span>
            <input id="timetable-search" type="search" value="${escapeAttr(state.timetableQuery)}" placeholder="Search course or room" aria-label="Search timetable">
          </label>
        </div>
        <div id="progress-wrap" class="progress-wrap" ${state.running ? "" : "hidden"}>
          <div class="progress-details">
            <div class="progress-heading">
              <strong id="progress-percent">0%</strong>
              <span id="progress-count">0 of ${formatNumber(state.included.size)} placed</span>
            </div>
            <div
              id="progress-track"
              class="progress-track"
              role="progressbar"
              aria-label="Timetable generation progress"
              aria-valuemin="0"
              aria-valuemax="100"
              aria-valuenow="0"
            >
              <span id="progress-bar" class="progress-bar" style="width:0%"></span>
            </div>
            <span id="progress-copy" class="progress-copy">Preparing course constraints…</span>
          </div>
          <button id="cancel-generation" class="text-button" type="button">Cancel</button>
        </div>
        <div id="grid-scroll" class="grid-scroll"></div>
      </section>`;
    renderTimetableGrid();
  }

  function segmentButton(value, label) {
    return `<button class="segment${state.timetableView === value ? " is-active" : ""}" type="button" data-timetable-view="${value}">${label}</button>`;
  }

  function refreshActiveSegments() {
    document.querySelectorAll("[data-timetable-view]").forEach((button) => {
      button.classList.toggle("is-active", button.dataset.timetableView === state.timetableView);
    });
  }

  function renderTimetableGrid() {
    const host = document.querySelector("#grid-scroll");
    if (!host) return;
    const dayIndexes =
      state.timetableView === "all" ? [0, 1, 2, 3, 4] : [Number(state.timetableView)];
    const columns = dayIndexes.length;
    host.innerHTML = `
      <div class="schedule" style="--day-count:${columns}">
        <div class="schedule-head">
          <div>Time</div>
          ${dayIndexes.map((day) => `<div>${DAYS[day]}</div>`).join("")}
        </div>
        <div class="time-axis">
          ${Array.from({ length: 14 }, (_, hour) => `<span class="time-label" style="top:${hour * 4 * ROW_HEIGHT}px">${String(hour + 9).padStart(2, "0")}:00</span>`).join("")}
        </div>
        <div id="grid-canvas" class="grid-canvas">
          ${renderCourseBlocks(dayIndexes)}
          ${state.assignments.length ? "" : `<div class="empty-grid"><strong>No generated timetable yet.</strong><br>Select Auto-generate to run the WASM search engine on ${formatNumber(state.included.size)} included courses.</div>`}
        </div>
      </div>`;
    filterTimetableBlocks();
  }

  function renderCourseBlocks(dayIndexes) {
    const blocks = [];
    state.assignments
      .filter((assignment) => assignment.assigned)
      .forEach((assignment) => {
        splitIntoBlocks(assignment.classtimeSlots || []).forEach((block) => {
          if (!dayIndexes.includes(block.day)) return;
          blocks.push({ ...block, assignment });
        });
      });

    dayIndexes.forEach((day) => {
      const dayBlocks = blocks
        .filter((block) => block.day === day)
        .sort((a, b) => a.start - b.start || b.length - a.length);
      let cluster = [];
      let clusterEnd = -1;
      const finishCluster = () => {
        if (!cluster.length) return;
        const laneEnds = [];
        cluster.forEach((block) => {
          let lane = laneEnds.findIndex((end) => end <= block.start);
          if (lane === -1) lane = laneEnds.length;
          laneEnds[lane] = block.start + block.length;
          block.lane = lane;
        });
        cluster.forEach((block) => {
          block.laneCount = laneEnds.length;
        });
      };
      dayBlocks.forEach((block) => {
        if (cluster.length && block.start >= clusterEnd) {
          finishCluster();
          cluster = [];
          clusterEnd = -1;
        }
        cluster.push(block);
        clusterEnd = Math.max(clusterEnd, block.start + block.length);
      });
      finishCluster();
    });

    return blocks
      .map((block) => {
        const column = dayIndexes.indexOf(block.day);
        const dayWidth = 100 / dayIndexes.length;
        const laneWidth = dayWidth / (block.laneCount || 1);
        const left = column * dayWidth + (block.lane || 0) * laneWidth;
        const top = block.start * ROW_HEIGHT + 2;
        const height = Math.max(22, block.length * ROW_HEIGHT - 4);
        const color = BLOCK_COLORS[colorIndex(block.assignment.courseID)];
        const course = state.courses.find((item) => item.id === block.assignment.courseID);
        const searchText = [
          block.assignment.code,
          block.assignment.name,
          block.assignment.nameKorean,
          block.assignment.classroom,
          ...(course?.instructors || []),
        ].join(" ");
        return `
          <button
            class="course-block"
            type="button"
            data-course-block="${escapeAttr(block.assignment.courseID)}"
            data-search-text="${escapeAttr(searchText.toLowerCase())}"
            title="${escapeAttr(`${courseLabel(block.assignment)} · ${block.assignment.name || ""} · ${block.assignment.classroom || ""}`)}"
            style="top:${top}px;height:${height}px;left:calc(${left}% + 2px);width:calc(${laneWidth}% - 4px);background:${color}"
          >
            <strong>${escapeHtml(courseLabel(block.assignment))}</strong>
            <span>${escapeHtml(block.assignment.name || block.assignment.nameKorean || "")}</span>
            <span>${escapeHtml(block.assignment.classroom || "")} · ${escapeHtml(slotRange(block))}</span>
          </button>`;
      })
      .join("");
  }

  function splitIntoBlocks(slots) {
    const sorted = Array.from(slots, Number).sort((a, b) => a - b);
    const blocks = [];
    let current = null;
    sorted.forEach((absolute) => {
      const day = Math.floor(absolute / SLOTS_PER_DAY);
      const slot = absolute % SLOTS_PER_DAY;
      if (day < 0 || day > 4) return;
      if (!current || current.day !== day || slot !== current.start + current.length) {
        current = { day, start: slot, length: 1 };
        blocks.push(current);
      } else {
        current.length += 1;
      }
    });
    return blocks;
  }

  function slotRange(block) {
    return `${slotTime(block.start)}–${slotTime(block.start + block.length)}`;
  }

  function slotTime(slot) {
    const minutes = 9 * 60 + slot * 15;
    return `${String(Math.floor(minutes / 60)).padStart(2, "0")}:${String(minutes % 60).padStart(2, "0")}`;
  }

  function filterTimetableBlocks() {
    const query = state.timetableQuery.trim().toLowerCase();
    document.querySelectorAll(".course-block").forEach((block) => {
      block.classList.toggle("is-dimmed", Boolean(query) && !block.dataset.searchText.includes(query));
    });
  }

  async function generateTimetable() {
    if (state.running || !state.included.size) return;
    state.running = true;
    renderTimetablePage();
    setEngineStatus(
      `Searching placements for ${formatNumber(state.included.size)} included courses in this browser…`,
    );

    try {
      const result = await window.DataWeb.generateTimetable({
        includeCourseIDs: Array.from(state.included),
        onProgress(progress) {
          const bar = document.querySelector("#progress-bar");
          const track = document.querySelector("#progress-track");
          const percentLabel = document.querySelector("#progress-percent");
          const count = document.querySelector("#progress-count");
          const copy = document.querySelector("#progress-copy");
          if (!bar || !track || !percentLabel || !count || !copy) return;

          const percent = Math.min(100, Math.max(0, Number(progress.percent) || 0));
          const completed = Math.max(0, Number(progress.completed) || 0);
          const total = Math.max(0, Number(progress.total) || state.included.size);
          const elapsed = formatElapsed(progress.elapsedMs);

          bar.style.width = `${percent}%`;
          track.setAttribute("aria-valuenow", String(percent));
          percentLabel.textContent = `${percent}%`;
          count.textContent = `${formatNumber(completed)} of ${formatNumber(total)} placed`;

          if (progress.phase === "backtracking") {
            const step = Math.max(1, Number(progress.backStep) || 1);
            copy.textContent = `Revising ${formatNumber(step)} earlier ${step === 1 ? "placement" : "placements"} near ${progress.course || "a conflict"} · ${formatNumber(progress.backtracks)} backtracks · ${elapsed}`;
          } else if (progress.phase === "complete") {
            copy.textContent = `All included courses allocated · ${elapsed}`;
          } else if (progress.course) {
            copy.textContent = `Considering ${progress.course} · ${formatNumber(progress.attempts)} search updates · ${elapsed}`;
          } else {
            copy.textContent = `Preparing course constraints · ${elapsed}`;
          }
        },
      });
      state.assignments = Array.from(result.assignments || []);
      state.assignmentByCourse = new Map(
        state.assignments.map((assignment) => [assignment.courseID, assignment]),
      );
      if (!result.ok && result.error) {
        setEngineStatus(`Generation stopped: ${result.error}`, "error");
      } else {
        const assigned = state.assignments.filter((item) => item.assigned).length;
        setEngineStatus(
          `Timetable generated locally — ${formatNumber(assigned)} of ${formatNumber(state.included.size)} included courses allocated.`,
          "ready",
        );
      }
    } catch (error) {
      setEngineStatus(readError(error), "error");
    } finally {
      state.running = false;
      if (state.activeTab === "timetable") renderTimetablePage();
    }
  }

  function downloadTimetable() {
    const rows = [
      ["Course", "Section", "Name", "Classroom", "Classtime"],
      ...state.assignments
        .filter((item) => item.assigned)
        .map((item) => [item.code, item.section, item.name, item.classroom, item.classtime]),
    ];
    const csv = rows
      .map((row) => row.map((value) => `"${text(value).replaceAll('"', '""')}"`).join(","))
      .join("\n");
    const url = URL.createObjectURL(new Blob([csv], { type: "text/csv;charset=utf-8" }));
    const link = document.createElement("a");
    link.href = url;
    link.download = "DATA-browser-demo-timetable.csv";
    link.click();
    URL.revokeObjectURL(url);
  }

  function renderSettings() {
    app.innerHTML = `
      ${pageHeading("Settings", "Review the active simulation profile and browser runtime.", "Application")}
      <section class="cards-grid">
        <article class="settings-card">
          <p class="detail-overline">Semester</p>
          <h2>2026 spring</h2>
          <p>The bundled workbook is configured for the spring semester.</p>
        </article>
        <article class="settings-card">
          <p class="detail-overline">Search engine</p>
          <h2>${escapeHtml(window.DataWeb?.version || "WebAssembly")}</h2>
          <p>The Go scheduling engine runs entirely inside this browser tab.</p>
        </article>
        <article class="settings-card">
          <p class="detail-overline">Data handling</p>
          <h2>Local processing</h2>
          <p>Imported workbooks are held in memory. They are not sent to a server or saved by this demo.</p>
        </article>
      </section>`;
  }

  function renderHelp() {
    app.innerHTML = `
      ${pageHeading("Help", "A compact guide to the DATA browser simulation.", "Documentation")}
      <section class="cards-grid">
        <article class="help-card algorithm-note">
          <p class="detail-overline">Core search idea</p>
          <h2>Adaptive backtracking</h2>
          <p>DATA builds a timetable through trial and error. When a placement creates a dead end, it goes backward and tries another valid choice. The backward step is adjusted as the search develops, helping it escape repeating patterns instead of cycling indefinitely.</p>
          <div class="algorithm-steps">
            <div class="algorithm-step">1. Try a valid room and time for the next course.</div>
            <div class="algorithm-step">2. Detect conflicts with rooms, instructors, and student plans.</div>
            <div class="algorithm-step">3. Vary how far the search steps back, then continue from a different state.</div>
          </div>
        </article>
        <article class="help-card">
          <p class="detail-overline">Getting started</p>
          <h2>Explore the sample</h2>
          <p>The protected sample loads automatically. Inspect the records, choose courses, then open Timetable and select Auto-generate.</p>
        </article>
        <article class="help-card">
          <p class="detail-overline">Your own data</p>
          <h2>Open an Excel workbook</h2>
          <p>Use Open data file on any data page. A workbook matching the DATA format replaces the sample in this tab.</p>
        </article>
        <article class="help-card">
          <p class="detail-overline">Project page</p>
          <h2>Return to the explanation</h2>
          <p><a href="../../explore/automated-timetabling-algorithm/">Read the project description and algorithm overview.</a></p>
        </article>
      </section>`;
  }

  function courseSort(a, b) {
    return (
      (Number(a.year) || 99) - (Number(b.year) || 99) ||
      text(a.code).localeCompare(text(b.code), undefined, { numeric: true }) ||
      text(a.section).localeCompare(text(b.section), undefined, { numeric: true })
    );
  }

  function courseLabel(course) {
    const section = text(course.sectionDisplay || course.section);
    return `${text(course.code) || "Course"}${section ? `-${section}` : ""}`;
  }

  function entityKey(type, record) {
    if (type === "instructors") return text(record.name);
    if (type === "students") return text(record.id);
    return text(record.name);
  }

  function courseNames(ids) {
    if (!ids?.length) return "No courses";
    return ids
      .map((id) => {
        const course = state.courses.find((item) => item.id === id);
        return course ? courseLabel(course) : id;
      })
      .join(", ");
  }

  function colorIndex(value) {
    let hash = 0;
    for (const char of text(value)) hash = (hash * 31 + char.charCodeAt(0)) | 0;
    return Math.abs(hash) % BLOCK_COLORS.length;
  }

  function joinOr(values, fallback) {
    return values?.length ? values.join(", ") : fallback;
  }

  function formatNumber(value) {
    return new Intl.NumberFormat("en").format(Number(value) || 0);
  }

  function formatElapsed(value) {
    const milliseconds = Math.max(0, Number(value) || 0);
    if (milliseconds < 1000) return "under 1s";
    return `${(milliseconds / 1000).toFixed(milliseconds < 10000 ? 1 : 0)}s`;
  }

  function text(value) {
    return value == null ? "" : String(value);
  }

  function readError(error) {
    if (typeof error === "string") return error;
    return error?.error || error?.message || "An unexpected error occurred.";
  }

  function escapeHtml(value) {
    return text(value)
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#039;");
  }

  function escapeAttr(value) {
    return escapeHtml(value).replaceAll("`", "&#096;");
  }

  renderLoading();
  initialize();
})();
