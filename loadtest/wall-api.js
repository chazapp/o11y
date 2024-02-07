import http from 'k6/http'
import { check, sleep } from 'k6'
import { SharedArray } from 'k6/data';

const data = new SharedArray('some quotes', function () {
  // All heavy work (opening and processing big files for example) should be done inside here.
  // This way it will happen only once and the result will be shared between all VUs, saving time and memory.
  const f = JSON.parse(open('./quotes.json'));
  console.log("Created Shared Array !")
  return f; // f must be an array
});



export default function () {
  const element = data[Math.floor(Math.random() * data.length)];
  const { text, from } = element;
  const payload = { username: from, message: text }
  let res = http.post(`${__ENV.API_URL}/message`, JSON.stringify(payload))
  check(res, { 'Success: Post Message': (r) => r.status === 201 })
  res = http.get(`${__ENV.API_URL}/message/${res.json("id")}`)
  check(res, {'Sucess: Get Message': (r) => r.status === 200 })
  sleep(1)
}