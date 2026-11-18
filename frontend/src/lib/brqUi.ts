// BRQ interaction layer: theme marker class and click ripples.
import { currentThemeId } from "./theme/stores";
import { DEFAULT_THEME_ID } from "./theme/types";

/** Set on <html> while the BRQ (default) theme is active; brq-ui.scss keys off it. */
export const BRQ_THEME_CLASS = "brq-theme";

const RIPPLE_SELECTOR = [
  ".actionbar button",
  ".brq-btn",
  ".button",
  ".settings-actions button",
  ".settings-card__row > button",
  ".modal-inline-actions button",
  ".dropdown-toggle",
  ".platform_list_item",
].join(",");

function motionAllowed(): boolean {
  const root = document.documentElement;
  if (root.classList.contains("animations-disabled")) return false;
  if (document.querySelector(".container.animations-disabled")) return false;
  return !window.matchMedia?.("(prefers-reduced-motion: reduce)").matches;
}

function spawnRipple(event: PointerEvent): void {
  if (event.button !== 0 || !motionAllowed()) return;
  if (!document.documentElement.classList.contains(BRQ_THEME_CLASS)) return;
  const target = event.target instanceof Element ? event.target.closest<HTMLElement>(RIPPLE_SELECTOR) : null;
  if (!target || (target as HTMLButtonElement).disabled) return;

  const rect = target.getBoundingClientRect();
  const size = Math.hypot(rect.width, rect.height) * 2;
  const ripple = document.createElement("span");
  ripple.className = "brq-ripple";
  ripple.setAttribute("aria-hidden", "true");
  ripple.style.width = ripple.style.height = `${size}px`;
  ripple.style.left = `${event.clientX - rect.left}px`;
  ripple.style.top = `${event.clientY - rect.top}px`;
  if (getComputedStyle(target).position === "static") target.style.position = "relative";
  target.appendChild(ripple);
  ripple.addEventListener("animationend", () => ripple.remove(), { once: true });
  // Safety net if the animation never runs (element hidden mid-click).
  window.setTimeout(() => ripple.remove(), 1200);
}

export function installBrqUi(): () => void {
  const root = document.documentElement;
  const offTheme = currentThemeId.subscribe((id) => {
    root.classList.toggle(BRQ_THEME_CLASS, id === DEFAULT_THEME_ID);
  });
  window.addEventListener("pointerdown", spawnRipple, true);
  return () => {
    offTheme();
    window.removeEventListener("pointerdown", spawnRipple, true);
  };
}
