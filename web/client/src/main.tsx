import { lift } from '@grammarly/focal'
import {
  Chip,
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

const Tcp = ({ item }: { item: Stream.Item }) => (
  <Chip label={item.tcp} color={Stream.Item.getTcpStatus(item)} />
)

export const Main = () => (
  <main>
    <TableContainer component={Paper}>
      <Table stickyHeader sx={{ minWidth: 650 }} size="small" aria-label="simple table">
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
                <TableRow
                  key={item.id}
                  sx={{
                    '&:last-child td, &:last-child th': { border: 0 },
                    background: Stream.Item.getStatusColor(item)
                  }}
                >
                  <TableCell width={150} component="th" scope="row">
                    <strong>{item.id}</strong>
                  </TableCell>
                  <TableCell align="right">
                    <Tcp item={item} />
                  </TableCell>
                  <TableCell align="right">{item.httpResponse}</TableCell>
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
