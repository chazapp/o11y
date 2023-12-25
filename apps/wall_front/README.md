# Wall_Front

A Front-End application written in React that implements the WallAPI.
Features:

- Form creating Messages for the WallAPI
- WebSocket listening for created messages and drawing them on a canvas element via Konva

## Usage

```bash
$ yarn install
yarn install v1.22.19
[1/4] Resolving packages...
...
$ echo "REACT_APP_API_URL=<http://your.api.endpoint>" >> .env
...
$ yarn start
Compiled successfully!

You can now view wall_front in the browser.

  Local:            http://localhost:3000
  On Your Network:  http://192.168.1.10:3000

Note that the development build is not optimized.
To create a production build, use yarn build.

webpack compiled successfully
No issues found.
```

## End-to-End Tests

A Playwright test suite is available to E2E test both front-end client and API.
Start up a database, API, `npm start`, then run the test suite

```bash
$ npx playwright test
<...>

```

## Production build

This repository contains a Dockerfile that will make a production build
of the application then run it in Nginx. The build step is not aware of
where to find the API that this client implements.
You must mount the `env.js` file at runtime in the application build directory:

```js
// build/env.js
window.env = {
    "API_URL": "<your WallAPI endpoint>"
}
```
