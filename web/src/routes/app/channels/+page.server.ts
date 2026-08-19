import { api } from '$lib/server/api';
import { channelPayload } from '$lib/server/forms';
import { toFormError } from '$lib/server/formError';
import type { Channel } from '$lib/types';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async () => ({
	channels: await api<Channel[]>('/api/channels')
});

export const actions: Actions = {
	create: async ({ request }) => {
		try {
			await api('/api/channels', {
				method: 'POST',
				body: JSON.stringify(channelPayload(await request.formData()))
			});
		} catch (err) {
			return toFormError(err);
		}
		return { message: 'Channel added.' };
	},

	delete: async ({ request }) => {
		const id = String((await request.formData()).get('id') ?? '');
		try {
			await api(`/api/channels/${id}`, { method: 'DELETE' });
		} catch (err) {
			return toFormError(err);
		}
		return { message: 'Channel removed.' };
	}
};
