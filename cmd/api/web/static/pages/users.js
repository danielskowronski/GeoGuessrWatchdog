import { getUsers } from "../api.js";
import { 
  setStatus,
  linkFormatter
} from "../table-utils.js";

let table;

function normalizeUsers(userMap) {
  return Object.entries(userMap).map(([id, name]) => ({ id, name }));
}

async function load() {
  setStatus("Loading users...");
  const rows = normalizeUsers(await getUsers());

  if (!table) {
    table = new Tabulator("#users-table", {
      data: rows,
      layout: "fitData",
      movableColumns: true,
      pagination: true,
      paginationSize: 50,
      initialSort: [{ column: "name", dir: "asc" }],
      columns: [
        {
          title: "Name",
          field: "name",
          sorter: "string",
          headerFilter: "input",
          formatter: linkFormatter(
            row => `/user/${encodeURIComponent(row.id)}`,
            row => row.name || row.id,
          ),
        },
        {
          title: "User ID",
          field: "id",
          sorter: "string",
          headerFilter: "input",
          formatter: linkFormatter(
            row => `/user/${encodeURIComponent(row.id)}`,
            row => row.id,
          ),
        },
      ],
    });
  } else {
    table.replaceData(rows);
  }

  setStatus(`Loaded ${rows.length} users`);
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
