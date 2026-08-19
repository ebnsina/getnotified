import { env } from '$env/dynamic/private';

const KEYS = ['API_URL', 'API_KEY', 'AUTH_SECRET', 'AUTH_PASSWORD_HASH', 'SITE_ORIGIN'] as const;

type Key = (typeof KEYS)[number];

let loaded: Record<Key, string> | null = null;

/**
 * Read on first use, not at import, so a build does not need the values.
 * Nothing has a default: a missing value is a misconfigured deployment.
 */
function load(): Record<Key, string> {
	if (loaded) return loaded;

	const missing = KEYS.filter((key) => !env[key]);
	if (missing.length > 0) {
		throw new Error(`Missing required environment variables: ${missing.join(', ')}`);
	}
	loaded = Object.fromEntries(KEYS.map((key) => [key, env[key]])) as Record<Key, string>;
	return loaded;
}

export const config = {
	get apiUrl() {
		return load().API_URL;
	},
	get apiKey() {
		return load().API_KEY;
	},
	get authSecret() {
		return load().AUTH_SECRET;
	},
	get authPasswordHash() {
		return load().AUTH_PASSWORD_HASH;
	},
	get publicOrigin() {
		return load().SITE_ORIGIN;
	}
} as const;
