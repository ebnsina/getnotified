import { createHmac, randomBytes, scryptSync, timingSafeEqual } from 'node:crypto';
import { config } from './env';

export const COOKIE = 'gn_session';
const TTL_MS = 7 * 24 * 60 * 60 * 1000;
const KEY_LENGTH = 64;

/** `salt:hash`, produced by scripts/hash-password.js. */
export function hashPassword(password: string): string {
	const salt = randomBytes(16).toString('hex');
	return `${salt}:${scryptSync(password, salt, KEY_LENGTH).toString('hex')}`;
}

export function verifyPassword(password: string): boolean {
	const [salt, hash] = config.authPasswordHash.split(':');
	if (!salt || !hash) return false;
	const expected = Buffer.from(hash, 'hex');
	return equal(expected, scryptSync(password, salt, expected.length));
}

/**
 * Stateless session: `<expiry>.<hmac>`. There is one account, so rotating
 * AUTH_SECRET is what revokes every cookie.
 */
export function issueSession(): string {
	const expiry = String(Date.now() + TTL_MS);
	return `${expiry}.${sign(expiry)}`;
}

export function validSession(token: string | undefined): boolean {
	const [expiry, mac] = token?.split('.') ?? [];
	if (!expiry || !mac || Number(expiry) < Date.now()) return false;
	return equal(Buffer.from(sign(expiry)), Buffer.from(mac));
}

function sign(value: string): string {
	return createHmac('sha256', config.authSecret).update(value).digest('hex');
}

function equal(a: Buffer, b: Buffer): boolean {
	return a.length === b.length && timingSafeEqual(a, b);
}

export const cookieOpts = {
	path: '/',
	httpOnly: true,
	sameSite: 'lax',
	secure: process.env.NODE_ENV === 'production',
	maxAge: TTL_MS / 1000
} as const;
