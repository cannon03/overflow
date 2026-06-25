<pre>
███████ ███████ ███████ ███████ ███████ ███████ ███████
██   ██ ██   ██ ██   ██ ██   ██ ██   ██ ██     ██   ██
██   ██ ██   ██ ███████ ███████ ███████ ██     ███████
██   ██ ██   ██ ██   ██ ██   ██ ██   ██ ██     ██   ██
███████ ███████ ██   ██ ██   ██ ██   ██ ███████ ██   ██

                     ███████████████
                     █  BULLET HELL  █
                     ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
</pre>

# OVERFLOW — ASCII Bullet Hell

A fast-paced terminal bullet hell game built in Go, rendered with ANSI escape codes and Unicode block graphics. Playable natively in your terminal and in the browser via WebAssembly + xterm.js.

## Quick Start

```bash
# Run natively (Linux/macOS terminal)
make native && ./overflow

# Run in browser (requires Python 3)
make serve
# → http://localhost:8080
```

## Controls

| Key | Action |
|---|---|
| `WASD` / `Arrow Keys` | Move |
| `Space` | Shoot |
| `P` | Pause |
| `ESC` | Quit |
| `R` | Restart (game over) |

## Gameplay

Survive increasingly difficult waves of enemies. Each wave brings more enemies with faster movement, quicker fire rates, and deadlier bullet patterns.

- **Basic enemies** — Slow, single-shot attackers
- **Fast enemies** — Quick movers with spread patterns
- **Tank enemies** — Heavy, multi-hit enemies
- **Boss waves (every 5th wave)** — Massive enemies with burst fire

## Building

```bash
# Native binary
make native

# WASM + web assets
make web

# Everything
make all

# Clean artifacts
make clean
```

## Deployment (Netlify)

The project is pre-configured for Netlify. Connect your Netlify site to this repository and it will automatically build and deploy the WASM version.

Or deploy manually:

```bash
# Build everything
make web

# The `web/` directory is ready to deploy
# Point your Netlify site to publish from `web/`
```

## Project Structure

```
├── Makefile              # Build system
├── netlify.toml          # Netlify deployment config
├── README.md
├── cmd/
│   └── game/main.go      # Entry point
├── internal/
│   ├── engine/           # Game loop, terminal setup, WASM support
│   ├── render/           # Cell, Framebuffer, Sprite, ANSI renderer
│   ├── input/            # Keyboard input (Unix + WASM)
│   ├── entity/           # Player, Enemy, Bullet, Particle
│   ├── world/            # Wave system, collision, game state
│   ├── ui/               # HUD (HP, score, wave, FPS)
│   └── assets/           # Sprite definitions
└── web/
    ├── index.html        # Browser entry with xterm.js
    ├── main.wasm         # Compiled WASM binary (generated)
    └── wasm_exec.js      # Go WASM runtime (generated)
```

## Architecture

```
                    ┌────────────────────┐
                    │    Game Loop (60fps)│
                    └─────────┬──────────┘
                              │
             ┌────────────────┼────────────────┐
             │                │                │
      Input System      World Update      Rendering
             │                │                │
             │         Entity Manager          │
             │                │                │
             │    ┌────┬────┬─┴──┬────┬────┐   │
             │    │Plyr│Enmy│Bull│Part│Wave│   │
             │    └────┴────┴────┴────┴────┘   │
             │                │                │
             └────────────────┼────────────────┘
                              │
                    ┌─────────┴─────────┐
                    │   ANSI Renderer   │
                    │  (Framebuffer →   │
                    │   terminal escape │
                    │      codes)       │
                    └───────────────────┘
```

## Tech Stack

- **Language:** Go
- **Rendering:** ANSI escape sequences + 24-bit RGB color
- **Graphics:** Unicode block characters
- **Loop:** Fixed timestep (60 FPS)
- **Browser:** WebAssembly + xterm.js

## License

MIT
