import React from 'react'
import Box from '@mui/material/Box'
import PersonOutlineIcon from '@mui/icons-material/PersonOutline'
import { Paper, Typography } from '@mui/material'
import { type WallMessage } from './types'

export default function MessageBox (props: { message: WallMessage }): JSX.Element {
  const { username, message } = props.message
  return (
        <Box>
            <Paper id="backBox" elevation={3} sx={{
              display: 'flex',
              flexDirection: 'row',
              margin: 1,
              backgroundColor: 'darkgrey'
            }}>
                <Box id="picture" sx={{
                  padding: '3em'
                }}>
                    <PersonOutlineIcon fontSize="large"/>
                </Box>

                <Box id="content" sx={{
                  margin: 'auto',
                  marginLeft: 0,
                  display: 'flex',
                  flexDirection: 'column',
                  textAlign: 'left'
                }}>
                    <Typography>{username}: </Typography>
                    <Typography>{message}</Typography>
                </Box>
            </Paper>
        </Box>
  )
}
