<script lang="ts">
	import Meta from '$lib/Meta.svelte';
	import Icon from '$lib/Icon.svelte';
	import FlapDemo from '$lib/landing/FlapDemo.svelte';
	import Shortcuts from '$lib/landing/Shortcuts.svelte';
	import { ArrowRight02Icon } from '@hugeicons/core-free-icons';
	import { CHANNELS, INCLUDED, STEPS, WATCHES } from '$lib/landing/content';

	let { data } = $props();
	const entrance = '/app';

	const principles = [
		{
			title: 'It waits before it wakes you',
			body: 'One failed check is usually the network. An incident opens only after the number of consecutive failures you set — two, by default.'
		},
		{
			title: 'Every channel goes on its own',
			body: 'Each message is sent separately with its own retries, so one slow channel never holds up the rest.'
		},
		{
			title: 'Nothing is measured but you',
			body: 'No third-party scripts anywhere, including on the status pages your customers read. Your data stays on your server.'
		}
	];
</script>

<Meta
	title="GetNotified — uptime monitoring that tells you on WhatsApp"
	description="Checks your sites, servers and certificates, and tells you the moment one stops answering — on WhatsApp, Slack, iMessage, text message or email. Once when it breaks, once when it returns."
	origin={data.origin}
/>

<Shortcuts keys={{ d: entrance, s: '/status/default' }} />

<div class="landing min-h-screen bg-night text-mid antialiased">
	<div class="mx-auto max-w-6xl border-x border-rule/60">
		<header class="flex flex-wrap items-center justify-between gap-3 border-b border-rule px-6 py-4">
			<a href="/" class="flex items-baseline gap-2">
				<span class="size-2 translate-y-[-1px] rounded-full bg-up"></span>
				<span class="font-display text-lg text-bright italic">GetNotified</span>
			</a>

			<nav class="flex flex-wrap items-center gap-4 font-mono text-xs sm:gap-6">
				<a href="#watches" class="text-dim hover:text-bright">What it checks</a>
				<a href="#channels" class="text-dim hover:text-bright">Where it tells you</a>
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

			<p class="mx-auto mt-6 max-w-2xl text-base leading-relaxed text-mid sm:text-lg">
				GetNotified checks your sites, servers and certificates around the clock, and tells you the
				moment one stops answering — on WhatsApp, in Slack, or as a text message. Once when it
				breaks, once when it comes back. Nothing in between.
			</p>

			<div class="mt-10 flex flex-wrap items-center justify-center gap-3">
				<a href={entrance} class="btn btn-primary">
					Open the dashboard
					<kbd class="rounded bg-night/20 px-1.5 py-0.5 font-mono text-[11px]">D</kbd>
				</a>
				<a href="/status/default" class="btn btn-secondary">
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

		<section id="watches" class="border-b border-rule px-6 py-10">
			<h2 class="font-display text-3xl font-light text-bright italic">What it checks</h2>
			<p class="mt-2 max-w-lg text-sm text-dim">Three kinds of thing, on whatever schedule you set.</p>
		</section>

		<section class="grid border-b border-rule sm:grid-cols-3">
			{#each WATCHES as watch, index (watch.title)}
				<article
					class="border-rule p-6 {index < WATCHES.length - 1 ? 'sm:border-r' : ''} {index > 0
						? 'border-t sm:border-t-0'
						: ''}"
				>
					<Icon icon={watch.icon} size={22} strokeWidth={1.5} class="text-up" />
					<h3 class="mt-3 font-medium text-bright">{watch.title}</h3>
					<p class="mt-2 text-sm leading-relaxed text-dim">{watch.body}</p>
					<p class="mt-3 font-mono text-xs text-dim">{watch.example}</p>
				</article>
			{/each}
		</section>

		<section id="channels" class="border-b border-rule px-6 py-10">
			<h2 class="font-display text-3xl font-light text-bright italic">Where it tells you</h2>
			<p class="mt-2 max-w-xl text-sm text-dim">
				Most monitoring tools stop at Slack and email. GetNotified reaches the places people
				actually look — attach as many as you like to any monitor.
			</p>
		</section>

		<section class="grid border-b border-rule sm:grid-cols-2 lg:grid-cols-3">
			{#each CHANNELS as channel, index (channel.name)}
				<article
					class="flex items-start gap-3 border-t border-rule p-6 first:border-t-0 sm:border-t
					{index % 2 === 0 ? 'sm:border-r' : ''} lg:border-t
					{index % 3 !== 2 ? 'lg:border-r' : 'lg:border-r-0'}"
				>
					<Icon icon={channel.icon} size={20} strokeWidth={1.5} class="mt-0.5 text-up" />
					<div>
						<h3 class="font-medium text-bright">{channel.name}</h3>
						<p class="mt-1 text-sm text-dim">{channel.note}</p>
					</div>
				</article>
			{/each}
		</section>

		<section class="border-b border-rule px-6 py-10">
			<h2 class="font-display text-3xl font-light text-bright italic">How it works</h2>
			<p class="mt-2 max-w-lg text-sm text-dim">
				Three things happen, in this order, on every interval.
			</p>
		</section>

		<section class="grid border-b border-rule sm:grid-cols-3">
			{#each STEPS as step, index (step.label)}
				<article
					class="border-rule p-6 {index < STEPS.length - 1 ? 'sm:border-r' : ''} {index > 0
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

		<section class="border-b border-rule px-6 py-10">
			<h2 class="font-display text-3xl font-light text-bright italic">What you get</h2>
		</section>

		<section class="grid border-b border-rule sm:grid-cols-2 lg:grid-cols-3">
			{#each INCLUDED as [title, body], index (title)}
				<article
					class="border-t border-rule p-6 first:border-t-0 sm:border-t
					{index % 2 === 0 ? 'sm:border-r' : ''} lg:border-t
					{index % 3 !== 2 ? 'lg:border-r' : 'lg:border-r-0'}"
				>
					<h3 class="font-medium text-bright">{title}</h3>
					<p class="mt-2 text-sm leading-relaxed text-dim">{body}</p>
				</article>
			{/each}
		</section>

		<section class="flex flex-wrap items-center justify-between gap-6 border-b border-rule px-6 py-12">
			<div>
				<h2 class="font-display text-3xl font-light text-bright italic">
					A status page you can hand to anyone
				</h2>
				<p class="mt-3 max-w-md text-sm text-dim">
					Open to anyone and carrying no tracking of any kind. It shows what is up, what is not, and
					what happened before.
				</p>
			</div>
			<a href="/status/default" class="btn btn-secondary">
				Look at a live one
				<Icon icon={ArrowRight02Icon} size={16} strokeWidth={1.5} />
			</a>
		</section>

		<section class="border-b border-rule px-6 py-12">
			<h2 class="font-display text-3xl font-light text-bright italic">Everything here is an API</h2>
			<p class="mt-3 max-w-xl text-sm text-dim">
				The dashboard has no private back door — it asks the same questions you can. Add a monitor
				from a deploy script, or pull uptime into a report of your own.
			</p>

			<pre class="mt-6 overflow-x-auto rounded-xl border border-rule bg-night-2 p-5 font-mono text-xs leading-relaxed text-mid"><code
					>curl -X POST https://your-server/api/monitors \
  -H "Authorization: Bearer $API_KEY" \
  -d '&lbrace;"name": "Checkout", "target": "https://example.com/health"&rbrace;'</code
				></pre>
		</section>

		<footer class="flex flex-wrap items-center gap-4 px-6 py-8 font-mono text-xs text-dim">
			<span>GetNotified</span>
			<a href="#watches" class="hover:text-bright">What it checks</a>
			<a href="#channels" class="hover:text-bright">Where it tells you</a>
			<a href="/status/default" class="hover:text-bright">Status</a>
			<a href="/login" class="hover:text-bright">Sign in</a>
		</footer>
	</div>
</div>
