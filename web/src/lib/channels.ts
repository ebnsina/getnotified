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
		{ name: 'url', label: 'Where should we send it?', placeholder: 'https://example.com/alerts' },
		{
			name: 'secret',
			label: 'Shared secret (optional)',
			placeholder: 'Sent with every message so you can tell it came from us'
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
	{ value: 'webhook', label: 'Your own address' },
	{ value: 'sms', label: 'Text message' },
	{ value: 'whatsapp', label: 'WhatsApp' },
	{ value: 'imessage', label: 'iMessage' }
] as const;
