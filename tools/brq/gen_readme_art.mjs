// Generates the README banner and buttons in .github/brq/.
//   node tools/brq/gen_readme_art.mjs
// SVGs shown through <img> can't load other files, so the badge icon is
// embedded as base64 and text uses common system fonts.
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "../..");
const outDir = resolve(root, ".github/brq");
mkdirSync(outDir, { recursive: true });

const icon = readFileSync(resolve(root, "build/trayicon.png")).toString("base64");
const FONT = `'Segoe UI Black','Segoe UI','Arial Black',Arial,sans-serif`;

// ---------------------------------------------------------------- banner
const banner = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="900" height="220" viewBox="0 0 900 220">
  <defs>
    <linearGradient id="bg" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#2a1546"/>
      <stop offset="0.55" stop-color="#160c26"/>
      <stop offset="1" stop-color="#0d0717"/>
    </linearGradient>
    <radialGradient id="glow" cx="0.2" cy="0.5" r="0.45">
      <stop offset="0" stop-color="#9b6ddc" stop-opacity="0.55"/>
      <stop offset="1" stop-color="#9b6ddc" stop-opacity="0"/>
    </radialGradient>
    <linearGradient id="silk" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0" stop-color="#6b3fb0" stop-opacity="0"/>
      <stop offset="0.5" stop-color="#6b3fb0" stop-opacity="0.35"/>
      <stop offset="1" stop-color="#6b3fb0" stop-opacity="0"/>
    </linearGradient>
    <linearGradient id="metal" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="#ffffff"/>
      <stop offset="0.5" stop-color="#e2cff8"/>
      <stop offset="1" stop-color="#a985d8"/>
    </linearGradient>
    <linearGradient id="shine" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0" stop-color="#fff" stop-opacity="0"/>
      <stop offset="0.5" stop-color="#fff" stop-opacity="0.35"/>
      <stop offset="1" stop-color="#fff" stop-opacity="0"/>
    </linearGradient>
    <clipPath id="frame"><path d="M22 0H900V198L878 220H0V22Z"/></clipPath>
  </defs>
  <g clip-path="url(#frame)">
    <rect width="900" height="220" fill="url(#bg)"/>
    <rect width="900" height="220" fill="url(#glow)"/>
    <path d="M-50 170 C 200 90, 420 250, 950 120 L950 160 C 430 290, 200 130, -50 210Z" fill="url(#silk)"/>
    <path d="M-50 40 C 260 120, 520 -20, 950 60 L950 80 C 520 0, 260 140, -50 64Z" fill="url(#silk)" opacity="0.6"/>
    <rect x="-300" y="0" width="160" height="220" fill="url(#shine)" transform="skewX(-20)">
      <animate attributeName="x" values="-300;1200;1200" keyTimes="0;0.45;1" dur="6s" repeatCount="indefinite"/>
    </rect>
  </g>
  <path d="M22 1H899V197.5L877.5 219H1V22.5Z" fill="none" stroke="#b996f0" stroke-opacity="0.55" stroke-width="2"/>
  <image x="48" y="30" width="160" height="160" xlink:href="data:image/png;base64,${icon}" href="data:image/png;base64,${icon}"/>
  <text x="246" y="112" font-family="${FONT}" font-size="64" font-weight="900" fill="url(#metal)" letter-spacing="4">BRQ</text>
  <text x="248" y="156" font-family="${FONT}" font-size="30" font-weight="800" fill="#e9dcf7" letter-spacing="3">ACCOUNT SWITCHER</text>
  <text x="249" y="186" font-family="'Segoe UI',Arial,sans-serif" font-size="15" fill="#b6a8cb" letter-spacing="1">Switch accounts in one click · Steam · Epic · Riot · Ubisoft · Discord · Battle.net</text>
</svg>
`;
writeFileSync(resolve(outDir, "banner.svg"), banner);

// ---------------------------------------------------------------- buttons
// Each button: chamfered frame, coloured glow edge, icon disc, label.
const ICONS = {
  download: `<path d="M0 -14 V6 M-9 -3 L0 6 L9 -3" fill="none" stroke="#fff" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round"/><path d="M-12 11 H12" stroke="#fff" stroke-width="3.5" stroke-linecap="round"/>`,
  info: `<circle r="2.6" cy="-9" fill="#fff"/><path d="M0 -2 V11" stroke="#fff" stroke-width="4" stroke-linecap="round"/>`,
  bug: `<ellipse rx="7" ry="9" cy="2" fill="none" stroke="#fff" stroke-width="3"/><path d="M0 -7 V11 M-7 0 H-13 M7 0 H13 M-6 -6 L-11 -11 M6 -6 L11 -11 M-6 8 L-11 12 M6 8 L11 12" stroke="#fff" stroke-width="2.6" stroke-linecap="round"/>`,
  code: `<path d="M-6 -9 L-14 0 L-6 9 M6 -9 L14 0 L6 9" fill="none" stroke="#fff" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round"/>`,
};

function button(file, label, sub, color, iconKey, primary = false) {
  const w = 400, h = 84, c = 14;
  const shape = `M${c} 0H${w}V${h - c}L${w - c} ${h}H0V${c}Z`;
  const shapeIn = `M${c} 1H${w - 1}V${h - c - 0.4}L${w - c - 0.4} ${h - 1}H1V${c + 0.4}Z`;
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${w}" height="${h}" viewBox="0 0 ${w} ${h}">
  <defs>
    <linearGradient id="fill" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="${primary ? "#3b1f66" : "#241536"}"/>
      <stop offset="1" stop-color="${primary ? "#1b0e2e" : "#120a1d"}"/>
    </linearGradient>
    <radialGradient id="halo" cx="0.14" cy="0.5" r="0.5">
      <stop offset="0" stop-color="${color}" stop-opacity="0.45"/>
      <stop offset="1" stop-color="${color}" stop-opacity="0"/>
    </radialGradient>
    <linearGradient id="edge" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0" stop-color="${color}"/>
      <stop offset="0.5" stop-color="#ffffff" stop-opacity="0.9"/>
      <stop offset="1" stop-color="${color}"/>
      ${primary ? `<animateTransform attributeName="gradientTransform" type="translate" values="-1 0;1 0" dur="3s" repeatCount="indefinite"/>` : ""}
    </linearGradient>
    <linearGradient id="sweep" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0" stop-color="#fff" stop-opacity="0"/>
      <stop offset="0.5" stop-color="#fff" stop-opacity="0.22"/>
      <stop offset="1" stop-color="#fff" stop-opacity="0"/>
    </linearGradient>
    <clipPath id="clip"><path d="${shape}"/></clipPath>
  </defs>
  <g clip-path="url(#clip)">
    <rect width="${w}" height="${h}" fill="url(#fill)"/>
    <rect width="${w}" height="${h}" fill="url(#halo)"/>
    ${primary ? `<rect x="-200" y="0" width="90" height="${h}" fill="url(#sweep)" transform="skewX(-20)"><animate attributeName="x" values="-200;520;520" keyTimes="0;0.5;1" dur="3.5s" repeatCount="indefinite"/></rect>` : ""}
  </g>
  <path d="${shapeIn}" fill="none" stroke="url(#edge)" stroke-width="2"/>
  <g transform="translate(52 42)">
    <circle r="27" fill="${color}" fill-opacity="0.18" stroke="${color}" stroke-width="2"/>
    ${ICONS[iconKey]}
  </g>
  <text x="98" y="${sub ? 44 : 53}" font-family="${FONT}" font-size="${primary ? 30 : 26}" font-weight="900" fill="#ffffff" letter-spacing="2">${label}</text>
  ${sub ? `<text x="99" y="66" font-family="'Segoe UI',Arial,sans-serif" font-size="14" fill="#cbbbe0" letter-spacing="0.5">${sub}</text>` : ""}
</svg>
`;
  writeFileSync(resolve(outDir, file), svg);
}

button("btn-download.svg", "DOWNLOAD", "Latest Windows release", "#b996f0", "download", true);
button("btn-howto.svg", "HOW IT WORKS", "What BRQ does with your logins", "#6fb7ff", "info");
button("btn-bug.svg", "REPORT A BUG", "Open an issue on GitHub", "#ff7a9c", "bug");
button("btn-build.svg", "SOURCE CODE", "Build it yourself · Go · Wails", "#7dffb0", "code");

console.log("written to", outDir);
