<script lang="ts">
	import { page } from '$app/state';
	import Icon from '$lib/Icon.svelte';
	import { Alert02Icon } from '@hugeicons/core-free-icons';

	// The API writes these messages, so they are already plain English.
	const heading = $derived(page.status === 404 ? 'Nothing here' : 'That did not work');
	const home = $derived(page.url.pathname.startsWith('/app') ? '/app' : '/');
</script>

<svelte:head>
	<title>{heading} · GetNotified</title>
	<meta name="robots" content="noindex" />
</svelte:head>

<div class="min-h-screen bg-night text-bright">
	<div class="mx-auto max-w-md px-6 pt-24 text-center">
	<Icon
		icon={Alert02Icon}
		size={32}
		strokeWidth={1.5}
		class="mx-auto text-dim"
	/>
	<h1 class="mt-4 text-xl font-semibold tracking-tight">{heading}</h1>
	<p class="mt-2 text-mid">{page.error?.message}</p>
		<a
			href={home}
			class="mt-6 inline-block rounded-md bg-up px-4 py-2 text-sm font-medium text-night hover:bg-up/90"
		>
			{home === '/app' ? 'Back to monitors' : 'Back to the home page'}
		</a>
	</div>
</div>
