import { createContext } from "react";

type AppContext = {
  authUrl: string,
}

let AUTH_URL = import.meta.env.VITE_AUTH_URL;

if (AUTH_URL === undefined) {
  AUTH_URL = window.env.AUTH_URL;
}

const ctx: AppContext = {
  authUrl: AUTH_URL 
}

export const appContext = createContext<AppContext>(ctx);