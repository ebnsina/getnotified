declare global {
	namespace App {
		// Mirrors the API's error envelope so pages render one shape.
		interface Error {
			code: string;
			message: string;
		}

		interface Locals {
			authenticated: boolean;
			locale: string;
		}
	}
}

export {};
