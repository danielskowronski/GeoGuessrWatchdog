import { getFetchStatuses } from "../api.js";
import {
  setStatus,
  linkFormatter,
  dateTimeFormatter,
  sortDateTime,
  dateTimeRelativeFormatter,
  wireToggleButtons,
} from "../table-utils.js";

let table;

async function load() {
  setStatus("Loading fetch statuses...");
  const rows = await getFetchStatuses();

  if (!table) {
    table = new Tabulator("#fetch-statuses-table", {
      data: rows,
      layout: "fitData",
      movableColumns: true,
      pagination: true,
      paginationSize: 50,
      initialSort: [{ column: "name", dir: "asc" }],
      columns: [
        {
          title: "Type",
          field: "fetchType",
          sorter: "string",
          headerFilter: "input",
          visible: false,
        },
        {
          title: "Description",
          field: "description",
          sorter: "string",
          headerFilter: "input",
        },
        {
          title: "Last Success",
          field: "lastSuccess",
          sorter: sortDateTime,
          formatter: dateTimeFormatter,
        },
        {
          title: "Last Success",
          field: "lastSuccess",
          sorter: sortDateTime,
          formatter: dateTimeRelativeFormatter,
          formatterParams: {
            thresholdHours: 24,
          },
        },
      ],
    });
    wireToggleButtons(table);
  } else {
    table.replaceData(rows);
  }

  setStatus(`Loaded ${rows.length} fetch statuses.`);
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
