import { CHANNEL_FIELDS } from '$lib/channels';

/**
 * Builds the monitor payload. Absent keys stay undefined so a PATCH leaves the
 * stored value alone.
 */
export function monitorPayload(fd: FormData) {
	const text = (key: string) => (fd.has(key) ? String(fd.get(key)).trim() : undefined);
	const number = (key: string) => (text(key) ? Number(text(key)) : undefined);
	const list = (key: string) =>
		fd.has(key)
			? String(fd.get(key))
					.split(',')
					.map((item) => item.trim())
					.filter(Boolean)
			: undefined;

	return {
		name: text('name'),
		type: text('type'),
		target: text('target'),
		interval_seconds: number('interval_seconds'),
		timeout_seconds: number('timeout_seconds'),
		failure_threshold: number('failure_threshold'),
		ssl_warn_days: number('ssl_warn_days'),
		expected_status: list('expected_status')?.map(Number),
		tags: list('tags') ?? []
	};
}

/** Keeps only the config keys the chosen channel type actually uses. */
export function channelPayload(fd: FormData) {
	const type = String(fd.get('type') ?? '');
	const config: Record<string, string> = {};

	for (const field of CHANNEL_FIELDS[type] ?? []) {
		const value = String(fd.get(field.name) ?? '').trim();
		if (value) config[field.name] = value;
	}
	return { name: String(fd.get('name') ?? '').trim(), type, config };
}
