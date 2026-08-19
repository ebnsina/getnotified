<script lang="ts">
	// A scripted run, not a random one: it always makes the same argument —
	// one failure is ignored, two in a row open an incident.
	const BLIP = 13;
	const OUTAGE = [27, 28, 29];
	const TOTAL = 44;
	const STEP_MS = 260;

	const failed = (index: number) => index === BLIP || OUTAGE.includes(index);

	const NOTES: Record<number, { text: string; tone: 'quiet' | 'down' | 'up' }> = {
		0: { text: 'Checking every 30 seconds.', tone: 'quiet' },
		[BLIP]: { text: 'One check failed. That is usually the network — nothing sent.', tone: 'quiet' },
		[OUTAGE[0]]: { text: 'One check failed.', tone: 'quiet' },
		[OUTAGE[1]]: { text: 'Two in a row. Incident opened, every channel notified.', tone: 'down' },
		[OUTAGE[2] + 1]: { text: 'Back up. One recovery notice, then quiet again.', tone: 'up' }
	};

	let cursor = $state(TOTAL);

	const note = $derived.by(() => {
		const reached = Object.keys(NOTES)
			.map(Number)
			.filter((key) => key <= cursor);
		return NOTES[Math.max(...reached, 0)];
	});

	$effect(() => {
		if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return;

		cursor = 0;
		const timer = setInterval(() => {
			cursor = cursor >= TOTAL + 8 ? 0 : cursor + 1;
		}, STEP_MS);
		return () => clearInterval(timer);
	});
</script>

<figure class="border-t border-rule">
	<figcaption class="flex items-baseline justify-between gap-4 border-b border-rule px-5 py-3">
		<span class="font-mono text-xs tracking-widest text-dim uppercase">api.example.com</span>
		<span class="font-mono text-xs text-dim">Opens an incident after 2 failures</span>
	</figcaption>

	<div class="px-5 py-10">
		<div class="flex h-14 items-end gap-[3px]" aria-hidden="true">
			{#each { length: TOTAL } as _, index (index)}
				<span
					class="w-full rounded-[1px] transition-[height] duration-200 {index > cursor
						? 'h-8 bg-night-2'
						: failed(index)
							? 'h-14 bg-down'
							: 'h-8 bg-up'}"
				></span>
			{/each}
		</div>
	</div>

	<p class="flex min-h-12 items-center border-t border-rule px-5 py-3 font-mono text-xs">
		<span class={note.tone === 'down' ? 'text-down' : note.tone === 'up' ? 'text-up' : 'text-dim'}>
			{note.text}
		</span>
	</p>
</figure>

<p class="sr-only">
	A demonstration: a single failed check is ignored, two consecutive failures open an incident, and
	recovery sends one notice.
</p>
