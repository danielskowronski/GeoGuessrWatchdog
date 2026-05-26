export function setStatus(text) {
  const el = document.getElementById("status");
  if (el) el.textContent = text;
}

export function toggleColumn(table, field) {
  const col = table.getColumn(field);
  if (!col) return;
  if (col.isVisible()) col.hide();
  else col.show();
}

export function wireToggleButtons(table, root = document) {
  root.querySelectorAll("[data-toggle-column]").forEach(btn => {
    btn.addEventListener("click", () => {
      for (const field of btn.dataset.toggleColumn.split(",")) {
        toggleColumn(table, field.trim());
      }
    });
  });
}

export function fmtCoordinates(lat, lng, options = {}) {
  const { digits = 5, empty = "" } = options;

  if (
    lat === null || lat === undefined || Number.isNaN(Number(lat)) ||
    lng === null || lng === undefined || Number.isNaN(Number(lng))
  ) {
    return empty;
  }

  const latStr = Number(Math.abs(lat)).toFixed(digits);
  const latHemisphere = lat >= 0 ? "N" : "S";

  const lngStr = Number(Math.abs(lng)).toFixed(digits);
  const lngHemisphere = lng >= 0 ? "E" : "W";

  return `${latStr}°${latHemisphere}, ${lngStr}°${lngHemisphere}`;
}
export function coordinatesFormatter(cell) {
  if (!cell) return "";
  console.log(cell.getValue());
  return fmtCoordinates(cell.getValue().lat, cell.getValue().lng);
}

export function fmtDate(value) {
  if (!value) return "";
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return String(value);
  return d.toLocaleString();
}

export function fmtDelta(value, digits = 0) {
  if (value === null || value === undefined || Number.isNaN(value)) return "";
  const n = Number(value);
  const body = digits > 0 ? n.toFixed(digits) : String(n);
  return n > 0 ? `+${body}` : body;
}

export function deltaFormatter(cell, params = { digits: 0, multiplier: 1 }) {
  const digits = params.digits ?? 0;
  const multiplier = params.multiplier ?? 1;
  const v = cell.getValue();
  const el = cell.getElement();

  el.classList.remove("delta-positive", "delta-negative", "delta-zero");

  if (v === null || v === undefined || Number.isNaN(v)) return "";

  if (v > 0) el.classList.add("delta-positive");
  else if (v < 0) el.classList.add("delta-negative");
  else el.classList.add("delta-zero");

  return fmtDelta(v * multiplier, digits);
}

export function addDeltasAscending(rows, fields, timestampField) {
  rows.sort((a, b) => {
    const at = new Date(a[timestampField] || 0).getTime();
    const bt = new Date(b[timestampField] || 0).getTime();
    return at - bt;
  });

  for (let i = 0; i < rows.length; i++) {
    const prev = i > 0 ? rows[i - 1] : null;

    for (const field of fields) {
      const current = rows[i][field.source];
      const previous = prev ? prev[field.source] : null;

      rows[i][field.delta] =
        current !== null &&
          current !== undefined &&
          previous !== null &&
          previous !== undefined
          ? current - previous
          : null;
    }
  }

  rows.reverse();
  return rows;
}

export function linkFormatter(hrefFn, textFn) {
  return function (cell) {
    const row = cell.getRow().getData();
    const a = document.createElement("a");
    a.href = hrefFn(row, cell);
    a.textContent = textFn(row, cell);
    return a;
  };
}

export function sortDateTime(a, b) {
  const at = Date.parse(a || "");
  const bt = Date.parse(b || "");

  if (Number.isNaN(at) && Number.isNaN(bt)) return 0;
  if (Number.isNaN(at)) return -1;
  if (Number.isNaN(bt)) return 1;

  return at - bt;
}

export function fmtDateTime(value) {
  if (!value) return "";

  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return String(value);

  return new Intl.DateTimeFormat('en-GB', {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  }).format(d);
}
export function fmtDateTimeDateOnly(value) {
  if (!value) return "";

  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return String(value);

  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(d);
}

export function gameModeFormatter(cell) {
  const v = cell.getValue();
  if (v === "standardDuels") return "Moving";
  if (v === "noMoveDuels") return "No Move";
  if (v === "nmpzDuels") return "NMPZ";
  return String(v);
}

export function dateTimeFormatter(cell) {
  return fmtDateTime(cell.getValue());
}
export function dateFormatter(cell) {
  return fmtDateTimeDateOnly(cell.getValue());
}

export function weekdayFromNumber(value, offset = 0) {
  if (value === null || value === undefined) return "";
  const n = Number(value);
  if (Number.isNaN(n) || n < 1 || n > 7) return String(value);
  const weekdays = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];
  return weekdays[(n - 1 + offset + 7) % 7];
}
export function weekdayNumberFromDate(value) {
  if (!value) return null;
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return null;
  return d.getDay() === 0 ? 7 : d.getDay();
}

export function numberText(value, options = {}) {
  const {
    digits = 0,
    multiplier = 1,
    suffix = "",
    empty = "",
    fillSpaces = 0,
  } = options;
  if (value === null || value === undefined || Number.isNaN(Number(value))) {
    return empty;
  }
  const v = Number(value) * multiplier;
  const body = digits > 0 ? v.toFixed(digits) : String(Math.round(v));
  return `${body.padStart(fillSpaces, " ")}${suffix}`;
}

export function deltaText(value, options = {}) {
  const {
    digits = 0,
    multiplier = 1,
    suffix = "",
    showZero = false,
    empty = "",
    intDigits = 2,
  } = options;

  if (value === null || value === undefined || Number.isNaN(Number(value))) {
    return empty;
  }

  const v = Number(value) * multiplier;

  if (v === 0 && !showZero) {
    return empty;
  }

  const body = digits > 0 ? v.toFixed(digits) : String(Math.round(v));
  const sign = v > 0 ? "+" : "";

  return (`${sign}${body}${suffix}`).padStart(1 + digits + suffix.length + intDigits, " ");
}

export function deltaClass(value) {
  if (value === null || value === undefined || Number.isNaN(Number(value))) {
    return "";
  }

  const v = Number(value);

  if (v > 0) return "delta-positive";
  if (v < 0) return "delta-negative";
  return "delta-zero";
}

export function valueWithDeltaFormatter(options = {}) {
  const {
    deltaField,
    valueDigits = 0,
    deltaDigits = 0,
    valueMultiplier = 1,
    deltaMultiplier = 1,
    suffix = "",
    deltaSuffix = suffix,
    showZeroDelta = false,
    separator = " ",
    fillSpaces = 0
  } = options;

  return function (cell) {
    const row = cell.getRow().getData();

    const value = cell.getValue();
    const delta = row[deltaField];

    const valueOut = numberText(value, {
      digits: valueDigits,
      multiplier: valueMultiplier,
      suffix,
      fillSpaces: fillSpaces,
    });

    if (!valueOut) {
      return "";
    }

    const deltaOut = deltaText(delta, {
      digits: deltaDigits,
      multiplier: deltaMultiplier,
      suffix: deltaSuffix,
      showZero: showZeroDelta,
    });

    if (!deltaOut) {
      return valueOut;
    }

    const cls = deltaClass(delta);

    return `${valueOut}${separator}<span class="${cls}">${deltaOut}</span>`;
  };
}

export function validCoordPoint(point) {
  return (
    point &&
    Number.isFinite(Number(point.lat)) &&
    Number.isFinite(Number(point.lng))
  );
}

export function bboxGeoJSONUrl(boundsMin, boundsMax, name = "Bounds") {
  if (!validCoordPoint(boundsMin) || !validCoordPoint(boundsMax)) {
    return "";
  }

  const minLat = Math.min(Number(boundsMin.lat), Number(boundsMax.lat));
  const maxLat = Math.max(Number(boundsMin.lat), Number(boundsMax.lat));
  const minLng = Math.min(Number(boundsMin.lng), Number(boundsMax.lng));
  const maxLng = Math.max(Number(boundsMin.lng), Number(boundsMax.lng));

  const geojson = {
    type: "FeatureCollection",
    features: [
      {
        type: "Feature",
        properties: { name },
        geometry: {
          type: "Polygon",
          coordinates: [[
            [minLng, minLat],
            [maxLng, minLat],
            [maxLng, maxLat],
            [minLng, maxLat],
            [minLng, minLat],
          ]],
        },
      },
    ],
  };

  return `https://geojson.io/#data=data:application/json,${encodeURIComponent(JSON.stringify(geojson))}`;
}
export function escapeHTML(value) {

  if (value === null || value === undefined) return "";

  return String(value).replace(/[&<>"']/g, ch => ({

    "&": "&amp;",

    "<": "&lt;",

    ">": "&gt;",

    '"': "&quot;",

    "'": "&#39;",

  }[ch]));

}