import React, { useEffect, useState } from 'react'
import './App.css'
import { Box, Typography } from '@mui/material'
import Form from './Form'
import { type WallMessage } from './types'
import { httpUrlToWebSocketUrl } from './utils'
import MessageStream from './MessageStream'
import BuildInfo from './BuildInfo'
import WallAppBar from './WallAppBar'

function App (props: { apiUrl: string, clientVersion: string }): JSX.Element {
  const { apiUrl, clientVersion } = props
  const [wallMessages, setWallMessages] = useState<WallMessage[]>([])
  const [webSocket, setWebSocket] = useState<WebSocket | undefined>(undefined)
  const [errMessage, setErrMessage] = useState('')

  useEffect(() => {
    if (webSocket !== undefined) {
      return
    }
    const ws = new WebSocket(`${httpUrlToWebSocketUrl(apiUrl)}/ws`)
    ws.onerror = (event) => {
      setErrMessage('WebSocket error !')
    }
    ws.onmessage = (event) => {
      console.log(event.data)
      const { message, username, id } = JSON.parse(event.data as string)

      console.warn('Received event message!')
      const newMessage = {
        message,
        username,
        id
      }
      // Add the message to the state while retaining max 50 elements
      setWallMessages((oldArray) => [...oldArray.slice(wallMessages.length - 49), newMessage])
    }
    setWebSocket(ws)
  }, [apiUrl, webSocket, wallMessages])

  return (
    <div className="App">
      <WallAppBar />
      <Box sx={{
        display: 'flex',
        flexDirection: 'row'
      }}>
        <Box sx={{
          marginTop: '1%',
          width: '500px',
          textAlign: 'center',
          margin: '10px'
        }}>
            <Form />
            {errMessage !== ''
              ? <Typography color={'red'}>{errMessage}</Typography>
              : <Typography variant="h1">{wallMessages.length}</Typography>
            }
        </Box>
        <Box sx={{
          marginTop: '1%',
          width: '50%'
        }}>
          <MessageStream messages={wallMessages}/>
        </Box>
      </Box>
      <BuildInfo clientVersion={clientVersion}/>
    </div>
  )
}

export default App
