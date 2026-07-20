// Generates one marketing card per Goca CLI command, in 3 formats (square,
// landscape, portrait) x 2 themes (dark, light) = 6 files/command.
// Run: node marketing/scripts/generate-commands.mjs
import { writeFileSync, mkdirSync } from 'fs';
import { fileURLToPath } from 'url';
import path from 'path';
import {
  FONT_SANS, FONT_MONO, THEMES,
  bgDefs, brandLockup, iconBadge, tagPillAuto, cmdChip, escXml, footerUrl,
} from './brand.mjs';
import { COMMANDS } from './commands-data.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const OUT_DIR = path.resolve(__dirname, '..', 'svg', 'commands');
mkdirSync(OUT_DIR, { recursive: true });

// ---- Format renderers ----

function renderSquare(c, t) {
  const W = 1080, H = 1080;
  return `<svg width="${W}" height="${H}" xmlns="http://www.w3.org/2000/svg">
  ${bgDefs('50%', '35%', '55%', t)}
  <rect width="${W}" height="${H}" fill="url(#bgGrad)"/>
  <rect width="${W}" height="${H}" fill="url(#glow)"/>
  ${brandLockup(88, 86, 26, 32, t)}
  ${tagPillAuto(990, 62, 44, c.tag, 17, t)}
  ${iconBadge(540, 430, 120, 90, c.icon, t)}
  <text x="540" y="612" font-family="${FONT_SANS}" font-size="58" font-weight="bold" fill="${t.headlineMain}" text-anchor="middle">${escXml(c.headline[0])}</text>
  <text x="540" y="678" font-family="${FONT_SANS}" font-size="58" font-weight="bold" fill="${t.headlineAccent}" text-anchor="middle">${escXml(c.headline[1])}</text>
  <text x="540" y="738" font-family="${FONT_SANS}" font-size="27" fill="${t.subtext}" fill-opacity="${t.subtextOpacity}" text-anchor="middle">${escXml(c.sub[0])}</text>
  <text x="540" y="772" font-family="${FONT_SANS}" font-size="27" fill="${t.subtext}" fill-opacity="${t.subtextOpacity}" text-anchor="middle">${escXml(c.sub[1])}</text>
  ${cmdChip(90, 850, 900, 72, c.cmd, 24, t, 'middle')}
  ${footerUrl(540, 1010, 22, t, 'middle')}
</svg>
`;
}

function renderPortrait(c, t) {
  const W = 1080, H = 1920;
  return `<svg width="${W}" height="${H}" xmlns="http://www.w3.org/2000/svg">
  ${bgDefs('50%', '28%', '45%', t)}
  <rect width="${W}" height="${H}" fill="url(#bgGrad)"/>
  <rect width="${W}" height="${H}" fill="url(#glow)"/>
  ${brandLockup(90, 170, 28, 34, t)}
  ${tagPillAuto(990, 148, 46, c.tag, 18, t)}
  ${iconBadge(540, 560, 140, 105, c.icon, t)}
  <text x="540" y="792" font-family="${FONT_SANS}" font-size="66" font-weight="bold" fill="${t.headlineMain}" text-anchor="middle">${escXml(c.headline[0])}</text>
  <text x="540" y="866" font-family="${FONT_SANS}" font-size="66" font-weight="bold" fill="${t.headlineAccent}" text-anchor="middle">${escXml(c.headline[1])}</text>
  <text x="540" y="928" font-family="${FONT_SANS}" font-size="30" fill="${t.subtext}" fill-opacity="${t.subtextOpacity}" text-anchor="middle">${escXml(c.sub[0])}</text>
  <text x="540" y="966" font-family="${FONT_SANS}" font-size="30" fill="${t.subtext}" fill-opacity="${t.subtextOpacity}" text-anchor="middle">${escXml(c.sub[1])}</text>
  ${cmdChip(90, 1040, 900, 80, c.cmd, 26, t, 'middle')}
  <path d="M525 1740 L540 1720 L555 1740" fill="none" stroke="${t.footer}" stroke-opacity="${t.footerOpacity}" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
  ${footerUrl(540, 1800, 28, t, 'middle')}
</svg>
`;
}

function renderLandscape(c, t) {
  const W = 1200, H = 630;
  return `<svg width="${W}" height="${H}" xmlns="http://www.w3.org/2000/svg">
  ${bgDefs('22%', '45%', '55%', t)}
  <rect width="${W}" height="${H}" fill="url(#bgGrad)"/>
  <rect width="${W}" height="${H}" fill="url(#glow)"/>
  ${brandLockup(66, 56, 22, 26, t)}
  ${tagPillAuto(1140, 40, 36, c.tag, 15, t)}
  ${iconBadge(190, 350, 110, 85, c.icon, t)}
  <text x="360" y="260" font-family="${FONT_SANS}" font-size="44" font-weight="bold" fill="${t.headlineMain}">${escXml(c.headline[0])}</text>
  <text x="360" y="316" font-family="${FONT_SANS}" font-size="44" font-weight="bold" fill="${t.headlineAccent}">${escXml(c.headline[1])}</text>
  <text x="360" y="372" font-family="${FONT_SANS}" font-size="21" fill="${t.subtext}" fill-opacity="${t.subtextOpacity}">${escXml(c.sub[0])}</text>
  <text x="360" y="404" font-family="${FONT_SANS}" font-size="21" fill="${t.subtext}" fill-opacity="${t.subtextOpacity}">${escXml(c.sub[1])}</text>
  ${cmdChip(360, 460, 780, 64, c.cmd, 20, t, 'start')}
  ${footerUrl(1140, 592, 17, t, 'end')}
</svg>
`;
}

const FORMATS = { square: renderSquare, landscape: renderLandscape, portrait: renderPortrait };

let count = 0;
for (const cmd of COMMANDS) {
  for (const [fmtName, renderFn] of Object.entries(FORMATS)) {
    for (const themeName of Object.keys(THEMES)) {
      const svg = renderFn(cmd, THEMES[themeName]);
      const file = path.join(OUT_DIR, `${cmd.id}-${fmtName}-${themeName}.svg`);
      writeFileSync(file, svg, 'utf-8');
      count++;
    }
  }
}

console.log(`Generated ${count} SVG files for ${COMMANDS.length} commands in ${OUT_DIR}`);
