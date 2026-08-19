<script lang="ts">
	import '../app.css';
	import { page } from '$app/state';
	import { HugeiconsIcon } from '@hugeicons/svelte';
	import { DashboardSquare01Icon, Notification03Icon, Logout03Icon } from '@hugeicons/core-free-icons';

	let { children, data } = $props();

	const nav = [
		{ href: '/', label: 'Monitors', icon: DashboardSquare01Icon },
		{ href: '/channels', label: 'Notifications', icon: Notification03Icon }
	];
</script>

<div class="min-h-screen bg-stone-50 text-stone-900">
	{#if data.authenticated}
		<header class="border-b border-stone-200 bg-white">
			<div class="mx-auto flex max-w-5xl items-center gap-6 px-6 py-4">
				<a href="/" class="font-semibold tracking-tight">GetNotified</a>
				<nav class="flex gap-4 text-sm">
					{#each nav as item (item.href)}
						<a
							href={item.href}
							aria-current={page.url.pathname === item.href ? 'page' : undefined}
							class="flex items-center gap-1.5 hover:text-stone-900 {page.url.pathname === item.href
								? 'text-stone-900'
								: 'text-stone-500'}"
						>
							<HugeiconsIcon icon={item.icon} size={16} strokeWidth={1.5} />
							{item.label}
						</a>
					{/each}
				</nav>
				<form method="POST" action="/logout" class="ml-auto">
					<button class="flex items-center gap-1.5 text-sm text-stone-500 hover:text-stone-900">
						<HugeiconsIcon icon={Logout03Icon} size={16} strokeWidth={1.5} />
						Sign out
					</button>
				</form>
			</div>
		</header>
	{/if}

	<main class="mx-auto max-w-5xl px-6 py-8">
		{@render children()}
	</main>
</div>
