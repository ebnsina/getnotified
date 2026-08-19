import { publicStatus } from '$lib/server/api';
import type { StatusPage } from '$lib/types';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, setHeaders }) => {
	setHeaders({ 'cache-control': 'public, max-age=30' });
	return { status: await publicStatus<StatusPage>(params.slug) };
};
