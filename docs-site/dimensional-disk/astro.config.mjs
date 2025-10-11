// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
    site: 'https://5xxxx.github.io',
    base: '/pie',
    output: 'static',
	integrations: [
		starlight({
			title: 'Pie Documentation',
            defaultLocale: 'en',
            locales: {
                en: { label: 'English', lang: 'en' },
                zh: { label: '简体中文', lang: 'zh-CN' }
            },
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/5xxxx/pie' }],
            sidebar: [
                { 
                    label: 'Overview', 
                    link: '/',
                    translations: { 'zh-CN': '概览' }
                },
                { 
                    label: 'Getting Started', 
                    link: '/getting-started/',
                    translations: { 'zh-CN': '快速开始' }
                },
                {
                    label: 'Core Features',
                    translations: { 'zh-CN': '核心功能' },
                    items: [
                        { 
                            label: 'Query Builder', 
                            link: '/core-features/query-builder/',
                            translations: { 'zh-CN': '查询构建器' }
                        },
                        { 
                            label: 'Struct Query', 
                            link: '/core-features/struct-query/',
                            translations: { 'zh-CN': '结构体查询' }
                        },
                        { 
                            label: 'Pagination', 
                            link: '/core-features/pagination/',
                            translations: { 'zh-CN': '分页查询' }
                        },
                        { 
                            label: 'Cursor Operations', 
                            link: '/core-features/cursor/',
                            translations: { 'zh-CN': '游标操作' }
                        },
                        { 
                            label: 'Bulk Operations', 
                            link: '/core-features/bulk-operations/',
                            translations: { 'zh-CN': '批量操作' }
                        },
                        { 
                            label: 'Aggregation', 
                            link: '/core-features/aggregation/',
                            translations: { 'zh-CN': '聚合查询' }
                        },
                        { 
                            label: 'Transactions', 
                            link: '/core-features/transactions/',
                            translations: { 'zh-CN': '事务管理' }
                        },
                    ]
                },
                {
                    label: 'Advanced Features',
                    translations: { 'zh-CN': '高级功能' },
                    items: [
                        { 
                            label: 'Cache Support', 
                            link: '/advanced/cache/',
                            translations: { 'zh-CN': '缓存支持' }
                        },
                        { 
                            label: 'Hook System', 
                            link: '/advanced/hooks/',
                            translations: { 'zh-CN': '钩子系统' }
                        },
                        { 
                            label: 'Soft Delete', 
                            link: '/advanced/soft-delete/',
                            translations: { 'zh-CN': '软删除' }
                        },
                        { 
                            label: 'Index Management', 
                            link: '/advanced/indexes/',
                            translations: { 'zh-CN': '索引管理' }
                        },
                        { 
                            label: 'Change Streams', 
                            link: '/advanced/change-streams/',
                            translations: { 'zh-CN': '变更流' }
                        },
                        { 
                            label: 'Query Scopes', 
                            link: '/advanced/scopes/',
                            translations: { 'zh-CN': '查询作用域' }
                        },
                        { 
                            label: 'Logging & Monitoring', 
                            link: '/advanced/logging/',
                            translations: { 'zh-CN': '日志监控' }
                        },
                        { 
                            label: 'Advanced Aggregation', 
                            link: '/advanced/aggregation-advanced/',
                            translations: { 'zh-CN': '高级聚合' }
                        },
                    ]
                },
                {
                    label: 'Reference',
                    translations: { 'zh-CN': '参考文档' },
                    items: [
                        { 
                            label: 'Configuration', 
                            link: '/reference/configuration/',
                            translations: { 'zh-CN': '配置选项' }
                        },
                        { 
                            label: 'Name Mappers', 
                            link: '/reference/mappers/',
                            translations: { 'zh-CN': '命名映射' }
                        },
                        { 
                            label: 'Error Handling', 
                            link: '/reference/error-handling/',
                            translations: { 'zh-CN': '错误处理' }
                        },
                        { 
                            label: 'Performance', 
                            link: '/reference/performance/',
                            translations: { 'zh-CN': '性能优化' }
                        },
                    ]
                },
                { 
                    label: 'Best Practices', 
                    link: '/best-practices/',
                    translations: { 'zh-CN': '最佳实践' }
                },
            ],
            lastUpdated: true,
		}),
	],
});
