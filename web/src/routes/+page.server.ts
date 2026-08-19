import { config } from '$lib/server/env';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals }) => ({
	origin: config.publicOrigin,
	authenticated: locals.authenticated
});
