#!/bin/sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
wasm_exec=$(go env GOROOT)/lib/wasm/wasm_exec.js

if [ ! -f "$wasm_exec" ]; then
  echo "Could not find wasm_exec.js at $wasm_exec" >&2
  exit 1
fi

for game_name in fill predator rocket rotate throw virus stone-skipping math-arcade lanchester
do
  source_dir="$repo_root/game-source/$game_name"
  output_dir="$repo_root/game/$game_name"
  mkdir -p "$output_dir"

  (
    cd "$source_dir"
    GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o "$output_dir/$game_name.wasm" .
  )

  cp "$wasm_exec" "$output_dir/wasm_exec.js"
  echo "Built $game_name"
done

for game_name in electric-circuit double-spring
do
  source_dir="$repo_root/game-source/differential-labs"
  output_dir="$repo_root/game/$game_name"
  mkdir -p "$output_dir"

  (
    cd "$source_dir"
    GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o "$output_dir/$game_name.wasm" .
  )

  cp "$wasm_exec" "$output_dir/wasm_exec.js"
  echo "Built $game_name"
done
