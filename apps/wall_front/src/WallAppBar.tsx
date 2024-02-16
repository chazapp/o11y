import React from 'react'
import { AppBar, Box, IconButton, Toolbar, Typography } from '@mui/material'
import MenuIcon from '@mui/icons-material/Menu'
import MoreVert from '@mui/icons-material/MoreVert'

export default function WallAppBar (): JSX.Element {
  return (
      <Box sx={{ flexGrow: 1 }}>
        <AppBar position="static">
          <Toolbar>
            <IconButton
              size="large"
              edge="start"
              color="inherit"
              aria-label="menu"
              sx={{ mr: 2 }}
            >
              <MenuIcon />
            </IconButton>
            <Typography variant="h6" component="div" sx={{ flexGrow: 1 }} noWrap>
              The Wall Application
            </Typography>
            <IconButton color="inherit">
              <MoreVert fontSize="large"/>
            </IconButton>
          </Toolbar>
        </AppBar>
      </Box>
  )
}
