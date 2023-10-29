

import React, { useEffect, useState } from 'react';
import './App.css';
import axios from "axios";
import { Box } from "@mui/material";
import Form from './Form';
import { Stage, Layer, Text } from 'react-konva';
import { WallMessage } from './types';
import { getRandomInt, httpUrlToWebSocketUrl } from './utils';

function App() {
  const [wallMessages, setWallMessages] = useState<WallMessage[]>([]);
  const [webSocket, setWebSocket] = useState<WebSocket | undefined>(undefined);
  const API_URL = window.env && window.env.API_URL ? window.env.API_URL : process.env.REACT_APP_API_URL;
  if (API_URL === undefined) {
    throw new Error("API URL is undefined !")
  }

  axios.defaults.baseURL = API_URL;
  useEffect(() => {
    if (webSocket) {
      return;
    }
    const ws = new WebSocket(`${httpUrlToWebSocketUrl(API_URL)}/ws`);
    ws.onmessage = (event) => {
      const { message, username, id } = JSON.parse(event.data);
      const wallMessage = {
        message,
        username,
        id,
        posX: getRandomInt(0, window.innerWidth),
        posY: getRandomInt(0, window.innerHeight),
        opacity: 1,
      };

      // Add the message to the state
      setWallMessages((oldArray) => [...oldArray, wallMessage]);

      // Schedule a setInterval to gradually reduce the opacity over 10 seconds
      const interval = 100; // Update opacity every 100 milliseconds
      const duration = 10000; 

      const opacityInterval = setInterval(() => {
        setWallMessages((oldMessages) =>
          oldMessages.map((oldMessage) =>
            oldMessage.id === wallMessage.id
              ? { ...oldMessage, opacity: oldMessage.opacity - 0.1  }
              : oldMessage
          )
        )
      }, interval);

      // Clear the interval once opacity reaches 0
      setTimeout(() => {
        clearInterval(opacityInterval);
        setWallMessages((oldMessages) =>
          oldMessages.filter((oldMessage) => oldMessage.id !== wallMessage.id)
        );
      }, duration);
    
    }
    setWebSocket(ws);
  }, [API_URL, webSocket]);
  
  return (
    <div className="App">
      <Stage style={{position: "absolute"}} width={window.innerWidth} height={window.innerHeight}>
        <Layer >
          {
            wallMessages.map((wallMessage) => {
              return (
                <Text
                  key={wallMessage.id} 
                  text={wallMessage.message}
                  x={wallMessage.posX}
                  y={window.innerHeight/2}
                  opacity={wallMessage.opacity < 0 ? 0 : wallMessage.opacity}
                  fill="red"
                  fontSize={50}
              />
              )
            })
          }
        </Layer>
      </Stage>
      <Box sx={{
          marginTop: "1%",
          width: "500px",
        }}>
          <Form />
      </Box>
    </div>
  );
}

export default App;
