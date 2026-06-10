# OpenLucky Clean-Room Reverse-Engineering Record

## Why This Document Exists

OpenLucky is not a Lucky fork. It is a clean-room implementation that uses observed behavior, route names, endpoint names, and risk taxonomy as requirements input for a new open-source codebase.

This document exists so contributors can understand:

- what was actually observed from an authorized Lucky installation,
- how that information was collected,
- which parts of the observed product were only inventoried,
- which parts were exercised safely,
- how those observations map to the current OpenLucky MVP,
- and where the project deliberately stops short of feature parity.

## Clean-Room Rules

The implementation repository must not import or copy:

- Lucky source code,
- compiled Lucky JavaScript bundles,
- Vue or Quasar component code from Lucky,
- Lucky CSS class names as a design source,
- Lucky icons, screenshots, images, or copied UI copy,
- Lucky private persistence formats,
- raw authenticated response bodies, tokens, secrets, or credentials.

Allowed inputs for OpenLucky design work are:

- route names,
- endpoint names,
- HTTP methods,
- coarse response shape,
- behavioral grouping,
- risk classification,
- and manually written notes about what a screen or module is for.

## Research Workflow Summary

The reverse-engineering process followed three distinct stages.

### 1. Static bundle inventory

An authorized Lucky admin instance was loaded and its frontend JavaScript chunks were inspected offline.

That stage produced:

- route inventory,
- endpoint wrapper inventory,
- HTTP method classification,
- parameterized vs non-parameterized path separation,
- and a first-pass risk split between safe read-only candidates and dangerous/state-changing candidates.

Key static findings:

- 33 observed hash routes
- 458 unique API wrappers in downloaded frontend chunks
- 214 `GET`, 123 `POST`, 74 `PUT`, 47 `DELETE`
- 173 read-only or redacted-read execution candidates
- 285 skipped candidates because they were mutating, parameterized, or action-like

### 2. Route and control inventory

Browser automation logged in and visited the observed routes to inventory visible controls and passive network behavior.

This stage was intentionally treated as **inventory only**, not “full button testing”. It showed which pages and controls existed, but it did not justify claims like “everything was clicked” or “the backend was fully tested”.

### 3. Safe live audit

A later audit pass used an explicit safety harness:

- log in,
- enumerate controls,
- classify controls before clicking,
- block external requests,
- block mutating methods,
- block action-like `GET` requests,
- dismiss dialogs,
- store only status/content-type/shape metadata,
- and never persist raw response bodies, headers, tokens, or secrets.

This safe live audit produced the final evidence that OpenLucky currently references.

## What Was Actually Verified

The final safe live audit verified:

- 33 frontend routes
- 1,571 discovered visible controls
- 1,232 controls clicked safely or under network guard
- 28 controls that attempted mutation/external activity and were blocked
- 39 controls blocked because they were secret-related
- 109 controls blocked because labels indicated dangerous actions
- 38 external links blocked from navigation
- 10 disabled controls
- 115 state-toggle controls blocked from execution
- 173 authenticated read-only/redacted/public backend candidates executed
- 285 backend candidates skipped with explicit reasons
- no remaining frontend `error` status in the final safe audit artifact

Observed backend status results in the safe run:

- 162 HTTP `200`
- 11 HTTP `404`

Those numbers describe the **safe audit coverage**, not feature completeness.

## Observed Route Surface

The observed Lucky surface included these route families:

- core: `/about`, `/status`, `/logscenter`, `/set`
- MVP-shaped service modules: `/ddns`, `/web`, `/portforward`, `/ssl`, `/cron`
- infrastructure and storage modules: `/stun`, `/cloudflared`, `/frp`, `/docker`, `/webterminal`, `/webdav`, `/smb`, `/ftpserver`, `/filebrowser`, `/dlnaservice`, `/storagemanagement`, `/rclone`
- security modules: `/ipfilter`, `/securitygroups`, `/ipdb`, `/coraza`, `/thirdPartyAuthManager`
- other admin surfaces: `/wol`, `/dbbackup`, `/wanjiadmin/*`

OpenLucky does **not** claim that all of these are implemented.

## Mapping From Research To OpenLucky

The current repository turns that reverse-engineering record into a conservative MVP:

### Implemented in the MVP

- `/healthz`
- `/lucky/` static admin shell
- `POST /api/login/challenge`
- `POST /api/login`
- `POST /api/logout`
- `GET /api/status/host-overview`
- `GET /api/status/module-overview`
- `GET /api/logscenter/query`
- `GET /api/baseconfigure`
- `PUT /api/baseconfigure`
- `GET /api/modules/list`
- `GET /api/ddnstasklist`
- `GET /api/webservice/rules`
- `GET /api/portforwards`
- `GET /api/ssl`
- `GET /api/cron/list`

### Present but intentionally stubbed

- Docker
- Web terminal / SFTP
- WebDAV
- SMB
- FTP server
- File browser
- DLNA
- Storage management
- Rclone
- Cloudflared
- FRP
- STUN
- IP database
- Coraza WAF
- Third-party auth
- Wanji admin surfaces

These routes return explicit `501 Not Implemented` responses rather than pretending to work.

## Why OpenLucky Does Not Mirror All Official Routes Yet

The reverse-engineering record showed many high-risk operations:

- file manipulation,
- shell and terminal access,
- Docker control,
- process kill and restart actions,
- config restore/update/reboot paths,
- WOL shutdown/wakeup flows,
- certificate sync and renewal paths,
- upload/download surfaces,
- and tunnel or proxy control flows.

Shipping those as fake or weakly protected MVP handlers would be worse than leaving them unimplemented.

So the current project explicitly chooses:

- honest stubs,
- typed module metadata,
- route visibility,
- and clear documentation of the boundary.

## Superpowers Planning Trail

The OpenLucky MVP was not produced as a one-shot code dump. It came from a structured planning trail that used:

- route and endpoint inventory,
- safety classification,
- parallel context gathering,
- implementation sequencing,
- verification before completion,
- and explicit clean-room constraints.

That planning trail is why the repository is split into:

- a small Hertz backend core,
- a dependency-free frontend shell,
- typed module registry and stubs,
- test-backed auth/config/status/logs primitives,
- and documentation that separates “observed” from “implemented”.

## Current Parity Statement

The correct parity statement for the current repository is:

> OpenLucky is informed by a clean-room reverse-engineering process against an authorized Lucky instance, but it currently implements only a conservative MVP subset of the observed route and API surface.

That means:

- it does **not** claim official route-for-route parity,
- it does **not** claim backend feature parity,
- it does **not** claim Lucky data format compatibility,
- and it does **not** claim that all observed modules are safe to ship yet.

## What Contributors Should Do Next

When extending OpenLucky, contributors should preserve the same discipline:

1. Start from behavior requirements, not copied Lucky code.
2. Keep dangerous modules stubbed until they have their own threat model.
3. Document whether a route is observed, implemented, or intentionally stubbed.
4. Add tests before widening capability.
5. Prefer small verified expansions over broad parity claims.

## Related Files

- `README.md` - public project overview and MVP scope
- `internal/app/app.go` - current backend route surface
- `internal/modules/registry.go` - current module registry and stub split
- `tests/browser/openlucky-smoke.spec.mjs` - current browser smoke flow
