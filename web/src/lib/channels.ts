export interface ChannelField {
	name: string;
	label: string;
	placeholder: string;
}

/** What each channel type asks for, shared by the form and the payload builder. */
export const CHANNEL_FIELDS: Record<string, ChannelField[]> = {
	slack: [
		{
			name: 'webhook_url',
			label: 'Incoming webhook URL',
			placeholder: 'https://hooks.slack.com/services/…'
		}
	],
	webhook: [
		{ name: 'url', label: 'Where should we post?', placeholder: 'https://example.com/alerts' },
		{
			name: 'secret',
			label: 'Shared secret (optional)',
			placeholder: 'Sent as the X-GetNotified-Secret header'
		}
	],
	email: [{ name: 'to', label: 'Send to', placeholder: 'you@example.com' }],
	sms: [{ name: 'to', label: 'Send to', placeholder: '+15551234567' }],
	whatsapp: [{ name: 'to', label: 'Send to', placeholder: '+15551234567' }],
	imessage: [{ name: 'to', label: 'Send to', placeholder: '+15551234567 or an Apple ID' }]
};

export const CHANNEL_TYPES = [
	{ value: 'slack', label: 'Slack' },
	{ value: 'email', label: 'Email' },
	{ value: 'webhook', label: 'Webhook' },
	{ value: 'sms', label: 'SMS' },
	{ value: 'whatsapp', label: 'WhatsApp' },
	{ value: 'imessage', label: 'iMessage' }
] as const;
