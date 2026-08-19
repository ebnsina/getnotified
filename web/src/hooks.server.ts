import { redirect, type Handle, type HandleServerError } from '@sveltejs/kit';
import { COOKIE, validSession } from '$lib/server/auth';
import { fromAcceptLanguage } from '$lib/locale';

// Public status pages stay reachable without a session — that is the point of them.
const isPublic = (path: string) => path === '/login' || path.startsWith('/status/');

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.authenticated = validSession(event.cookies.get(COOKIE));
	event.locals.locale = fromAcceptLanguage(event.request.headers.get('accept-language'));

	if (!event.locals.authenticated && !isPublic(event.url.pathname)) {
		redirect(303, `/login?next=${encodeURIComponent(event.url.pathname)}`);
	}
	return resolve(event, {
		transformPageChunk: ({ html }) => html.replace('%lang%', event.locals.locale)
	});
};

// Unexpected failures are logged in full and shown as one plain sentence.
export const handleError: HandleServerError = ({ error, event, status, message }) => {
	console.error('[getnotified]', event.request.method, event.url.pathname, error);
	return {
		code: 'unexpected_error',
		message:
			status === 404
				? 'We could not find that page.'
				: `Something went wrong on our side. ${message}`.trim()
	};
};
