import { apiFetch, clearToken, getTokenExpiry, hasToken } from './api.js';
import { navigate, registerRoute, startRouter } from './router.js';
import { renderLogin } from './views/login.js';
import {
  MODULE_ROUTES,
  renderLogs,
  renderModules,
  renderSettings,
  renderStatus,
  renderStubModule
} from './views/dashboard.js';

const root = document.querySelector('[data-app-root]');
const savedTheme = localStorage.getItem('openlucky.theme') || 'harbor-light';

document.body.dataset.theme = savedTheme;

registerRoute('/login', renderLogin);
registerRoute('/status', renderStatus);
registerRoute('/dashboard', renderStatus);
registerRoute('/logs', renderLogs);
registerRoute('/logscenter', renderLogs);
registerRoute('/settings', renderSettings);
registerRoute('/set', renderSettings);
registerRoute('/modules', renderModules);
registerRoute('/stub', renderStubModule);

for (const route of MODULE_ROUTES) {
  registerRoute(route.path, renderStubModule);
}

document.addEventListener('click', async event => {
  const logoutButton = event.target.closest('[data-action="logout"]');
  if (!logoutButton) {
    return;
  }
  logoutButton.setAttribute('disabled', 'disabled');
  try {
    await apiFetch('/api/logout', { method: 'POST' });
  } catch (error) {
    console.info('Logout completed locally after server response failed.', error);
  } finally {
    clearToken();
    navigate('/login');
  }
});

document.addEventListener('change', event => {
  const picker = event.target.closest('[data-theme-picker]');
  if (!picker) {
    return;
  }
  document.body.dataset.theme = picker.value;
  localStorage.setItem('openlucky.theme', picker.value);
});

window.addEventListener('openlucky:login', () => navigate('/status'));

if (!location.hash) {
  location.hash = hasToken() ? '#/status' : '#/login';
}

window.openLucky = Object.freeze({
  tokenPresent: hasToken,
  tokenExpiresAt: getTokenExpiry
});

startRouter(root);
