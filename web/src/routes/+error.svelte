<script lang="ts">
	import { page } from '$app/state';
	import { HugeiconsIcon } from '@hugeicons/svelte';
	import { Alert02Icon } from '@hugeicons/core-free-icons';

	// The API writes these messages, so they are already plain English.
	const heading = $derived(page.status === 404 ? 'Nothing here' : 'That did not work');
</script>

<svelte:head>
	<title>{heading} · GetNotified</title>
	<meta name="robots" content="noindex" />
</svelte:head>

<div class="mx-auto max-w-md pt-16 text-center">
	<HugeiconsIcon
		icon={Alert02Icon}
		size={32}
		strokeWidth={1.5}
		class="mx-auto text-stone-400"
	/>
	<h1 class="mt-4 text-xl font-semibold tracking-tight">{heading}</h1>
	<p class="mt-2 text-stone-600">{page.error?.message}</p>
	<a
		href="/"
		class="mt-6 inline-block rounded-md bg-stone-900 px-4 py-2 text-sm text-white hover:bg-stone-700"
	>
		Back to monitors
	</a>
</div>
