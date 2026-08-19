import { config } from '$lib/server/env';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = ({ locals }) => ({
	authenticated: locals.authenticated,
	locale: locals.locale,
	origin: config.publicOrigin
});
