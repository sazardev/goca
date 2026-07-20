// Shared Goca brand tokens + SVG helper builders, used by both the static
// command-card generator and the video-frame generator. Keep this the single
// source of truth for colors, fonts and icon paths.

export const FONT_SANS = 'Liberation Sans, Arial, sans-serif';
export const FONT_MONO = 'Liberation Mono, monospace';

// -- Icons: inner shapes only (24x24 viewBox), stroked by the wrapping <g> --
export const ICONS = {
  package: '<path d="m7.5 4.27 9 5.15"/><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 22V12"/>',
  rocket: '<path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 0 0-2.91-.09z"/><path d="m12 15-3-3a22 22 0 0 1 2-3.95A12.88 12.88 0 0 1 22 2c0 2.72-.78 7.5-6 11a22.35 22.35 0 0 1-4 2z"/><path d="M9 12H4s.55-3.03 2-4c1.62-1.08 5 0 5 0"/><path d="M12 15v5s3.03-.55 4-2c1.08-1.62 0-5 0-5"/>',
  layout: '<rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9h18"/><path d="M9 21V9"/>',
  refresh: '<path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/>',
  database: '<ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5v14a9 3 0 0 0 18 0V5"/>',
  zap: '<path d="M13 2 3 14h9l-1 8 10-12h-9l1-8z"/>',
  target: '<circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="6"/><circle cx="12" cy="12" r="2"/>',
  book: '<path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/>',
  flask: '<path d="M10 2v7.527a2 2 0 0 1-.211.896L4.72 20.55a1 1 0 0 0 .9 1.45h12.76a1 1 0 0 0 .9-1.45l-5.069-10.127A2 2 0 0 1 14 9.527V2"/><path d="M8.5 2h7"/><path d="M7 16h10"/>',
  plug: '<path d="M12 22v-5"/><path d="M9 8V2"/><path d="M15 8V2"/><path d="M18 8v5a4 4 0 0 1-4 4h-4a4 4 0 0 1-4-4V8Z"/>',
  gitbranch: '<line x1="6" x2="6" y1="3" y2="15"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><path d="M18 9a9 9 0 0 1-9 9"/>',
  shield: '<path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.79 17 5 19 5a1 1 0 0 1 1 1z"/>',
  activity: '<path d="M22 12h-4l-3 9L9 3l-3 9H2"/>',
  search: '<circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>',
  settings: '<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1Z"/>',
  copy: '<rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>',
  arrowupcircle: '<circle cx="12" cy="12" r="10"/><path d="m16 12-4-4-4 4"/><path d="M12 16V8"/>',
  bot: '<path d="M12 8V4H8"/><rect width="16" height="12" x="4" y="8" rx="2"/><path d="M2 14h2"/><path d="M20 14h2"/><path d="M15 13v2"/><path d="M9 13v2"/>',
};

// -- Theme tokens --
export const THEMES = {
  dark: {
    bgFrom: '#0a0a0a', bgTo: '#1a1a2e', glowOpacity: 0.30,
    brandText: '#ffffff', headlineMain: '#ffffff', headlineAccent: '#00D9FF',
    subtext: '#ffffff', subtextOpacity: 0.72,
    chipFill: '#000000', chipFillOpacity: 0.35, chipStroke: '#00ADD8', chipText: '#00D9FF',
    footer: '#ffffff', footerOpacity: 0.5,
    badgeRingFill: '#00ADD8', badgeRingFillOpacity: 0.12,
    tagFill: '#00ADD8', tagFillOpacity: 0.14, tagText: '#00D9FF',
    termBg: '#0d1117', termBar: '#161b22', termMuted: '#8b949e', termText: '#c9d1d9', termGreen: '#3fb950',
  },
  light: {
    bgFrom: '#ffffff', bgTo: '#e8f7fb', glowOpacity: 0.14,
    brandText: '#0f172a', headlineMain: '#0f172a', headlineAccent: '#00758F',
    subtext: '#334155', subtextOpacity: 1,
    chipFill: '#0f172a', chipFillOpacity: 1, chipStroke: '#00ADD8', chipText: '#00D9FF',
    footer: '#334155', footerOpacity: 0.6,
    badgeRingFill: '#00ADD8', badgeRingFillOpacity: 0.10,
    tagFill: '#00ADD8', tagFillOpacity: 0.10, tagText: '#00758F',
    termBg: '#0d1117', termBar: '#161b22', termMuted: '#8b949e', termText: '#c9d1d9', termGreen: '#3fb950',
  },
};

export function bgDefs(glowCx, glowCy, glowR, t, id = 'bgGrad', glowId = 'glow') {
  return `<defs>
    <linearGradient id="${id}" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="${t.bgFrom}"/>
      <stop offset="100%" stop-color="${t.bgTo}"/>
    </linearGradient>
    <radialGradient id="${glowId}" cx="${glowCx}" cy="${glowCy}" r="${glowR}">
      <stop offset="0%" stop-color="#00ADD8" stop-opacity="${t.glowOpacity}"/>
      <stop offset="100%" stop-color="#00ADD8" stop-opacity="0"/>
    </radialGradient>
  </defs>`;
}

export function brandLockup(x, y, r, fontSize, t, opacity = 1) {
  const barW = r * 0.92, barH = r * 0.12, barX = x - barW / 2, barGap = r * 0.135;
  const bars = [0.9, 0.7, 0.5, 0.3].map((op, i) =>
    `<rect x="${barX.toFixed(1)}" y="${(y - r * 0.4 + i * barGap).toFixed(1)}" width="${barW.toFixed(1)}" height="${barH.toFixed(1)}" rx="1.5" fill="#ffffff" fill-opacity="${op}"/>`
  ).join('');
  return `<g opacity="${opacity}"><circle cx="${x}" cy="${y}" r="${r}" fill="#00ADD8"/>${bars}
  <text x="${x + r + 18}" y="${y + fontSize * 0.32}" font-family="${FONT_SANS}" font-size="${fontSize}" font-weight="bold" fill="${t.brandText}">Goca</text></g>`;
}

export function iconBadge(cx, cy, rOuter, rInner, iconKey, t, opacity = 1, scaleMul = 1) {
  const s = ((rInner * 2 * 0.55) / 24) * scaleMul;
  const tx = cx - 12 * s, ty = cy - 12 * s;
  return `<g opacity="${opacity}"><circle cx="${cx}" cy="${cy}" r="${rOuter}" fill="${t.badgeRingFill}" fill-opacity="${t.badgeRingFillOpacity}" stroke="#00ADD8" stroke-width="2"/>
  <circle cx="${cx}" cy="${cy}" r="${rInner}" fill="#00ADD8"/>
  <g transform="translate(${tx.toFixed(1)},${ty.toFixed(1)}) scale(${s.toFixed(3)})" fill="none" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">${ICONS[iconKey]}</g></g>`;
}

// Right-aligns a pill of auto-computed width to rightEdge, sized to fit `text`.
export function tagPillAuto(rightEdge, y, h, text, fontSize, t, opacity = 1) {
  const w = Math.round(text.length * (fontSize * 0.62 + 0.5) + 44);
  const x = rightEdge - w;
  return `<g opacity="${opacity}"><rect x="${x}" y="${y}" width="${w}" height="${h}" rx="${h / 2}" fill="${t.tagFill}" fill-opacity="${t.tagFillOpacity}" stroke="#00ADD8" stroke-width="1.3"/>
  <text x="${x + w / 2}" y="${y + h / 2 + fontSize * 0.32}" font-family="${FONT_SANS}" font-size="${fontSize}" font-weight="700" letter-spacing="0.5" fill="${t.tagText}" text-anchor="middle">${text}</text></g>`;
}

export function cmdChip(x, y, w, h, text, fontSize, t, anchor = 'middle', opacity = 1) {
  const tx = anchor === 'middle' ? x + w / 2 : x + 24;
  return `<g opacity="${opacity}"><rect x="${x}" y="${y}" width="${w}" height="${h}" rx="14" fill="${t.chipFill}" fill-opacity="${t.chipFillOpacity}" stroke="${t.chipStroke}" stroke-width="1.5" stroke-opacity="0.5"/>
  <text x="${tx}" y="${y + h / 2 + fontSize * 0.32}" font-family="${FONT_MONO}" font-size="${fontSize}" fill="${t.chipText}" text-anchor="${anchor}">${escXml(text)}</text></g>`;
}

export function escXml(s) {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

export function footerUrl(x, y, fontSize, t, anchor = 'middle', opacity = 1) {
  return `<text x="${x}" y="${y}" font-family="${FONT_SANS}" font-size="${fontSize}" fill="${t.footer}" fill-opacity="${t.footerOpacity * opacity}" text-anchor="${anchor}">sazardev.github.io/goca</text>`;
}
