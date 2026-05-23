import { getDivisions } from "../api.js";
import {
  setStatus,
  wireToggleButtons,
  linkFormatter,
  dateTimeFormatter,
  sortDateTime,
  gameModeFormatter
} from "../table-utils.js";

let table;

function normalize(row) {
  return {
    divisionName: row.divisionName || "",
    gameMode: row.gameMode || "",
    mapId: row.mapId || "",
    mapName: row.mapName || "",
    lastChanged: row.lastChanged || "",
  };
}

async function load() {
  setStatus("Loading divisions...");
  const started = performance.now();

  const rows = (await getDivisions()).map(normalize);

  if (!table) {
    table = new Tabulator("#divisions-table", {
      data: rows,
      layout: "fitData",
      movableColumns: true,
      pagination: false,
      /*initialSort: [
        { column: "divisionName", dir: "asc" },
        { column: "gameMode", dir: "asc" },
      ],*/ // FIXME: DDB should store division ID and that should be base for sorting
      columns: [
        {
          title: "Division",
          field: "divisionName",
          sorter: "string",
          headerFilter: "input",
        },
        {
          title: "Mode",
          field: "gameMode",
          sorter: "string",
          formatter: gameModeFormatter,
        },
        {
          title: "Map",
          field: "mapName",
          sorter: "string",
          headerFilter: "input",
          formatter: linkFormatter(
            row => `/map/${encodeURIComponent(row.mapId)}`,
            row => row.mapName || row.mapId,
          ),
        },
        {
          title: "Map ID",
          field: "mapId",
          sorter: "string",
          headerFilter: "input",
          formatter: linkFormatter(
            row => `/map/${encodeURIComponent(row.mapId)}`,
            row => row.mapId,
          ),
        },
        {
          title: "Last Changed",
          field: "lastChanged",
          sorter: sortDateTime,
          formatter: dateTimeFormatter,
        },
      ],
    });

    wireToggleButtons(table);
  } else {
    table.replaceData(rows);
  }

  setStatus(`Loaded ${rows.length} divisions in ${(performance.now() - started).toFixed(1)} ms`);
}

document.getElementById("reload")?.addEventListener("click", () => {
  load().catch(err => {
    console.error(err);
    setStatus(`ERROR: ${err.stack || err.message}`);
  });
});

load().catch(err => {
  console.error(err);
  setStatus(`ERROR: ${err.stack || err.message}`);
});
