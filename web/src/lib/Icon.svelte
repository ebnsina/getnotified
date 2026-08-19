<script lang="ts">
	// Hugeicons' own component builds its SVG in onMount, so it never server
	// renders. The icon data is plain markup, so we render it directly.
	type IconNode = readonly (readonly [string, { readonly [key: string]: string | number }])[];

	let {
		icon,
		size = 24,
		strokeWidth = 1.5,
		class: className = ''
	}: { icon: IconNode; size?: number; strokeWidth?: number; class?: string } = $props();
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	width={size}
	height={size}
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-width={strokeWidth}
	class={className}
	aria-hidden="true"
>
	{#each icon as [tag, attrs], index (index)}
		<svelte:element this={tag} {...attrs} />
	{/each}
</svg>
