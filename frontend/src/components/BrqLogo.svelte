<script lang="ts">
  import {
    BRQ_BOLT_PATH,
    BRQ_LETTER_PATHS,
    BRQ_MARK_TRANSFORM,
    BRQ_RING_PATH,
  } from "../lib/brqLogoPaths";

  /** "wordmark" = hexagon mark + BRQ letters; "mark" = hexagon mark only. */
  export let variant: "wordmark" | "mark" = "wordmark";
  /** Plays the periodic shine sweep and bolt flicker. */
  export let animated = true;
  let className = "";
  export { className as class };

  // Gradient ids must be unique per instance or the first <defs> wins for all.
  const uid = `brq${Math.random().toString(36).slice(2, 9)}`;
  $: viewBox = variant === "mark" ? "0 0 100 100" : "0 0 768 264";
  $: markTransform = variant === "mark" ? "" : BRQ_MARK_TRANSFORM;
  $: shineWidth = variant === "mark" ? 40 : 220;
  $: shineHeight = variant === "mark" ? 100 : 264;
</script>

<svg
  class="brq-logo {className}"
  class:brq-logo--animated={animated}
  style="--brq-shine-dist: {variant === 'mark' ? 190 : 1100}px"
  {viewBox}
  xmlns="http://www.w3.org/2000/svg"
  aria-hidden="true"
>
  <defs>
    <linearGradient id="{uid}-metal" x1="0" y1="0" x2="0.35" y2="1">
      <stop offset="0" stop-color="#f6eefe" />
      <stop offset="0.38" stop-color="#d5bdf0" />
      <stop offset="0.62" stop-color="#9a78c4" />
      <stop offset="1" stop-color="#e2d0f6" />
    </linearGradient>
    <linearGradient id="{uid}-bolt" x1="0.2" y1="0" x2="0.8" y2="1">
      <stop offset="0" stop-color="#ffffff" />
      <stop offset="0.45" stop-color="#e3d2f8" />
      <stop offset="1" stop-color="#a888d4" />
    </linearGradient>
    <linearGradient id="{uid}-shine" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0" stop-color="#fff" stop-opacity="0" />
      <stop offset="0.5" stop-color="#fff" stop-opacity="0.85" />
      <stop offset="1" stop-color="#fff" stop-opacity="0" />
    </linearGradient>
    <clipPath id="{uid}-clip">
      <g transform={markTransform}>
        <path d={BRQ_RING_PATH} />
        <path d={BRQ_BOLT_PATH} />
      </g>
      {#if variant === "wordmark"}
        {#each BRQ_LETTER_PATHS as d}
          <path {d} clip-rule="evenodd" fill-rule="evenodd" />
        {/each}
      {/if}
    </clipPath>
  </defs>

  <g transform={markTransform}>
    <path class="brq-logo__ring" d={BRQ_RING_PATH} fill="url(#{uid}-metal)" />
    <path class="brq-logo__bolt" d={BRQ_BOLT_PATH} fill="url(#{uid}-bolt)" />
  </g>
  {#if variant === "wordmark"}
    {#each BRQ_LETTER_PATHS as d}
      <path {d} fill-rule="evenodd" fill="url(#{uid}-metal)" />
    {/each}
  {/if}

  {#if animated}
    <g clip-path="url(#{uid}-clip)">
      <rect
        class="brq-logo__shine"
        x={-shineWidth}
        y="0"
        width={shineWidth}
        height={shineHeight}
        fill="url(#{uid}-shine)"
        transform="skewX(-20)"
      />
    </g>
  {/if}
</svg>

<style>
  .brq-logo {
    display: block;
    overflow: visible;
    filter: drop-shadow(0 0 4px rgba(185, 150, 240, 0.35));
  }
  .brq-logo--animated .brq-logo__shine {
    animation: brq-logo-shine 5.5s cubic-bezier(0.45, 0, 0.2, 1) infinite;
  }
  .brq-logo--animated .brq-logo__bolt {
    animation: brq-logo-bolt 5.5s linear infinite;
  }
  /* The sweep crosses once, then rests; the bolt flickers as the light reaches it. */
  @keyframes brq-logo-shine {
    0% { transform: skewX(-20deg) translateX(0); }
    35%, 100% { transform: skewX(-20deg) translateX(var(--brq-shine-dist, 1100px)); }
  }
  @keyframes brq-logo-bolt {
    0%, 5%, 13%, 100% { opacity: 1; }
    7% { opacity: 0.55; }
    9% { opacity: 1; }
    11% { opacity: 0.7; }
  }
</style>
