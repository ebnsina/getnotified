<script lang="ts">
	import { goto } from '$app/navigation';

	// The badges on the buttons are real keys, not decoration.
	let { keys }: { keys: Record<string, string> } = $props();

	function onKeydown(event: KeyboardEvent) {
		if (event.metaKey || event.ctrlKey || event.altKey) return;
		const target = event.target as HTMLElement | null;
		if (target && ['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName)) return;

		const destination = keys[event.key.toLowerCase()];
		if (destination) goto(destination);
	}
</script>

<svelte:window onkeydown={onKeydown} />
