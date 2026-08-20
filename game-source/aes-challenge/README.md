# AES Challenge WebAssembly simulator

This simulator uses the 16 strings from the Fall 2023 BS203 Linear Algebra
assignment. Go owns the round order, classification feedback, hidden-key
extraction, and AES-128 block decryption. The browser files under
`game/aes-challenge/` only load the compiled WebAssembly program and style its
accessible HTML interface.

## Test and build

From this directory:

```sh
go test ./...
GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" \
  -o ../../game/aes-challenge/aes-challenge.wasm .
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" \
  ../../game/aes-challenge/wasm_exec.js
```

From the repository root, `./scripts/build_games.sh` performs the same build
along with the other Go games.

WebAssembly must be served over HTTP rather than opened as a local file. For a
quick preview from the repository root:

```sh
python3 -m http.server 4000
```

Then open <http://localhost:4000/game/aes-challenge/main.html>.

## Publish with GitHub Pages

Commit both this source directory and the generated files in
`game/aes-challenge/`, including `aes-challenge.wasm` and `wasm_exec.js`. The
AES learning page embeds `game/aes-challenge/main.html` through
`_data/project_media.yml`. Pushing the repository's main publishing branch
triggers `.github/workflows/deploy.yml`, whose Jekyll build copies the static
game directory into the GitHub Pages site.

After deployment, verify these two URLs:

- `/learn/linear-algebra/aes-encryption/`
- `/game/aes-challenge/main.html`
