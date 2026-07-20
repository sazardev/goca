// Generates animated SVG frames for a short (~8s) vertical video per Goca CLI
// command: terminal typing + explainer card. One subfolder per command under
// video/frames-svg/<command>/. Run: node marketing/scripts/generate-command-videos.mjs
import { writeFileSync, mkdirSync, rmSync } from 'fs';
import { fileURLToPath } from 'url';
import path from 'path';
import { FONT_SANS, FONT_MONO, THEMES, bgDefs, brandLockup, iconBadge, escXml } from './brand.mjs';
import { COMMANDS } from './commands-data.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const FRAMES_ROOT = path.resolve(__dirname, '..', 'video', 'frames-svg');

const FPS = 30;
const W = 1080, H = 1920;
const T = THEMES.dark;
const CF = 0.3;

const clamp = (v, lo, hi) => Math.max(lo, Math.min(hi, v));
const lerp = (a, b, p) => a + (b - a) * p;
const easeOutCubic = (p) => 1 - Math.pow(1 - p, 3);
const easeInCubic = (p) => p * p * p;
function fadeP(t, start, dur, ease = easeOutCubic) { return ease(clamp((t - start) / dur, 0, 1)); }

function group(inner, { opacity = 1, dy = 0, dx = 0, scale = 1, cx, cy } = {}) {
  const transforms = [];
  if (dx || dy) transforms.push(`translate(${dx.toFixed(2)},${dy.toFixed(2)})`);
  if (scale !== 1) {
    if (cx !== undefined) transforms.push(`translate(${cx},${cy}) scale(${scale.toFixed(4)}) translate(${-cx},${-cy})`);
    else transforms.push(`scale(${scale.toFixed(4)})`);
  }
  return `<g opacity="${opacity.toFixed(3)}" ${transforms.length ? `transform="${transforms.join(' ')}"` : ''}>${inner}</g>`;
}
function bg(glowCy) {
  return `${bgDefs('50%', glowCy, '55%', T)}<rect width="${W}" height="${H}" fill="url(#bgGrad)"/><rect width="${W}" height="${H}" fill="url(#glow)"/>`;
}

const TERM_FONT = 22;
const TYPE_START = 0.35;
const typeDurFor = (c) => clamp(c.cmd.length * 0.045, 1.1, 2.6);
// Time at which the "Ready." line has fully faded in (see stagger below).
const readyDoneFor = (c) => TYPE_START + typeDurFor(c) + 1.55;

function renderTerminalScene(c, t) {
  const winP = fadeP(t, 0.0, 0.35);
  const typeStart = TYPE_START;
  const typeDur = typeDurFor(c);
  const typed = c.cmd.slice(0, Math.round(c.cmd.length * clamp((t - typeStart) / typeDur, 0, 1)));
  const cursorX = 100 + typed.length * (TERM_FONT * 0.62);
  const showCursor = t < typeStart + typeDur + 0.3 && Math.floor(t / 0.25) % 2 === 0;

  const afterType = typeStart + typeDur + 0.15;
  const tagP = fadeP(t, afterType, 0.3);
  const tagW = Math.round(c.tag.length * (16 * 0.62 + 0.5) + 32);
  const tag = group(`
    <rect x="100" y="${462 - 30}" width="${tagW}" height="38" rx="9" fill="#00ADD8" fill-opacity="0.14" stroke="#00ADD8" stroke-width="1.2"/>
    <text x="${100 + tagW / 2}" y="${462 - 5}" font-family="${FONT_SANS}" font-size="15" font-weight="700" letter-spacing="0.4" fill="#00D9FF" text-anchor="middle">${c.tag}</text>
  `, { opacity: tagP, dy: lerp(10, 0, tagP) });

  const line1P = fadeP(t, afterType + 0.3, 0.35);
  const line2P = fadeP(t, afterType + 0.65, 0.35);
  const comments = group(
    `<text x="100" y="560" font-family="${FONT_MONO}" font-size="23" fill="#8b949e">// ${escXml(c.sub[0])}</text>`,
    { opacity: line1P, dy: lerp(10, 0, line1P) }
  ) + group(
    `<text x="100" y="596" font-family="${FONT_MONO}" font-size="23" fill="#8b949e">// ${escXml(c.sub[1])}</text>`,
    { opacity: line2P, dy: lerp(10, 0, line2P) }
  );

  const readyP = fadeP(t, afterType + 1.05, 0.35);
  const ready = group(
    `<text x="100" y="668" font-family="${FONT_MONO}" font-size="30" font-weight="bold" fill="#3fb950">&#10003; Ready.</text>`,
    { opacity: readyP, scale: lerp(0.92, 1, readyP), cx: 100, cy: 668 }
  );

  const term = group(`
    <rect x="60" y="300" width="960" height="1220" rx="20" fill="#0d1117" stroke="#00ADD8" stroke-opacity="0.35" stroke-width="1.5"/>
    <rect x="60" y="300" width="960" height="56" rx="20" fill="#161b22"/>
    <rect x="60" y="336" width="960" height="20" fill="#161b22"/>
    <circle cx="102" cy="328" r="9" fill="#ff5f56"/>
    <circle cx="132" cy="328" r="9" fill="#ffbd2e"/>
    <circle cx="162" cy="328" r="9" fill="#27c93f"/>
    <text x="540" y="335" font-family="${FONT_MONO}" font-size="19" fill="#8b949e" text-anchor="middle">goca &#8212; zsh</text>
    <text x="100" y="430" font-family="${FONT_MONO}" font-size="${TERM_FONT}" fill="#c9d1d9"><tspan fill="#27c93f">$</tspan> ${escXml(typed.replace(/^\$ /, ''))}</text>
    ${showCursor ? `<rect x="${cursorX.toFixed(1)}" y="404" width="12" height="28" fill="#00D9FF"/>` : ''}
    ${tag}
    ${comments}
    ${ready}
  `, { opacity: winP, dy: lerp(40, 0, winP) });

  return `${bg('14%')}${brandLockup(90, 170, 26, 32, T)}${term}`;
}

function renderCardScene(c, t) {
  const p = fadeP(t, 0.0, 0.4);
  const icon = group(iconBadge(540, 620, 140, 105, c.icon, T), { opacity: p, scale: lerp(0.88, 1, p), cx: 540, cy: 620 });
  const heads = group(
    `<text x="540" y="852" font-family="${FONT_SANS}" font-size="62" font-weight="bold" fill="#ffffff" text-anchor="middle">${escXml(c.headline[0])}</text>
     <text x="540" y="924" font-family="${FONT_SANS}" font-size="62" font-weight="bold" fill="#00D9FF" text-anchor="middle">${escXml(c.headline[1])}</text>`,
    { opacity: p, dy: lerp(16, 0, p) }
  );
  const chipP = fadeP(t, 0.35, 0.4);
  const cmdShort = c.cmd.length > 46 ? c.cmd.slice(0, 44) + '…' : c.cmd;
  const chip = group(`
    <rect x="130" y="1000" width="820" height="70" rx="14" fill="#000000" fill-opacity="0.35" stroke="#00ADD8" stroke-width="1.5" stroke-opacity="0.5"/>
    <text x="540" y="1044" font-family="${FONT_MONO}" font-size="23" fill="#00D9FF" text-anchor="middle">${escXml(cmdShort)}</text>
  `, { opacity: chipP, dy: lerp(12, 0, chipP) });

  const footP = fadeP(t, 0.55, 0.4);
  const foot = group(
    `<text x="540" y="1800" font-family="${FONT_SANS}" font-size="26" fill="#ffffff" fill-opacity="0.5" text-anchor="middle">sazardev.github.io/goca</text>`,
    { opacity: footP }
  );

  return `${bg('30%')}${brandLockup(90, 170, 26, 32, T)}${icon}${heads}${chip}${foot}`;
}

const CARD_DUR = 3.6;
const INTRO_FADE = 0.3, OUTRO_FADE = 0.5;
// Hold the fully-typed "Ready." line on screen for this long before crossfading out.
const READY_HOLD = 0.5;

function blackOverlay(t, total) {
  if (t < INTRO_FADE) return 1 - easeOutCubic(t / INTRO_FADE);
  if (t > total - OUTRO_FADE) return easeInCubic((t - (total - OUTRO_FADE)) / OUTRO_FADE);
  return 0;
}

function composeFrame(c, t, termDur, total) {
  let inner;
  if (t >= termDur - CF) {
    const p = clamp((t - (termDur - CF)) / CF, 0, 1);
    inner = group(renderTerminalScene(c, clamp(t, 0, termDur)), { opacity: 1 - p }) +
            group(renderCardScene(c, clamp(t - termDur, 0, CARD_DUR)), { opacity: p });
  } else {
    inner = renderTerminalScene(c, t);
  }
  const overlay = blackOverlay(t, total);
  return `<svg width="${W}" height="${H}" xmlns="http://www.w3.org/2000/svg">${inner}${overlay > 0 ? `<rect width="${W}" height="${H}" fill="#000000" opacity="${overlay.toFixed(3)}"/>` : ''}</svg>`;
}

let grand = 0;
for (const c of COMMANDS) {
  const termDur = readyDoneFor(c) + READY_HOLD;
  const total = termDur + CARD_DUR;
  const totalFrames = Math.round(total * FPS);

  const dir = path.join(FRAMES_ROOT, c.id);
  rmSync(dir, { recursive: true, force: true });
  mkdirSync(dir, { recursive: true });
  for (let f = 0; f < totalFrames; f++) {
    const svg = composeFrame(c, f / FPS, termDur, total);
    writeFileSync(path.join(dir, `frame_${String(f).padStart(5, '0')}.svg`), svg, 'utf-8');
  }
  grand += totalFrames;
  console.log(`  ${c.id}: ${totalFrames} frames (${total.toFixed(2)}s)`);
}
console.log(`Generated ${grand} frames across ${COMMANDS.length} command videos (~8s each) in ${FRAMES_ROOT}/<command>/`);
