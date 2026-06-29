# TODO

## Pending

### Site Audit — Remove Dead/Broken Sites

Site status as of 2026-04-07:

| Site | Domain | Status | Verdict |
|------|--------|--------|---------|
| readcomicsonline.ru | readcomicsonline.ru | ✅ HTTP 200, serves comics | **KEEP** |
| mangadex | api.mangadex.org | ✅ HTTP 200, API responds `pong` | **KEEP** |
| readallcomics | readallcomics.com | ✅ HTTP 200 (with browser UA) | **KEEP** |
| mangatown | mangatown.com | ✅ HTTP 200, site live | **KEEP** |
| comicextra | comicextra.com | ❌ Domain parked (Sedo parking page, not the comic site) | **REMOVE** |
| mangareader | mangareader.tv | ❌ DNS resolution failure — no address associated with hostname | **REMOVE** |
| mangakakalot | mangakakalot.com | ❌ HTTP 522 (Cloudflare origin timeout) — server unreachable | **REMOVE** |
| manganato | manganato.com | ❌ HTTP 522 (Cloudflare origin timeout) — server unreachable | **REMOVE** |

- [ ] Remove `comicextra` scraper: `pkg/sites/comicextra.go`, `pkg/sites/comicextra_test.go`, `pkg/sites/comicextra_deobfuscate_test.go`, `pkg/sites/deobfuscate.go` (if comicextra-only), update `pkg/sites/loader.go` switch, update README
- [ ] Remove `mangareader` scraper: `pkg/sites/mangareader.go`, `pkg/sites/mangareader_test.go`, update `pkg/sites/loader.go` switch, update README
- [ ] Remove `mangakakalot` scraper: `pkg/sites/mangakakalot.go`, `pkg/sites/mangakakalot_test.go`, update `pkg/sites/loader.go` switch, update README
- [ ] Remove `manganato` scraper: `pkg/sites/manganato.go`, `pkg/sites/manganato_test.go`, update `pkg/sites/loader.go` switch, update README
- [ ] Verify `deobfuscate.go` and `common.go` are not shared with kept sites before deleting
- [ ] Run `go test ./...` and `go build ./...` after removals to confirm clean compile
- [ ] Update README supported-sites table to reflect only active sites

## Completed
- [x] Add `readcomicsonline.ru` scraper (separate from `readcomiconline.li`)
- [x] Remove `readcomiconline.li` scraper: site uses server-side encrypted obfuscation, permanently broken (`pkg/sites/readcomicsonline.go`) with data-src primary path, JS `var pages` fallback, issue listing, `--all`, `--last` support, and full test coverage
- [x] Convert Makefile to Justfile
- [x] Add GUI unit tests (`cmd/gui/gui_test.go`) and `test-gui` Justfile recipe
