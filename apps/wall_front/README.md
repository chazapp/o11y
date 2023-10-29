# Wall_Front

A Front-End application written in React that implements the WallAPI.
Features:
 - Form creating Messages for the WallAPI
 - WebSocket listening for created messages and drawing them on a canvas element via Konva


## Usage

```
$ yarn install
$ echo "REACT_APP_API_URL=<http://your.api.endpoint>" >> .env
$ yarn start
```


## Production build
This repository contains a Dockerfile that will make a production build
of the application then run it in Nginx. The build step is not aware of
where to find the API that this client implements.
You must mount the `env.js` file at runtime in the application build directory:

```
// build/env.js
window.env = {
    "API_URL": "<your WallAPI endpoint>"
}
```
