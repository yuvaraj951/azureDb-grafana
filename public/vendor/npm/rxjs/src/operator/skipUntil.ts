import {Operator} from '../Operator';
import {Subscriber} from '../Subscriber';
import {Observable} from '../Observable';

import {OuterSubscriber} from '../OuterSubscriber';
import {InnerSubscriber} from '../InnerSubscriber';
import {subscribeToResult} from '../util/subscribeToResult';

/**
 * Returns an Observable that skips items emitted by the source Observable until a second Observable emits an item.
 *
 * <img src="./img/skipUntil.png" width="100%">
 *
 * @param {Observable} the second Observable that has to emit an item before the source Observable's elements begin to
 * be mirrored by the resulting Observable.
 * @return {Observable<T>} an Observable that skips items from the source Observable until the second Observable emits
 * an item, then emits the remaining items.
 * @method skipUntil
 * @owner Observable
 */
export function skipUntil<T>(notifier: Observable<any>): Observable<T> {
  return this.lift(new SkipUntilOperator(notifier));
}

export interface SkipUntilSignature<T> {
  (notifier: Observable<any>): Observable<T>;
}

class SkipUntilOperator<T> implements Operator<T, T> {
  constructor(private notifier: Observable<any>) {
  }

  call(subscriber: Subscriber<T>): Subscriber<T> {
    return new SkipUntilSubscriber(subscriber, this.notifier);
  }
}

class SkipUntilSubscriber<T, R> extends OuterSubscriber<T, R> {

  private hasValue: boolean = false;
  private isInnerStopped: boolean = false;

  constructor(destination: Subscriber<any>,
              notifier: Observable<any>) {
    super(destination);
    this.add(subscribeToResult(this, notifier));
  }

  protected _next(value: T) {
    if (this.hasValue) {
      super._next(value);
    }
  }

  protected _complete() {
    if (this.isInnerStopped) {
      super._complete();
    } else {
      this.unsubscribe();
    }
  }

  notifyNext(outerValue: T, innerValue: R,
             outerIndex: number, innerIndex: number,
             innerSub: InnerSubscriber<T, R>): void {
    this.hasValue = true;
  }

  notifyComplete(): void {
    this.isInnerStopped = true;
    if (this.isStopped) {
      super._complete();
    }
  }
}
