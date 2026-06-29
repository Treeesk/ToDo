import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

const backendOrigin = process.env.BACKEND_ORIGIN || "http://localhost:8080";

/**
 * Vite development configuration.
 *
 * The proxy lets the React app call /api routes as same-origin requests, so
 * the Go backend can keep its HttpOnly cookie auth without adding CORS changes.
 */
export default defineConfig({
  plugins: [react()],
  server: {
    port: Number(process.env.FRONTEND_PORT || 5173),
    proxy: {
      "/api": {
        target: backendOrigin,
        changeOrigin: true,
        configure(proxy) {
          proxy.on("proxyRes", (proxyRes) => {
            const cookies = proxyRes.headers["set-cookie"];

            if (Array.isArray(cookies)) {
              proxyRes.headers["set-cookie"] = cookies.map((cookie) =>
                cookie.replace(/;\s*Secure/gi, ""),
              );
            }
          });
        },
      },
    },
  },
});
