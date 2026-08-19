<script lang="ts">
	import Meta from '$lib/Meta.svelte';
	import { dateTime, duration, percent, pillClass, relative, statusLabel } from '$lib/format';

	let { data } = $props();

	const status = $derived(data.status);
	const listed = $derived(status.monitors.filter((monitor) => !monitor.paused));
	const summary = $derived(
		status.overall === 'operational'
			? 'All monitored services are responding normally.'
			: 'Some monitored services are not responding.'
	);
</script>

<Meta
	title="{status.org.name} status"
	description={summary}
	origin={data.origin}
/>

<div class="min-h-screen bg-night text-bright">
<div class="mx-auto max-w-3xl px-6 py-12">

<h1 class="font-display text-2xl font-normal text-bright italic">{status.org.name} status</h1>
<p class="mt-1 text-sm text-dim">
	{summary} Last checked {relative(data.locale, status.as_of)}.
</p>

<div class="mt-6 divide-y divide-rule/60 panel">
	{#each listed as monitor (monitor.id)}
		<div class="flex flex-wrap items-center justify-between gap-2 px-4 py-3 text-sm">
			<span class="text-bright">{monitor.name}</span>
			<span class="flex items-center gap-3">
				<span class="numeric text-dim">{percent(data.locale, monitor.up_30d)} over 30 days</span>
				<span class={pillClass(monitor.status)}>{statusLabel(monitor.status)}</span>
			</span>
		</div>
	{:else}
		<p class="px-4 py-3 text-sm text-dim">Nothing is being monitored yet.</p>
	{/each}
</div>

<section class="mt-8">
	<h2 class="text-sm font-medium text-dim">Past incidents</h2>
	{#if status.incidents.length === 0}
		<p class="mt-2 text-sm text-dim">No incidents so far.</p>
	{:else}
		<ul class="mt-2 space-y-2 text-sm">
			{#each status.incidents as incident (incident.id)}
				<li class="panel px-4 py-3">
					<span class="font-medium">
						{status.monitors.find((monitor) => monitor.id === incident.monitor_id)?.name ??
							'A service'}
					</span>
					<span class="text-dim">
						· {incident.resolved_at ? 'resolved' : 'still open'} ·
						<span class="numeric">
							{duration(data.locale, incident.started_at, incident.resolved_at)}
						</span>
					</span>
					<div class="text-xs text-dim">
						{dateTime(data.locale, incident.started_at)}
					</div>
				</li>
			{/each}
		</ul>
	{/if}
</section>

<p class="mt-10 text-xs text-dim">
	Status page by <a href="/" class="underline">GetNotified</a>. This page does not track you.
</p>

</div>
</div>
