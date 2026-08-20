document
  .querySelectorAll("details.project-code[data-code-url]")
  .forEach((panel) => {
    const code = panel.querySelector("code");

    const loadCode = async () => {
      if (!code || panel.dataset.codeState) return;

      panel.dataset.codeState = "loading";
      code.textContent = "Loading MATLAB code…";

      try {
        const response = await fetch(panel.dataset.codeUrl);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);

        code.textContent = await response.text();
        panel.dataset.codeState = "loaded";
      } catch (_error) {
        code.textContent =
          "The code could not be loaded. Please use the download link above.";
        panel.dataset.codeState = "error";
      }
    };

    panel.addEventListener("toggle", () => {
      if (panel.open) loadCode();
    });

    if (panel.open) loadCode();
  });
