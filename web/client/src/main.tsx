import { F, lift, ReadOnlyAtom } from '@grammarly/focal'
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
import { map, Observable } from 'rxjs'
import { Stream } from './stream'

const Body = lift(TableBody)

const Tcp = ({ item }: { item: Stream.Item }) => (
  <Chip label={item.tcp} color={Stream.Item.getTcpStatus(item)} />
)

const Row = ({ item }: { item: Observable<Stream.Item> }) => (
  <F.Fragment>
    {item.pipe(
      map(item => (
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
    )}
  </F.Fragment>
)

export const Main = () => {
  const { ids, getItem } = Stream.create()
  return (
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
          <Body>{ids.pipe(map(ids => ids.map(id => <Row key={id} item={getItem(id)} />)))}</Body>
        </Table>
      </TableContainer>
    </main>
  )
}
