(() => {
  const root = document.getElementById("aes-challenge");

  const showError = (message) => {
    root.innerHTML = `
      <main class="loading-card error-card" role="alert">
        <p class="eyebrow">Unable to start</p>
        <h1>The AES challenge could not load.</h1>
        <p>${message}</p>
        <button class="primary-action" type="button" onclick="window.location.reload()">Try again</button>
      </main>`;
  };

  const load = async () => {
    if (!window.WebAssembly) {
      showError("This browser does not support WebAssembly.");
      return;
    }

    try {
      const go = new Go();
      let result;
      try {
        result = await WebAssembly.instantiateStreaming(
          fetch("aes-challenge.wasm"),
          go.importObject,
        );
      } catch {
        const response = await fetch("aes-challenge.wasm");
        if (!response.ok)
          throw new Error(`WebAssembly request failed (${response.status})`);
        const bytes = await response.arrayBuffer();
        result = await WebAssembly.instantiate(bytes, go.importObject);
      }
      go.run(result.instance).catch((error) => showError(error.message));
    } catch (error) {
      showError(
        error instanceof Error ? error.message : "Unknown loading error",
      );
    }
  };

  load();
})();
