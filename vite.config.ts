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
                board: 'src/components/board/board.scss',
            },
            output: {
                entryFileNames: 'js/[name].[hash].js',
                chunkFileNames: 'js/[name].[hash].js',
                assetFileNames: (assetInfo) => {
                    const originalFileName = assetInfo.originalFileName || '';
                    if (originalFileName.match(/.*css/)) {
                        if (originalFileName.match(/components/)) {
                            return 'css/components/[name][extname]'
                        }
                        return 'css/[name].[hash][extname]'
                    }
                    return 'assets/[name][extname]'
                }
            }
        }
    }
})
