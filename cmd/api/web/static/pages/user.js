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


        { title: "Countries<br />Best", field: "bestCountries", sorter: "string", visible: false, cssClass: "wrap-cell" },
        { title: "Countries<br />Worst", field: "worstCountries", sorter: "string", visible: false, cssClass: "wrap-cell" },
      ],
    });

    wireToggleButtons(table);
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
