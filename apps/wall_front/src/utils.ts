export function getRandomInt (min: number, max: number): number {
  min = Math.ceil(min)
  max = Math.floor(max)
  return Math.floor(Math.random() * (max - min + 1)) + min
}

export function httpUrlToWebSocketUrl (httpUrl: string): string {
  if (httpUrl.startsWith('http://')) {
    return `ws://${httpUrl.slice(7)}`
  } else if (httpUrl.startsWith('https://')) {
    return `wss://${httpUrl.slice(8)}`
  } else {
    throw new Error('Invalid URL. Must start with "http://" or "https://".')
  }
}
