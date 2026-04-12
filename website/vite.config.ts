import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
    plugins: [tailwindcss(), sveltekit()], server: {
        proxy: {
            // 匹配所有以 /api 开头的请求
            '/api': {
                target: 'http://127.0.0.1:8090', // PocketBase 默认地址
                changeOrigin: true,
                // 如果 PocketBase 没改 API 路径，通常它期望的是 /api/...
                // 所以这里不需要 rewrite，除非你想用不同的前缀
            }
        }
    }
});
