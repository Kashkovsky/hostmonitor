import { lift } from '@grammarly/focal'
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow
} from '@mui/material'
import * as React from 'react'
import { map } from 'rxjs'
import { Stream } from './stream'

const Body = lift(TableBody)

export const Main = () => (
  <main>
    <TableContainer component={Paper}>
      <Table sx={{ minWidth: 650 }} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>Address</TableCell>
            <TableCell align="right">TCP status</TableCell>
            <TableCell align="right">HTTP status</TableCell>
            <TableCell align="right">Request duration</TableCell>
          </TableRow>
        </TableHead>
        <Body>
          {Stream.create().pipe(
            map(items =>
              items.map(item => (
                <TableRow key={item.id} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
                  <TableCell component="th" scope="row">
                    {item.id}
                  </TableCell>
                  <TableCell align="right">{item.tcp}</TableCell>
                  <TableCell align="right">{item.httpStatus}</TableCell>
                  <TableCell align="right">{item.duration}</TableCell>
                </TableRow>
              ))
            )
          )}
        </Body>
      </Table>
    </TableContainer>
  </main>
)
