<script lang="ts">
	import { enhance } from '$app/forms';
	import Meta from '$lib/Meta.svelte';
	import MonitorForm from '$lib/MonitorForm.svelte';
	import Icon from '$lib/Icon.svelte';
	import { PlayIcon, PauseIcon, Delete02Icon } from '@hugeicons/core-free-icons';
	import { dateTime, duration, milliseconds, percent, pillClass, relative, statusLabel } from '$lib/format';
	import { validateMonitorForm } from '$lib/schemas';

	let { data, form } = $props();
	let clientErrors = $state<Record<string, string>>({});

	const monitor = $derived(data.monitor);

	const stats = $derived([
		['24 hours', percent(data.locale, monitor.up_24h)],
		['7 days', percent(data.locale, monitor.up_7d)],
		['30 days', percent(data.locale, monitor.up_30d)],
		['Latency', milliseconds(data.locale, monitor.latency_ms)]
	]);
	const errors = $derived({ ...clientErrors, ...(form?.errors ?? {}) });
	const timeline = $derived([...data.checks].reverse());

	const submit = ({ formData, cancel }: { formData: FormData; cancel: () => void }) => {
		clientErrors = validateMonitorForm(formData);
		if (Object.keys(clientErrors).length > 0) cancel();
	};
</script>

<Meta
	title="{monitor.name} · GetNotified"
	description="Status, incident history and settings for {monitor.name}."
	origin={data.origin}
	noindex
/>

<nav class="mb-4 font-mono text-xs text-dim">
	<a href="/app" class="hover:text-bright">Monitors</a>
	<span class="px-1.5">/</span>
	<span class="text-mid">{monitor.name}</span>
</nav>

<div class="flex items-start justify-between gap-4">
	<div>
		<h1 class="page-title">{monitor.name}</h1>
		<p class="mt-2 flex flex-wrap items-center gap-2 text-sm text-dim">
			<span class={pillClass(monitor.status, monitor.paused)}>
				{statusLabel(monitor.status, monitor.paused)}
			</span>
			<span class="font-mono text-xs">{monitor.target}</span>
		</p>
	</div>
	<form method="POST" action="?/pause" use:enhance>
		<input type="hidden" name="paused" value={String(!monitor.paused)} />
		<button
			class="btn btn-secondary btn-sm"
		>
			<Icon icon={monitor.paused ? PlayIcon : PauseIcon} size={16} strokeWidth={1.5} />
			{monitor.paused ? 'Resume checks' : 'Pause checks'}
		</button>
	</form>
</div>

{#if form?.message}
	<p class="mt-4 panel px-4 py-2.5 text-sm text-mid">{form.message}</p>
{/if}

<section class="panel mt-6 grid grid-cols-2 divide-rule/60 sm:grid-cols-4 sm:divide-x">
	{#each stats as [label, value] (label)}
		<div class="px-5 py-4">
			<p class="section-label">{label}</p>
			<p class="numeric mt-1 text-lg text-bright">{value}</p>
		</div>
	{/each}
</section>

<section class="mt-8">
	<h2 class="section-label">Recent checks</h2>
	<div class="mt-2 flex flex-wrap gap-1">
		{#each timeline as check (check.id)}
			<span
				class="h-8 w-2 rounded-sm {check.ok ? 'bg-up' : 'bg-down'}"
				title="{dateTime(data.locale, check.checked_at)} — {check.ok
					? milliseconds(data.locale, check.latency_ms)
					: (check.error ?? 'no response')}"
			></span>
		{:else}
			<p class="text-sm text-dim">Nothing checked yet. The first result appears shortly.</p>
		{/each}
	</div>
</section>

<section class="mt-8">
	<h2 class="section-label">Incidents</h2>
	{#if data.incidents.length === 0}
		<p class="mt-2 text-sm text-dim">No incidents so far.</p>
	{:else}
		<ul class="mt-2 divide-y divide-rule/60 panel">
			{#each data.incidents as incident (incident.id)}
				<li class="flex flex-wrap items-baseline justify-between gap-2 px-4 py-3 text-sm">
					<div>
						<span class="font-medium">{incident.resolved_at ? 'Resolved' : 'Still open'}</span>
						<span class="text-dim"> · {incident.cause ?? 'no response'}</span>
					</div>
					<div class="text-dim">
						{relative(data.locale, incident.started_at)} · lasted
						<span class="numeric">{duration(data.locale, incident.started_at, incident.resolved_at)}</span>
					</div>
				</li>
			{/each}
		</ul>
	{/if}
</section>

<section class="mt-8 max-w-2xl">
	<h2 class="section-label">Who hears about it</h2>
	{#if data.channels.length === 0}
		<p class="mt-2 text-sm text-dim">
			No channels set up yet. <a href="/app/channels" class="underline">Add one</a>.
		</p>
	{:else}
		<form method="POST" action="?/channels" use:enhance class="mt-2 space-y-2">
			{#each data.channels as channel (channel.id)}
				<label class="flex items-center gap-2 text-sm">
					<input
						class="check"
						type="checkbox"
						name="channel_ids"
						value={channel.id}
						checked={data.attached.includes(channel.id)}
					/>
					{channel.name} <span class="text-dim">({channel.type})</span>
				</label>
			{/each}
			<button class="btn btn-secondary btn-sm">
				Save channels
			</button>
		</form>
	{/if}
</section>

<section class="mt-8 max-w-2xl">
	<h2 class="section-label">Settings</h2>
	<form method="POST" action="?/update" use:enhance={submit} class="mt-2">
		<MonitorForm {monitor} {errors} />
		{#if errors.form}<p class="mt-4 text-sm text-down">{errors.form}</p>{/if}

		<button class="btn btn-primary mt-6">
			Save changes
		</button>
	</form>
</section>

<section class="mt-12 max-w-2xl border-t border-rule pt-6">
	<form method="POST" action="?/delete">
		<button class="btn btn-danger btn-sm">
			<Icon icon={Delete02Icon} size={16} strokeWidth={1.5} />
			Delete this monitor
		</button>
		<p class="mt-1 text-xs text-dim">
			This also removes its check history and past incidents.
		</p>
	</form>
</section>
