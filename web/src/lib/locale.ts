const DEFAULT_LOCALE = 'en';

/** Picks the visitor's preferred locale, ignoring quality weights. */
export function fromAcceptLanguage(header: string | null): string {
	const tag = header?.split(',')[0]?.split(';')[0]?.trim();
	if (!tag) return DEFAULT_LOCALE;
	try {
		return Intl.getCanonicalLocales(tag)[0] ?? DEFAULT_LOCALE;
	} catch {
		return DEFAULT_LOCALE;
	}
}
