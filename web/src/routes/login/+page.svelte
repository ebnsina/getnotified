<script lang="ts">
	import { enhance } from '$app/forms';
	import Meta from '$lib/Meta.svelte';

	let { data, form } = $props();
</script>

<Meta
	title="Sign in · GetNotified"
	description="Sign in to manage your monitors and notification channels."
	origin={data.origin}
	noindex
/>

<div class="min-h-screen bg-night text-bright">
	<div class="mx-auto max-w-sm px-6 pt-24">
	<h1 class="font-display text-2xl font-normal text-bright">Sign in</h1>
	<p class="mt-1 text-sm text-dim">GetNotified is for the account owner only.</p>

	<form method="POST" use:enhance class="mt-6 space-y-4">
		<div>
			<label for="password" class="block text-sm font-medium">Password</label>
			<input
				id="password"
				name="password"
				type="password"
				autocomplete="current-password"
				required
				aria-describedby={form?.message ? 'password-error' : undefined}
				class="mt-1 w-full rounded-md border border-rule bg-night px-3 py-2 text-bright focus:border-dim focus:outline-none"
			/>
		</div>
		{#if form?.message}
			<p id="password-error" class="text-sm text-down">{form.message}</p>
		{/if}
		<button class="w-full rounded-md bg-up px-4 py-2 font-medium text-night hover:bg-up/90">
			Sign in
		</button>
		</form>

		<a href="/" class="mt-6 inline-block text-sm text-dim hover:text-bright">
			Back to the home page
		</a>
	</div>
</div>
