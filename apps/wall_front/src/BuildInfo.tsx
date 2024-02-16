import React, { useEffect, useState } from 'react'
import { Box, ClickAwayListener, IconButton, Paper, Table, TableBody, TableCell, TableContainer, TableRow, Typography } from '@mui/material'
import HelpIcon from '@mui/icons-material/Help'
import axios from 'axios'
import { type ApiStatus } from './types'

export default function BuildInfo (props: { clientVersion: string }): JSX.Element {
  const { clientVersion } = props
  const [open, setOpen] = useState(false)
  const [error, setError] = useState<string>('')
  const [apiStatus, setApiStatus] = useState<ApiStatus | undefined>(undefined)

  const handleClose = (): void => {
    setOpen(false)
  }

  const handleOpen = (): void => {
    setOpen(true)
  }

  useEffect(() => {
    axios.get('/status').then((resp) => {
      setApiStatus(resp.data as ApiStatus)
    }).catch((err) => {
      setError(`Unable to show status: ${err}`)
      console.error(err)
    })
  }, [open])

  return (
    <Box sx={{
      position: 'fixed',
      bottom: '2%',
      right: '2%'
    }}>
        {open && (
            <ClickAwayListener onClickAway={handleClose}>
                <Paper elevation={3} sx={{
                  position: 'fixed',
                  zIndex: 2,
                  bottom: '2%',
                  right: '2%',
                  width: 'auto',
                  height: 'auto'
                }}>
                    {error !== '' && <Typography sx={{ color: 'red' }}>{error}</Typography>}
                    {apiStatus !== undefined && (
                        <TableContainer>
                            <Table size="small" aria-label="status-data-table">
                                <TableBody>
                                    <TableRow key="client-version">
                                        <TableCell component="th" scope="row">Client Version</TableCell>
                                        <TableCell component="th" scope="row">{clientVersion}</TableCell>
                                    </TableRow>
                                    <TableRow key="api-version">
                                        <TableCell component="th" scope="row">API Version</TableCell>
                                        <TableCell component="th" scope="row">{apiStatus.version}</TableCell>
                                    </TableRow>
                                    <TableRow key="go-version">
                                        <TableCell component="th" scope="row">Go Version</TableCell>
                                        <TableCell component="th" scope="row">{apiStatus.goVersion}</TableCell>
                                    </TableRow>
                                    <TableRow key="connect-ws">
                                        <TableCell component="th" scope="row">Connected WS</TableCell>
                                        <TableCell component="th" scope="row">{apiStatus.connectedWS}</TableCell>
                                    </TableRow>
                                    <TableRow key="total-messages">
                                        <TableCell component="th" scope="row">Total Message Count</TableCell>
                                        <TableCell component="th" scope="row">{apiStatus.messagesCount}</TableCell>
                                    </TableRow>
                                </TableBody>
                            </Table>
                        </TableContainer>
                    )}
                </Paper>
            </ClickAwayListener>
        )}
        <IconButton aria-label="build-info" size="large" onClick={handleOpen}>
            <HelpIcon fontSize="large"/>
        </IconButton>
    </Box>
  )
}
