import { getMapHistory } from "../api.js";
import {
  setStatus,
  addDeltasAscending,
  wireToggleButtons,
  sortDateTime,
  dateTimeFormatter,
  coordinatesFormatter,
  bboxGeoJSONUrl,
  valueWithDeltaFormatter,
  escapeHTML,
} from "../table-utils.js";

let table;

function pointOrNull(lat, lng) {
  const latNum = Number(lat);
  const lngNum = Number(lng);

  if (!Number.isFinite(latNum) || !Number.isFinite(lngNum)) {
    return null;
  }

  return {
    lat: latNum,
    lng: lngNum,
  };
}

function normalizeEntry(entry) {
  const info = entry.info || {};

  const boundsMin = pointOrNull(info.boundsMinLat, info.boundsMinLng);
  const boundsMax = pointOrNull(info.boundsMaxLat, info.boundsMaxLng);

  const name = info.name || "";
  const timestamp = entry.timestamp || "";
  const updatedAt = info.updatedAt || "";

  return {
    timestamp,
    updatedAt,

    id: info.id || "",
    name,
    description: info.description || "",

    coordinateCount: info.coordinateCount ?? null,
    maxErrorDistance: info.maxErrorDistance ?? null,

    boundsMin,
    boundsMax,

    boundsMinLat: boundsMin?.lat ?? null,
    boundsMinLng: boundsMin?.lng ?? null,
    boundsMaxLat: boundsMax?.lat ?? null,
    boundsMaxLng: boundsMax?.lng ?? null,

    geojsonUrl: bboxGeoJSONUrl(
      boundsMin,
      boundsMax,
      name || timestamp || "Bounds",
    ),
  };
}

function prepareRows(entries) {
  const rows = entries.map(normalizeEntry);

  return addDeltasAscending(rows, [
    { source: "coordinateCount", delta: "coordinateCountDelta" },
    { source: "maxErrorDistance", delta: "maxErrorDistanceDelta" },
  ], "timestamp");
}

async function load() {
  const mapId = window.GGWD_PAGE?.id || "";
  setStatus(`Loading map ${mapId}...`);

  const result = await getMapHistory(mapId);

  document.getElementById("map-description").textContent = result.description || "";

  const mapUrl = document.getElementById("map-url");
  mapUrl.href = result.url || "#";
  mapUrl.textContent = result.url || "";

  document.querySelector("h1").textContent = result.name || mapId;

  const rows = prepareRows(result.entries || []);

  if (!table) {
    table = new Tabulator("#map-history-table", {
      data: rows,
      layout: "fitData",
      movableColumns: true,
      pagination: true,
      paginationSize: 50,

      initialSort: [
        { column: "timestamp", dir: "desc" },
      ],

      columns: [
        {
          title: "Fetch<br />Timestamp",
          field: "timestamp",
          sorter: sortDateTime,
          formatter: dateTimeFormatter,
          visible: false,
        },
        {
          title: "Updated<br />At",
          field: "updatedAt",
          sorter: sortDateTime,
          formatter: dateTimeFormatter,
        },
        {
          title: "Locations<br />Count",
          field: "coordinateCount",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "coordinateCountDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
        },
        {
          title: "Max Error<br />Distance",
          field: "maxErrorDistance",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "maxErrorDistanceDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "Bounds",
          field: "geojsonUrl",
          sorter: "string",
          formatter: cell => {
            const url = cell.getValue();
            if (!url) return "";

            return `<a href="${escapeHTML(url)}" target="_blank" rel="noreferrer">GeoJSON</a>`;
          },
        },
        {
          title: "Bounds<br />Minimum",
          field: "boundsMin",
          formatter: coordinatesFormatter,
          visible: false,
        },
        {
          title: "Bounds<br />Maximum",
          field: "boundsMax",
          formatter: coordinatesFormatter,
          visible: false,
        },
      ],
    });

    wireToggleButtons(table);
  } else {
    table.replaceData(rows);
  }

  setStatus(`Loaded ${rows.length} map history entries`);
}

document.getElementById("reload")?.addEventListener("click", () => {
  load().catch(err => {
    console.error(err);
    setStatus(`ERROR: ${err.stack || err.message}`);
  });
});

requestAnimationFrame(() => {
  load().catch(err => {
    console.error(err);
    setStatus(`ERROR: ${err.stack || err.message}`);
  });
});