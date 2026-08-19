import * as v from 'valibot';

// Mirrors the API's rules so the form can answer before a round trip. The API
// stays authoritative; this only saves the person a wasted submit.
export const MonitorSchema = v.object({
	name: v.pipe(
		v.string(),
		v.trim(),
		v.minLength(1, 'Give the monitor a name so you can recognise it later.')
	),
	type: v.picklist(['http', 'tcp', 'ssl_expiry'], 'Choose one of HTTP, TCP port, or SSL expiry.'),
	target: v.pipe(v.string(), v.trim(), v.minLength(1, 'Tell us what to check — a URL or a host.')),
	interval_seconds: v.pipe(
		v.number('Enter how often to check, in seconds.'),
		v.minValue(10, 'Check no more often than once every 10 seconds.')
	),
	timeout_seconds: v.pipe(
		v.number('Enter a timeout in seconds.'),
		v.minValue(1, 'Set a timeout between 1 and 120 seconds.'),
		v.maxValue(120, 'Set a timeout between 1 and 120 seconds.')
	),
	failure_threshold: v.pipe(
		v.number('Enter how many failures should open an incident.'),
		v.minValue(1, 'It takes at least one failure to open an incident.')
	),
	ssl_warn_days: v.optional(
		v.pipe(
			v.number('Enter a number of days.'),
			v.minValue(1, 'Warn at least one day before the certificate expires.')
		)
	)
});

export const HttpTargetSchema = v.pipe(
	v.string(),
	v.url('An HTTP monitor needs a full address, like https://example.com.')
);

export const ChannelSchema = v.object({
	name: v.pipe(v.string(), v.trim(), v.minLength(1, 'Give the channel a name, such as “Ops Slack”.')),
	type: v.picklist(['slack', 'email', 'webhook', 'sms', 'whatsapp', 'imessage'])
});

export const EmailSchema = v.pipe(v.string(), v.email('Enter a valid email address.'));
export const UrlSchema = v.pipe(v.string(), v.url('Enter a full address, starting with https://.'));

/** Maps valibot issues to `{ field: message }` for rendering next to inputs. */
export function fieldErrors(
	issues: [v.BaseIssue<unknown>, ...v.BaseIssue<unknown>[]]
): Record<string, string> {
	const flat = v.flatten(issues).nested ?? {};
	return Object.fromEntries(
		Object.entries(flat).map(([field, messages]) => [field, messages?.[0] ?? ''])
	);
}

/** Validates the monitor form before it is submitted. Empty means it is fine. */
export function validateMonitorForm(fd: FormData): Record<string, string> {
	const number = (key: string) => Number(fd.get(key) ?? Number.NaN);
	const type = String(fd.get('type') ?? 'http');

	const result = v.safeParse(MonitorSchema, {
		name: String(fd.get('name') ?? ''),
		type,
		target: String(fd.get('target') ?? ''),
		interval_seconds: number('interval_seconds'),
		timeout_seconds: number('timeout_seconds'),
		failure_threshold: number('failure_threshold'),
		ssl_warn_days: fd.has('ssl_warn_days') ? number('ssl_warn_days') : undefined
	});

	const errors = result.success ? {} : fieldErrors(result.issues);
	if (!errors.target && type === 'http') {
		const target = v.safeParse(HttpTargetSchema, String(fd.get('target') ?? ''));
		if (!target.success) errors.target = target.issues[0].message;
	}
	return errors;
}

/** Validates the channel form, including the destination its type requires. */
export function validateChannelForm(fd: FormData): Record<string, string> {
	const type = String(fd.get('type') ?? '');
	const result = v.safeParse(ChannelSchema, { name: String(fd.get('name') ?? ''), type });
	const errors = result.success ? {} : fieldErrors(result.issues);

	if (type === 'slack') assign(errors, 'webhook_url', fd, UrlSchema);
	if (type === 'webhook') assign(errors, 'url', fd, UrlSchema);
	if (type === 'email') assign(errors, 'to', fd, EmailSchema);
	return errors;
}

function assign(
	errors: Record<string, string>,
	field: string,
	fd: FormData,
	schema: v.BaseSchema<string, string, v.BaseIssue<unknown>>
) {
	const result = v.safeParse(schema, String(fd.get(field) ?? ''));
	if (!result.success) errors[field] = result.issues[0].message;
}
