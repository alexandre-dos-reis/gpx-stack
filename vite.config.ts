import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [tailwindcss()],
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
    minify: true,
  },
});
