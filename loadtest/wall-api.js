import http from 'k6/http'

import { check, sleep } from 'k6'


export default function () {
  const data = { username: 'shad', message: "Hello World !" }
  let res = http.post('https://wall-api.local/message', JSON.stringify(data))
  check(res, { 'Success: Message': (r) => r.status === 201 })
  sleep(0.3)
}