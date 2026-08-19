<script lang="ts">
	import Meta from '$lib/Meta.svelte';
	import Icon from '$lib/Icon.svelte';
	import { PlusSignIcon } from '@hugeicons/core-free-icons';
	import { dotClass, milliseconds, percent, statusLabel } from '$lib/format';

	let { data } = $props();
</script>

<Meta
	title="Monitors · GetNotified"
	description="Every monitor you are running, with current status and uptime."
	origin={data.origin}
	noindex
/>

<div class="flex items-center justify-between">
	<h1 class="font-display text-2xl font-normal text-bright italic">Monitors</h1>
	<a
		href="/app/monitors/new"
		class="btn btn-primary btn-sm"
	>
		<Icon icon={PlusSignIcon} size={16} strokeWidth={2} />
		New monitor
	</a>
</div>

{#if data.monitors.length === 0}
	<p class="mt-8 text-dim">
		No monitors yet. <a href="/app/monitors/new" class="underline">Add the first one</a>.
	</p>
{:else}
	<div class="mt-6 overflow-x-auto rounded-lg border border-rule bg-night-2">
		<table class="w-full text-sm">
			<thead class="border-b border-rule text-left text-xs text-dim">
				<tr>
					<th scope="col" class="px-4 py-2 font-medium">Monitor</th>
					<th scope="col" class="px-4 py-2 font-medium">Status</th>
					<th scope="col" class="px-4 py-2 font-medium">24 hours</th>
					<th scope="col" class="px-4 py-2 font-medium">7 days</th>
					<th scope="col" class="px-4 py-2 font-medium">30 days</th>
					<th scope="col" class="px-4 py-2 font-medium">Latency</th>
				</tr>
			</thead>
			<tbody>
				{#each data.monitors as monitor (monitor.id)}
					<tr class="border-b border-rule/60 last:border-0 hover:bg-night-2">
						<td class="px-4 py-3">
							<a href="/app/monitors/{monitor.id}" class="font-medium hover:underline">{monitor.name}</a>
							<div class="truncate text-xs text-dim">{monitor.target}</div>
						</td>
						<td class="px-4 py-3">
							<span class="inline-flex items-center gap-2">
								<span class="size-2 rounded-full {dotClass(monitor.status, monitor.paused)}"></span>
								{statusLabel(monitor.status, monitor.paused)}
							</span>
						</td>
						<td class="numeric px-4 py-3">{percent(data.locale, monitor.up_24h)}</td>
						<td class="numeric px-4 py-3">{percent(data.locale, monitor.up_7d)}</td>
						<td class="numeric px-4 py-3">{percent(data.locale, monitor.up_30d)}</td>
						<td class="numeric px-4 py-3 text-dim">
							{milliseconds(data.locale, monitor.latency_ms)}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}
