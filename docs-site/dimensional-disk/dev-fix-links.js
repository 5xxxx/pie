#!/usr/bin/env node

// 开发环境链接修复脚本
// 这个脚本会监听文件变化并自动修复链接

import { watch } from 'fs';
import { readFileSync, writeFileSync } from 'fs';
import { join } from 'path';

const distDir = './dist';

function fixLinksInFile(filePath) {
    try {
        let content = readFileSync(filePath, 'utf8');
        
        // 修复侧边栏中的链接，添加 /pie 前缀
        content = content.replace(
            /href="\/(en|zh)\/([^"]*?)"/g,
            'href="/pie/$1/$2"'
        );
        
        // 修复语言选择器中的链接
        content = content.replace(
            /value="\/(en|zh)\/"/g,
            'value="/pie/$1/"'
        );
        
        writeFileSync(filePath, content);
        console.log(`Fixed links in: ${filePath}`);
    } catch (error) {
        // 忽略文件不存在的错误
    }
}

function processDirectory(dir) {
    try {
        const items = require('fs').readdirSync(dir);
        
        for (const item of items) {
            const fullPath = join(dir, item);
            const stat = require('fs').statSync(fullPath);
            
            if (stat.isDirectory()) {
                processDirectory(fullPath);
            } else if (item.endsWith('.html')) {
                fixLinksInFile(fullPath);
            }
        }
    } catch (error) {
        // 忽略目录不存在的错误
    }
}

console.log('Starting dev link fixer...');

// 监听 dist 目录的变化
watch(distDir, { recursive: true }, (eventType, filename) => {
    if (filename && filename.endsWith('.html')) {
        const filePath = join(distDir, filename);
        setTimeout(() => fixLinksInFile(filePath), 100); // 延迟一点确保文件写入完成
    }
});

// 初始处理
setTimeout(() => {
    processDirectory(distDir);
}, 1000);
