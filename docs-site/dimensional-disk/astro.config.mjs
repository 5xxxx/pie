// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
    site: 'https://5xxxx.github.io',
    output: 'static',
	integrations: [
		starlight({
			title: 'Pie 文档',
            defaultLocale: 'root',
            locales: {
                root: { label: '简体中文', lang: 'zh-CN' },
                en:   { label: 'English' }
            },
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/5xxxx/pie' }],
            sidebar: [
                { label: '概览', link: '/' },
                { label: '快速开始', link: '/getting-started/' },
                { label: 'API 参考', link: '/api/' },
                { label: 'FAQ', link: '/faq/' }
            ],
            lastUpdated: true,
		}),
	],
});
