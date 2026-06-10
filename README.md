# OpenLucky

OpenLucky is a clean-room, open-source admin appliance inspired by the operational shape of Lucky-style home-lab gateways. It is not a Lucky fork and does not include Lucky source code, compiled frontend bundles, CSS, icons, images, UI copy, private persistence formats, or proprietary assets.

OpenLucky is also not a claim of official Lucky route-for-route parity. The repository contains an MVP implementation plus a public research trail explaining how the observed route/API surface was inventoried and how the current boundary was chosen.

## MVP Scope

- Single Go binary using CloudWeGo Hertz.
- Admin UI served from `/lucky/` with dependency-free vanilla HTML/CSS/JS.
- JSON APIs under `/api/*` with a normalized response envelope.
- Password login with short-lived admin tokens.
- MVP modules for status, logs, settings, module registry, DDNS/Web/port-forward/SSL/Cron list surfaces.
- Explicit `501 Not Implemented` stubs for high-risk modules such as Docker, terminal/SFTP, file services, tunnels, WAF, storage integrations, and Wanji-specific surfaces.

## Route Parity

The reverse-engineering pass observed a broad Lucky route surface, including core pages, network services, storage services, terminal/SFTP, Docker, tunnel managers, security modules, database backup, and Wanji admin surfaces.

The current OpenLucky MVP intentionally implements only the conservative subset:

- auth, config, status, logs, module registry, and static admin shell,
- read-only MVP list endpoints for DDNS, Web service rules, port forwards, SSL, and Cron,
- explicit stubs for modules that require separate threat models before they should execute real work.

That means OpenLucky shows the observed shape of the appliance without pretending that dangerous modules are complete or safe.

For the public clean-room research record, see `docs/research/clean-room-reverse-engineering.md`.

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

The public research record in `docs/research/` summarizes observed behavior, route names, API risk taxonomy, safety classification, and Superpowers-style planning decisions from an authorized audit. Those observations are used only as behavior requirements. OpenLucky implementation files are original code and should stay free of copied Lucky bundles, styles, images, icon packs, generated code, or private strings.

When adding modules, document whether each route is observed, implemented, or intentionally stubbed. Do not turn a stub into an active handler without tests and a module-specific threat model.
