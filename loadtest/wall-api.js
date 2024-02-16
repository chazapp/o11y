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

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 1000,
      timeUnit: '1s', // 1000 iterations per second, i.e. 1000 RPS
      duration: '30s',
      preAllocatedVUs: 200, // how large the initial pool of VUs would be
      maxVUs: 200, // if the preAllocatedVUs are not enough, we can initialize more
    },
  },
};

export default function () {
  const element = data[Math.floor(Math.random() * data.length)];
  const { text, from } = element;
  const payload = { username: from, message: text }
  let res = http.post(`${__ENV.API_URL}/message`, JSON.stringify(payload))
  check(res, { 'Success: Post Message': (r) => r.status === 201 })
}