import { map, Observable, scan } from 'rxjs'
import { WebSocketSubject } from 'rxjs/webSocket'

export type Stream = Observable<Stream.Item>

export namespace Stream {
  export interface Item {
    id: string
    inProgress: boolean
    tcp: string
    httpStatus: string
    duration: string
  }
  export const create = () => {
    const sock = new WebSocketSubject<Stream.Item>(`ws://${location.host}/ws`)
    return sock.pipe(
      scan<Stream.Item, Map<string, Stream.Item>>(
        (a, c) => (a.set(c.id, c), a),
        new Map<string, Stream.Item>()
      ),
      map(items => [...items.values()])
    )
  }
}
