import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
    export default defineConfig({
        plugins: [
            vue(),
            vueDevTools(),
        ],
        server :{
            proxy: {
                "/upload" :{
                    target: 'http://localhost:8000',
                    changeOrigin: true,
                    secure: false,
                }
            }
        }
    })
