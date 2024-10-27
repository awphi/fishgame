import { GameServerUpdate, PlayerAction } from './fishgame.pb';

export interface FishGameClientOptions {
	host: string;
	reconnectTimeout?: number;
}

export class FishGameClient {
	private ws: WebSocket | undefined;

	private get readyState(): number {
		return this.ws?.readyState ?? 3;
	}

	constructor(private opts: FishGameClientOptions) {
		this.connect();
	}

	private connect(): void {
		if (this.ws !== undefined && this.ws.readyState <= 1) {
			return;
		}

		const onMessage = this.onMessage.bind(this);
		const connect = this.connect.bind(this);
		const onClose = () => {
			setTimeout(connect, this.opts.reconnectTimeout ?? 1000);
		};

		this.ws = new WebSocket(this.opts.host);
		this.ws.binaryType = 'arraybuffer';
		this.ws.addEventListener('message', onMessage);
		this.ws.addEventListener('close', onClose);

		this.ws.addEventListener('open', () => {
			const action = PlayerAction.encode({ ping: { id: 90 } }).finish();
			this.ws!.send(action);
		});
	}

	private onMessage(ev: MessageEvent<ArrayBuffer>): void {
		console.log(ev, GameServerUpdate.decode(new Uint8Array(ev.data)));
	}

	// TODO add ping functionality
}

// single global client instance for convenience
export const client = new FishGameClient({ host: 'ws://localhost:8081' });
