<script>
  import { Events } from "@wailsio/runtime";
  import { TimerService } from "../bindings/changeme/backend/index.js";

  let formattedTime = "25:00";
  let state = "stopped";
  let sessionType = "focus";
  let progressPercentage = 100;
  let totalSeconds = 1500;

  Events.On('timer_tick', (e) => {
    console.log("JS timer_tick received:", e);
    const tick = e.data[0] || e.data; // fallback for non-array?
    if (tick && typeof tick === 'object') {
      formattedTime = tick.formatted_time || tick.FormattedTime;
      state = tick.state || tick.State;
      sessionType = tick.session_type || tick.SessionType;
      totalSeconds = tick.total_seconds || tick.TotalSeconds;
      const remainingSeconds = tick.remaining_seconds || tick.RemainingSeconds;
      if (totalSeconds > 0) {
        progressPercentage = (remainingSeconds / totalSeconds) * 100;
      } else {
        progressPercentage = 0;
      }
    }
  });

  function play() { console.log("play"); TimerService.Play().catch(e => console.error("play err:", e)); }
  function pause() { console.log("pause"); TimerService.Pause().catch(e => console.error("pause err:", e)); }
  function reset() { console.log("reset"); TimerService.Reset().catch(e => console.error("reset err:", e)); }
  function skip() { console.log("skip"); TimerService.Skip().catch(e => console.error("skip err:", e)); }
</script>

<div class="app-container">
  <div class="header">
    <h2 class="drag-handle">Foculist</h2>
    <span class="session-badge">{sessionType.replace('_', ' ')}</span>
  </div>

  <div class="timer-display">
    <div class="progress-ring">
      <svg viewBox="0 0 120 120">
        <circle class="ring-bg" cx="60" cy="60" r="54"></circle>
        <circle class="ring-progress" cx="60" cy="60" r="54"
                style="stroke-dashoffset: {339.292 * (1 - progressPercentage / 100)};"></circle>
      </svg>
      <div class="time-text">{formattedTime}</div>
    </div>
  </div>

  <div class="controls">
    {#if state === "running"}
      <button on:click={pause} class="btn icon-btn" title="Pause">⏸</button>
    {:else}
      <button on:click={play} class="btn icon-btn play-btn" title="Play">▶</button>
    {/if}
    <button on:click={reset} class="btn icon-btn" title="Reset">⟲</button>
    <button on:click={skip} class="btn icon-btn" title="Skip">⏭</button>
  </div>
</div>

<style>
  :global(html, body, #app) {
    margin: 0;
    padding: 0;
    height: 100%;
    width: 100%;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    color: white;
    background: transparent !important;
    user-select: none;
    -webkit-user-select: none;
  }

  .app-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
    height: 100%;
    padding: 20px;
    box-sizing: border-box;
    background: rgba(30, 30, 30, 0.3); /* Subtle dark overlay, maintains vibrancy */
    border-radius: 12px;
  }

  .header {
    display: flex;
    justify-content: space-between;
    width: 100%;
    align-items: center;
  }

  .drag-handle {
    margin: 0;
    font-size: 1.2rem;
    font-weight: 500;
    opacity: 0.9;
    --wails-draggable: drag; /* Wails 3 native drag area */
    cursor: default;
  }

  .session-badge {
    background: rgba(255, 255, 255, 0.15);
    padding: 4px 10px;
    border-radius: 12px;
    font-size: 0.8rem;
    text-transform: capitalize;
  }

  .timer-display {
    position: relative;
    width: 220px;
    height: 220px;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .progress-ring {
    position: absolute;
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .progress-ring svg {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    transform: rotate(-90deg);
  }

  .ring-bg {
    fill: none;
    stroke: rgba(255, 255, 255, 0.1);
    stroke-width: 6;
  }

  .ring-progress {
    fill: none;
    stroke: #ff9f0a; /* iOS Orange */
    stroke-width: 6;
    stroke-linecap: round;
    transition: stroke-dashoffset 1s linear;
    stroke-dasharray: 339.292;
  }

  .time-text {
    font-size: 4rem;
    font-weight: 200;
    font-variant-numeric: tabular-nums;
    letter-spacing: -2px;
    z-index: 10;
  }

  .controls {
    display: flex;
    gap: 20px;
    margin-bottom: 20px;
  }

  .btn {
    background: rgba(255, 255, 255, 0.1);
    border: none;
    color: white;
    font-size: 1.2rem;
    width: 50px;
    height: 50px;
    border-radius: 25px;
    cursor: pointer;
    display: flex;
    justify-content: center;
    align-items: center;
    transition: background 0.2s, transform 0.1s;
  }

  .btn:hover {
    background: rgba(255, 255, 255, 0.2);
  }

  .btn:active {
    transform: scale(0.95);
  }

  .play-btn {
    background: rgba(255, 159, 10, 0.8);
  }

  .play-btn:hover {
    background: rgba(255, 159, 10, 1);
  }
</style>
