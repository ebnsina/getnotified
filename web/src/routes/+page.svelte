<script lang="ts">
	import Meta from '$lib/Meta.svelte';
	import FlapDemo from '$lib/landing/FlapDemo.svelte';
	import Shortcuts from '$lib/landing/Shortcuts.svelte';
	import Icon from '$lib/Icon.svelte';
	import { ArrowRight01Icon } from '@hugeicons/core-free-icons';

	let { data } = $props();

	const entrance = '/app';

	const principles = [
		{
			title: 'It waits before it wakes you',
			body: 'One failed check is usually the network. An incident opens only after the number of consecutive failures you set — two, by default.'
		},
		{
			title: 'Every channel goes on its own',
			body: 'Slack, email, text message, WhatsApp, iMessage, or your own address. Each is sent separately, so a slow one never holds up the others.'
		},
		{
			title: 'Status pages that watch nobody',
			body: 'Open to anyone, quick to load, and carrying no tracking of any kind. Your visitors read the page; nothing on it reads them.'
		}
	];

	const steps = [
		{
			label: 'Check',
			body: 'A web address, a port on a server, or a security certificate’s expiry date — as often as you like.'
		},
		{
			label: 'Decide',
			body: 'Failures are counted in a row. Cross your threshold and an incident opens; recover and it closes.'
		},
		{
			label: 'Tell',
			body: 'Each place you chose gets its own message, and we keep trying if one does not arrive. Once when it breaks, once when it returns.'
		}
	];
</script>

<Meta
	title="GetNotified — uptime monitoring that stays quiet"
	description="Keeps an eye on your sites and services, and tells you once when something breaks and once when it comes back. Runs on your own server, with status pages that carry no tracking."
	origin={data.origin}
/>

<Shortcuts keys={{ d: entrance, s: '/status/default' }} />

<div class="min-h-screen bg-night text-mid antialiased">
	<div class="mx-auto max-w-6xl border-x border-rule/60">
		<header class="flex flex-wrap items-center justify-between gap-3 border-b border-rule px-6 py-4">
			<a href="/" class="flex items-baseline gap-2">
				<span class="size-2 translate-y-[-1px] rounded-full bg-up"></span>
				<span class="font-semibold tracking-tight text-bright">GetNotified</span>
			</a>

			<nav class="flex flex-wrap items-center gap-4 font-mono text-xs sm:gap-6">
				<a href="#how" class="text-dim hover:text-bright">How it works</a>
				<a href="/status/default" class="text-dim hover:text-bright">Live status page</a>
				<a href={entrance} class="text-bright hover:text-up">
					{data.authenticated ? 'Dashboard' : 'Sign in'}
				</a>
			</nav>
		</header>

		<section class="border-b border-rule px-6 py-16 text-center sm:py-20">
			<p class="font-mono text-xs tracking-[0.2em] text-dim uppercase">
				Uptime monitoring, on your own server
			</p>

			<h1
				class="mx-auto mt-6 max-w-3xl font-display text-4xl leading-[1.05] font-light text-bright italic sm:text-6xl"
			>
				Quiet until it matters
			</h1>

			<p class="mx-auto mt-6 max-w-xl text-base leading-relaxed text-mid sm:text-lg">
				GetNotified keeps an eye on your sites and services, and tells you once when something
				breaks, and once when it comes back. No countdowns. No badgering.
			</p>

			<div class="mt-10 flex flex-wrap items-center justify-center gap-3">
				<a
					href={entrance}
					class="btn btn-primary"
				>
					Open the dashboard
					<kbd class="rounded bg-night/20 px-1.5 py-0.5 font-mono text-[11px]">D</kbd>
				</a>
				<a
					href="/status/default"
					class="btn btn-secondary"
				>
					See a status page
					<kbd class="rounded border border-rule px-1.5 py-0.5 font-mono text-[11px] text-dim">S</kbd>
				</a>
			</div>

			<p class="mt-6 font-mono text-xs text-dim">Runs on a server you control.</p>
		</section>

		<section class="border-b border-rule">
			<FlapDemo />
		</section>

		<section class="grid border-b border-rule sm:grid-cols-3">
			{#each principles as principle, index (principle.title)}
				<article
					class="border-rule p-6 {index < principles.length - 1 ? 'sm:border-r' : ''} {index > 0
						? 'border-t sm:border-t-0'
						: ''}"
				>
					<h2 class="font-display text-xl text-bright italic">{principle.title}</h2>
					<p class="mt-3 text-sm leading-relaxed text-dim">{principle.body}</p>
				</article>
			{/each}
		</section>

		<section id="how" class="border-b border-rule px-6 py-10">
			<h2 class="font-display text-3xl font-light text-bright italic">How it works</h2>
			<p class="mt-2 max-w-lg text-sm text-dim">
				Three things happen, in this order, on every interval.
			</p>
		</section>

		<section class="grid border-b border-rule sm:grid-cols-3">
			{#each steps as step, index (step.label)}
				<article
					class="border-rule p-6 {index < steps.length - 1 ? 'sm:border-r' : ''} {index > 0
						? 'border-t sm:border-t-0'
						: ''}"
				>
					<p class="font-mono text-xs tracking-widest text-up uppercase">
						{String(index + 1).padStart(2, '0')} · {step.label}
					</p>
					<p class="mt-3 text-sm leading-relaxed text-mid">{step.body}</p>
				</article>
			{/each}
		</section>

		<section class="flex flex-wrap items-center justify-between gap-6 border-b border-rule px-6 py-12">
			<div>
				<h2 class="font-display text-3xl font-light text-bright italic">
					A status page you can hand to anyone
				</h2>
				<p class="mt-3 max-w-md text-sm text-dim">
					Open to anyone and carrying no tracking of any kind. It shows what is
					up, what is not, and what happened before.
				</p>
			</div>
			<a
				href="/status/default"
				class="btn btn-secondary"
			>
				Look at a live one
				<Icon icon={ArrowRight01Icon} size={16} strokeWidth={1.5} />
			</a>
		</section>

		<footer class="flex flex-wrap items-center gap-4 px-6 py-8 font-mono text-xs text-dim">
			<span>GetNotified</span>
			<a href="/status/default" class="hover:text-bright">Status</a>
			<a href="/login" class="hover:text-bright">Sign in</a>
		</footer>
	</div>
</div>
