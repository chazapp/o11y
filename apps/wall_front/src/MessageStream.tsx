import React from 'react'
import Box from '@mui/material/Box'
import { type WallMessage } from './types'
import MessageBox from './MessageBox'

export default function MessageStream (props: { messages: WallMessage[] }): JSX.Element {
  const { messages } = props

  const orderedMessages = messages.reverse()
  return (
        <Box>
        {
            orderedMessages.map((message, index) => {
              return (
                    <MessageBox key={index} message={message}></MessageBox>
              )
            })
        }
        </Box>
  )
}
