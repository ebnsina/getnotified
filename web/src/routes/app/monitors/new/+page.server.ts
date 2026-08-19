import { redirect } from '@sveltejs/kit';
import { api } from '$lib/server/api';
import { monitorPayload } from '$lib/server/forms';
import { toFormError } from '$lib/server/formError';
import type { Monitor } from '$lib/types';
import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ request }) => {
		const payload = monitorPayload(await request.formData());

		let created: Monitor;
		try {
			created = await api<Monitor>('/api/monitors', {
				method: 'POST',
				body: JSON.stringify(payload)
			});
		} catch (err) {
			return toFormError(err);
		}
		redirect(303, `/monitors/${created.id}`);
	}
};
