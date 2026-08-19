<script lang="ts">
	// The last handful of results, oldest on the left. Reads at a glance without
	// asking anyone to open the monitor.
	let { recent, height = 'h-5' }: { recent: boolean[]; height?: string } = $props();
</script>

{#if recent.length > 0}
	<span class="flex items-end gap-px {height}" aria-hidden="true">
		{#each recent as ok, index (index)}
			<span class="w-1 rounded-[1px] {ok ? 'h-full bg-up/70' : 'h-full bg-down'}"></span>
		{/each}
	</span>
	<span class="sr-only">
		{recent.filter(Boolean).length} of the last {recent.length} checks succeeded.
	</span>
{:else}
	<span class="text-xs text-dim">Not checked yet</span>
{/if}
