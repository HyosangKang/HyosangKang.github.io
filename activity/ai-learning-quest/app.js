(() => {
  "use strict";

  const config = window.LEARNING_QUEST_CONFIG || {};
  const state = {
    intakeOpen: Boolean(config.demoMode),
    currentLevel: 1,
    requestId: "",
  };

  const startScreen = document.getElementById("start-screen");
  const startForm = document.getElementById("start-form");
  const startMessage = document.getElementById("start-message");
  const questGame = document.getElementById("quest-game");
  const questForm = document.getElementById("quest-form");
  const completeScreen = document.getElementById("complete-screen");
  const status = document.getElementById("intake-status");
  const statusLabel = document.getElementById("status-label");
  const submitButton = document.getElementById("submit-quest");
  const submissionError = document.getElementById("submission-error");
  const sessionCode = document.getElementById("session-code");
  const submissionCode = document.getElementById("submission-code");
  const requestIdInput = document.getElementById("request-id");
  const coachLink = document.getElementById("open-coach");

  const makeRequestId = () => {
    if (window.crypto && crypto.randomUUID) return crypto.randomUUID();
    return `${Date.now()}-${Math.random().toString(16).slice(2)}`;
  };

  const setStatus = (kind, label) => {
    status.classList.remove("is-open", "is-closed", "is-error");
    if (kind) status.classList.add(`is-${kind}`);
    statusLabel.textContent = label;
  };

  const readStatus = () => {
    if (config.demoMode || !config.endpoint) {
      state.intakeOpen = true;
      setStatus("open", "Demo intake open");
      return;
    }

    const callbackName = `learningQuestStatus_${Date.now()}`;
    const script = document.createElement("script");
    const timer = window.setTimeout(() => {
      cleanup();
      state.intakeOpen = false;
      setStatus("error", "Status unavailable");
    }, 8000);

    const cleanup = () => {
      window.clearTimeout(timer);
      delete window[callbackName];
      script.remove();
    };

    window[callbackName] = (payload) => {
      state.intakeOpen = Boolean(payload && payload.open);
      setStatus(
        state.intakeOpen ? "open" : "closed",
        state.intakeOpen
          ? `${payload.sessionTitle || "Class intake"} · Open`
          : "Intake closed",
      );
      cleanup();
    };

    script.onerror = () => {
      state.intakeOpen = false;
      setStatus("error", "Status unavailable");
      cleanup();
    };
    script.src = `${config.endpoint}?action=status&callback=${callbackName}&t=${Date.now()}`;
    document.head.append(script);
  };

  const showLevel = (level) => {
    state.currentLevel = level;
    document.querySelectorAll("[data-level]").forEach((panel) => {
      panel.hidden = Number(panel.dataset.level) !== level;
    });
    document.querySelectorAll("[data-progress]").forEach((item) => {
      const itemLevel = Number(item.dataset.progress);
      item.classList.toggle("is-current", itemLevel === level);
      item.classList.toggle("is-complete", itemLevel < level);
    });
    document.getElementById("control-points").textContent = `${level - 1} / 4`;
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const validateLevel = (level) => {
    const panel = document.querySelector(`[data-level="${level}"]`);
    const controls = panel.querySelectorAll("input, textarea");
    for (const control of controls) {
      if (!control.checkValidity()) {
        control.reportValidity();
        return false;
      }
    }
    return true;
  };

  const selectedMission = () =>
    questForm.querySelector('input[name="mission"]:checked')?.value ||
    "this question";

  const starterPrompt =
    () => `You are my learning coach and discussion partner for Engineering Mathematics I. Help me explore: ${selectedMission()}

Do not decide the answer for me. Begin with three different directions and a fourth option that lets me reject them and propose my own. Ask one question at a time, let me disagree or redirect you, and distinguish facts from assumptions. Focus on mathematical reasoning and engineering applications rather than merely obtaining correct answers. Do not draft a conclusion until I ask. When I do, use only ideas I have accepted and ask me what I want to keep, reject, or rewrite.`;

  startForm.addEventListener("submit", (event) => {
    event.preventDefault();
    startMessage.classList.remove("is-error");
    if (!state.intakeOpen) {
      startMessage.textContent =
        "This class intake is closed. Please wait for the lecturer to open it.";
      startMessage.classList.add("is-error");
      return;
    }
    const code = sessionCode.value.trim().toUpperCase();
    if (!code) return;
    submissionCode.value = code;
    startScreen.hidden = true;
    questGame.hidden = false;
    showLevel(1);
  });

  document.querySelectorAll(".next-button").forEach((button) => {
    button.addEventListener("click", () => {
      if (!validateLevel(state.currentLevel)) return;
      showLevel(Math.min(4, state.currentLevel + 1));
    });
  });

  document.querySelectorAll(".back-button").forEach((button) => {
    button.addEventListener("click", () =>
      showLevel(Math.max(1, state.currentLevel - 1)),
    );
  });

  document.querySelectorAll("textarea").forEach((textarea) => {
    const counter = document.querySelector(`[data-count-for="${textarea.id}"]`);
    textarea.addEventListener("input", () => {
      counter.textContent = String(textarea.value.length);
    });
  });

  document
    .getElementById("copy-prompt")
    .addEventListener("click", async (event) => {
      const button = event.currentTarget;
      try {
        await navigator.clipboard.writeText(starterPrompt());
        button.textContent = "Prompt copied ✓";
      } catch {
        button.textContent = "Copy unavailable";
      }
      window.setTimeout(() => {
        button.textContent = "Copy starter prompt";
      }, 1800);
    });

  coachLink.href = config.coachUrl || "https://chatgpt.com/";

  questForm.addEventListener("submit", (event) => {
    if (!validateLevel(4)) {
      event.preventDefault();
      return;
    }

    submissionError.hidden = true;
    state.requestId = makeRequestId();
    requestIdInput.value = state.requestId;
    submitButton.disabled = true;
    submitButton.textContent = "Submitting…";

    if (config.demoMode || !config.endpoint) {
      event.preventDefault();
      const formData = Object.fromEntries(new FormData(questForm).entries());
      const responses = JSON.parse(
        localStorage.getItem("learningQuestDemoResponses") || "[]",
      );
      responses.push({ submittedAt: new Date().toISOString(), ...formData });
      localStorage.setItem(
        "learningQuestDemoResponses",
        JSON.stringify(responses),
      );
      window.setTimeout(() => completeSubmission(), 450);
      return;
    }

    questForm.action = config.endpoint;
    window.setTimeout(() => {
      if (submitButton.disabled) {
        failSubmission(
          "The submission took too long. Check your connection and try again.",
        );
      }
    }, 12000);
  });

  const completeSubmission = () => {
    submitButton.disabled = false;
    submitButton.textContent = "Submit contribution";
    questGame.hidden = true;
    completeScreen.hidden = false;
    document.getElementById("control-points").textContent = "4 / 4";
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const failSubmission = (message) => {
    submitButton.disabled = false;
    submitButton.textContent = "Submit contribution";
    submissionError.textContent = message;
    submissionError.hidden = false;
  };

  window.addEventListener("message", (event) => {
    const payload = event.data;
    if (!payload || payload.type !== "learning-quest-submit") return;
    if (payload.requestId !== state.requestId) return;
    if (payload.ok) completeSubmission();
    else
      failSubmission(payload.message || "The response could not be submitted.");
  });

  document.getElementById("new-response").addEventListener("click", () => {
    questForm.reset();
    document.querySelectorAll("[data-count-for]").forEach((counter) => {
      counter.textContent = "0";
    });
    completeScreen.hidden = true;
    startScreen.hidden = false;
    submitButton.disabled = false;
    submissionError.hidden = true;
    readStatus();
  });

  readStatus();
})();
