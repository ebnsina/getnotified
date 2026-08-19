import { redirect } from '@sveltejs/kit';
import { COOKIE, cookieOpts } from '$lib/server/auth';
import type { Actions } from './$types';

export const actions: Actions = {
	default: ({ cookies }) => {
		cookies.delete(COOKIE, cookieOpts);
		redirect(303, '/login');
	}
};
