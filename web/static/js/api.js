const TOKEN_KEY = 'openlucky.adminToken';
const TOKEN_EXPIRES_KEY = 'openlucky.adminTokenExpiresAt';

export class APIError extends Error {
  constructor(message, status, payload) {
    super(message);
    this.name = 'APIError';
    this.status = status;
    this.payload = payload;
  }
}

export function getToken() {
  return sessionStorage.getItem(TOKEN_KEY) || '';
}

export function setToken(token, expiresAt = '') {
  sessionStorage.setItem(TOKEN_KEY, token);
  if (expiresAt) {
    sessionStorage.setItem(TOKEN_EXPIRES_KEY, expiresAt);
  } else {
    sessionStorage.removeItem(TOKEN_EXPIRES_KEY);
  }
}

export function getTokenExpiry() {
  return sessionStorage.getItem(TOKEN_EXPIRES_KEY) || '';
}

export function clearToken() {
  sessionStorage.removeItem(TOKEN_KEY);
  sessionStorage.removeItem(TOKEN_EXPIRES_KEY);
}

export function hasToken() {
  return Boolean(getToken());
}

export async function apiFetch(path, options = {}) {
  const headers = new Headers(options.headers || {});
  headers.set('Accept', 'application/json');

  const body = prepareBody(options.body, headers);
  const token = getToken();
  if (token) {
    headers.set('OpenLucky-Admin-Token', token);
    headers.set('Lucky-Admin-Token', token);
  }

  const response = await fetch(path, {
    ...options,
    body,
    headers,
    credentials: options.credentials || 'same-origin'
  });

  const payload = await parsePayload(response);
  if (response.status === 401) {
    clearToken();
    window.dispatchEvent(new CustomEvent('openlucky:unauthorized'));
  }
  if (!response.ok) {
    throw new APIError(errorMessage(payload, response.status), response.status, payload);
  }

  return unwrapPayload(payload);
}

function prepareBody(body, headers) {
  if (body === undefined || body === null) {
    return body;
  }

  const isFormData = typeof FormData !== 'undefined' && body instanceof FormData;
  const isBlob = typeof Blob !== 'undefined' && body instanceof Blob;
  if (typeof body === 'string' || isFormData || isBlob) {
    return body;
  }

  if (!headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }
  return JSON.stringify(body);
}

async function parsePayload(response) {
  const text = await response.text();
  if (!text) {
    return null;
  }
  const contentType = response.headers.get('Content-Type') || response.headers.get('content-type') || '';
  if (!contentType.includes('application/json')) {
    return text;
  }
  try {
    return JSON.parse(text);
  } catch (error) {
    throw new APIError('The server returned invalid JSON.', response.status, text);
  }
}

function unwrapPayload(payload) {
  if (payload && typeof payload === 'object' && 'data' in payload) {
    return payload.data;
  }
  return payload;
}

function errorMessage(payload, status) {
  if (payload && typeof payload === 'object') {
    return payload.error?.message || payload.message || `Request failed with HTTP ${status}`;
  }
  if (typeof payload === 'string' && payload.trim()) {
    return payload;
  }
  return `Request failed with HTTP ${status}`;
}
