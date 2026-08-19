export type MonitorType = 'http' | 'tcp' | 'ssl_expiry';
export type Status = 'pending' | 'up' | 'down';

export interface MonitorSummary {
	id: string;
	name: string;
	type: MonitorType;
	target: string;
	status: Status;
	paused: boolean;
	tags: string[];
	up_24h: number | null;
	up_7d: number | null;
	up_30d: number | null;
	latency_ms: number | null;
}

export interface Monitor extends Omit<MonitorSummary, 'up_24h' | 'up_7d' | 'up_30d' | 'latency_ms'> {
	org_id: string;
	interval_seconds: number;
	timeout_seconds: number;
	expected_status: number[];
	ssl_warn_days: number;
	failure_threshold: number;
	consecutive_failures: number;
	next_check_at: string;
	created_at: string;
}

export interface Check {
	id: number;
	monitor_id: string;
	checked_at: string;
	ok: boolean;
	latency_ms: number;
	status_code: number | null;
	error: string | null;
	region: string;
}

export interface Incident {
	id: string;
	monitor_id: string;
	started_at: string;
	resolved_at: string | null;
	cause: string | null;
}

export type ChannelType = 'slack' | 'email' | 'webhook' | 'sms' | 'whatsapp' | 'imessage';

export interface Channel {
	id: string;
	org_id: string;
	name: string;
	type: ChannelType;
	config: Record<string, string>;
}

export interface StatusPage {
	org: { id: string; name: string; slug: string };
	monitors: MonitorSummary[];
	incidents: Incident[];
	overall: 'operational' | 'degraded';
	as_of: string;
}
