const routes = new Map();
let routeRoot = null;

export function registerRoute(path, render) {
  routes.set(normalizePath(path), render);
}

export function navigate(path) {
  const nextPath = normalizePath(path);
  if (currentPath() === nextPath) {
    return renderCurrentRoute(routeRoot);
  }
  location.hash = nextPath;
  return Promise.resolve();
}

export function currentPath() {
  return normalizePath(location.hash.replace(/^#/, '') || '/login');
}

export async function renderCurrentRoute(root = routeRoot) {
  if (!root) {
    throw new Error('OpenLucky router root is missing.');
  }

  const path = currentPath();
  const render = routes.get(path) || routes.get('/stub');
  root.innerHTML = renderLoading();

  try {
    root.innerHTML = await render(path);
    updateCurrentLinks(path);
    focusMainHeading(root);
  } catch (error) {
    root.innerHTML = renderError(error);
  }
}

export function startRouter(root) {
  routeRoot = root;
  window.addEventListener('hashchange', () => renderCurrentRoute(root));
  window.addEventListener('openlucky:unauthorized', () => {
    if (currentPath() !== '/login') {
      navigate('/login');
    }
  });
  return renderCurrentRoute(root);
}

function normalizePath(path) {
  const cleaned = String(path || '/login').trim().replace(/^#/, '');
  const withSlash = cleaned.startsWith('/') ? cleaned : `/${cleaned}`;
  return withSlash.replace(/\/+$/, '') || '/login';
}

function updateCurrentLinks(path) {
  for (const link of document.querySelectorAll('[data-route-link]')) {
    const href = link.getAttribute('href') || '';
    link.setAttribute('aria-current', href === `#${path}` ? 'page' : 'false');
  }
}

function focusMainHeading(root) {
  const heading = root.querySelector('main h1');
  if (heading) {
    heading.setAttribute('tabindex', '-1');
    heading.focus({ preventScroll: true });
  }
}

function renderLoading() {
  return `
    <main class="ol-loading" aria-live="polite">
      <section class="ol-card ol-loading__card">
        <p class="ol-kicker">OpenLucky</p>
        <h1>Preparing console</h1>
        <p class="ol-message">Loading the requested view.</p>
      </section>
    </main>`;
}

function renderError(error) {
  const message = escapeHTML(error?.message || 'The requested view could not be rendered.');
  return `
    <main class="ol-error" aria-live="assertive">
      <section class="ol-card ol-error__card">
        <p class="ol-kicker">OpenLucky</p>
        <h1>View unavailable</h1>
        <p class="ol-message ol-message--error">${message}</p>
        <a class="ol-button ol-button--accent" href="#/login">Return to login</a>
      </section>
    </main>`;
}

function escapeHTML(value) {
  return String(value).replace(/[&<>"]/g, char => ({
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;'
  })[char]);
}
