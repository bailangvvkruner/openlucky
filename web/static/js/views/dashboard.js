import { APIError, apiFetch, hasToken } from '../api.js';

export const MODULE_ROUTES = [
  { path: '/ddns', title: 'DDNS', group: 'MVP modules', phase: 'MVP list', endpoint: '/api/ddnstasklist', empty: 'No DDNS tasks have been defined yet.' },
  { path: '/web', title: 'Web service', group: 'MVP modules', phase: 'MVP list', endpoint: '/api/webservice/rules', empty: 'No web service rules have been defined yet.' },
  { path: '/portforward', title: 'Port forward', group: 'MVP modules', phase: 'MVP list', endpoint: '/api/portforwards', empty: 'No port forwarding rules have been defined yet.' },
  { path: '/ssl', title: 'SSL certificates', group: 'MVP modules', phase: 'MVP list', endpoint: '/api/ssl', empty: 'No certificate records have been defined yet.' },
  { path: '/cron', title: 'Cron jobs', group: 'MVP modules', phase: 'MVP list', endpoint: '/api/cron/list', empty: 'No scheduled jobs have been defined yet.' },
  { path: '/ipfilter', title: 'IP filter', group: 'Security stubs', phase: 'Stub' },
  { path: '/securitygroups', title: 'Security groups', group: 'Security stubs', phase: 'Stub' },
  { path: '/dbbackup', title: 'Database backup', group: 'Safety stubs', phase: 'Stub' },
  { path: '/stun', title: 'STUN', group: 'Deferred modules', phase: 'Stub' },
  { path: '/cloudflared', title: 'Cloudflared', group: 'Deferred modules', phase: 'Stub' },
  { path: '/frp', title: 'FRP', group: 'Deferred modules', phase: 'Stub' },
  { path: '/docker', title: 'Docker', group: 'Deferred modules', phase: 'Stub' },
  { path: '/webterminal', title: 'Web terminal', group: 'Deferred modules', phase: 'Stub' },
  { path: '/webdav', title: 'WebDAV', group: 'Deferred modules', phase: 'Stub' },
  { path: '/smb', title: 'SMB', group: 'Deferred modules', phase: 'Stub' },
  { path: '/ftpserver', title: 'FTP server', group: 'Deferred modules', phase: 'Stub' },
  { path: '/filebrowser', title: 'File browser', group: 'Deferred modules', phase: 'Stub' },
  { path: '/dlnaservice', title: 'DLNA service', group: 'Deferred modules', phase: 'Stub' },
  { path: '/storagemanagement', title: 'Storage management', group: 'Deferred modules', phase: 'Stub' },
  { path: '/rclone', title: 'Rclone', group: 'Deferred modules', phase: 'Stub' },
  { path: '/wol', title: 'Wake on LAN', group: 'Deferred modules', phase: 'Stub' },
  { path: '/ipdb', title: 'IP database', group: 'Deferred modules', phase: 'Stub' },
  { path: '/coraza', title: 'Coraza WAF', group: 'Deferred modules', phase: 'Stub' },
  { path: '/thirdPartyAuthManager', title: 'Third-party auth', group: 'Deferred modules', phase: 'Stub' }
];

const CORE_ROUTES = [
  { path: '/status', title: 'Dashboard', group: 'Core', phase: 'Live' },
  { path: '/logs', title: 'Logs', group: 'Core', phase: 'Live' },
  { path: '/settings', title: 'Settings', group: 'Core', phase: 'Live' },
  { path: '/modules', title: 'Modules', group: 'Core', phase: 'Live' }
];

const LOCAL_CATALOG = [...CORE_ROUTES, ...MODULE_ROUTES];

export async function renderStatus() {
  requireToken();
  const [host, moduleOverview, modules] = await Promise.all([
    apiFetch('/api/status/host-overview'),
    apiFetch('/api/status/module-overview'),
    loadModules()
  ]);

  const moduleItems = normalizeItems(moduleOverview);
  const body = `
    ${pageTitle('Dashboard', 'Watch service health, module readiness, and runtime identity from a compact control surface.')}
    <section class="ol-dashboard-grid" aria-label="Host overview">
      ${statCard('Service', host?.service || 'OpenLucky')}
      ${statCard('Version', host?.version || 'dev')}
      ${statCard('Runtime', host?.runtime || 'go+hertz')}
      ${statCard('Uptime', formatUptime(host?.uptimeSeconds))}
    </section>
    <section class="ol-section" aria-labelledby="module-health-title">
      <div class="ol-section__header">
        <div>
          <h2 id="module-health-title">Module health</h2>
          <p>${modules.length} modules are known to this console.</p>
        </div>
        <a class="ol-button ol-button--ghost" href="#/modules">Review modules</a>
      </div>
      ${renderModuleTable(moduleItems.length ? moduleItems : modules)}
    </section>`;
  return renderShell('/status', body);
}

export async function renderLogs() {
  requireToken();
  const logs = normalizeItems(await apiFetch('/api/logscenter/query'));
  const body = `
    ${pageTitle('Logs', 'Inspect recent appliance events from the in-memory MVP log stream.')}
    <section class="ol-section" aria-labelledby="logs-title">
      <div class="ol-section__header">
        <div>
          <h2 id="logs-title">Recent events</h2>
          <p>Entries are read-only in the frontend MVP.</p>
        </div>
      </div>
      ${renderLogsList(logs)}
    </section>`;
  return renderShell('/logs', body);
}

export async function renderSettings() {
  requireToken();
  const config = await apiFetch('/api/baseconfigure');
  setTimeout(bindSettingsForm, 0);
  const body = `
    ${pageTitle('Settings', 'Review and update the JSON configuration exposed by the MVP backend.')}
    <section class="ol-card" aria-labelledby="settings-title">
      <div class="ol-card__body">
        <div>
          <h2 id="settings-title">Base configuration</h2>
          <p>Edits are sent to <code>/api/baseconfigure</code> as JSON.</p>
        </div>
        <form id="settings-form" class="ol-form">
          <div class="ol-field">
            <label class="ol-label" for="settings-json">Configuration JSON</label>
            <textarea class="ol-textarea" id="settings-json" name="config" spellcheck="false" required>${escapeHTML(JSON.stringify(config || {}, null, 2))}</textarea>
          </div>
          <button class="ol-button ol-button--accent" type="submit">Save settings</button>
          <p id="settings-status" class="ol-message" role="status" aria-live="polite"></p>
        </form>
      </div>
    </section>`;
  return renderShell('/settings', body);
}

export async function renderModules() {
  requireToken();
  const modules = await loadModules();
  const body = `
    ${pageTitle('Modules', 'Track the live MVP modules and the explicit stubs reserved for later security reviews.')}
    <section class="ol-section" aria-labelledby="module-list-title">
      <div class="ol-section__header">
        <div>
          <h2 id="module-list-title">Module registry</h2>
          <p>${modules.length} entries are available from the registry or local route catalog.</p>
        </div>
      </div>
      ${renderModuleTable(modules)}
    </section>`;
  return renderShell('/modules', body);
}

export async function renderStubModule(path) {
  requireToken();
  const route = MODULE_ROUTES.find(item => item.path === path) || { path, title: titleFromPath(path), phase: 'Stub' };
  let listContent = '';

  if (route.endpoint) {
    try {
      const result = await apiFetch(route.endpoint);
      listContent = renderDataPreview(result, route.empty);
    } catch (error) {
      listContent = renderEndpointError(error, route.endpoint);
    }
  } else {
    listContent = renderPlannedStub(route);
  }

  const body = `
    ${pageTitle(route.title, 'This route is present so operators can see the MVP boundary without unsafe placeholder actions.')}
    <section class="ol-card" aria-labelledby="stub-title">
      <div class="ol-card__body">
        <div class="ol-section__header">
          <div>
            <h2 id="stub-title">${escapeHTML(route.title)} route</h2>
            <p>${escapeHTML(route.endpoint ? 'Read-only list endpoint connected.' : 'Implementation reserved for a future module plan.')}</p>
          </div>
          <span class="ol-pill ${route.endpoint ? 'ol-pill--ok' : 'ol-pill--warn'}">${escapeHTML(route.phase)}</span>
        </div>
        ${listContent}
      </div>
    </section>`;
  return renderShell(path, body);
}

function renderShell(currentPath, body) {
  return `
    <div class="ol-shell ol-view">
      <a class="ol-skip-link" href="#main-content">Skip to content</a>
      <aside class="ol-sidebar" aria-label="OpenLucky navigation">
        <a class="ol-brand" href="#/status" aria-label="OpenLucky dashboard">
          <span class="ol-brand-mark" aria-hidden="true">OL</span>
          <span class="ol-brand__text">
            <span class="ol-brand__name">OpenLucky</span>
            <span class="ol-brand__caption">Admin console</span>
          </span>
        </a>
        <nav class="ol-nav" aria-label="Primary">
          ${renderNav(currentPath)}
        </nav>
        <div class="ol-sidebar__footer">
          <label class="ol-field">
            <span class="ol-label">Theme</span>
            <select class="ol-select" data-theme-picker aria-label="Theme">
              ${themeOption('harbor-light', 'Harbor light')}
              ${themeOption('night-watch', 'Night watch')}
            </select>
          </label>
          <button class="ol-button ol-button--ghost" type="button" data-action="logout">Log out</button>
        </div>
      </aside>
      <main id="main-content" class="ol-main">
        ${body}
      </main>
    </div>`;
}

function renderNav(currentPath) {
  const groups = groupBy(LOCAL_CATALOG, item => item.group);
  return Object.entries(groups).map(([group, items]) => `
    <div class="ol-nav__group">
      <span class="ol-nav__title">${escapeHTML(group)}</span>
      ${items.map(item => navLink(item, currentPath)).join('')}
    </div>`).join('');
}

function navLink(item, currentPath) {
  const active = equivalentPath(item.path, currentPath);
  return `
    <a class="ol-nav__link" href="#${item.path}" data-route-link aria-current="${active ? 'page' : 'false'}">
      <span>${escapeHTML(item.title)}</span>
      <span class="ol-nav__badge">${escapeHTML(item.phase)}</span>
    </a>`;
}

function equivalentPath(routePath, currentPath) {
  if (routePath === currentPath) {
    return true;
  }
  return (routePath === '/logs' && currentPath === '/logscenter') || (routePath === '/settings' && currentPath === '/set');
}

function pageTitle(title, description) {
  return `
    <header class="ol-page-title">
      <p class="ol-kicker">OpenLucky</p>
      <h1>${escapeHTML(title)}</h1>
      <p>${escapeHTML(description)}</p>
    </header>`;
}

function statCard(label, value) {
  return `
    <article class="ol-card ol-stat">
      <span class="ol-stat__label">${escapeHTML(label)}</span>
      <strong class="ol-stat__value">${escapeHTML(value)}</strong>
    </article>`;
}

function renderModuleTable(modules) {
  const rows = normalizeItems(modules).map(item => ({
    name: item.title || item.name || titleFromPath(item.route || item.path),
    route: item.route || item.path || `/${item.name || ''}`,
    phase: item.phase || (item.implemented === false ? 'Stub' : 'MVP'),
    state: stateLabel(item)
  }));

  if (!rows.length) {
    return emptyState('No modules are currently reported by the backend.');
  }

  return renderTable(
    ['Name', 'Route', 'Phase', 'State'],
    rows.map(row => [row.name, row.route, row.phase, row.state])
  );
}

function renderLogsList(logs) {
  if (!logs.length) {
    return emptyState('No log entries are available yet.');
  }

  return `
    <ul class="ol-log-list">
      ${logs.map(entry => `
        <li class="ol-log-entry">
          <span class="ol-pill ${levelClass(entry.level)}">${escapeHTML(entry.level || 'info')}</span>
          <time class="ol-log-entry__time" datetime="${escapeHTML(entry.time || '')}">${escapeHTML(formatTime(entry.time))}</time>
          <strong>${escapeHTML(entry.module || 'core')}</strong>
          <p class="ol-log-entry__message">${escapeHTML(entry.message || 'Event recorded.')}</p>
        </li>`).join('')}
    </ul>`;
}

function renderDataPreview(result, emptyMessage) {
  const items = normalizeItems(result);
  if (!items.length) {
    return emptyState(emptyMessage || 'The endpoint returned an empty list.');
  }

  const first = items[0];
  if (first && typeof first === 'object' && !Array.isArray(first)) {
    const headers = Object.keys(first).slice(0, 5);
    return renderTable(headers, items.map(item => headers.map(header => formatValue(item[header]))));
  }

  return `<pre class="ol-code-block">${escapeHTML(JSON.stringify(result, null, 2))}</pre>`;
}

function renderEndpointError(error, endpoint) {
  const message = error instanceof APIError && error.status === 404
    ? 'The backend has not exposed this MVP list endpoint yet.'
    : error.message;
  return `
    <div class="ol-empty">
      <span class="ol-pill ol-pill--warn">Endpoint</span>
      <h3>${escapeHTML(endpoint)}</h3>
      <p>${escapeHTML(message)}</p>
    </div>`;
}

function renderPlannedStub(route) {
  return `
    <div class="ol-empty">
      <span class="ol-pill ol-pill--warn">Stub</span>
      <h3>${escapeHTML(route.title)} is intentionally inactive.</h3>
      <p>This clean-room MVP avoids unsafe placeholder controls for modules that need their own runtime and security plan.</p>
    </div>`;
}

function renderTable(headers, rows) {
  return `
    <div class="ol-table-wrap">
      <table class="ol-table">
        <thead>
          <tr>${headers.map(header => `<th scope="col">${escapeHTML(header)}</th>`).join('')}</tr>
        </thead>
        <tbody>
          ${rows.map(row => `
            <tr>${row.map((cell, index) => `<td data-label="${escapeHTML(headers[index])}">${escapeHTML(cell)}</td>`).join('')}</tr>`).join('')}
        </tbody>
      </table>
    </div>`;
}

function emptyState(message) {
  return `
    <div class="ol-empty">
      <span class="ol-pill">Empty</span>
      <p>${escapeHTML(message)}</p>
    </div>`;
}

function bindSettingsForm() {
  const form = document.querySelector('#settings-form');
  if (!form) {
    return;
  }

  form.addEventListener('submit', async event => {
    event.preventDefault();
    const status = form.querySelector('#settings-status');
    const submit = form.querySelector('button[type="submit"]');
    status.className = 'ol-message';
    status.textContent = 'Saving settings.';
    submit.disabled = true;

    try {
      const payload = JSON.parse(form.elements.config.value);
      await apiFetch('/api/baseconfigure', { method: 'PUT', body: payload });
      status.className = 'ol-message ol-message--success';
      status.textContent = 'Settings saved.';
    } catch (error) {
      status.className = 'ol-message ol-message--error';
      status.textContent = error.message;
    } finally {
      submit.disabled = false;
    }
  });
}

async function loadModules() {
  try {
    const modules = normalizeItems(await apiFetch('/api/modules/list'));
    return modules.length ? modules : LOCAL_CATALOG;
  } catch (error) {
    if (error instanceof APIError && error.status === 401) {
      throw error;
    }
    return LOCAL_CATALOG;
  }
}

function normalizeItems(value) {
  if (Array.isArray(value)) {
    return value;
  }
  if (!value || typeof value !== 'object') {
    return [];
  }
  for (const candidate of ['items', 'list', 'modules', 'entries', 'logs', 'data']) {
    if (Array.isArray(value[candidate])) {
      return value[candidate];
    }
  }
  const arrayValue = Object.values(value).find(Array.isArray);
  return arrayValue || [];
}

function requireToken() {
  if (!hasToken()) {
    location.hash = '#/login';
    throw new Error('Sign in to continue.');
  }
}

function groupBy(items, keyFn) {
  return items.reduce((groups, item) => {
    const key = keyFn(item);
    groups[key] = groups[key] || [];
    groups[key].push(item);
    return groups;
  }, {});
}

function themeOption(value, label) {
  const selected = document.body.dataset.theme === value ? ' selected' : '';
  return `<option value="${escapeHTML(value)}"${selected}>${escapeHTML(label)}</option>`;
}

function stateLabel(item) {
  if (item.implemented === false || item.phase === 'Stub') {
    return 'Stubbed';
  }
  if (item.enabled === false) {
    return 'Disabled';
  }
  return 'Available';
}

function levelClass(level = '') {
  const normalized = String(level).toLowerCase();
  if (['error', 'fatal', 'panic'].includes(normalized)) {
    return 'ol-pill--danger';
  }
  if (['warn', 'warning'].includes(normalized)) {
    return 'ol-pill--warn';
  }
  return 'ol-pill--ok';
}

function formatUptime(seconds) {
  const total = Number(seconds || 0);
  if (!Number.isFinite(total) || total <= 0) {
    return 'Starting';
  }
  const hours = Math.floor(total / 3600);
  const minutes = Math.floor((total % 3600) / 60);
  if (hours) {
    return `${hours}h ${minutes}m`;
  }
  return `${minutes || 1}m`;
}

function formatTime(value) {
  if (!value) {
    return 'Just now';
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString();
}

function formatValue(value) {
  if (value === null || value === undefined || value === '') {
    return 'None';
  }
  if (typeof value === 'object') {
    return JSON.stringify(value);
  }
  return value;
}

function titleFromPath(path) {
  return String(path || 'module')
    .replace(/^\//, '')
    .replace(/[-_]/g, ' ')
    .replace(/\b\w/g, char => char.toUpperCase());
}

function escapeHTML(value) {
  return String(value).replace(/[&<>"]/g, char => ({
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;'
  })[char]);
}
