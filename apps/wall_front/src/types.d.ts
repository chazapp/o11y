export interface WallMessage {
  id: number
  message: string
  username: string
}

export interface ApiStatus {
  goVersion: string
  version: string
  connectedWS: number
  messagesCount: number
}

declare global {
  interface Window {
    env: {
      API_URL: string
      FARO_URL?: string
      VERSION?: string
    }
  }
}
