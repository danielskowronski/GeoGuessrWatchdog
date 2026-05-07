# Competitive Duels - Maps

Goal of watchdog: detect when GeoGuessr changes map assigned to player current division and when that map gets updated. This is needed, because GeoGuessr doesn't announce those changes and maps differ in play style. Map owners are not prohibited from changing them during competitive week and some of them update them daily, which impacts play.

Watcher fetches data about all division+mode maps, notifier is responsible for filtering info.

## Watcher

### API

Two calls are needed to fetch data:

1. `https://www.geoguessr.com/api/v4/ranked-system/divisions`

this returns:

```json
{
  "divisions": [
    {
      "divisionNumber": 2,
      "divisionRank": 2,
      "tier": "Master",
      "name": "Master I",
      "pointCategories": {},
      "promotionPercentage": 0.1,
      "demotionPercentage": 0.2,
      "doublePromotionPercentage": 0,
      "bucketSortedBy": "Rating",
      "gameModes": [
        "StandardDuels",
        "NoMoveDuels",
        "NmpzDuels"
      ],
      "numberOfGames": 0,
      "maps": {
        "standardDuels": {
          "mapId": "676340ae2f718dbabdf30331",
          "mapName": "An Informed World"
        },
        "noMoveDuels": {
          "mapId": "65c86935d327035509fd616f",
          "mapName": "A Rainbolt World"
        },
        "nmpzDuels": {
          "mapId": "643dbc7ccc47d3a344307998",
          "mapName": "An Arbitrary Rural World"
        }
      }
    },
    // rest of divisions
  ]
}
```

- `divisions` is list of objects, where `name` is unique (this is useful for later filtering)
- `gameModes` list defines which modes will appear in `maps` dict (this is useful for later filtering)
- each game mode in `maps` dict have `mapId` needed for next step

2. `https://www.geoguessr.com/api/maps/{mapId}`

this returns:

```json
{
  "id": "676340ae2f718dbabdf30331",
  "name": "An Informed World",
  "slug": "676340ae2f718dbabdf30331",
  "description": "Locations maximum 2 km from an area considered realitively dense in coverage. 109034 locations. See https://bit.ly/slashp for distribution and other maps.",
  "url": "/maps/676340ae2f718dbabdf30331",
  "playUrl": "/676340ae2f718dbabdf30331/play",
  "published": true,
  "banned": false,
  "images": {
    "backgroundLarge": null,
    "incomplete": true
  },
  "bounds": {
    "min": {
      "lat": -54.84040253403922,
      "lng": -170.73234474472065
    },
    "max": {
      "lat": 78.21578286107786,
      "lng": 178.38109969334502
    }
  },
  "customCoordinates": null,
  "coordinateCount": "100K+",
  "regions": null,
  "creator": {
    // info about creator
  },
  "createdAt": "2024-12-18T21:37:52.9780000Z",
  "updatedAt": "2026-05-07T02:15:00.2040000Z",
  "numFinishedGames": 141045,
  "likedByUser": null,
  "averageScore": 17495,
  "avatar": {
    "background": "night",
    "decoration": "cactus",
    "ground": "water",
    "landscape": "skyline"
  },
  "difficulty": "EASY",
  "difficultyLevel": 2,
  "highscore": null,
  "isUserMap": true,
  "highlighted": false,
  "deleted": false,
  "free": false,
  "panoramaProvider": "StreetView",
  "inExplorerMode": false,
  "maxErrorDistance": 18499075,
  "likes": 2388,
  "locationSelectionMode": 1,
  "tags": [
    "Generated",
    "Official Cov.",
    "Global",
    "Urban"
  ],
  "collaborators": null,
  "flair": 0,
  "mapSize": {
    "coordinateCount": 109034,
    "label": "L"
  }
}
```

only following data is useful:

- `id` = `mapId`
- `name` - map name as appears in GeoGuessr, mandatory for reasonable notifications
- `updatedAt` - ISO time format of last time it was updated, mandatory for notifications about updates
- `mapSize.coordinateCount` - exact numer of locations, may be interesting to compare between map updates
- `bounds` - in future, may make sense to compare

### Flow

TBD

### Data structure in code

TBD

#### Key / watcher selector

TBD

Both division name and game mode key must be stored. It cannot be safely simplified to `C`/`M1`/`M2` because GeoGuessr may rename divisions at any given time. Game modes also cannot be assumed to be set in stone forever. 

This data is not expected to change every minute, so any decent DB engine will handle it well for years.

#### Map info

TBD

### DB schema

TBD

---

## Notifier

TBD
