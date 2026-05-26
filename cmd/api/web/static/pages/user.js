import { getUserHistory } from "../api.js";
import {
  setStatus,
  wireToggleButtons,
  deltaFormatter,
  addDeltasAscending,
  sortDateTime,
  dateTimeFormatter,
  dateFormatter,
  weekdayFromNumber,
  weekdayNumberFromDate,
  valueWithDeltaFormatter,
} from "../table-utils.js";

let table;

function prepareRows(rows) {
  rows = rows.map(row => ({
    ...row,
    fetchWeekday: weekdayNumberFromDate(row.fetchTimestamp),
  }));

  return addDeltasAscending(rows, [
    { source: "ratingOverall", delta: "ratingOverallDelta" },
    { source: "ratingMoving", delta: "ratingMovingDelta" },
    { source: "ratingNomove", delta: "ratingNomoveDelta" },
    { source: "ratingNmpz", delta: "ratingNmpzDelta" },
    { source: "guessedFirst", delta: "guessedFirstDelta" },

    { source: "rankedSoloMovingGames", delta: "rankedSoloMovingGamesDelta" },
    { source: "rankedSoloMovingWins", delta: "rankedSoloMovingWinsDelta" },
    { source: "rankedSoloNomoveGames", delta: "rankedSoloNomoveGamesDelta" },
    { source: "rankedSoloNomoveWins", delta: "rankedSoloNomoveWinsDelta" },
    { source: "rankedSoloNmpzGames", delta: "rankedSoloNmpzGamesDelta" },
    { source: "rankedSoloNmpzWins", delta: "rankedSoloNmpzWinsDelta" },
  ], "fetchTimestamp");

  // TODO: add daily win ratios
}

function dailyRatioFormatter({ winsDeltaField = "", gamesDeltaField = "", digits = 0 } = {}) {
  return function (cell) {
    const row = cell.getRow().getData();
    const el = cell.getElement();

    el.classList.remove("delta-positive", "delta-negative", "delta-zero");

    const wins = Number(row[winsDeltaField]);
    const games = Number(row[gamesDeltaField]);

    if (!Number.isFinite(wins) || !Number.isFinite(games) || games === 0) {
      return "";
    }

    const ratio = wins / games;

    if (ratio >= 0.5) {
      el.classList.add("delta-positive");
    } else {
      el.classList.add("delta-negative");
    }

    const padDigits = digits > 0 ? digits + 4 : 3; // 4 for "100.0", 3 for "100"
    return (ratio * 100).toFixed(digits).padStart(padDigits, " ") + "%";
  };
}

async function load() {
  const userId = window.GGWD_PAGE?.id || "";
  setStatus(`Loading user ${userId}...`);

  const rows = prepareRows(await getUserHistory(userId));

  if (!table) {
    table = new Tabulator("#user-history-table", {
      data: rows,
      layout: "fitData",
      movableColumns: true,
      pagination: true,
      paginationSize: 50,
      initialSort: [{ column: "fetchTimestamp", dir: "asc" }],
      columns: [
        { title: "Fetch ID", field: "fetchID", sorter: "number", visible: false },
        {
          title: "Fetch<br />Date",
          field: "fetchTimestamp",
          sorter: sortDateTime,
          formatter: dateFormatter,
        },
        {
          title: "Prev<br />Day",
          field: "fetchWeekday",
          sorter: "number",
          formatter: cell => {
            const v = cell.getValue();
            if (v === null || v === undefined) return "";
            return weekdayFromNumber(v, -1);
          },
          headerFilter: "input",
          headerFilterFunc: (headerValue, rowValue) => {
            if (!headerValue) return true;
            if (rowValue === null || rowValue === undefined) return false;
            const needle = String(headerValue).toLowerCase();
            const weekday = weekdayFromNumber(rowValue, -1).toLowerCase();
            return weekday.includes(needle);
          }
        },
        { title: "Division", field: "divisionName", sorter: "string", headerFilter: "input", visible: false },
        {
          title: "ELO<br />Overall",
          field: "ratingOverall",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "ratingOverallDelta",
            valueDigits: 0,
            deltaDigits: 0,
            fillSpaces: 4
          }),
        },
        {
          title: "ELO<br />Moving",
          field: "ratingMoving",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "ratingMovingDelta",
            valueDigits: 0,
            deltaDigits: 0,
            fillSpaces: 4
          }),
          visible: false,
        },
        {
          title: "ELO<br />No Move",
          field: "ratingNomove",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "ratingNomoveDelta",
            valueDigits: 0,
            deltaDigits: 0,
            fillSpaces: 4
          }),
          visible: false,
        },
        {
          title: "ELO<br />NMPZ",
          field: "ratingNmpz",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "ratingNmpzDelta",
            valueDigits: 0,
            deltaDigits: 0,
            fillSpaces: 4
          }),
          visible: false,
        },
        {
          title: "Guessed<br />First",
          field: "guessedFirst",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "guessedFirstDelta",
            valueMultiplier: 100,
            deltaMultiplier: 100,
            valueDigits: 2,
            deltaDigits: 2,
            suffix: "%",
            fillSpaces: 5,
          }),
        },

        {
          title: "# Games<br />Ranked<br />Moving",
          field: "rankedSoloMovingGames",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedSoloMovingGamesDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "# Wins<br />Ranked<br />Moving",
          field: "rankedSoloMovingWins",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedSoloMovingWinsDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "Ratio<br />Ranked<br />Moving",
          field: "rankedSoloMovingGamesDaily",
          sorter: "number",
          formatter: dailyRatioFormatter({
            winsDeltaField: "rankedSoloMovingWinsDelta",
            gamesDeltaField: "rankedSoloMovingGamesDelta",
          }),
          visible: false,
        },
        {
          title: "# Games<br />Ranked<br />No Move",
          field: "rankedSoloNomoveGames",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedSoloNomoveGamesDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "# Wins<br />Ranked<br />No Move",
          field: "rankedSoloNomoveWins",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedSoloNomoveWinsDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "Ratio<br />Ranked<br />No Move",
          field: "rankedSoloNomoveGamesDaily",
          sorter: "number",
          formatter: dailyRatioFormatter({
            winsDeltaField: "rankedSoloNomoveWinsDelta",
            gamesDeltaField: "rankedSoloNomoveGamesDelta",
          }),
          visible: false,
        },

        {
          title: "# Games<br />Ranked<br />NMPZ",
          field: "rankedSoloNmpzGames",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedSoloNmpzGamesDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "# Wins<br />Ranked<br />NMPZ",
          field: "rankedSoloNmpzWins",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedSoloNmpzWinsDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "Ratio<br />Ranked<br />NMPZ",
          field: "rankedSoloNmpzGamesDaily",
          sorter: "number",
          formatter: dailyRatioFormatter({
            winsDeltaField: "rankedSoloNmpzWinsDelta",
            gamesDeltaField: "rankedSoloNmpzGamesDelta",
          }),
          visible: false,
        },

        {
          title: "# Games<br />Team<br />Moving",
          field: "rankedTeamMovingGames",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedTeamMovingGamesDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "# Wins<br />Team<br />Moving",
          field: "rankedTeamMovingWins",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedTeamMovingWinsDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "Ratio<br />Team<br />Moving",
          field: "rankedTeamMovingGamesDaily",
          sorter: "number",
          formatter: dailyRatioFormatter({
            winsDeltaField: "rankedTeamMovingWinsDelta",
            gamesDeltaField: "rankedTeamMovingGamesDelta",
          }),
          visible: false,
        },

        {
          title: "# Games<br />Team<br />No Move",
          field: "rankedTeamNomoveGames",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedTeamNomoveGamesDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "# Wins<br />Team<br />No Move",
          field: "rankedTeamNomoveWins",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedTeamNomoveWinsDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "Ratio<br />Team<br />No Move",
          field: "rankedTeamNomoveGamesDaily",
          sorter: "number",
          formatter: dailyRatioFormatter({
            winsDeltaField: "rankedTeamNomoveWinsDelta",
            gamesDeltaField: "rankedTeamNomoveGamesDelta",
          }),
          visible: false,
        },

        {
          title: "# Games<br />Team<br />NMPZ",
          field: "rankedTeamNmpzGames",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedTeamNmpzGamesDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "# Wins<br />Team<br />NMPZ",
          field: "rankedTeamNmpzWins",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "rankedTeamNmpzWinsDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "Ratio<br />Team<br />NMPZ",
          field: "rankedTeamNmpzGamesDaily",
          sorter: "number",
          formatter: dailyRatioFormatter({
            winsDeltaField: "rankedTeamNmpzWinsDelta",
            gamesDeltaField: "rankedTeamNmpzGamesDelta",
          }),
          visible: false,
        },

        {
          title: "# Games<br />Unranked<br />Moving",
          field: "unrankedSoloMovingGames",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "unrankedSoloMovingGamesDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "# Wins<br />Unranked<br />Moving",
          field: "unrankedSoloMovingWins",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "unrankedSoloMovingWinsDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "Ratio<br />Unranked<br />Moving",
          field: "unrankedSoloMovingGamesDaily",
          sorter: "number",
          formatter: dailyRatioFormatter({
            winsDeltaField: "unrankedSoloMovingWinsDelta",
            gamesDeltaField: "unrankedSoloMovingGamesDelta",
          }),
          visible: false,
        },

        {
          title: "# Games<br />Unranked<br />No Move",
          field: "unrankedSoloNomoveGames",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "unrankedSoloNomoveGamesDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "# Wins<br />Unranked<br />No Move",
          field: "unrankedSoloNomoveWins",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "unrankedSoloNomoveWinsDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "Ratio<br />Unranked<br />No Move",
          field: "unrankedSoloNomoveGamesDaily",
          sorter: "number",
          formatter: dailyRatioFormatter({
            winsDeltaField: "unrankedSoloNomoveWinsDelta",
            gamesDeltaField: "unrankedSoloNomoveGamesDelta",
          }),
          visible: false,
        },

        {
          title: "# Games<br />Unranked<br />NMPZ",
          field: "unrankedSoloNmpzGames",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "unrankedSoloNmpzGamesDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "# Wins<br />Unranked<br />NMPZ",
          field: "unrankedSoloNmpzWins",
          sorter: "number",
          formatter: valueWithDeltaFormatter({
            deltaField: "unrankedSoloNmpzWinsDelta",
            valueDigits: 0,
            deltaDigits: 0,
          }),
          visible: false,
        },
        {
          title: "Ratio<br />Unranked<br />NMPZ",
          field: "unrankedSoloNmpzGamesDaily",
          sorter: "number",
          formatter: dailyRatioFormatter({
            winsDeltaField: "unrankedSoloNmpzWinsDelta",
            gamesDeltaField: "unrankedSoloNmpzGamesDelta",
          }),
          visible: false,
        },


        { title: "Countries<br />Best", field: "bestCountries", sorter: "string", visible: false, cssClass: "wrap-cell" },
        { title: "Countries<br />Worst", field: "worstCountries", sorter: "string", visible: false, cssClass: "wrap-cell" },
      ],
    });

    wireToggleButtons(table);
    setTimeout(() => initTable(table), 1); 
  } else {
    table.replaceData(rows);
  }

  document.querySelector("h1").innerHTML = `User history: <a href="https://www.geoguessr.com/user/${userId}">${userId}</a>`;
  // TODO: get user name and show that instead of ID
  setStatus(`Loaded ${rows.length} user history rows`);
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

function readTableSettings() {
  const defaultSettings = {
    "gamesAndWins": {
      "rankedSoloMoving": true,
      "rankedSoloNomove": false,
      "rankedSoloNmpz": false,
      "rankedTeamMoving": false,
      "rankedTeamNomove": false,
      "rankedTeamNmpz": false,
      "unrankedSoloMoving": false,
      "unrankedSoloNomove": false,
      "unrankedSoloNmpz": false,
    },
    "rating": {
      "Moving": true,
      "Nomove": false,
      "Nmpz": false,
    }
  };
  const settings = localStorage.getItem("userHistoryTableSettings");
  if (!settings) {
    localStorage.setItem("userHistoryTableSettings", JSON.stringify(defaultSettings));
    return defaultSettings;
  }
  else {
    try {
      return JSON.parse(settings);
    } catch (err) {
      localStorage.setItem("userHistoryTableSettings", JSON.stringify(defaultSettings));
      return defaultSettings;
    }
  }
}
function updateTableSettings(table, group, field, value) {
  const settings = readTableSettings();
  settings[group][field] = value;
  localStorage.setItem("userHistoryTableSettings", JSON.stringify(settings));
  if (group === "rating") {
    const col = table.getColumn("rating" + field);
    if (col) {
      if (value) col.show();
      else col.hide();
    }
  }
  else if (group === "gamesAndWins") {
    const columnSuffixes = ["Games", "Wins", "GamesDaily"];
    for (const suffix of columnSuffixes) {
      const col = table.getColumn(field + suffix);
      if (col) {
        if (value) col.show();
        else col.hide();
      }
    }
  }
}
function initTable(table) {
  const settings = readTableSettings();

  document.querySelectorAll("#userModesTable [data-select-type]").forEach(checkbox => {
    const group = checkbox.dataset.selectGroup;
    const type = checkbox.dataset.selectType;
    checkbox.checked = settings[group][type] || false;
    updateTableSettings(table, group, type, checkbox.checked);
    checkbox.addEventListener("change", () => {
      updateTableSettings(table, group, type, checkbox.checked);
    });
  });
}
