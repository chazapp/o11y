import React, { useState } from 'react'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import Box from '@mui/material/Box'
import axios from 'axios'

function Form (): JSX.Element {
  const [message, setMessage] = useState<string>('')
  const [username, setUsername] = useState<string>('')

  const handleSubmit = (event: React.FormEvent): void => {
    event.preventDefault()

    if (message.trim() === '' || username.trim() === '') {
      console.log('Bad!')
      return
    }
    // Handle the new message (e.g., send it to a server or update your state)
    axios.post('/message', {
      message,
      username
    }).then(() => {
      setMessage('')
      setUsername('')
    }).catch((err) => {
      console.error(err)
    })
  }

  return (
        <form onSubmit={handleSubmit}>
          <Box display="flex" flexDirection="column">
            <TextField
              label="Nick"
              variant="outlined"
              value={username}
              onChange={(e) => { setUsername(e.target.value) }}
              inputProps={{
                autoComplete: 'off'
              }}
              required
            />
            <TextField
              label="Message"
              variant="outlined"
              value={message}
              onChange={(e) => { setMessage(e.target.value) }}
              multiline
              rows={1}
              inputProps={{
                autoComplete: 'off'
              }}
              required
            />
            <Button type="submit" variant="contained" color="primary">
              Send
            </Button>
          </Box>
        </form>
  )
}

export default Form
