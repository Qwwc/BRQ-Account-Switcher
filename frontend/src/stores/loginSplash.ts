import { writable } from "svelte/store";

/** What the BRQ sign-in splash shows after a successful account switch. */
export type LoginSplash = {
  /** Display name of the account that was signed in. */
  accountName: string;
  /** Platform the account belongs to, e.g. "Steam". */
  platformName: string;
  /** Avatar image URL, or "" to show the BRQ bolt instead. */
  avatarSrc: string;
  /** Changes on every show so the same account can replay the splash. */
  nonce: number;
};

export const loginSplash = writable<LoginSplash | null>(null);

let nonce = 0;

export function showLoginSplash(splash: Omit<LoginSplash, "nonce">): void {
  nonce += 1;
  loginSplash.set({ ...splash, nonce });
}

export function hideLoginSplash(): void {
  loginSplash.set(null);
}
