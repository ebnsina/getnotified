import {
	AppleIcon,
	Certificate01Icon,
	GlobeIcon,
	Mail01Icon,
	Message01Icon,
	ServerStackIcon,
	SlackIcon,
	WebhookIcon,
	WhatsappIcon
} from '@hugeicons/core-free-icons';

export const WATCHES = [
	{
		icon: GlobeIcon,
		title: 'A web address',
		body: 'Your site, a page, an API. GetNotified asks for it and grades the reply — you decide which responses count as healthy.',
		example: 'https://example.com/health'
	},
	{
		icon: ServerStackIcon,
		title: 'A port on a server',
		body: 'Opens a connection, checks something answers, closes it again. For databases, mail servers, anything without a web page.',
		example: 'db.example.com:5432'
	},
	{
		icon: Certificate01Icon,
		title: 'A security certificate',
		body: 'Warns you before a certificate runs out, as many days ahead as you ask for. Nobody should learn about this from a visitor.',
		example: 'example.com · 30 days ahead'
	}
];

export const CHANNELS = [
	{ icon: WhatsappIcon, name: 'WhatsApp', note: 'The group chat your team is already in' },
	{ icon: SlackIcon, name: 'Slack', note: 'Into any channel, through a webhook' },
	{ icon: Message01Icon, name: 'Text message', note: 'For when nobody is at a desk' },
	{ icon: AppleIcon, name: 'iMessage', note: 'Through a relay you run on a Mac' },
	{ icon: Mail01Icon, name: 'Email', note: 'Over your own mail server' },
	{ icon: WebhookIcon, name: 'Your own address', note: 'The whole incident, as JSON' }
];

export const INCLUDED = [
	['Uptime that means something', 'Measured over 24 hours, 7 days and 30 days, from every check actually made.'],
	['A timeline of every check', 'See exactly which ones failed, how slow each was, and what came back.'],
	['Incident history', 'When it started, when it ended, how long it lasted, and what caused it.'],
	['Public status pages', 'One address you can hand to customers. No sign-in, no tracking.'],
	['Pause without deleting', 'Working on something? Pause the monitor and keep its history.'],
	['Labels', 'Group monitors however you think about them.']
];

export const STEPS = [
	{
		label: 'Check',
		body: 'A web address, a port on a server, or a certificate’s expiry date — as often as you like.'
	},
	{
		label: 'Decide',
		body: 'Failures are counted in a row. Cross your threshold and an incident opens; recover and it closes.'
	},
	{
		label: 'Tell',
		body: 'Each place you chose gets its own message, and we keep trying if one does not arrive. Once when it breaks, once when it returns.'
	}
];
