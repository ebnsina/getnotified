<script lang="ts">
	import { enhance } from '$app/forms';
	import Meta from '$lib/Meta.svelte';
	import MonitorForm from '$lib/MonitorForm.svelte';
	import { validateMonitorForm } from '$lib/schemas';

	let { data, form } = $props();
	let clientErrors = $state<Record<string, string>>({});

	const errors = $derived({ ...clientErrors, ...(form?.errors ?? {}) });

	// Catches what we can before the round trip; the API still decides.
	const submit = ({ formData, cancel }: { formData: FormData; cancel: () => void }) => {
		clientErrors = validateMonitorForm(formData);
		if (Object.keys(clientErrors).length > 0) cancel();
	};
</script>

<Meta
	title="New monitor · GetNotified"
	description="Add something new to keep an eye on."
	origin={data.origin}
	noindex
/>

<h1 class="font-display text-2xl font-normal text-bright italic">New monitor</h1>

<form method="POST" use:enhance={submit} class="mt-6 max-w-2xl">
	<MonitorForm {errors} />
	{#if errors.form}<p class="mt-4 text-sm text-down">{errors.form}</p>{/if}

	<div class="mt-6 flex gap-3">
		<button class="btn btn-primary">
			Create monitor
		</button>
		<a href="/app" class="btn btn-secondary">Cancel</a>
	</div>
</form>
