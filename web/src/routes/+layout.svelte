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
		<header class="sticky top-0 z-10 border-b border-rule bg-night-2/95 backdrop-blur">
			<div class="mx-auto flex max-w-5xl items-center gap-6 px-6 py-3">
				<a href="/app" class="flex items-baseline gap-2">
					<span class="size-1.5 translate-y-[-1px] rounded-full bg-up"></span>
					<span class="font-medium tracking-tight text-bright">GetNotified</span>
				</a>
				<nav class="flex gap-1 text-sm">
					{#each nav as item (item.href)}
						<a
							href={item.href}
							aria-current={page.url.pathname === item.href ? 'page' : undefined}
							class="flex items-center gap-1.5 rounded-lg px-2.5 py-1.5 transition-colors {page.url.pathname ===
							item.href
								? 'bg-night text-bright'
								: 'text-dim hover:text-bright'}"
						>
							<Icon icon={item.icon} size={16} strokeWidth={1.5} />
							{item.label}
						</a>
					{/each}
				</nav>
				<form method="POST" action="/logout" class="ml-auto">
					<button class="btn btn-secondary btn-sm">
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
