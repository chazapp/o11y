import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import App from './App'
import reportWebVitals from './reportWebVitals'
import { getWebInstrumentations, initializeFaro, ReactIntegration } from '@grafana/faro-react'
import { TracingInstrumentation } from '@grafana/faro-web-tracing'
import axios from 'axios'

let API_URL = process.env.REACT_APP_API_URL
if (API_URL === undefined) {
  API_URL = window.env?.API_URL
}

let FARO_URL = process.env.REACT_APP_FARO_URL
if (FARO_URL === undefined) {
  FARO_URL = window.env?.FARO_URL
}

let VERSION = 'dev'
if (window.env?.VERSION !== undefined) {
  VERSION = window.env.VERSION
}

if (API_URL === undefined) {
  throw new Error('API URL is undefined !')
}
axios.defaults.baseURL = API_URL

if (FARO_URL !== undefined) {
  initializeFaro({
    url: `${FARO_URL}/collect`,
    app: {
      name: 'wall-browser',
      version: VERSION
    },
    instrumentations: [
      // Load the default Web instrumentations
      ...getWebInstrumentations(),
      // Tracing Instrumentation is needed if you want to use the React Profiler
      new TracingInstrumentation({
        instrumentationOptions: {
          // Requests to these URLs will have tracing headers attached.
          propagateTraceHeaderCorsUrls: [new RegExp(`${API_URL}/*`)]
        }
      }),
      new ReactIntegration({
      })
    ]
  })
}

const rootDiv = document.getElementById('root')
if (rootDiv === null) {
  throw new Error('No root <div> found !')
}

const root = ReactDOM.createRoot(
  rootDiv
)
root.render(
  <App apiUrl={API_URL}/>
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
