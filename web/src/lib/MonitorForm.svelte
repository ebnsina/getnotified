<script lang="ts">
	import type { Monitor } from '$lib/types';

	let {
		monitor = null,
		errors = {}
	}: { monitor?: Monitor | null; errors?: Record<string, string> } = $props();

	// Seeded from the monitor once; after that the select owns it.
	const initialType = monitor?.type ?? 'http';
	let type = $state(initialType);

	const field = 'field mt-1';
</script>

{#snippet error(name: string)}
	{#if errors[name]}
		<p id="{name}-error" class="mt-1 text-sm text-down">{errors[name]}</p>
	{/if}
{/snippet}

<div class="space-y-4">
	<div>
		<label for="name" class="block text-sm font-medium">Name</label>
		<input
			id="name"
			name="name"
			required
			value={monitor?.name ?? ''}
			aria-invalid={errors.name ? 'true' : undefined}
			class={field}
		/>
		{@render error('name')}
	</div>

	<div class="grid gap-4 sm:grid-cols-2">
		<div>
			<label for="type" class="block text-sm font-medium">Check type</label>
			<select id="type" name="type" bind:value={type} class="field field-select mt-1">
				<option value="http">A web address</option>
				<option value="tcp">A port on a server</option>
				<option value="ssl_expiry">A security certificate</option>
			</select>
		</div>
		<div>
			<label for="target" class="block text-sm font-medium">What to check</label>
			<input
				id="target"
				name="target"
				required
				value={monitor?.target ?? ''}
				placeholder={type === 'http' ? 'https://example.com/health' : 'example.com:443'}
				aria-invalid={errors.target ? 'true' : undefined}
				class="{field} font-mono"
			/>
			{@render error('target')}
		</div>
	</div>

	<div class="grid gap-4 sm:grid-cols-3">
		<div>
			<label for="interval_seconds" class="block text-sm font-medium">Check every (seconds)</label>
			<input
				id="interval_seconds"
				name="interval_seconds"
				type="number"
				min="10"
				value={monitor?.interval_seconds ?? 60}
				class="{field} numeric"
			/>
			{@render error('interval_seconds')}
		</div>
		<div>
			<label for="timeout_seconds" class="block text-sm font-medium">Give up after (seconds)</label>
			<input
				id="timeout_seconds"
				name="timeout_seconds"
				type="number"
				min="1"
				max="120"
				value={monitor?.timeout_seconds ?? 10}
				class="{field} numeric"
			/>
			{@render error('timeout_seconds')}
		</div>
		<div>
			<label for="failure_threshold" class="block text-sm font-medium">Failures before alert</label>
			<input
				id="failure_threshold"
				name="failure_threshold"
				type="number"
				min="1"
				value={monitor?.failure_threshold ?? 2}
				class="{field} numeric"
			/>
			<p class="mt-1 text-xs text-dim">
				How many checks in a row must fail before anyone hears about it.
			</p>
			{@render error('failure_threshold')}
		</div>
	</div>

	{#if type === 'http'}
		<div>
			<label for="expected_status" class="block text-sm font-medium">Replies that count as healthy</label>
			<input
				id="expected_status"
				name="expected_status"
				value={(monitor?.expected_status ?? [200]).join(', ')}
				class="{field} numeric"
			/>
		</div>
	{:else if type === 'ssl_expiry'}
		<div>
			<label for="ssl_warn_days" class="block text-sm font-medium">Warn this many days ahead</label>
			<input
				id="ssl_warn_days"
				name="ssl_warn_days"
				type="number"
				min="1"
				value={monitor?.ssl_warn_days ?? 14}
				class="{field} numeric"
			/>
			{@render error('ssl_warn_days')}
		</div>
	{/if}

	<div>
		<label for="tags" class="block text-sm font-medium">Labels</label>
		<input id="tags" name="tags" value={(monitor?.tags ?? []).join(', ')} class={field} />
	</div>
</div>
