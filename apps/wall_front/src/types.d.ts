export interface WallMessage {
    id: number,
    message: string,
    username: string
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