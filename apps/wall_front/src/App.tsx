

import React, { useEffect, useState } from 'react';
import './App.css';
import { Box, Typography } from "@mui/material";
import Form from './Form';
import { WallMessage } from './types';
import { getRandomInt, httpUrlToWebSocketUrl } from './utils';
import MessageStream from './MessageStream';

function App(props: {apiUrl: string}) {
  const { apiUrl } = props;
  const [wallMessages, setWallMessages] = useState<WallMessage[]>([]);
  const [webSocket, setWebSocket] = useState<WebSocket | undefined>(undefined);
  const [errMessage, setErrMessage] = useState("")

  useEffect(() => {
    if (webSocket) {
      return;
    }
    const ws = new WebSocket(`${httpUrlToWebSocketUrl(apiUrl)}/ws`);
    if (!ws) {
      setErrMessage("Could not connect WebSocket !")
      return
    }
    ws.onerror = (event) => {
      setErrMessage("WebSocket error !")
    }
    ws.onmessage = (event) => {
      const { message, username, id } = JSON.parse(event.data);

      console.error("Received event message!");
      const newMessage = {
        message,
        username,
        id,
      };
      // Add the message to the state while retaining max 50 elements
      setWallMessages((oldArray) => [...oldArray.slice(wallMessages.length - 49), newMessage])
    }
    setWebSocket(ws);
  }, [apiUrl, webSocket, wallMessages]);

  return (
    <div className="App">
      <Box sx={{
        display: "flex",
        flexDirection: "row"
      }}>
        <Box sx={{
            marginTop: "1%",
            width: "500px",
          }}>
            <Form />
            {errMessage !== "" ?
              <Typography color={"red"}>{errMessage}</Typography> 
             : 
              <Typography variant="h1">{wallMessages.length}</Typography>
            }
        </Box>
        <Box sx={{
          marginTop: "1%",
          width: "50%"
        }}>
          <MessageStream messages={wallMessages}/>
        </Box>
      </Box>
      
    </div>
  );
}

export default App;
