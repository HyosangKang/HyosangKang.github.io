(() => {
  "use strict";

  const state = {
    engineReady: false,
    courses: [],
    selected: new Set(),
    assignments: [],
    filteredCourses: [],
    running: false,
  };

  const elements = {
    statusBanner: document.querySelector("#status-banner"),
    engineStatus: document.querySelector("#engine-status"),
    runtimeBadge: document.querySelector(".runtime-badge"),
    loadSample: document.querySelector("#load-sample"),
    workbookFile: document.querySelector("#workbook-file"),
    importSummary: document.querySelector("#import-summary"),
    summaryCourses: document.querySelector("#summary-courses"),
    summaryInstructors: document.querySelector("#summary-instructors"),
    summaryClassrooms: document.querySelector("#summary-classrooms"),
    summaryStudents: document.querySelector("#summary-students"),
    courseSearch: document.querySelector("#course-search"),
    courseList: document.querySelector("#course-list"),
    selectedCount: document.querySelector("#selected-count"),
    selectVisible: document.querySelector("#select-visible"),
    clearVisible: document.querySelector("#clear-visible"),
    generate: document.querySelector("#generate"),
    runProgress: document.querySelector("#run-progress"),
    progressLabel: document.querySelector("#progress-label"),
    progressPercent: document.querySelector("#progress-percent"),
    progressBar: document.querySelector("#progress-bar"),
    progressCourse: document.querySelector("#progress-course"),
    cancelRun: document.querySelector("#cancel-run"),
    resultSummary: document.querySelector("#result-summary"),
    assignedCount: document.querySelector("#assigned-count"),
    roomCount: document.querySelector("#room-count"),
    dayCount: document.querySelector("#day-count"),
    resultsEmpty: document.querySelector("#results-empty"),
    resultsTableWrap: document.querySelector("#results-table-wrap"),
    resultsList: document.querySelector("#results-list"),
    resultsDescription: document.querySelector("#results-description"),
    resultSearch: document.querySelector("#result-search"),
    downloadResults: document.querySelector("#download-results"),
    newRun: document.querySelector("#new-run"),
  };

  function setStatus(message, type = "loading") {
    elements.engineStatus.textContent = message;
    elements.statusBanner.classList.toggle("is-ready", type === "ready");
    elements.statusBanner.classList.toggle("is-error", type === "error");
  }

  async function initializeEngine() {
    try {
      if (typeof Go !== "function") {
        throw new Error("The WebAssembly support script could not be loaded.");
      }

      const go = new Go();
      const response = await fetch("dataweb.wasm");
      if (!response.ok) {
        throw new Error(
          `The timetabling engine could not be downloaded (${response.status}).`,
        );
      }

      let result;
      try {
        result = await WebAssembly.instantiateStreaming(
          response.clone(),
          go.importObject,
        );
      } catch {
        result = await WebAssembly.instantiate(
          await response.arrayBuffer(),
          go.importObject,
        );
      }

      go.run(result.instance);
      await waitForEngineApi();
      await window.DataWeb.ready;

      state.engineReady = true;
      elements.loadSample.disabled = false;
      elements.workbookFile.disabled = false;
      elements.runtimeBadge.classList.add("is-ready");
      setStatus(
        "The timetabling engine is ready. Choose a workbook to begin.",
        "ready",
      );
    } catch (error) {
      setStatus(readError(error), "error");
    }
  }

  function waitForEngineApi() {
    return new Promise((resolve, reject) => {
      let attempts = 0;
      const check = () => {
        if (window.DataWeb) {
          resolve();
          return;
        }
        attempts += 1;
        if (attempts > 200) {
          reject(new Error("The timetabling engine did not start."));
          return;
        }
        window.setTimeout(check, 25);
      };
      check();
    });
  }

  async function loadSampleWorkbook() {
    elements.loadSample.disabled = true;
    setStatus("Loading the sample workbook…");
    try {
      const response = await fetch("sample-data.xlsx");
      if (!response.ok) {
        throw new Error(
          `The sample workbook could not be loaded (${response.status}).`,
        );
      }
      await importWorkbook(
        new Uint8Array(await response.arrayBuffer()),
        "Sample workbook",
      );
    } catch (error) {
      setStatus(readError(error), "error");
    } finally {
      elements.loadSample.disabled = !state.engineReady;
    }
  }

  async function loadChosenWorkbook(event) {
    const [file] = event.target.files;
    if (!file) return;

    setStatus(`Reading ${file.name}…`);
    try {
      await importWorkbook(new Uint8Array(await file.arrayBuffer()), file.name);
    } catch (error) {
      setStatus(readError(error), "error");
    } finally {
      event.target.value = "";
    }
  }

  async function importWorkbook(bytes, sourceName) {
    if (!state.engineReady) {
      throw new Error("The timetabling engine is still loading.");
    }

    await new Promise((resolve) => window.setTimeout(resolve, 20));
    const result = window.DataWeb.importWorkbook(bytes);
    if (!result.ok) {
      throw new Error(
        (result.errors || ["The workbook could not be imported."]).join(" · "),
      );
    }

    state.courses = Array.from(result.courses || []);
    state.selected = new Set(state.courses.map((course) => course.id));
    state.assignments = [];

    elements.summaryCourses.textContent = formatNumber(result.summary.courses);
    elements.summaryInstructors.textContent = formatNumber(
      result.summary.instructors,
    );
    elements.summaryClassrooms.textContent = formatNumber(
      result.summary.classrooms,
    );
    elements.summaryStudents.textContent = formatNumber(
      result.summary.students,
    );
    elements.importSummary.hidden = false;

    renderCourses();
    enableStep("courses-panel");
    setStatus(
      `${sourceName} imported successfully. ${state.courses.length} courses are ready.`,
      "ready",
    );
    showStep("courses-panel");
  }

  function renderCourses() {
    const query = elements.courseSearch.value.trim().toLowerCase();
    state.filteredCourses = state.courses.filter((course) => {
      const haystack =
        `${course.code} ${course.section} ${course.name} ${course.nameKorean}`.toLowerCase();
      return haystack.includes(query);
    });

    const fragment = document.createDocumentFragment();
    state.filteredCourses.forEach((course) => {
      const row = document.createElement("tr");
      const checked = state.selected.has(course.id);
      row.innerHTML = `
        <td><input class="course-check" type="checkbox" aria-label="Include ${escapeHtml(course.code)}" data-course-id="${escapeHtml(course.id)}" ${checked ? "checked" : ""}></td>
        <td class="course-code">${escapeHtml(course.code)}${course.section ? `-${escapeHtml(course.section)}` : ""}</td>
        <td>${escapeHtml(course.name || course.nameKorean || "Untitled course")}</td>
        <td class="muted-cell">${escapeHtml(String(course.credit ?? "—"))}</td>
        <td class="muted-cell">${formatNumber(course.students || 0)}</td>
      `;
      fragment.appendChild(row);
    });

    elements.courseList.replaceChildren(fragment);
    updateSelectionCount();
  }

  function updateSelectionCount() {
    elements.selectedCount.textContent = formatNumber(state.selected.size);
    elements.generate.disabled = state.running || state.selected.size === 0;
  }

  function changeVisibleSelection(selected) {
    state.filteredCourses.forEach((course) => {
      if (selected) state.selected.add(course.id);
      else state.selected.delete(course.id);
    });
    renderCourses();
  }

  async function generateTimetable() {
    if (state.running || state.selected.size === 0) return;

    state.running = true;
    updateSelectionCount();
    enableStep("results-panel");
    showStep("results-panel");
    elements.newRun.disabled = true;
    elements.runProgress.hidden = false;
    elements.resultSummary.hidden = true;
    elements.resultsEmpty.hidden = true;
    elements.resultsTableWrap.hidden = true;
    elements.progressBar.classList.add("is-searching");
    elements.progressBar.style.width = "";
    elements.progressPercent.textContent = "Searching";
    elements.progressLabel.textContent = "Trying course placements…";
    elements.progressCourse.textContent = "Preparing the search order";
    setStatus("Generating a timetable in this browser…");

    try {
      const result = await window.DataWeb.generateTimetable({
        includeCourseIDs: Array.from(state.selected),
        onProgress(progress) {
          const percent = Math.max(
            0,
            Math.min(99, Number(progress.percent) || 0),
          );
          elements.progressCourse.textContent = progress.course
            ? `Considering ${progress.course}`
            : "Revising the partial timetable";
          if (percent > 2) {
            elements.progressBar.classList.remove("is-searching");
            elements.progressBar.style.width = `${percent}%`;
            elements.progressPercent.textContent = `${percent}%`;
          }
        },
      });

      if (!result.ok) {
        if (String(result.error).toLowerCase().includes("cancel")) {
          throw new Error("The generation run was cancelled.");
        }
        throw new Error(
          result.error || "The timetable could not be completed.",
        );
      }

      state.assignments = Array.from(result.assignments || []).filter(
        (assignment) => assignment.assigned,
      );
      elements.progressBar.classList.remove("is-searching");
      elements.progressBar.style.width = "100%";
      elements.progressPercent.textContent = "100%";
      elements.progressLabel.textContent = "Timetable complete";
      elements.progressCourse.textContent = `${state.assignments.length} courses received a time and room.`;
      elements.cancelRun.hidden = true;
      renderResults();
      setStatus("Timetable generated successfully.", "ready");
    } catch (error) {
      elements.progressBar.classList.remove("is-searching");
      elements.progressBar.style.width = "0";
      elements.progressPercent.textContent = "Stopped";
      elements.progressLabel.textContent = "Generation stopped";
      elements.progressCourse.textContent = readError(error);
      elements.resultsEmpty.hidden = false;
      setStatus(readError(error), "error");
    } finally {
      state.running = false;
      elements.cancelRun.hidden = false;
      elements.newRun.disabled = false;
      updateSelectionCount();
    }
  }

  function renderResults() {
    const query = elements.resultSearch.value.trim().toLowerCase();
    const assignments = state.assignments
      .filter((assignment) => {
        const haystack =
          `${assignment.code} ${assignment.section} ${assignment.name} ${assignment.nameKorean} ${assignment.classtime} ${assignment.classroom}`.toLowerCase();
        return haystack.includes(query);
      })
      .sort((a, b) => {
        const firstA = Math.min(...Array.from(a.classtimeSlots || [9999]));
        const firstB = Math.min(...Array.from(b.classtimeSlots || [9999]));
        return firstA - firstB || String(a.code).localeCompare(String(b.code));
      });

    const fragment = document.createDocumentFragment();
    assignments.forEach((assignment) => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td class="course-code">${escapeHtml(assignment.code || assignment.courseID)}${assignment.section ? `-${escapeHtml(assignment.section)}` : ""}</td>
        <td>${escapeHtml(assignment.name || assignment.nameKorean || "Untitled course")}</td>
        <td>${escapeHtml(assignment.classtime || "—")}</td>
        <td>${escapeHtml(assignment.classroom || "—")}</td>
      `;
      fragment.appendChild(row);
    });
    elements.resultsList.replaceChildren(fragment);

    const rooms = new Set(
      state.assignments.map((item) => item.classroom).filter(Boolean),
    );
    const days = new Set();
    state.assignments.forEach((item) => {
      Array.from(item.classtimeSlots || []).forEach((slot) =>
        days.add(Math.floor(slot / 52)),
      );
    });

    elements.assignedCount.textContent = formatNumber(state.assignments.length);
    elements.roomCount.textContent = formatNumber(rooms.size);
    elements.dayCount.textContent = formatNumber(days.size);
    elements.resultSummary.hidden = false;
    elements.resultsEmpty.hidden = true;
    elements.resultsTableWrap.hidden = false;
    elements.resultsDescription.textContent =
      "Each row is an assignment returned directly by the WebAssembly timetabling engine.";
  }

  function downloadCsv() {
    if (!state.assignments.length) return;

    const rows = [
      ["Course", "Section", "Name", "Time", "Room"],
      ...state.assignments.map((item) => [
        item.code || item.courseID,
        item.section || "",
        item.name || item.nameKorean || "",
        item.classtime || "",
        item.classroom || "",
      ]),
    ];
    const csv = rows.map((row) => row.map(csvCell).join(",")).join("\n");
    const blob = new Blob([`\uFEFF${csv}`], { type: "text/csv;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "data-timetable.csv";
    link.click();
    URL.revokeObjectURL(url);
  }

  function showStep(panelId) {
    document.querySelectorAll(".step-panel").forEach((panel) => {
      panel.classList.toggle("is-visible", panel.id === panelId);
    });
    document.querySelectorAll(".step").forEach((step) => {
      const active = step.dataset.stepTarget === panelId;
      step.classList.toggle("is-active", active);
      if (active) step.setAttribute("aria-current", "step");
      else step.removeAttribute("aria-current");
    });
  }

  function enableStep(panelId) {
    const step = document.querySelector(`[data-step-target="${panelId}"]`);
    if (step) step.disabled = false;
  }

  function readError(error) {
    if (!error) return "An unexpected error occurred.";
    if (typeof error === "string") return error;
    if (typeof error.error === "string") return error.error;
    if (error.message) return error.message;
    return "An unexpected error occurred.";
  }

  function formatNumber(value) {
    return new Intl.NumberFormat().format(Number(value) || 0);
  }

  function escapeHtml(value) {
    return String(value ?? "")
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#039;");
  }

  function csvCell(value) {
    return `"${String(value ?? "").replaceAll('"', '""')}"`;
  }

  document.querySelectorAll(".step").forEach((step) => {
    step.addEventListener("click", () => {
      if (!step.disabled) showStep(step.dataset.stepTarget);
    });
  });

  document.querySelectorAll("[data-go-step]").forEach((button) => {
    button.addEventListener("click", () => showStep(button.dataset.goStep));
  });

  elements.loadSample.addEventListener("click", loadSampleWorkbook);
  elements.workbookFile.addEventListener("change", loadChosenWorkbook);
  elements.courseSearch.addEventListener("input", renderCourses);
  elements.selectVisible.addEventListener("click", () =>
    changeVisibleSelection(true),
  );
  elements.clearVisible.addEventListener("click", () =>
    changeVisibleSelection(false),
  );
  elements.courseList.addEventListener("change", (event) => {
    const checkbox = event.target.closest("[data-course-id]");
    if (!checkbox) return;
    if (checkbox.checked) state.selected.add(checkbox.dataset.courseId);
    else state.selected.delete(checkbox.dataset.courseId);
    updateSelectionCount();
  });
  elements.generate.addEventListener("click", generateTimetable);
  elements.cancelRun.addEventListener("click", () => {
    if (state.running) window.DataWeb.cancel();
  });
  elements.newRun.addEventListener("click", () => showStep("courses-panel"));
  elements.resultSearch.addEventListener("input", renderResults);
  elements.downloadResults.addEventListener("click", downloadCsv);

  initializeEngine();
})();
