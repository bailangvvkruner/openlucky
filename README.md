# OpenLucky

OpenLucky is a clean-room, open-source admin appliance inspired by the operational shape of Lucky-style home-lab gateways. It is not a Lucky fork and does not include Lucky source code, compiled frontend bundles, CSS, icons, images, UI copy, private persistence formats, or proprietary assets.

## MVP Scope

- Single Go binary using CloudWeGo Hertz.
- Admin UI served from `/lucky/` with dependency-free vanilla HTML/CSS/JS.
- JSON APIs under `/api/*` with a normalized response envelope.
- Password login with short-lived admin tokens.
- MVP modules for status, logs, settings, module registry, DDNS/Web/port-forward/SSL/Cron list surfaces.
- Explicit `501 Not Implemented` stubs for high-risk modules such as Docker, terminal/SFTP, file services, tunnels, WAF, storage integrations, and Wanji-specific surfaces.

## Safety Model

OpenLucky binds to `127.0.0.1:16601` by default. Public exposure should be an explicit operator choice behind a trusted reverse proxy or network boundary.

Set a real admin password before running outside local development:

```sh
export OPENLUCKY_ADMIN_USER=admin
read -r -s OPENLUCKY_ADMIN_PASSWORD
export OPENLUCKY_ADMIN_PASSWORD
export OPENLUCKY_ADDR=127.0.0.1:16601
go run ./cmd/openlucky
```

Open the UI at `http://127.0.0.1:16601/lucky/`.

## Development

The current workspace may not include Go. With Go 1.23+ installed:

```sh
go mod tidy
go test ./...
go run ./cmd/openlucky
```

Frontend files are native ES modules and can be syntax-checked with Node:

```sh
node --check web/static/js/api.js
node --check web/static/js/router.js
node --check web/static/js/app.js
node --check web/static/js/views/login.js
node --check web/static/js/views/dashboard.js
```

## Clean-Room Notes

The planning artifacts in `docs/superpowers/plans/` summarize observed behavior, route names, and API risk taxonomy from an authorized audit. Those observations are used only as behavior requirements. OpenLucky implementation files are original code and should stay free of copied Lucky bundles, styles, images, icon packs, generated code, or private strings.
