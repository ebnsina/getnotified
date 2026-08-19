<script lang="ts">
	import Meta from '$lib/Meta.svelte';
	import Icon from '$lib/Icon.svelte';
	import { PlusSignIcon, DashboardSquare01Icon } from '@hugeicons/core-free-icons';
	import { milliseconds, percent, pillClass, statusLabel } from '$lib/format';

	let { data } = $props();

	const down = $derived(data.monitors.filter((m) => !m.paused && m.status === 'down').length);
</script>

<Meta
	title="Monitors · GetNotified"
	description="Every monitor you are running, with current status and uptime."
	origin={data.origin}
	noindex
/>

<div class="flex flex-wrap items-end justify-between gap-4">
	<div>
		<h1 class="font-display text-2xl font-normal text-bright italic">Monitors</h1>
		<p class="mt-1 text-sm text-dim">
			{#if data.monitors.length === 0}
				Nothing is being checked yet.
			{:else if down === 0}
				All {data.monitors.length} responding normally.
			{:else}
				{down} of {data.monitors.length} not responding.
			{/if}
		</p>
	</div>
	<a href="/app/monitors/new" class="btn btn-primary btn-sm">
		<Icon icon={PlusSignIcon} size={16} strokeWidth={2} />
		New monitor
	</a>
</div>

{#if data.monitors.length === 0}
	<div class="panel mt-8 px-6 py-16 text-center">
		<Icon icon={DashboardSquare01Icon} size={28} strokeWidth={1.5} class="mx-auto text-dim" />
		<h2 class="mt-4 font-medium text-bright">Add the first thing to watch</h2>
		<p class="mx-auto mt-2 max-w-sm text-sm text-dim">
			A web address, a port on a server, or a certificate. You choose how often it is checked and
			who hears about it.
		</p>
		<a href="/app/monitors/new" class="btn btn-primary mt-6">Create a monitor</a>
	</div>
{:else}
	<div class="panel mt-6 overflow-x-auto">
		<table class="w-full text-sm">
			<thead class="text-left text-xs text-dim">
				<tr class="border-b border-rule">
					<th scope="col" class="px-5 py-3 font-normal">Monitor</th>
					<th scope="col" class="px-5 py-3 font-normal">Status</th>
					<th scope="col" class="px-5 py-3 text-right font-normal">24 hours</th>
					<th scope="col" class="px-5 py-3 text-right font-normal">7 days</th>
					<th scope="col" class="px-5 py-3 text-right font-normal">30 days</th>
					<th scope="col" class="px-5 py-3 text-right font-normal">Latency</th>
				</tr>
			</thead>
			<tbody>
				{#each data.monitors as monitor (monitor.id)}
					<tr class="border-b border-rule/60 transition-colors last:border-0 hover:bg-night">
						<td class="px-5 py-4">
							<a href="/app/monitors/{monitor.id}" class="font-medium text-bright hover:text-up">
								{monitor.name}
							</a>
							<div class="mt-0.5 truncate font-mono text-xs text-dim">{monitor.target}</div>
						</td>
						<td class="px-5 py-4">
							<span class={pillClass(monitor.status, monitor.paused)}>
								{statusLabel(monitor.status, monitor.paused)}
							</span>
						</td>
						<td class="numeric px-5 py-4 text-right text-mid">{percent(data.locale, monitor.up_24h)}</td>
						<td class="numeric px-5 py-4 text-right text-mid">{percent(data.locale, monitor.up_7d)}</td>
						<td class="numeric px-5 py-4 text-right text-mid">{percent(data.locale, monitor.up_30d)}</td>
						<td class="numeric px-5 py-4 text-right text-dim">
							{milliseconds(data.locale, monitor.latency_ms)}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}
