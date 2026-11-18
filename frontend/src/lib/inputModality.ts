export type InputModality = "keyboard" | "pointer" | "controller";
export type ControllerGlyphScheme = "xbox" | "playstation" | "nintendo" | "generic";

const MODALITY_CLASSES: Record<InputModality, string> = {
  keyboard: "input-modality-keyboard",
  pointer: "input-modality-pointer",
  controller: "input-modality-controller",
};

const KEYBOARD_FOCUS_KEYS = new Set([
  "ArrowUp",
  "ArrowDown",
  "ArrowLeft",
  "ArrowRight",
  "Home",
  "End",
  "PageUp",
  "PageDown",
  "Enter",
  " ",
  "Escape",
]);
const SYNTHETIC_CONTROLLER_KEY_EVENT = "__tcnoControllerKeyEvent";

let currentModality: InputModality = "pointer";
let currentControllerGlyphScheme: ControllerGlyphScheme = "xbox";

function rootElement(): HTMLElement | null {
  return document.documentElement ?? document.body ?? null;
}

export function setInputModality(modality: InputModality): void {
  currentModality = modality;
  const root = rootElement();
  if (!root) return;
  // Every pointerdown re-asserts "pointer". Rewriting the same class churns the attribute for
  // nothing, and shows up as a change to anything watching the DOM.
  if (root.dataset.inputModality !== modality) {
    root.classList.remove(...Object.values(MODALITY_CLASSES));
    root.classList.add(MODALITY_CLASSES[modality]);
    root.dataset.inputModality = modality;
  }
  if (!root.dataset.controllerScheme) {
    root.dataset.controllerScheme = currentControllerGlyphScheme;
  }
}

export function getInputModality(): InputModality {
  return currentModality;
}

export function getControllerGlyphScheme(): ControllerGlyphScheme {
  return currentControllerGlyphScheme;
}

export function setControllerGlyphScheme(scheme: ControllerGlyphScheme): void {
  currentControllerGlyphScheme = scheme;
  const root = rootElement();
  if (!root) return;
  root.dataset.controllerScheme = scheme;
}

export function inferControllerGlyphScheme(id: string | null | undefined): ControllerGlyphScheme {
  const normalized = (id ?? "").toLowerCase();
  if (normalized.includes("playstation") || normalized.includes("dualsense") || normalized.includes("dualshock") || normalized.includes("sony")) {
    return "playstation";
  }
  if (normalized.includes("nintendo") || normalized.includes("switch") || normalized.includes("joy-con") || normalized.includes("joycon")) {
    return "nintendo";
  }
  if (normalized.includes("xbox") || normalized.includes("xinput") || normalized.includes("microsoft")) {
    return "xbox";
  }
  return "generic";
}

export function markControllerInput(scheme?: ControllerGlyphScheme): void {
  if (scheme) {
    setControllerGlyphScheme(scheme);
  }
  setInputModality("controller");
}

export function markSyntheticControllerKeyEvent(event: KeyboardEvent): void {
  Object.defineProperty(event, SYNTHETIC_CONTROLLER_KEY_EVENT, {
    value: true,
    configurable: true,
  });
}

export function isSyntheticControllerKeyEvent(event: KeyboardEvent): boolean {
  return (event as KeyboardEvent & Record<string, unknown>)[SYNTHETIC_CONTROLLER_KEY_EVENT] === true;
}

export function installInputModalityTracking(): () => void {
  if (typeof window === "undefined" || typeof document === "undefined") {
    return () => {};
  }

  setInputModality(currentModality);

  // Moving or re-activating the frameless window hands focus back to the
  // WebView, and WebView2 lands it on the first tabbable element: the first
  // platform tile or account card. To a mouse user that reads as the app
  // selecting something on its own. Focus that arrives without a click, key or
  // controller press while in pointer mode is dropped.
  let lastPointerDownAt = 0;
  const STRAY_FOCUS_SELECTOR = "[data-dnd-cell], .acc_list_item, .platform_list_item";
  const onFocusIn = (event: FocusEvent): void => {
    if (currentModality !== "pointer") return;
    if (performance.now() - lastPointerDownAt < 400) return;
    const target = event.target;
    if (target instanceof HTMLElement && target.matches(STRAY_FOCUS_SELECTOR)) {
      target.blur();
    }
  };

  const onPointerInput = (): void => {
    lastPointerDownAt = performance.now();
    setInputModality("pointer");
  };
  const onKeydown = (event: KeyboardEvent): void => {
    // BRQ: the Tab key does nothing. Tabbing walked a highlight across tiles
    // and buttons, which read as the app selecting things by itself.
    if (event.key === "Tab" && !event.ctrlKey && !event.altKey && !event.metaKey) {
      event.preventDefault();
      event.stopImmediatePropagation();
      return;
    }
    if (event.metaKey || event.ctrlKey || event.altKey) return;
    if (isSyntheticControllerKeyEvent(event)) return;
    if (KEYBOARD_FOCUS_KEYS.has(event.key)) {
      setInputModality("keyboard");
    }
  };

  window.addEventListener("keydown", onKeydown, true);
  window.addEventListener("mousedown", onPointerInput, true);
  window.addEventListener("pointerdown", onPointerInput, true);
  window.addEventListener("touchstart", onPointerInput, true);
  window.addEventListener("focusin", onFocusIn, true);

  return () => {
    window.removeEventListener("focusin", onFocusIn, true);
    window.removeEventListener("keydown", onKeydown, true);
    window.removeEventListener("mousedown", onPointerInput, true);
    window.removeEventListener("pointerdown", onPointerInput, true);
    window.removeEventListener("touchstart", onPointerInput, true);
  };
}
