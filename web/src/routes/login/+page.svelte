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

<div class="mx-auto max-w-sm pt-16">
	<h1 class="text-xl font-semibold tracking-tight">Sign in</h1>
	<p class="mt-1 text-sm text-stone-500">GetNotified is for the account owner only.</p>

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
				class="mt-1 w-full rounded-md border border-stone-300 px-3 py-2 focus:border-stone-500 focus:outline-none"
			/>
		</div>
		{#if form?.message}
			<p id="password-error" class="text-sm text-down">{form.message}</p>
		{/if}
		<button class="w-full rounded-md bg-stone-900 px-4 py-2 text-white hover:bg-stone-700">
			Sign in
		</button>
	</form>
</div>
