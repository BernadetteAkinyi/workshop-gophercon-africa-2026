# Wild Africa: design

Warm and quiet: sunlit paper, clay and olive, savanna gold for the loudest
moments, and a soft serif for every headline. Editorial rather than poster — no
weight above 500 anywhere on the site.

(Replaces the earlier bold "deep charcoal and magenta" look, and before that the dark
"Night Safari" and neon "electric dusk" attempts.)

## Colours (`web/assets/css/app.css`, `:root`)

Every colour is a token; change one there and the whole site follows.

| Token | Hex | Used for |
|-------|-----|----------|
| `--paper` | `#FBF6ED` | the page itself |
| `--surface` | `#FFFFFF` | cards, panels, filter, map popups |
| `--surface-2` | `#F1E7D6` | sand wells: ticker, chips, chart tracks, weather pill |
| `--line` | `#E4D5BC` | hairlines between steps, card borders |
| `--ink` | `#2A2018` | headings and body text (14.8:1 on `--paper`) |
| `--ink-deep` | `#211A12` | footer, the numbers band, the closing panel, photo washes |
| `--muted` | `#5F5040` | secondary text (7.2:1) |
| `--muted-2` | `#6F5F4B` | small labels and captions (5.0:1 on `--surface-2`) |
| `--white` | `#FFFDF8` | type on dark, globe dots (16.9:1 on `--ink-deep`) |
| `--white-soft` | `#D9CDBB` | secondary type on dark (10.8:1) |
| `--accent` | `#B4552F` | clay: buttons, markers, the odd border — **fills only** |
| `--accent-deep` | `#96401F` | clay for text and links (6.4:1 on `--paper`) |
| `--accent-2` | `#55663C` | olive: focus rings, secondary links, every chart line |
| `--ochre` | `#A97613` | marks on light — the headline full stop, step numerals, icons |
| `--ochre-light` | `#E0A93B` | gold for dark: the band's figures, globe dots, footer links |

Two rules keep the warmth legible. Clay at full strength is only 4.56:1 on paper, so it is
never text: `--accent-deep` carries links and labels instead. And the ochres are cut in two —
`--ochre` on light, `--ochre-light` on dark — because no single amber clears 4.5:1 on both.
Focus rings are olive, except on the numbers band and the photo banner, where they switch to
`--white`.

## Lettering

- Headlines, figures and the ticker: **Fraunces** (400–500, optical size axis on). Warm,
  slightly wonky old-style serif; it replaces Syne's poster caps.
- Body: **Manrope** (400–500).
- **Nothing is bold**: headings 500, body 400, small uppercase labels and table headers 500.
  `b`, `strong` and `th` are pinned to 500 so the browser's own 700 never leaks in.
- Hierarchy through size and colour, not weight: `h1` `clamp(2.6rem, 6.5vw, 5rem)`, `h2` up to
  3.1rem, line-height just over 1, letter-spacing `-0.01em`; eyebrows 0.75rem uppercase, 500,
  0.2em tracking. Every major headline is signed off with a gold full stop.

## Shape, spacing, motion

- Radius `--radius: 14px`, `--radius-small: 8px`; buttons and chips are pills.
- Shadows are brown-tinted and shallow (`rgba(74, 52, 32, …)`), paired with a `--line`
  hairline rather than a heavy drop shadow.
- Roomy: sections 4.5–8rem apart, width 76rem.
- Motion is gentle: the banner photo settles in over 40s, the hero copy rises once on load,
  sections fade in as they arrive, a slow ticker of animal names runs under the banner, cards
  and tiles lift a few pixels on hover. All of it stops under reduced motion.

## The home page, top to bottom

1. **Banner**: full-width wildlife photo (`HeroPhoto`) under a warm brown wash, editorial
   headline in cream, a clay button plus a light secondary, photo credit.
2. **Ticker**: animal names drifting across a sand strip. Decorative, `aria-hidden`.
3. **Map first**: the map leads on OpenStreetMap's own light tiles with a light sepia wash,
   clay clusters; then the filter, then every sighting as a card underneath — one section,
   `#sightings`.
4. **The idea**: intro and three steps separated by hairlines, numbered in ochre.
5. **Numbers band**: full-bleed deep-brown panel with a clay glow, gold figures, gold sparkline.
6. **Meet the animals**: staggered photo wall — every third tile drops on wide screens.
7. **Globe**: sand panel with an ochre glow and a clay halo; dots in clay, gold and cream
   (`assets/js/scenes/globe.js` reads `--accent`, `--ochre-light`, `--white`).
8. **Gallery teaser**: eight square photos.
9. **Closing call**: deep-brown panel with clay and gold glows, inviting new sightings.

The home page shows nine animals — the Big Five (lion, leopard, elephant, rhino, buffalo) plus
giraffe, zebra, hippo and cheetah — set by `featuredWildlife` in `web/handlers/home.go`. Other
pages keep the full twenty-one.

## Pages

`/` home · `/sightings/{id}` · `/animals/{animal}` · `/countries/{code}` · `/stats` · `/gallery`

## Avoid

- Dark page backgrounds; anything off-token in components.
- Clay as text on paper (`--accent-deep` instead), and amber text on light (`--ochre` is for marks).
- Font weights above 500, and all-caps headlines — the uppercase is reserved for 0.75rem eyebrows
  and labels.
- A fourth hue: clay, olive and the two ochres are the whole palette.
- Glassmorphism, hard black shadows, and 3D outside the globe.
