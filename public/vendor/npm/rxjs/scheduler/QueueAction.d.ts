import { Action } from './Action';
import { FutureAction } from './FutureAction';
export declare class QueueAction<T> extends FutureAction<T> {
    protected _schedule(state?: T, delay?: number): Action;
}
