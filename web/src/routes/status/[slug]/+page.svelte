<script lang="ts">
	import Meta from '$lib/Meta.svelte';
	import { dateTime, dotClass, duration, percent, relative, statusLabel } from '$lib/format';

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

<h1 class="text-xl font-semibold tracking-tight">{status.org.name} status</h1>
<p class="mt-1 text-sm text-stone-500">
	{summary} Last checked {relative(data.locale, status.as_of)}.
</p>

<div class="mt-6 divide-y divide-stone-100 rounded-lg border border-stone-200 bg-white">
	{#each listed as monitor (monitor.id)}
		<div class="flex flex-wrap items-center justify-between gap-2 px-4 py-3 text-sm">
			<span class="flex items-center gap-2">
				<span class="size-2 rounded-full {dotClass(monitor.status)}"></span>
				{monitor.name}
			</span>
			<span class="text-stone-500">
				{statusLabel(monitor.status)} ·
				<span class="numeric">{percent(data.locale, monitor.up_30d)}</span> over 30 days
			</span>
		</div>
	{:else}
		<p class="px-4 py-3 text-sm text-stone-500">Nothing is being monitored yet.</p>
	{/each}
</div>

<section class="mt-8">
	<h2 class="text-sm font-medium text-stone-500">Past incidents</h2>
	{#if status.incidents.length === 0}
		<p class="mt-2 text-sm text-stone-500">No incidents so far.</p>
	{:else}
		<ul class="mt-2 space-y-2 text-sm">
			{#each status.incidents as incident (incident.id)}
				<li class="rounded-lg border border-stone-200 bg-white px-4 py-3">
					<span class="font-medium">
						{status.monitors.find((monitor) => monitor.id === incident.monitor_id)?.name ??
							'A service'}
					</span>
					<span class="text-stone-500">
						· {incident.resolved_at ? 'resolved' : 'still open'} ·
						<span class="numeric">
							{duration(data.locale, incident.started_at, incident.resolved_at)}
						</span>
					</span>
					<div class="text-xs text-stone-500">
						{dateTime(data.locale, incident.started_at)}
					</div>
				</li>
			{/each}
		</ul>
	{/if}
</section>

<p class="mt-10 text-xs text-stone-400">
	Status page by GetNotified. No third-party trackers on this page.
</p>
