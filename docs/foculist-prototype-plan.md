# Foculist Prototype - Implementation Plan

This document outlines the step-by-step plan to build **Foculist**, the macOS-first Pomodoro application using Go and Wails v3, based on the project requirements.

## Phase 1: Project Initialization & Setup
*   **Initialize Repository:** Create the project directory and initialize a new Git repository.
*   **Wails v3 Initialization:** Scaffold a new Wails v3 project. We will use a lightweight frontend framework like Svelte to keep the footprint minimal.
*   **Directory Structure Definition:** Set up the required architecture:
    *   `main.go`: Application entry point and Wails configuration.
    *   `backend/timer.go`: The core timer domain logic.
    *   `frontend/`: The frontend application directory.
    *   `build/`: For tray icons and assets.
*   **Asset Preparation:** Create a `TrayIconTemplate.png` to ensure native macOS color inversion support for Light and Dark modes.

## Phase 2: Core Backend Logic (`timer.go`)
*   **Timer Struct & State:** Define the Go service with states (Running, Paused, Stopped) and durations (Focus: 25m, Short Break: 5m, Long Break: 15m).
*   **Ticker Implementation:** Use a `time.Ticker` as the single source of truth for the countdown.
*   **Event Broadcasting:** Integrate Wails Events API to push timer ticks (remaining time, progress percentage, current state) to the frontend every second.
*   **Control Methods:** Implement Go methods for `Play`, `Pause`, `Reset`, and `Skip` that can be invoked from the frontend or via global hotkeys.

## Phase 3: Wails Application & macOS Integration (`main.go`)
*   **Window Configuration:** Configure the Wails app with macOS specific settings to achieve the vibrant, native look:
    *   `TitleBarHidden: true`
    *   `FullSizeContent: true`
    *   `WebviewIsTransparent: true` (crucial for macOS vibrancy).
*   **Menu Bar / Tray Setup:** 
    *   Configure the tray icon using the prepared template asset.
    *   Implement logic to dynamically update the Tray Label with the remaining countdown time (e.g., "25:00").
    *   Set the tray icon click behavior to toggle the visibility of the application window, attaching it to the icon's position.
*   **System Notifications:** Integrate native OS alerts to trigger when a session (Focus or Break) completes.
*   **Global Hotkeys:** Bind a configurable shortcut (default `Cmd+Shift+P`) to toggle the timer's play/pause state globally.

## Phase 4: Frontend Implementation
*   **Minimalist UI:** Develop a simple UI using Svelte (or chosen framework) adhering to Apple's Human Interface Guidelines (system fonts, rounded corners).
*   **Event Listening:** Make the frontend "dumb" by having it solely react to the events pushed by the Go backend to display the countdown and a progress bar/circle.
*   **Vibrancy Application:** Apply CSS properties that allow the macOS vibrancy effect to bleed through the transparent webview.
*   **Controls Integration:** Add buttons for Play, Pause, and Reset that call the corresponding Go backend methods.

## Phase 5: Optimization & Review
*   **Resource Footprint Check:** Profile the application to ensure it meets the constraint of using < 40MB RAM.
*   **Battery Efficiency Check:** Verify that the frontend is not continuously polling and is efficiently updating only on backend events.
*   **UI/UX Polish:** Ensure smooth transitions, correct window positioning near the menu bar, and proper contrast in both Light and Dark modes.
