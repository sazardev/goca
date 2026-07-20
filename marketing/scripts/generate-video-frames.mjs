// Generates one PNG-ready SVG frame per video frame for the "Goca in 23
// seconds" TikTok/Reels/Shorts explainer (1080x1920, vertical).
// Run: node marketing/scripts/generate-video-frames.mjs
// Then: render-video.sh converts frames to PNG and encodes the MP4.
import { writeFileSync, mkdirSync, readdirSync, unlinkSync } from 'fs';
import { fileURLToPath } from 'url';
import path from 'path';
import { FONT_SANS, FONT_MONO, THEMES, bgDefs, brandLockup, iconBadge, escXml } from './brand.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const OUT_DIR = path.resolve(__dirname, '..', 'video', 'frames-svg');
mkdirSync(OUT_DIR, { recursive: true });
for (const f of readdirSync(OUT_DIR)) unlinkSync(path.join(OUT_DIR, f)); // clean stale frames

const FPS = 30;
const W = 1080, H = 1920;
const T = THEMES.dark;
const CF = 0.35; // crossfade seconds between scenes

// ---- easing / helpers ----
const clamp = (v, lo, hi) => Math.max(lo, Math.min(hi, v));
const lerp = (a, b, p) => a + (b - a) * p;
const easeOutCubic = (p) => 1 - Math.pow(1 - p, 3);
const easeInCubic = (p) => p * p * p;
// progress 0..1 of a fade that starts at `start` and lasts `dur`, eased
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

// ==================== SCENE 0: INTRO ====================
const INTRO_DUR = 2.8;
function renderIntro(t) {
  const logoP = fadeP(t, 0.1, 0.6);
  const wordP = fadeP(t, 0.5, 0.6);
  const tagP = fadeP(t, 1.0, 0.6);
  const logo = group(iconBadge(540, 640, 130, 98, 'layout', T), { opacity: logoP, dy: lerp(30, 0, logoP), scale: lerp(0.85, 1, logoP), cx: 540, cy: 640 });
  const word = group(
    `<text x="540" y="880" font-family="${FONT_SANS}" font-size="150" font-weight="bold" fill="#00D9FF" text-anchor="middle">Goca</text>`,
    { opacity: wordP, dy: lerp(24, 0, wordP) }
  );
  const tag = group(
    `<text x="540" y="950" font-family="${FONT_SANS}" font-size="38" font-weight="600" fill="#ffffff" fill-opacity="0.85" text-anchor="middle">Go Clean Architecture,</text>
     <text x="540" y="1000" font-family="${FONT_SANS}" font-size="38" font-weight="600" fill="#ffffff" fill-opacity="0.85" text-anchor="middle">generated in seconds.</text>`,
    { opacity: tagP, dy: lerp(18, 0, tagP) }
  );
  return `${bg('26%')}${logo}${word}${tag}`;
}

// ==================== SCENE 1: TERMINAL ====================
const TERM_DUR = 7.2;
const CMD = '$ goca feature Order --fields "id:string,total:float64"';
const CHECKLIST = [
  '✓ Domain layer generated',
  '✓ Use case layer generated',
  '✓ Repository generated',
  '✓ HTTP handler generated',
  '✓ Integration tests generated',
];
function renderTerminal(t) {
  const winP = fadeP(t, 0.0, 0.4);
  const typeStart = 0.4, typeDur = 2.2;
  const typeProgress = clamp((t - typeStart) / typeDur, 0, 1);
  const chars = Math.round(CMD.length * typeProgress);
  const typed = CMD.slice(0, chars);
  const cmdFontSize = 25;
  const cursorX = 100 + typed.length * (cmdFontSize * 0.62);
  const cursorOn = Math.floor(t / 0.25) % 2 === 0;
  const showCursor = t < typeStart + typeDur + 0.3;

  const checklistStart = typeStart + typeDur + 0.15;
  const stagger = 0.45;
  const checklistLines = CHECKLIST.map((line, i) => {
    const p = fadeP(t, checklistStart + i * stagger, 0.35);
    return group(
      `<text x="100" y="${520 + i * 56}" font-family="${FONT_MONO}" font-size="27" fill="#00D9FF">${escXml(line)}</text>`,
      { opacity: p, dy: lerp(14, 0, p) }
    );
  }).join('');

  const noteStart = checklistStart + CHECKLIST.length * stagger + 0.1;
  const noteP = fadeP(t, noteStart, 0.4);
  const note = group(
    `<text x="100" y="850" font-family="${FONT_SANS}" font-size="22" fill="#8b949e">Dependency injection wired.</text>
     <text x="100" y="884" font-family="${FONT_SANS}" font-size="22" fill="#8b949e">Clean Architecture rules validated.</text>`,
    { opacity: noteP, dy: lerp(10, 0, noteP) }
  );

  const successStart = noteStart + 0.55;
  const successP = fadeP(t, successStart, 0.4);
  const success = group(
    `<text x="100" y="970" font-family="${FONT_MONO}" font-size="32" font-weight="bold" fill="#3fb950">Feature 'Order' ready in 1.4s</text>`,
    { opacity: successP, scale: lerp(0.92, 1, successP), cx: 100, cy: 970 }
  );

  const term = group(`
    <rect x="60" y="300" width="960" height="1220" rx="20" fill="#0d1117" stroke="#00ADD8" stroke-opacity="0.35" stroke-width="1.5"/>
    <rect x="60" y="300" width="960" height="56" rx="20" fill="#161b22"/>
    <rect x="60" y="336" width="960" height="20" fill="#161b22"/>
    <circle cx="102" cy="328" r="9" fill="#ff5f56"/>
    <circle cx="132" cy="328" r="9" fill="#ffbd2e"/>
    <circle cx="162" cy="328" r="9" fill="#27c93f"/>
    <text x="540" y="335" font-family="${FONT_MONO}" font-size="19" fill="#8b949e" text-anchor="middle">goca &#8212; zsh</text>
    <text x="100" y="430" font-family="${FONT_MONO}" font-size="${cmdFontSize}" fill="#c9d1d9"><tspan fill="#27c93f">$</tspan> ${escXml(typed.replace(/^\$ /, ''))}</text>
    ${showCursor && cursorOn ? `<rect x="${cursorX.toFixed(1)}" y="404" width="13" height="30" fill="#00D9FF"/>` : ''}
    ${checklistLines}
    ${note}
    ${success}
  `, { opacity: winP, dy: lerp(40, 0, winP) });

  return `${bg('12%')}${brandLockup(90, 170, 26, 32, T)}${term}`;
}

// ==================== SCENE 2: LAYERS ====================
const LAYERS_DUR = 4.5;
const LAYER_ROWS = [
  { label: 'Domain — pure entities', op: 0.3, textDark: true, y: 1300 },
  { label: 'Use Case — business rules', op: 0.5, textDark: true, y: 1170 },
  { label: 'Repository — persistence', op: 0.7, textDark: false, y: 1040 },
  { label: 'Handler — HTTP / gRPC / CLI', op: 0.9, textDark: false, y: 910 },
];
function renderLayers(t) {
  const titleP = fadeP(t, 0.0, 0.5);
  const title = group(
    `<text x="540" y="380" font-family="${FONT_SANS}" font-size="52" font-weight="bold" fill="#ffffff" text-anchor="middle">All layers.</text>
     <text x="540" y="446" font-family="${FONT_SANS}" font-size="52" font-weight="bold" fill="#00D9FF" text-anchor="middle">One command.</text>`,
    { opacity: titleP, dy: lerp(-16, 0, titleP) }
  );

  const bars = LAYER_ROWS.map((row, i) => {
    const start = 0.5 + i * 0.4;
    const p = fadeP(t, start, 0.4);
    const textFill = row.textDark ? '#04202b' : '#ffffff';
    return group(`
      <rect x="110" y="${row.y}" width="860" height="110" rx="16" fill="#00ADD8" fill-opacity="${row.op}"/>
      <text x="150" y="${row.y + 66}" font-family="${FONT_SANS}" font-size="30" font-weight="700" fill="${textFill}">${row.label}</text>
    `, { opacity: p, dy: lerp(30, 0, p) });
  }).join('');

  const subP = fadeP(t, 0.5 + LAYER_ROWS.length * 0.4 + 0.2, 0.4);
  const sub = group(
    `<text x="540" y="1500" font-family="${FONT_SANS}" font-size="28" fill="#ffffff" fill-opacity="0.72" text-anchor="middle">Dependencies always point inward.</text>`,
    { opacity: subP }
  );

  return `${bg('30%')}${brandLockup(90, 170, 26, 32, T)}${title}${bars}${sub}`;
}

// ==================== SCENE 3: VALUE CARDS ====================
const CARDS_DUR = 5.0;
const CARDS = [
  { icon: 'target', h1: 'Anti-patterns', h2: 'blocked by design.' },
  { icon: 'flask', h1: 'Tests generated,', h2: 'not an afterthought.' },
  { icon: 'plug', h1: 'Wired together,', h2: 'not by hand.' },
];
function renderOneCard(card, p) {
  const icon = group(iconBadge(540, 700, 150, 112, card.icon, T), { opacity: p, scale: lerp(0.9, 1, p), cx: 540, cy: 700 });
  const heads = group(
    `<text x="540" y="940" font-family="${FONT_SANS}" font-size="64" font-weight="bold" fill="#ffffff" text-anchor="middle">${escXml(card.h1)}</text>
     <text x="540" y="1014" font-family="${FONT_SANS}" font-size="64" font-weight="bold" fill="#00D9FF" text-anchor="middle">${escXml(card.h2)}</text>`,
    { opacity: p, dy: lerp(16, 0, p) }
  );
  return icon + heads;
}
function renderCards(t) {
  const n = CARDS.length;
  const each = CARDS_DUR / n;
  let idx = clamp(Math.floor(t / each), 0, n - 1);
  const localT = t - idx * each;
  const p = fadeP(localT, 0, 0.35);
  let inner = renderOneCard(CARDS[idx], p);
  // crossfade out near the end of this card's slot into the next
  if (idx < n - 1 && localT > each - CF) {
    const outP = 1 - clamp((localT - (each - CF)) / CF, 0, 1);
    const nextP = fadeP(0, 0, 0.35); // next card's intro fade hasn't started yet -> renders at p=0 (invisible-ish icon pop start)
    inner = group(renderOneCard(CARDS[idx], 1), { opacity: outP }) + group(renderOneCard(CARDS[idx + 1], 0.15), { opacity: 1 - outP });
  }
  const label = group(
    `<text x="540" y="1650" font-family="${FONT_SANS}" font-size="24" fill="#ffffff" fill-opacity="0.5" text-anchor="middle">goca feature &#8212; every time</text>`,
    { opacity: fadeP(t, 0.3, 0.4) }
  );
  return `${bg('34%')}${brandLockup(90, 170, 26, 32, T)}${inner}${label}`;
}

// ==================== SCENE 4: OUTRO ====================
const OUTRO_DUR = 3.5;
function renderOutro(t) {
  const p1 = fadeP(t, 0.0, 0.5);
  const p2 = fadeP(t, 0.25, 0.5);
  const p3 = fadeP(t, 0.5, 0.5);
  const p4 = fadeP(t, 0.75, 0.5);

  const logo = group(iconBadge(540, 560, 130, 98, 'layout', T), { opacity: p1, scale: lerp(0.85, 1, p1), cx: 540, cy: 560 });
  const word = group(`<text x="540" y="800" font-family="${FONT_SANS}" font-size="140" font-weight="bold" fill="#00D9FF" text-anchor="middle">Goca</text>`, { opacity: p2, dy: lerp(16, 0, p2) });
  const tag = group(`<text x="540" y="864" font-family="${FONT_SANS}" font-size="34" font-weight="600" fill="#ffffff" fill-opacity="0.85" text-anchor="middle">Go Clean Architecture Code Generator</text>`, { opacity: p2 });

  const ctaW = 480, ctaX = 540 - ctaW / 2;
  const cta = group(`
    <rect x="${ctaX}" y="1000" width="${ctaW}" height="96" rx="48" fill="#00ADD8"/>
    <text x="540" y="1060" font-family="${FONT_SANS}" font-size="34" font-weight="bold" fill="#ffffff" text-anchor="middle">Get Goca free</text>
  `, { opacity: p3, dy: lerp(16, 0, p3) });

  const cmd = group(
    `<text x="540" y="1160" font-family="${FONT_MONO}" font-size="26" fill="#00D9FF" text-anchor="middle">$ go install github.com/sazardev/goca@latest</text>`,
    { opacity: p4 }
  );
  const url = group(`<text x="540" y="1220" font-family="${FONT_SANS}" font-size="26" fill="#ffffff" fill-opacity="0.5" text-anchor="middle">sazardev.github.io/goca</text>`, { opacity: p4 });

  return `${bg('24%')}${logo}${word}${tag}${cta}${cmd}${url}`;
}

// ==================== TIMELINE ====================
const SCENES = [
  { name: 'intro', dur: INTRO_DUR, render: renderIntro },
  { name: 'terminal', dur: TERM_DUR, render: renderTerminal },
  { name: 'layers', dur: LAYERS_DUR, render: renderLayers },
  { name: 'cards', dur: CARDS_DUR, render: renderCards },
  { name: 'outro', dur: OUTRO_DUR, render: renderOutro },
];
const starts = [];
{ let acc = 0; for (const s of SCENES) { starts.push(acc); acc += s.dur; } }
const TOTAL = starts[starts.length - 1] + SCENES[SCENES.length - 1].dur;

const INTRO_FADE = 0.35, OUTRO_FADE = 0.6;
function blackOverlay(t) {
  if (t < INTRO_FADE) return 1 - easeOutCubic(t / INTRO_FADE);
  if (t > TOTAL - OUTRO_FADE) return easeInCubic((t - (TOTAL - OUTRO_FADE)) / OUTRO_FADE);
  return 0;
}

function composeFrame(t) {
  let idx = 0;
  for (let i = 0; i < SCENES.length; i++) if (starts[i] <= t) idx = i;
  const nextIdx = idx + 1 < SCENES.length ? idx + 1 : null;
  let inner;
  if (nextIdx !== null && t >= starts[nextIdx] - CF) {
    const p = clamp((t - (starts[nextIdx] - CF)) / CF, 0, 1);
    const curT = clamp(t - starts[idx], 0, SCENES[idx].dur);
    const nxtT = clamp(t - starts[nextIdx], 0, SCENES[nextIdx].dur);
    inner = group(SCENES[idx].render(curT), { opacity: 1 - p }) + group(SCENES[nextIdx].render(nxtT), { opacity: p });
  } else {
    inner = SCENES[idx].render(clamp(t - starts[idx], 0, SCENES[idx].dur));
  }
  const overlay = blackOverlay(t);
  return `<svg width="${W}" height="${H}" xmlns="http://www.w3.org/2000/svg">${inner}${overlay > 0 ? `<rect width="${W}" height="${H}" fill="#000000" opacity="${overlay.toFixed(3)}"/>` : ''}</svg>`;
}

const totalFrames = Math.round(TOTAL * FPS);
for (let f = 0; f < totalFrames; f++) {
  const t = f / FPS;
  const svg = composeFrame(t);
  writeFileSync(path.join(OUT_DIR, `frame_${String(f).padStart(5, '0')}.svg`), svg, 'utf-8');
}

console.log(`Generated ${totalFrames} frames (${TOTAL.toFixed(2)}s @ ${FPS}fps) in ${OUT_DIR}`);
