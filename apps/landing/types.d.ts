import "react-router";

declare global {
  interface Window {
    env: {
      AUTH_URL: string
      VERSION?: string
    }
  }
}
