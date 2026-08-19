<script lang="ts">
	import '../app.css';
	import { page } from '$app/state';
	import Icon from '$lib/Icon.svelte';
	import { DashboardSquare01Icon, Megaphone01Icon, Logout03Icon } from '@hugeicons/core-free-icons';

	let { children, data } = $props();

	// Only the dashboard wears the app chrome; the landing and status pages
	// bring their own.
	const inApp = $derived(page.url.pathname.startsWith('/app'));

	const nav = [
		{ href: '/app', label: 'Monitors', icon: DashboardSquare01Icon },
		{ href: '/app/channels', label: 'Notifications', icon: Megaphone01Icon }
	];
</script>

{#if inApp && data.authenticated}
	<div class="min-h-screen bg-night text-bright">
		<header class="border-b border-rule bg-night-2">
			<div class="mx-auto flex max-w-5xl items-center gap-6 px-6 py-4">
				<a href="/app" class="font-semibold tracking-tight">GetNotified</a>
				<nav class="flex gap-4 text-sm">
					{#each nav as item (item.href)}
						<a
							href={item.href}
							aria-current={page.url.pathname === item.href ? 'page' : undefined}
							class="flex items-center gap-1.5 hover:text-bright {page.url.pathname === item.href
								? 'text-bright'
								: 'text-dim'}"
						>
							<Icon icon={item.icon} size={16} strokeWidth={1.5} />
							{item.label}
						</a>
					{/each}
				</nav>
				<form method="POST" action="/logout" class="ml-auto">
					<button class="flex items-center gap-1.5 text-sm text-dim hover:text-bright">
						<Icon icon={Logout03Icon} size={16} strokeWidth={1.5} />
						Sign out
					</button>
				</form>
			</div>
		</header>

		<main class="mx-auto max-w-5xl px-6 py-8">
			{@render children()}
		</main>
	</div>
{:else}
	{@render children()}
{/if}
