import { fromEvent, map, mergeMap, Observable, retryWhen, scan, take, timer } from 'rxjs'
import { WebSocketSubject } from 'rxjs/webSocket'

export type Stream = Observable<Stream.Item>

export namespace Stream {
  const RETRY_DELAY = 5000
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
      retryWhen(errors =>
        errors.pipe(
          mergeMap(() => {
            if (window.navigator.onLine) {
              console.warn(`Retrying in ${RETRY_DELAY}ms.`)
              return timer(RETRY_DELAY)
            } else {
              return fromEvent(window, 'online').pipe(take(1))
            }
          })
        )
      ),
      scan<Stream.Item, Map<string, Stream.Item>>(
        (a, c) => (a.set(c.id, c), a),
        new Map<string, Stream.Item>()
      ),
      map(items => [...items.values()])
    )
  }
}
