import assert from 'node:assert/strict';
import { test } from 'node:test';
import { duration, percent } from './format.ts';

const from = (seconds: number) => new Date(Date.now() - seconds * 1000).toISOString();
const now = () => new Date().toISOString();

test('durations never render empty or with stray gaps', () => {
	for (const seconds of [0, 1, 45, 841, 3661, 7200]) {
		const text = duration('en', from(seconds), now());
		assert.notEqual(text.trim(), '', `${seconds}s rendered empty`);
		assert.doesNotMatch(text, /\s{2}/, `${seconds}s rendered with a double space: ${text}`);
	}
});

test('uptime with no data reads as a dash, not zero', () => {
	assert.equal(percent('en', null), '—');
	assert.equal(percent('en', 1), '100%');
});
