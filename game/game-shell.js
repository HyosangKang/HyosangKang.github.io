(() => {
  const body = document.body;
  const safeGameText = (value) =>
    String(value)
      .toUpperCase()
      .replace(/[^A-Z0-9 ]+/g, " ")
      .replace(/\s+/g, " ")
      .trim();
  const title = safeGameText(body.dataset.gameTitle || "Math Project");
  const subtitle = safeGameText(
    body.dataset.gameSubtitle || "Interactive mathematical model",
  );
  const wasm = body.dataset.gameWasm;
  const controls = (body.dataset.gameControls || "")
    .split("|")
    .map((item) => safeGameText(item))
    .filter(Boolean);

  body.innerHTML = `
    <div class="game-shell">
      <header class="game-topbar">
        <div class="game-heading">
          <p class="game-kicker">MATH PROJECT GAME</p>
          <div class="game-title-row">
            <h1 class="game-title">${title}</h1>
            <p class="game-subtitle">${subtitle}</p>
          </div>
        </div>
        <div class="game-actions">
          <span class="game-status" id="game-status">
            <span class="game-status-dot" aria-hidden="true"></span>
            <span class="game-status-text">LOADING</span>
          </span>
          <details class="game-help">
            <summary aria-label="Show controls" title="Controls">H</summary>
            <div class="game-help-panel">
              <strong>CONTROLS</strong>
              <ul>${controls.map((item) => `<li>${item}</li>`).join("")}</ul>
            </div>
          </details>
          <button class="game-action" id="game-restart" type="button" aria-label="Restart game" title="Restart">R</button>
          <button class="game-action" id="game-fullscreen" type="button" aria-label="Open full screen" title="Full screen">F</button>
        </div>
      </header>
      <main class="game-stage" id="game-stage" aria-label="${title} game canvas">
        <div class="game-loading" id="game-loading">
          <span class="game-loading-spinner" aria-hidden="true"></span>
          <span>LOADING GAME</span>
        </div>
      </main>
    </div>
  `;

  const stage = document.getElementById("game-stage");
  const loading = document.getElementById("game-loading");
  const status = document.getElementById("game-status");
  const statusText = status.querySelector(".game-status-text");

  document
    .getElementById("game-restart")
    .addEventListener("click", () => window.location.reload());
  document
    .getElementById("game-fullscreen")
    .addEventListener("click", async () => {
      if (document.fullscreenElement) {
        await document.exitFullscreen();
      } else {
        await document.documentElement.requestFullscreen();
      }
    });

  const mountCanvas = () => {
    const canvas = document.querySelector("body > canvas");
    if (!canvas) return false;
    stage.append(canvas);
    loading.remove();
    status.classList.add("is-ready");
    statusText.textContent = "READY";
    return true;
  };

  const canvasObserver = new MutationObserver(() => {
    if (mountCanvas()) canvasObserver.disconnect();
  });
  canvasObserver.observe(body, { childList: true });

  const showError = (message) => {
    loading.classList.add("game-error");
    loading.innerHTML = `<strong>GAME COULD NOT START</strong><span>${safeGameText(message)}</span>`;
    statusText.textContent = "UNAVAILABLE";
  };

  const start = async () => {
    try {
      if (!wasm) throw new Error("The WebAssembly file was not specified.");
      const go = new Go();
      let result;
      if (WebAssembly.instantiateStreaming) {
        try {
          result = await WebAssembly.instantiateStreaming(
            fetch(wasm),
            go.importObject,
          );
        } catch {
          const source = await (await fetch(wasm)).arrayBuffer();
          result = await WebAssembly.instantiate(source, go.importObject);
        }
      } else {
        const source = await (await fetch(wasm)).arrayBuffer();
        result = await WebAssembly.instantiate(source, go.importObject);
      }
      go.run(result.instance);
      mountCanvas();
    } catch (error) {
      showError(
        error instanceof Error ? error.message : "Unknown loading error",
      );
    }
  };

  start();
})();
