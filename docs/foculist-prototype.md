# PROMPT: macOS-First Go (Wails v3) Pomodoro Prototype

### CONTEXT
I am building a high-performance, minimalist desktop Pomodoro application named **Foculist**. The primary goal is to create a "Mac-first" experience that feels like a native utility, while maintaining the ability to support Windows in the future. I want to use **Go** for the backend logic and **Wails v3** for the desktop bridge and frontend.

### FEATURES
*   **Menu Bar Attachment:** The app should primarily live in the macOS Menu Bar. Clicking the tray icon should toggle a small, sleek window attached to the icon's position.
*   **Live Tray Label:** The Menu Bar icon must dynamically display the remaining countdown time (e.g., "25:00") directly in the bar.
*   **Vibrant UI:** The window must utilize macOS **Vibrancy** (frosted glass/blur effect) to match the system aesthetic.
*   **Timer Controls:** Play, Pause, and Reset functionality for 25-minute focus and 5/15-minute break sessions.
*   **System Notifications:** Native OS alerts when a session completes.
*   **Global Hotkeys:** A configurable shortcut (e.g., `Cmd+Shift+P`) to toggle the timer status from any app.

### TECHNICAL REQUIREMENTS
*   **Language & Framework:** Go with **Wails v3**.
*   **State Management:** Use a Go service with a `time.Ticker` as the single source of truth.
*   **Communication:** Use the Wails **Events API** to "push" timer updates from Go to the frontend every second.
*   **Window Configuration:** 
    *   `TitleBarHidden: true`
    *   `FullSizeContent: true`
    *   `WebviewIsTransparent: true` (to allow for vibrancy effects).
*   **Tray Handling:** Implement a "Template Icon" strategy for the macOS menu bar to ensure the icon automatically adapts to Light and Dark modes.

### CONSTRAINTS
*   **Minimal Footprint:** Target < 40MB RAM usage.
*   **Clean Architecture:** Separate the timer domain logic from the Wails-specific application setup.
*   **UI Guidelines:** The frontend should follow Apple’s Human Interface Guidelines (system fonts, rounded corners, subtle transitions).
*   **Efficiency:** The frontend should be "dumb"; it should only react to events pushed by Go to preserve battery life.

### PROJECT STRUCTURE
Please provide the following:
1.  **`main.go`**: Application entry point with specific `MacOptions` and `Tray` setup for macOS.
2.  **`timer.go`**: The backend Go service handling the ticker logic and event broadcasting.
3.  **Frontend Implementation**: A minimalist component (Svelte, React, or Vue) that displays the countdown and progress bar.
4.  **Asset Guidance**: Instructions on naming and placing tray icons (e.g., `TrayIconTemplate.png`) to ensure native macOS color inversion support.
