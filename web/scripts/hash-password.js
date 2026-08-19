#!/usr/bin/env node
// Usage: node scripts/hash-password.js 'your password'  ->  AUTH_PASSWORD_HASH value
import { randomBytes, scryptSync } from 'node:crypto';

const password = process.argv[2];
if (!password) {
	console.error("usage: node scripts/hash-password.js 'your password'");
	process.exit(1);
}
const salt = randomBytes(16).toString('hex');
console.log(`${salt}:${scryptSync(password, salt, 64).toString('hex')}`);
