# API definition

Api is not officially documented, but [github.com/teamcoltra/geoguessr-api-docs](https://github.com/teamcoltra/geoguessr-api-docs) contains great effort at reverse engineering.

## Auth

All endpoints, even read-only ones that contain public data like competitive game settings for everyone require auth.

Auth requires only `_ncfa` cookie. It is valid for a year and is set by server in normal login flow. This step cannot be automated (to only use email and password), as log in requires solving captcha. Therefore, value of that cookie is treated as API key and is required as config input (in file or env variable). User must extract it from browser themselves. For safety (to avoid ban), it's better to create fresh free account - it'll still have access to info about paid features like maps. To gather data about user (like game stats), it should be enough to add said account to friends.

For extra safety, optional proxy via any anonymizing VPN is possible. Apparently, GeoGuessr doesn't care if country changes.
