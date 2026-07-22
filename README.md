# Swarm Defense

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-or-square&logo=go)](https://go.dev)
[![Engine](https://img.shields.io/badge/Engine-Ebitengine-6c5ce7?style=flat-or-square)](https://ebitengine.org)

**Swarm Defense** is a feature-complete 2D tower defense survival game built entirely using the [Go programming language](https://go.dev) and the [Ebitengine](https://ebitengine.org) 2D game engine. 

What makes this game truly unique is its origin story: **every single line of code, sprite sheet, UI asset, save state utility, and sound effect was designed, developed, and compiled programmatically by an autonomous swarm of AI subagents** under the orchestration of a Swarm Coordinator.

> 📖 **Read the Story:** To understand the engineering architecture, agent dynamics, and full narrative behind how this game was built, read the published article: **[The Rise of the Subagents on danicat.dev](https://danicat.dev/posts/20260722-the-rise-of-the-subagents/)**.

---

## 🕹️ Game Features

*   **16-bit Retro Aesthetics:** Custom-designed 32x32 sprites with distinct 256-color palettes and multi-frame movement animations.
*   **Dual-Phase Gameplay:** 
    *   *Build Phase:* Strategically construct defensive structures and place tactical defender units.
    *   *Battle Phase:* Survive onslaughts of progressively harder waves of enemies.
*   **Diverse Entities:**
    *   **4 Playable Units** and **4 Buildable Defensive Structures**, each with unique properties and tactical positioning advantages.
    *   **8 Monster Wave Types**, concluding with a massive, challenging boss monster.
*   **Dynamic Audio Synthesis:** Sound effects (blips, shots, base damage, and enemy explosions) are synthesized **mathematically and programmatically** from pure sine, square, and triangle wave frequencies—no static audio assets required!
*   **Polished Interface:** Complete with a cinematic intro sequence, retro title screen, custom gameplay UI, victory screen, game-over screen, and local persistent high-score tracking.

---

## 🛠️ Codebase Architecture

The codebase was decomposed by the AI swarm into highly cohesive, specialized modules:

| File | Subagent / Specialist | Purpose |
| :--- | :--- | :--- |
| `main.go` | **System Architect** | Entry point that initializes the OS window, configures audio contexts, and launches the Ebitengine game loop. |
| `game.go` | **Gameplay Coordinator** | Controls core scene transitions, phase switching (build/battle), level updates, and frame-by-frame orchestration. |
| `entities.go` | **Backend Engineer** | Defines state models, attack patterns, stat blocks, and movement pathfinding for players, towers, and monsters. |
| `sprites.go` | **Art Specialist** | Houses custom binary/bitmap rendering definitions, programmatically designing every frame of sprite and tile animation on the fly. |
| `sound.go` | **Audio Engineer** | Contains mathematical sound synth generators that calculate and pipe live audio frequencies directly to the OS audio stream. |
| `ui.go` | **Frontend Engineer** | Renders HUD overlays, coin trackers, health bars, retro menus, and dialogs. |
| `save.go` | **Systems QA** | Handles local IO operation to safely serialize, persist, and load player high-scores. |

---

## 🚀 How to Run

### Prerequisites

To play Swarm Defense, you need to have **Go (1.26.0 or higher)** installed on your machine.

### Run Locally

1. Clone or download this repository.
2. Navigate to the project root directory.
3. Fetch dependencies and run the game directly:

```bash
# Tidy modules and fetch Ebitengine
go mod tidy

# Run the game
go run .
```

---

## 📄 License

This project is open-source and available under the MIT License.
