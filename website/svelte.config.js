import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter({
			pages: '../pb_public',
			assets: '../pb_public',
			fallback: 'index.html' // may differ from host to host
		})
	}
};

export default config;
