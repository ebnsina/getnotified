import type { HandleClientError } from '@sveltejs/kit';

export const handleError: HandleClientError = ({ error }) => {
	console.error('[getnotified]', error);
	return {
		code: 'unexpected_error',
		message: 'Something went wrong in your browser. Reloading the page usually fixes it.'
	};
};
