# Night Safari: design

Dark and glowing, like a night drive: a night-sky background, warm amber highlights like
fireflies or a torch beam, moonlit teal as the cool accent. Only African wildlife is shown.

## Colours (`web/assets/css/app.css`, `:root`)

| Token | Hex | Used for |
|-------|-----|----------|
| `--night` | `#0E1420` | page background, `theme-color` |
| `--night-2` | `#151D2E` | cards, panels |
| `--night-3` | `#1E2940` | raised surfaces, photo placeholders, hover |
| `--cream` | `#F2E9D8` | body text, headings |
| `--muted` | `#A9B0BF` | secondary text (8:1 on night) |
| `--amber` | `#F5B841` | links, buttons, figures, map dots, globe dots, charts |
| `--teal` | `#3FB8AF` | eyebrows, labels, globe halo, wet-weather icons, bars |
| `--line` | `rgba(242,233,216,.12)` | hairline borders |

The page has two faint glows in its background: teal top right, amber top left.

## Lettering

- Headings: Space Grotesk (500–700), modern and clean, slight negative tracking.
- Body: Inter (400–600).
- Small uppercase "eyebrow" labels in teal above page titles.

## Shape and spacing

- Radius `--radius: 14px` (cards, panels, map, photos), `--radius-small: 8px`; buttons and chips are pills.
- Roomy: sections `--section` apart (3–5.5rem), `--gap` 1–1.75rem, page width 76rem.

## Components

- **Header**: sticky, translucent night with blur; logo (amber full moon behind an acacia,
  `assets/images/logo.svg`) and "Night Safari" wordmark; nav pills, current page in amber.
- **Hero**: headline and summary on the left, a 3D globe (three.js, `assets/js/scenes/globe.js`) on
  the right: NASA Earth-at-night texture, amber glowing dot per sighting, teal halo, stars,
  slow auto-rotate, drag to turn. Stacks on phones.
- **Cards**: 4:3 photo on top (`object-fit: cover`), name, scientific name in italics, seen/where/
  count/by, weather readout at the bottom. Lift 3px with an amber border on hover.
- **Animal tiles**: square photo tiles with a dark gradient and the name over it.
- **Map**: Leaflet with dark (inverted) tiles, amber cluster bubbles with a glow.
- **Charts**: amber sparkline, amber gradient columns, teal bars.
- **Photo wall**: square tiles, lightbox on click, caption with animal and country.
- Everything that moves respects reduced motion.

## Pages

`/` home · `/sightings/{id}` · `/animals/{animal}` · `/countries/{code}` · `/stats` · `/gallery`

## Avoid

- Light backgrounds, or the generic indigo/violet the design search suggested.
- 3D per card (the attendee chose the globe only).
