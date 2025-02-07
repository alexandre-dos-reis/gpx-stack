import { defineConfig } from "vite";

export default defineConfig({
  server: {
    cors: {
      origin: "http://localhost:3000",
    },
  },
  build: {
    manifest: true,
    rollupOptions: {
      input: "./assets/index.ts",
    },
    outDir: "./public/assets",
    assetsDir: ".",
    minify: true,
  },
});
