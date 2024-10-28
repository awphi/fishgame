type Callback<T> = (payload: NonNullable<T>) => void;

export class EventBus<EventMap> {
	private events: Map<keyof EventMap, Callback<EventMap[keyof EventMap]>[]> = new Map();

	on<T extends keyof EventMap>(event: T, callback: Callback<EventMap[T]>): void {
		if (!this.events.has(event)) {
			this.events.set(event, []);
		}
		this.events.get(event)!.push(callback as Callback<EventMap[keyof EventMap]>);
	}

	off<T extends keyof EventMap>(event: T, callback: Callback<EventMap[T]>): void {
		const callbacks = this.events.get(event);
		if (callbacks !== undefined) {
			this.events.set(
				event,
				callbacks.filter((cb) => cb !== callback)
			);
		}
	}

	fire<T extends keyof EventMap>(event: T, payload: NonNullable<EventMap[T]>): void {
		const callbacks = this.events.get(event);
		if (callbacks !== undefined) {
			callbacks.forEach((callback) => callback(payload));
		}
	}
}
