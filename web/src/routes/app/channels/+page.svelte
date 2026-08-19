<script lang="ts">
	import { enhance } from '$app/forms';
	import Meta from '$lib/Meta.svelte';
	import Icon from '$lib/Icon.svelte';
	import { Delete02Icon, SentIcon } from '@hugeicons/core-free-icons';
	import { CHANNEL_FIELDS, CHANNEL_TYPES } from '$lib/channels';
	import { validateChannelForm } from '$lib/schemas';

	let { data, form } = $props();
	let type = $state('slack');
	let clientErrors = $state<Record<string, string>>({});

	const errors = $derived({ ...clientErrors, ...(form?.errors ?? {}) });
	const field = 'field mt-1';

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

<h1 class="page-title">Notifications</h1>
<p class="mt-1 text-sm text-dim">
	Each channel is sent on its own, so a slow one never holds up the rest. Send a test to check
	one works before you rely on it.
</p>

{#if form?.message}
	<p class="mt-4 panel px-4 py-2.5 text-sm text-mid">{form.message}</p>
{/if}

{#if data.channels.length > 0}
	<ul class="mt-6 divide-y divide-rule/60 panel">
		{#each data.channels as channel (channel.id)}
			<li class="flex items-center justify-between gap-4 px-4 py-3 text-sm">
				<div class="min-w-0">
					<span class="font-medium">{channel.name}</span>
					<span class="text-dim"> · {channel.type}</span>
					<div class="truncate font-mono text-xs text-dim">
						{Object.values(channel.config)[0] ?? ''}
					</div>
				</div>
				<div class="flex flex-none items-center gap-1">
					<form method="POST" action="?/test" use:enhance>
						<input type="hidden" name="id" value={channel.id} />
						<button aria-label="Send a test to {channel.name}" class="btn btn-secondary btn-sm">
							<Icon icon={SentIcon} size={16} strokeWidth={1.5} />
							Send a test
						</button>
					</form>
					<form method="POST" action="?/delete" use:enhance>
						<input type="hidden" name="id" value={channel.id} />
						<button aria-label="Remove {channel.name}" class="btn btn-secondary btn-sm">
							<Icon icon={Delete02Icon} size={16} strokeWidth={1.5} />
							Remove
						</button>
					</form>
				</div>
			</li>
		{/each}
	</ul>
{/if}

<form method="POST" action="?/create" use:enhance={submit} class="mt-8 max-w-lg space-y-4">
	<h2 class="section-label">Add a channel</h2>

	<div>
		<label for="name" class="block text-sm font-medium">Name</label>
		<input id="name" name="name" required placeholder="Ops Slack" class={field} />
		{#if errors.name}<p class="mt-1 text-sm text-down">{errors.name}</p>{/if}
	</div>

	<div>
		<label for="type" class="block text-sm font-medium">Type</label>
		<select id="type" name="type" bind:value={type} class="field field-select mt-1">
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

	<button class="btn btn-primary">
		Add channel
	</button>
	<p class="text-xs text-dim">
		Sign-in details for email and text messages live in your server settings, not here.
	</p>
</form>
