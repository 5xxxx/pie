#!/usr/bin/env node

import { readFileSync, writeFileSync, readdirSync, statSync } from 'fs';
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
        console.error(`Error processing ${filePath}:`, error.message);
    }
}

function processDirectory(dir) {
    const items = readdirSync(dir);
    
    for (const item of items) {
        const fullPath = join(dir, item);
        const stat = statSync(fullPath);
        
        if (stat.isDirectory()) {
            processDirectory(fullPath);
        } else if (item.endsWith('.html')) {
            fixLinksInFile(fullPath);
        }
    }
}

console.log('Fixing links in generated HTML files...');
processDirectory(distDir);
console.log('Link fixing completed!');
