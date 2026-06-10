import { apiFetch, hasToken, setToken } from '../api.js';

export async function renderLogin() {
  setTimeout(bindLoginForm, 0);
  return `
    <main class="ol-login ol-view">
      <section class="ol-login__hero" aria-labelledby="login-title">
        <div class="ol-brand-mark" aria-hidden="true">OL</div>
        <p class="ol-kicker">OpenLucky Admin</p>
        <h1 id="login-title">Operate the appliance from one quiet cockpit.</h1>
        <p class="ol-login__lede">A dependency-free console for status checks, logs, settings, and safe MVP module entry points.</p>
        <ul class="ol-login__note-list" aria-label="Console highlights">
          <li>Native browser modules and clean-room interface code.</li>
          <li>Token-protected API requests with OpenLucky header support.</li>
          <li>Responsive navigation designed for desktops, tablets, and phones.</li>
        </ul>
      </section>
      <section class="ol-card ol-card--raised" aria-labelledby="login-form-title">
        <div class="ol-card__body">
          <div>
            <p class="ol-kicker">Secure entry</p>
            <h2 id="login-form-title">Sign in</h2>
            <p class="ol-message">Use the administrator credentials configured for this OpenLucky instance.</p>
          </div>
          ${hasToken() ? '<p class="ol-message ol-message--success">A session token is already present. Signing in again replaces it.</p>' : ''}
          <form id="login-form" class="ol-form">
            <div class="ol-field">
              <label class="ol-label" for="login-username">Username</label>
              <input class="ol-input" id="login-username" name="username" autocomplete="username" required>
            </div>
            <div class="ol-field">
              <label class="ol-label" for="login-password">Password</label>
              <input class="ol-input" id="login-password" name="password" type="password" autocomplete="current-password" required>
            </div>
            <button class="ol-button ol-button--accent" type="submit">Sign in</button>
            <p id="login-status" class="ol-message" role="status" aria-live="polite"></p>
          </form>
        </div>
      </section>
    </main>`;
}

function bindLoginForm() {
  const form = document.querySelector('#login-form');
  if (!form) {
    return;
  }

  form.addEventListener('submit', async event => {
    event.preventDefault();
    const status = form.querySelector('#login-status');
    const submit = form.querySelector('button[type="submit"]');
    status.className = 'ol-message';
    status.textContent = 'Requesting challenge.';
    submit.disabled = true;

    try {
      const challenge = await apiFetch('/api/login/challenge', { method: 'POST' });
      const credentials = Object.fromEntries(new FormData(form));
      const result = await apiFetch('/api/login', {
        method: 'POST',
        body: {
          username: credentials.username,
          password: credentials.password,
          challengeId: challenge.id
        }
      });
      setToken(result.token, result.expiresAt);
      status.className = 'ol-message ol-message--success';
      status.textContent = 'Signed in. Opening dashboard.';
      window.dispatchEvent(new CustomEvent('openlucky:login'));
    } catch (error) {
      status.className = 'ol-message ol-message--error';
      status.textContent = error.message;
    } finally {
      submit.disabled = false;
    }
  });
}
