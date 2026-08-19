<script lang="ts">
	import { enhance } from '$app/forms';
	import Meta from '$lib/Meta.svelte';
	import { HugeiconsIcon } from '@hugeicons/svelte';
	import { Delete02Icon } from '@hugeicons/core-free-icons';
	import { CHANNEL_FIELDS, CHANNEL_TYPES } from '$lib/channels';
	import { validateChannelForm } from '$lib/schemas';

	let { data, form } = $props();
	let type = $state('slack');
	let clientErrors = $state<Record<string, string>>({});

	const errors = $derived({ ...clientErrors, ...(form?.errors ?? {}) });
	const field =
		'mt-1 w-full rounded-md border border-stone-300 px-3 py-2 focus:border-stone-500 focus:outline-none';

	const submit = ({ formData, cancel }: { formData: FormData; cancel: () => void }) => {
		clientErrors = validateChannelForm(formData);
		if (Object.keys(clientErrors).length > 0) cancel();
	};
</script>

<Meta
	title="Notifications · GetNotified"
	description="Choose where outage and recovery notices are sent."
	origin={data.origin}
	noindex
/>

<h1 class="text-xl font-semibold tracking-tight">Notifications</h1>
<p class="mt-1 text-sm text-stone-500">
	Each channel is sent on its own, so a slow one never holds up the rest.
</p>

{#if form?.message}
	<p class="mt-4 rounded-md bg-stone-100 px-3 py-2 text-sm text-stone-700">{form.message}</p>
{/if}

{#if data.channels.length > 0}
	<ul class="mt-6 divide-y divide-stone-100 rounded-lg border border-stone-200 bg-white">
		{#each data.channels as channel (channel.id)}
			<li class="flex items-center justify-between gap-4 px-4 py-3 text-sm">
				<div class="min-w-0">
					<span class="font-medium">{channel.name}</span>
					<span class="text-stone-500"> · {channel.type}</span>
					<div class="truncate font-mono text-xs text-stone-500">
						{Object.values(channel.config)[0] ?? ''}
					</div>
				</div>
				<form method="POST" action="?/delete" use:enhance>
					<input type="hidden" name="id" value={channel.id} />
					<button
						aria-label="Remove {channel.name}"
						class="flex items-center gap-1.5 text-stone-500 hover:text-down"
					>
						<HugeiconsIcon icon={Delete02Icon} size={16} strokeWidth={1.5} />
						Remove
					</button>
				</form>
			</li>
		{/each}
	</ul>
{/if}

<form method="POST" action="?/create" use:enhance={submit} class="mt-8 max-w-lg space-y-4">
	<h2 class="text-sm font-medium text-stone-500">Add a channel</h2>

	<div>
		<label for="name" class="block text-sm font-medium">Name</label>
		<input id="name" name="name" required placeholder="Ops Slack" class={field} />
		{#if errors.name}<p class="mt-1 text-sm text-down">{errors.name}</p>{/if}
	</div>

	<div>
		<label for="type" class="block text-sm font-medium">Type</label>
		<select id="type" name="type" bind:value={type} class={field}>
			{#each CHANNEL_TYPES as option (option.value)}
				<option value={option.value}>{option.label}</option>
			{/each}
		</select>
	</div>

	{#each CHANNEL_FIELDS[type] as channelField (channelField.name)}
		<div>
			<label for={channelField.name} class="block text-sm font-medium">{channelField.label}</label>
			<input
				id={channelField.name}
				name={channelField.name}
				placeholder={channelField.placeholder}
				aria-invalid={errors[channelField.name] ? 'true' : undefined}
				class={field}
			/>
			{#if errors[channelField.name]}
				<p class="mt-1 text-sm text-down">{errors[channelField.name]}</p>
			{/if}
		</div>
	{/each}

	{#if errors.form || errors.config}
		<p class="text-sm text-down">{errors.form ?? errors.config}</p>
	{/if}

	<button class="rounded-md bg-stone-900 px-4 py-2 text-sm text-white hover:bg-stone-700">
		Add channel
	</button>
	<p class="text-xs text-stone-500">
		Credentials for email, Twilio and the iMessage relay live in server settings, not here.
	</p>
</form>
