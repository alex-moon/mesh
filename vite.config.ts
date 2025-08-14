import { defineConfig } from 'vite'

export default defineConfig({
    base: '/static/',
    publicDir: 'src/assets',
    build: {
        outDir: 'static',
        emptyOutDir: true,
        manifest: true,
        rollupOptions: {
            input: {
                main: 'src/main.ts',
                styles: 'src/main.scss',
                app: 'src/components/app/app.scss',
            },
            output: {
                entryFileNames: 'js/[name].[hash].js',
                chunkFileNames: 'js/[name].[hash].js',
                assetFileNames: (assetInfo) => {
                    if ((assetInfo as any).name?.endsWith('.css')) {
                        const name = (assetInfo as any).name;
                        if (name === 'app.css') {
                            return 'css/components/[name].[hash][extname]'
                        }
                        return 'css/[name].[hash][extname]'
                    }
                    return 'assets/[name][extname]'
                }
            }
        }
    }
})
