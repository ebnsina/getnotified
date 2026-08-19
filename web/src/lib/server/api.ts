import { error } from '@sveltejs/kit';
import { config } from './env';

// The API is the single source of truth for failures: it returns a message
// written for a person, and we pass it through untouched.
interface ApiError {
	error: { code: string; message: string };
}

async function toError(res: Response): Promise<never> {
	const body = (await res.json().catch(() => null)) as ApiError | null;
	error(res.status, {
		code: body?.error.code ?? 'unexpected_error',
		message: body?.error.message ?? 'Something went wrong. Please try again.'
	});
}

async function request<T>(url: string, init: RequestInit): Promise<T> {
	let res: Response;
	try {
		res = await fetch(url, init);
	} catch {
		error(503, {
			code: 'service_unreachable',
			message: 'We cannot reach the monitoring service right now. Please try again shortly.'
		});
	}
	if (!res.ok) await toError(res);
	return res.status === 204 ? (undefined as T) : res.json();
}

/** Calls the private API. The key never leaves the server. */
export function api<T>(path: string, init: RequestInit = {}): Promise<T> {
	return request<T>(config.apiUrl + path, {
		...init,
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Bearer ${config.apiKey}`,
			...init.headers
		}
	});
}

/** Public status endpoint — no key, no session. */
export function publicStatus<T>(slug: string): Promise<T> {
	return request<T>(`${config.apiUrl}/status/${encodeURIComponent(slug)}`, {});
}
