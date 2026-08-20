# Math Projects game sources

This directory keeps the Go source for every WebAssembly game published on the
site.

- `fill`, `predator`, `rocket`, `rotate`, and `virus` were copied from the
  `multi-game-book` repository.
- `throw` was copied from its original teaching-game repository.
- `stone-skipping` is a new implementation based on
  `reference/StoneSkipping.m`.
- `lanchester` extends the earlier two-force Go model into a multi-team
  Ebitengine game with random laser targeting and tick-by-tick history.
- `differential-labs` contains the exact Laplace/modal models and shared
  Ebitengine interface used by the double-circuit and double-spring games.

The compiled `.wasm` files and their web interface live under `game/`. Run
`scripts/build_games.sh` from the repository root to rebuild every game.
