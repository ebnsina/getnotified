import { fail, redirect } from '@sveltejs/kit';
import { COOKIE, cookieOpts, issueSession, verifyPassword } from '$lib/server/auth';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals }) => {
	if (locals.authenticated) redirect(303, '/app');
};

export const actions: Actions = {
	default: async ({ request, cookies, url }) => {
		const password = String((await request.formData()).get('password') ?? '');
		if (!verifyPassword(password)) {
			return fail(401, { message: 'That password did not match. Please try again.' });
		}

		cookies.set(COOKIE, issueSession(), cookieOpts);
		const next = url.searchParams.get('next');
		redirect(303, next?.startsWith('/') ? next : '/app');
	}
};
