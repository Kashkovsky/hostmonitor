import { pipe } from 'fp-ts/es6/function'
import { isNumber } from 'fp-ts/es6/number'
import * as Ord from 'fp-ts/es6/Ord'
import {
  distinctUntilChanged,
  fromEvent,
  map,
  mergeMap,
  Observable,
  retryWhen,
  scan,
  share,
  take,
  timer
} from 'rxjs'
import { WebSocketSubject } from 'rxjs/webSocket'
import { boolean, string } from 'fp-ts'
import { Atom } from '@grammarly/focal'
import * as Eq from 'fp-ts/es6/Eq'
import * as A from 'fp-ts/es6/Array'

export type Stream = Observable<Stream.Item>

export namespace Stream {
  const RETRY_DELAY = 5000
  export interface Item {
    readonly id: string
    readonly inProgress: boolean
    readonly tcp: string
    readonly httpResponse: string
    readonly duration: string
    readonly status: Stream.Item.Status
  }

  export namespace Item {
    export enum Status {
      StatusOK = 'OK',
      StatusErr = 'Error',
      StatusErrResponse = 'ErrorResponse'
    }

    export const eq: Eq.Eq<Stream.Item> = Eq.struct({
      id: string.Eq,
      inProgress: boolean.Eq,
      tcp: string.Eq,
      httpResponse: string.Eq,
      duration: string.Eq,
      status: string.Eq
    })

    export const ord: Ord.Ord<Stream.Item> = pipe(
      Ord.contramap<string, Stream.Item>(x => x.status)(string.Ord),
      Ord.reverse
    )

    export const getStatusColor = (item: Stream.Item) => {
      switch (item.status) {
        case Stream.Item.Status.StatusOK:
          return '#a1ffc3'
        case Stream.Item.Status.StatusErr:
          return '#ebffa1'
        case Stream.Item.Status.StatusErrResponse:
          return '#fc8672'
        default:
          return undefined
      }
    }

    export const getTcpStatus = (item: Stream.Item) => {
      const success = +item.tcp.split('/')[0]
      if (isNumber(success)) {
        switch (true) {
          case success >= 6:
            return 'success'
          case success < 6 && success > 3:
            return 'warning'
          case success <= 3:
            return 'error'
          default:
            return 'default'
        }
      }

      return 'default'
    }
  }

  export const create = () => {
    const sock = new WebSocketSubject<Stream.Item>(`ws://${location.host}/ws`)
    const collection = sock.pipe(
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
      share()
    )

    const getItem = (id: string) =>
      collection.pipe(
        map(c => c.get(id)!),
        distinctUntilChanged(Stream.Item.eq.equals)
      )

    const ids = collection.pipe(
      map(c => [...c.keys()]),
      distinctUntilChanged((a, b) => JSON.stringify(a) === JSON.stringify(b))
    )

    return {
      ids,
      getItem
    }
  }
}
