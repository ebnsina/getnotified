import type { Status } from './types';

const MINUTE = 60;
const HOUR = 60 * MINUTE;
const DAY = 24 * HOUR;

// Formatters are cached per locale; building them is the expensive part.
const cache = new Map<string, Intl.NumberFormat | Intl.RelativeTimeFormat | Intl.DateTimeFormat>();

function memo<T>(key: string, build: () => T): T {
	if (!cache.has(key)) cache.set(key, build() as never);
	return cache.get(key) as T;
}

export function percent(locale: string, value: number | null): string {
	if (value === null) return '—';
	return memo(`p:${locale}`, () =>
		new Intl.NumberFormat(locale, { style: 'percent', maximumFractionDigits: 2 })
	).format(value);
}

export function milliseconds(locale: string, value: number | null): string {
	if (value === null) return '—';
	return memo(`ms:${locale}`, () =>
		new Intl.NumberFormat(locale, { style: 'unit', unit: 'millisecond', unitDisplay: 'short' })
	).format(value);
}

export function dateTime(locale: string, iso: string): string {
	return memo(`dt:${locale}`, () =>
		new Intl.DateTimeFormat(locale, { dateStyle: 'medium', timeStyle: 'short' })
	).format(new Date(iso));
}

export function relative(locale: string, iso: string): string {
	const seconds = Math.round((Date.parse(iso) - Date.now()) / 1000);
	const format = memo(`rt:${locale}`, () =>
		new Intl.RelativeTimeFormat(locale, { numeric: 'auto' })
	);

	const magnitude = Math.abs(seconds);
	if (magnitude >= DAY) return format.format(Math.round(seconds / DAY), 'day');
	if (magnitude >= HOUR) return format.format(Math.round(seconds / HOUR), 'hour');
	if (magnitude >= MINUTE) return format.format(Math.round(seconds / MINUTE), 'minute');
	return format.format(seconds, 'second');
}

export function duration(locale: string, from: string, to: string | null): string {
	const total = Math.max(
		0,
		Math.round(((to ? Date.parse(to) : Date.now()) - Date.parse(from)) / 1000)
	);
	// A duration of nothing still needs to read as a number, and DurationFormat
	// drops zero units entirely.
	if (total < MINUTE) {
		return memo(`s:${locale}`, () =>
			new Intl.NumberFormat(locale, { style: 'unit', unit: 'second', unitDisplay: 'narrow' })
		).format(total);
	}

	// Zero units are dropped too, otherwise they format as stray gaps.
	const parts = Object.fromEntries(
		Object.entries({
			hours: Math.floor(total / HOUR),
			minutes: Math.floor((total % HOUR) / MINUTE),
			seconds: total % MINUTE
		}).filter(([, value]) => value > 0)
	);
	return new Intl.DurationFormat(locale, { style: 'narrow' }).format(parts);
}

export function list(locale: string, items: string[]): string {
	return new Intl.ListFormat(locale, { style: 'short', type: 'conjunction' }).format(items);
}

export const pillClass = (status: Status, paused = false) =>
	paused ? 'pill pill-idle' : status === 'up' ? 'pill pill-up' : status === 'down' ? 'pill pill-down' : 'pill pill-idle';

export const statusLabel = (status: Status, paused = false) =>
	paused ? 'Paused' : status === 'up' ? 'Up' : status === 'down' ? 'Down' : 'Not checked yet';
