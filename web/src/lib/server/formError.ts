import { fail, isHttpError } from '@sveltejs/kit';

/**
 * Turns an API failure into field errors for the form. The message is the
 * API's, unchanged — it is the single source of truth for what went wrong.
 */
export function toFormError(err: unknown) {
	if (!isHttpError(err)) throw err;

	const { code, message } = err.body;
	const field = code.startsWith('invalid_') ? code.slice('invalid_'.length) : 'form';
	return fail(err.status, { errors: { [field]: message } });
}
