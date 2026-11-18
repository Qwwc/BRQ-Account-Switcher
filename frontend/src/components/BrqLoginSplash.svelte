<script lang="ts">
  import { onDestroy } from "svelte";
  import { fade } from "svelte/transition";
  import { loginSplash, hideLoginSplash, type LoginSplash } from "../stores/loginSplash";
  import { animationsEnabled } from "../stores/animationSettings";
  import { BRQ_BOLT_PATH, BRQ_RING_PATH } from "../lib/brqLogoPaths";

  /** How long the splash stays before it leaves on its own. */
  const SHOW_MS = 2100;
  const PARTICLES = 14;

  let current: LoginSplash | null = null;
  let timer: ReturnType<typeof setTimeout> | null = null;

  const unsub = loginSplash.subscribe((value) => {
    current = value;
    if (timer) clearTimeout(timer);
    timer = value ? setTimeout(hideLoginSplash, SHOW_MS) : null;
  });

  onDestroy(() => {
    unsub();
    if (timer) clearTimeout(timer);
  });

  function onKeydown(e: KeyboardEvent): void {
    if (current && (e.key === "Escape" || e.key === "Enter" || e.key === " ")) {
      e.preventDefault();
      hideLoginSplash();
    }
  }

  $: motion = $animationsEnabled !== false;
</script>

<svelte:window on:keydown={onKeydown} />

{#if current}
  {#key current.nonce}
    <!-- A click anywhere dismisses it early; Escape/Enter do the same via svelte:window. -->
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions a11y_no_noninteractive_element_interactions -->
    <div
      class="brq-splash"
      class:brq-splash--still={!motion}
      role="status"
      aria-live="polite"
      on:click={hideLoginSplash}
      in:fade={{ duration: motion ? 160 : 0 }}
      out:fade={{ duration: motion ? 320 : 0 }}
    >
      <div class="brq-splash__flash" aria-hidden="true"></div>

      <div class="brq-splash__stage">
        <!-- Shock rings that leave the emblem when it lands. -->
        <span class="brq-splash__shock" aria-hidden="true"></span>
        <span class="brq-splash__shock brq-splash__shock--2" aria-hidden="true"></span>
        <span class="brq-splash__shock brq-splash__shock--3" aria-hidden="true"></span>

        <!-- Sparks thrown outward. -->
        {#each Array(PARTICLES) as _, i}
          <span
            class="brq-splash__spark"
            style="--a: {(360 / PARTICLES) * i}deg; --d: {90 + (i % 3) * 28}px; --delay: {120 + (i % 4) * 30}ms"
            aria-hidden="true"
          ></span>
        {/each}

        <!-- Circling energy behind the hexagon. -->
        <span class="brq-splash__orbit" aria-hidden="true"></span>

        <div class="brq-splash__emblem">
          <svg class="brq-splash__ring" viewBox="0 0 100 100" aria-hidden="true">
            <defs>
              <linearGradient id="brq-splash-metal" x1="0" y1="0" x2="0.35" y2="1">
                <stop offset="0" stop-color="#f6eefe" />
                <stop offset="0.45" stop-color="#cdb2ee" />
                <stop offset="0.7" stop-color="#9a78c4" />
                <stop offset="1" stop-color="#e2d0f6" />
              </linearGradient>
            </defs>
            <path d={BRQ_RING_PATH} fill="url(#brq-splash-metal)" />
          </svg>

          <div class="brq-splash__avatar">
            {#if current.avatarSrc}
              <img src={current.avatarSrc} alt="" draggable="false" />
            {:else}
              <svg viewBox="0 0 100 100" aria-hidden="true"><path d={BRQ_BOLT_PATH} /></svg>
            {/if}
          </div>

          <span class="brq-splash__badge" aria-hidden="true">
            <svg viewBox="0 0 24 24"><path d="M9.5 16.6 4.9 12l-1.4 1.4 6 6 12-12-1.4-1.4z" /></svg>
          </span>
        </div>
      </div>

      <div class="brq-splash__text">
        <p class="brq-splash__title">SIGNED IN</p>
        <p class="brq-splash__name">{current.accountName}</p>
        <p class="brq-splash__platform">{current.platformName}</p>
      </div>
    </div>
  {/key}
{/if}

<style>
  .brq-splash {
    position: fixed;
    inset: 0;
    z-index: 100000;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1.4rem;
    cursor: pointer;
    overflow: hidden;
    background:
      radial-gradient(circle at 50% 44%, rgba(var(--accentInt, 185, 150, 240), 0.28), transparent 45%),
      rgba(9, 4, 16, 0.78);
    backdrop-filter: blur(10px) saturate(1.2);
    -webkit-backdrop-filter: blur(10px) saturate(1.2);
    user-select: none;
  }

  /* One bright frame when the switch lands. */
  .brq-splash__flash {
    position: absolute;
    inset: 0;
    pointer-events: none;
    background: radial-gradient(circle at 50% 44%, rgba(255, 255, 255, 0.55), transparent 55%);
    opacity: 0;
    animation: brq-splash-flash 650ms ease-out 80ms;
  }

  .brq-splash__stage {
    position: relative;
    width: 190px;
    height: 190px;
    display: grid;
    place-items: center;
  }

  .brq-splash__emblem {
    position: relative;
    width: 170px;
    height: 170px;
    display: grid;
    place-items: center;
    animation: brq-splash-land 700ms cubic-bezier(0.2, 1.4, 0.35, 1) both;
    filter: drop-shadow(0 0 18px rgba(var(--accentInt, 185, 150, 240), 0.7));
  }

  .brq-splash__ring {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    animation: brq-splash-ring-in 900ms cubic-bezier(0.2, 0.9, 0.3, 1) both;
  }

  /* Avatar clipped to the logo's hexagon. */
  .brq-splash__avatar {
    width: 66%;
    height: 66%;
    clip-path: polygon(50% 0, 100% 25%, 100% 75%, 50% 100%, 0 75%, 0 25%);
    background: linear-gradient(160deg, #2a1846, #120a1e);
    display: grid;
    place-items: center;
    animation: brq-splash-avatar 800ms cubic-bezier(0.2, 0.9, 0.3, 1) 120ms both;
  }
  .brq-splash__avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .brq-splash__avatar svg {
    width: 70%;
    height: 70%;
    fill: #e9dcf7;
  }

  .brq-splash__badge {
    position: absolute;
    right: 14%;
    bottom: 10%;
    width: 38px;
    height: 38px;
    border-radius: 50%;
    display: grid;
    place-items: center;
    background: linear-gradient(135deg, #8dffa0, #2fcf6a);
    box-shadow: 0 0 0 3px #140b20, 0 0 18px rgba(80, 255, 140, 0.7);
    animation: brq-splash-badge 520ms cubic-bezier(0.2, 1.6, 0.4, 1) 520ms both;
  }
  .brq-splash__badge svg {
    width: 24px;
    height: 24px;
    fill: #0b2a14;
  }

  .brq-splash__orbit {
    position: absolute;
    inset: -14px;
    border-radius: 50%;
    padding: 3px;
    background: conic-gradient(from 0deg,
        transparent 0deg,
        rgba(var(--accentInt, 185, 150, 240), 0.9) 60deg,
        #fff 80deg,
        transparent 120deg,
        transparent 180deg,
        rgba(var(--accentInt, 185, 150, 240), 0.9) 240deg,
        #fff 260deg,
        transparent 300deg);
    -webkit-mask: linear-gradient(#000 0 0) content-box, linear-gradient(#000 0 0);
    mask: linear-gradient(#000 0 0) content-box, linear-gradient(#000 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    opacity: 0.85;
    animation: brq-splash-spin 1.1s linear infinite, brq-splash-fade-in 400ms ease-out 200ms both;
  }

  .brq-splash__shock {
    position: absolute;
    inset: 10px;
    border-radius: 50%;
    border: 2px solid rgba(var(--accentInt, 185, 150, 240), 0.9);
    opacity: 0;
    animation: brq-splash-shock 900ms cubic-bezier(0.1, 0.7, 0.3, 1) 180ms;
  }
  .brq-splash__shock--2 { animation-delay: 300ms; border-color: rgba(255, 255, 255, 0.7); }
  .brq-splash__shock--3 { animation-delay: 420ms; }

  .brq-splash__spark {
    position: absolute;
    left: 50%;
    top: 50%;
    width: 6px;
    height: 6px;
    margin: -3px;
    border-radius: 50%;
    background: #fff;
    box-shadow: 0 0 10px rgba(var(--accentInt, 185, 150, 240), 1);
    opacity: 0;
    animation: brq-splash-spark 760ms cubic-bezier(0.1, 0.8, 0.3, 1) var(--delay);
  }

  .brq-splash__text {
    text-align: center;
    line-height: 1.15;
  }
  .brq-splash__text p {
    margin: 0;
  }
  .brq-splash__title {
    font-size: 0.95rem;
    font-weight: 800;
    letter-spacing: 0.5em;
    padding-left: 0.5em;
    color: #8dffa0;
    text-shadow: 0 0 12px rgba(80, 255, 140, 0.6);
    animation: brq-splash-text 520ms cubic-bezier(0.2, 0.9, 0.3, 1) 380ms both;
  }
  .brq-splash__name {
    margin-top: 0.4rem !important;
    font-size: 2rem;
    font-weight: 900;
    letter-spacing: 0.02em;
    max-width: 80vw;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    background: linear-gradient(180deg, #ffffff 20%, var(--accent, #b996f0));
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    filter: drop-shadow(0 3px 10px rgba(0, 0, 0, 0.6));
    animation: brq-splash-text 560ms cubic-bezier(0.2, 0.9, 0.3, 1) 460ms both;
  }
  .brq-splash__platform {
    margin-top: 0.35rem !important;
    font-size: 0.9rem;
    font-weight: 600;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: rgba(255, 255, 255, 0.6);
    animation: brq-splash-text 560ms cubic-bezier(0.2, 0.9, 0.3, 1) 540ms both;
  }

  /* With animations off, the splash simply appears: same content, no motion. */
  .brq-splash--still *,
  .brq-splash--still {
    animation: none !important;
  }
  .brq-splash--still .brq-splash__spark,
  .brq-splash--still .brq-splash__shock,
  .brq-splash--still .brq-splash__flash {
    display: none;
  }

  @keyframes brq-splash-flash {
    0% { opacity: 0; }
    20% { opacity: 1; }
    100% { opacity: 0; }
  }
  @keyframes brq-splash-land {
    0% { transform: scale(0.2) rotate(-25deg); opacity: 0; }
    60% { opacity: 1; }
    100% { transform: none; opacity: 1; }
  }
  @keyframes brq-splash-ring-in {
    from { transform: rotate(-120deg) scale(0.6); }
    to { transform: none; }
  }
  @keyframes brq-splash-avatar {
    from { transform: scale(0.4); filter: brightness(2.5) blur(4px); }
    to { transform: none; filter: none; }
  }
  @keyframes brq-splash-badge {
    from { transform: scale(0) rotate(-90deg); }
    to { transform: none; }
  }
  @keyframes brq-splash-spin {
    to { transform: rotate(360deg); }
  }
  @keyframes brq-splash-fade-in {
    from { opacity: 0; }
  }
  @keyframes brq-splash-shock {
    0% { opacity: 0.9; transform: scale(0.6); }
    100% { opacity: 0; transform: scale(2.4); }
  }
  @keyframes brq-splash-spark {
    0% { opacity: 1; transform: rotate(var(--a)) translateX(30px) scale(1); }
    100% { opacity: 0; transform: rotate(var(--a)) translateX(var(--d)) scale(0.2); }
  }
  @keyframes brq-splash-text {
    from { opacity: 0; transform: translateY(14px); letter-spacing: 0.9em; }
    to { opacity: 1; transform: none; }
  }
</style>
