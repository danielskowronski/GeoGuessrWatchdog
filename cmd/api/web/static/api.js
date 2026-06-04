export async function apiGet(path) {
  const res = await fetch(path, {
    headers: {
      "Accept": "application/json"
    }
  });

  if (!res.ok) {
    let detail = `${res.status} ${res.statusText}`;
    try {
      const body = await res.json();
      detail = body.detail || body.title || detail;
    } catch (_) { }
    throw new Error(detail);
  }

  return await res.json();
}

export async function getDivisions() {
  const body = await apiGet("/api/divisions");
  return body.divisions || [];
}

export async function getUsers() {
  const body = await apiGet("/api/users");
  return body.users || {};
}

export async function getFetchStatuses() {
  const body = await apiGet("/api/fetch_statuses");
  return body.fetches || [];
}

export async function getUserHistory(userId) {
  const body = await apiGet(`/api/user/${encodeURIComponent(userId)}`);
  return body.entries || [];
}

export async function getMapHistory(mapId) {
  const body = await apiGet(`/api/map/${encodeURIComponent(mapId)}`);
  return {
    name: body.name || "",
    description: body.description || "",
    url: body.url || "",
    entries: body.entries || [],
  };
}
