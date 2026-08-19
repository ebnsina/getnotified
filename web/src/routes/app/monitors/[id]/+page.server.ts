import { redirect } from '@sveltejs/kit';
import { api } from '$lib/server/api';
import { monitorPayload } from '$lib/server/forms';
import { toFormError } from '$lib/server/formError';
import type { Channel, Check, Incident, Monitor } from '$lib/types';
import type { Actions, PageServerLoad } from './$types';

const CHECK_HISTORY = 60;

export const load: PageServerLoad = async ({ params }) => {
	const [monitor, checks, incidents, channels, attached] = await Promise.all([
		api<Monitor>(`/api/monitors/${params.id}`),
		api<Check[]>(`/api/monitors/${params.id}/checks?limit=${CHECK_HISTORY}`),
		api<Incident[]>(`/api/monitors/${params.id}/incidents`),
		api<Channel[]>('/api/channels'),
		api<string[]>(`/api/monitors/${params.id}/channels`)
	]);
	return { monitor, checks, incidents, channels, attached };
};

export const actions: Actions = {
	update: async ({ request, params }) => {
		try {
			await api(`/api/monitors/${params.id}`, {
				method: 'PATCH',
				body: JSON.stringify(monitorPayload(await request.formData()))
			});
		} catch (err) {
			return toFormError(err);
		}
		return { message: 'Your changes are saved.' };
	},

	pause: async ({ request, params }) => {
		const paused = (await request.formData()).get('paused') === 'true';
		try {
			await api(`/api/monitors/${params.id}`, {
				method: 'PATCH',
				body: JSON.stringify({ paused })
			});
		} catch (err) {
			return toFormError(err);
		}
		return { message: paused ? 'Checks are paused.' : 'Checks have resumed.' };
	},

	channels: async ({ request, params }) => {
		const channel_ids = (await request.formData()).getAll('channel_ids').map(String);
		try {
			await api(`/api/monitors/${params.id}/channels`, {
				method: 'PUT',
				body: JSON.stringify({ channel_ids })
			});
		} catch (err) {
			return toFormError(err);
		}
		return { message: 'Notification channels updated.' };
	},

	delete: async ({ params }) => {
		try {
			await api(`/api/monitors/${params.id}`, { method: 'DELETE' });
		} catch (err) {
			return toFormError(err);
		}
		redirect(303, '/app');
	}
};
