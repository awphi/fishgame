import { GameServerUpdate, PlayerAction } from './fishgame.pb';
import { EventBus } from '../event-bus';

const HEARTBEAT_TIMER = 5000; // ms

export interface FishGameClientOptions {
	host: string;
	reconnectTimeout?: number;
}

export class FishGameClient {
	private ws: WebSocket | undefined;
	private eventBus = new EventBus<GameServerUpdate>();
	private pingTimeout: number = -1;
	private deferedMessages: Uint8Array<ArrayBufferLike>[] = [];

	public get readyState(): number {
		return this.ws?.readyState ?? 3;
	}

	constructor(private opts: FishGameClientOptions) {
		this.connect();

		// heartbeat handler
		this.on('pong', () => {
			// TODO record latency figures with performance.now()
			this.pingTimeout = setTimeout(() => this.ping(), HEARTBEAT_TIMER);
		});
	}

	// send with built-in message deferal
	private send(msg: Uint8Array<ArrayBufferLike>) {
		if (this.readyState !== 1) {
			this.deferedMessages.push(msg);
			return;
		}

		this.ws!.send(msg);
	}

	private ping(): void {
		console.log('ping');
		const action = PlayerAction.encode({ ping: {} }).finish();
		this.send(action);
	}

	private connect(): void {
		if (this.ws !== undefined && this.ws.readyState <= 1) {
			return;
		}

		const onMessage = this.onMessage.bind(this);
		const connect = this.connect.bind(this);
		const onClose = () => {
			clearTimeout(this.pingTimeout);
			this.pingTimeout = -1;
			setTimeout(connect, this.opts.reconnectTimeout ?? 1000);
		};

		this.ws = new WebSocket(this.opts.host);
		this.ws.binaryType = 'arraybuffer';
		this.ws.addEventListener('message', onMessage);
		this.ws.addEventListener('close', onClose);

		// start heartbeat
		this.ws.addEventListener('open', () => {
			this.ping();

			for (const msg of this.deferedMessages) {
				this.send(msg);
			}
			this.deferedMessages = [];
		});
	}

	private onMessage(ev: MessageEvent<ArrayBuffer>): void {
		const update = GameServerUpdate.decode(new Uint8Array(ev.data));
		const key = Object.keys(update)[0] as keyof GameServerUpdate;
		this.eventBus.fire(key, update[key]!);
	}

	on = this.eventBus.on.bind(this.eventBus);
	off = this.eventBus.off.bind(this.eventBus);
}

// single global client instance for convenience
export const client = new FishGameClient({ host: 'ws://localhost:8081' });
