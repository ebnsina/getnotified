import assert from 'node:assert/strict';
import { test } from 'node:test';
import { validateChannelForm, validateMonitorForm } from './schemas.ts';

const monitorForm = (overrides: Record<string, string> = {}) => {
	const fd = new FormData();
	const base = {
		name: 'Marketing site',
		type: 'http',
		target: 'https://example.com',
		interval_seconds: '60',
		timeout_seconds: '10',
		failure_threshold: '2',
		...overrides
	};
	for (const [key, value] of Object.entries(base)) fd.set(key, value);
	return fd;
};

test('a well formed monitor passes', () => {
	assert.deepEqual(validateMonitorForm(monitorForm()), {});
});

test('an http monitor must have a full URL', () => {
	const errors = validateMonitorForm(monitorForm({ target: 'example.com' }));
	assert.match(errors.target, /web address/);
});

test('a tcp monitor accepts a bare host and port', () => {
	assert.deepEqual(validateMonitorForm(monitorForm({ type: 'tcp', target: 'example.com:5432' })), {});
});

test('intervals below ten seconds are rejected', () => {
	const errors = validateMonitorForm(monitorForm({ interval_seconds: '5' }));
	assert.match(errors.interval_seconds, /10 seconds/);
});

test('timeouts outside one to 120 seconds are rejected', () => {
	assert.ok(validateMonitorForm(monitorForm({ timeout_seconds: '0' })).timeout_seconds);
	assert.ok(validateMonitorForm(monitorForm({ timeout_seconds: '300' })).timeout_seconds);
});

test('a slack channel needs a webhook URL', () => {
	const fd = new FormData();
	fd.set('name', 'Ops Slack');
	fd.set('type', 'slack');
	assert.ok(validateChannelForm(fd).webhook_url);

	fd.set('webhook_url', 'https://hooks.slack.com/services/abc');
	assert.deepEqual(validateChannelForm(fd), {});
});

test('an email channel needs a real address', () => {
	const fd = new FormData();
	fd.set('name', 'Me');
	fd.set('type', 'email');
	fd.set('to', 'not-an-email');
	assert.ok(validateChannelForm(fd).to);
});
