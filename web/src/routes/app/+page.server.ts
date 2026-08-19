import { api } from '$lib/server/api';
import type { MonitorSummary } from '$lib/types';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => ({
	monitors: await api<MonitorSummary[]>('/api/monitors')
});
