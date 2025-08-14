import { defineConfig } from 'vite'

export default defineConfig({
    publicDir: 'src/assets',
    build: {
        outDir: '../static',
        emptyOutDir: true,
        manifest: true,
        rollupOptions: {
            input: {
                main: 'src/main.ts',
                styles: 'src/styles/main.scss',
                app: 'src/components/app/app.scss',
            },
            output: {
                entryFileNames: 'js/[name].js',
                chunkFileNames: 'js/[name].js',
                assetFileNames: (assetInfo) => {
                    if ((assetInfo as any).name?.endsWith('.css')) {
                        const name = (assetInfo as any).name;
                        if (name === 'app.css') {
                            return 'css/components/[name][extname]'
                        }
                        return 'css/[name][extname]'
                    }
                    return 'assets/[name][extname]'
                }
            }
        }
    }
})
