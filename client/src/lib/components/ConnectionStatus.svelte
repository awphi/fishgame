<script lang="ts">
	import { client } from '$lib/game-client';
	import { writable } from 'svelte/store';

	export const updateRate = 250;

	// https://developer.mozilla.org/en-US/docs/Web/API/WebSocket/readyState
	const statusText = ['Connecting', 'Open', 'Closing', 'Closed'];
	const statusColors = ['bg-blue-500', 'bg-green-500', 'bg-orange-500', 'bg-red-500'];

	// hacky way to re-render properties from the client "reactively"
	const clientStore = writable(client);
	setInterval(() => clientStore.set(client), updateRate);
</script>

<div class=" bg-neutral-700 rounded-md px-4 py-1 select-none w-fit flex items-center gap-2">
	<div class={`w-2 h-2 rounded-full ${statusColors[$clientStore.readyState]}`}></div>
	<span>{statusText[$clientStore.readyState]}</span>
	<span>{$clientStore.latency.toFixed(0)}ms</span>
</div>
