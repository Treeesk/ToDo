/**
 * API client for the Go notes backend.
 *
 * The backend stores auth in HttpOnly cookies. Browser JavaScript cannot read
 * those cookies, so every request includes credentials and protected note
 * requests try one refresh before returning an unauthorized error.
 */

export class ApiError extends Error {
  constructor(message, status, payload = null) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.payload = payload;
  }
}

export function createApiClient({ baseUrl = "", fetcher = globalThis.fetch } = {}) {
  const request = (path, options = {}) => requestJson(fetcher, baseUrl, path, options);

  async function protectedRequest(path, options = {}) {
    try {
      return await request(path, options);
    } catch (error) {
      if (!(error instanceof ApiError) || error.status !== 401 || options.skipRefresh) {
        throw error;
      }

      await request("/api/refresh/", { method: "POST", skipRefresh: true });
      return request(path, { ...options, skipRefresh: true });
    }
  }

  return {
    getNotes: () => protectedRequest("/api/"),
    addNote: (text) => protectedRequest("/api/add/", { method: "POST", body: { text } }),
    editNote: (id, text) => protectedRequest("/api/edit/", { method: "PUT", body: { id, text } }),
    deleteNote: (id) => protectedRequest("/api/del/", { method: "DELETE", body: { id } }),
    login: (login, password) => request("/api/login/", { method: "POST", body: { login, password } }),
    register: (login, password) => request("/api/register/", { method: "POST", body: { login, password } }),
    logout: () => request("/api/logout/", { method: "POST" }),
  };
}

/**
 * Performs one JSON request and normalizes empty successful backend responses
 * into null.
 */
export async function requestJson(fetcher, baseUrl, path, options = {}) {
  const headers = new Headers(options.headers || {});
  const hasBody = Object.hasOwn(options, "body");

  if (hasBody) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetcher(`${baseUrl}${path}`, {
    method: options.method || "GET",
    credentials: "include",
    headers,
    body: hasBody ? JSON.stringify(options.body) : undefined,
  });

  const text = await response.text();
  const payload = text ? parseJson(text) : null;

  if (!response.ok) {
    throw new ApiError(payload?.message || "Не удалось выполнить запрос.", response.status, payload);
  }

  return payload;
}

function parseJson(text) {
  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}
