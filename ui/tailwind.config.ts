import type { Config } from "tailwindcss";
import tailwindTypography from "@tailwindcss/typography";
import daisyui from "daisyui";

const config = {
  content: ["./src/**/*.{html,js,svelte,ts}"],
  theme: {
    extend: {
      maxHeight: {
        select: "18rem",
      },
    },
    fontFamily: {
      sans: ["Inter", "ui-sans-serif"],
    },

    keyframes: {
      progress: {
        "0%": { transform: "translateX(-100%) scaleX(1)" },
        "40%": { transform: "translateX(10w) scaleX(0.5)" },
        "100%": { transform: "translateX(100vw) scaleX(0.7)" },
      },
    },

    animation: {
      progress: "progress 1s linear infinite",
      progressFast: "progress 0.8s linear infinite",
    },
  },
  plugins: [tailwindTypography, daisyui],
  daisyui: {
    themes: [
      {
        beans: {
          primary: "#2F855A",
          "primary-content": "#ffffff",
          secondary: "#CBD5E0",
          accent: "#9AE6B4",
          neutral: "#374151",
          error: "#9B2C2C",
          success: "#48BB78",
          warning: "#ECC94B",
          "base-100": "#f9fafb",
          "base-200": "#f3f4f6",
          "base-300": "#e5e7eb",

          "--rounded-btn": "0.25rem",
        },
      },
      {
        beansDark: {
          primary: "#307050",
          "primary-content": "#ffffff",
          secondary: "#343b48",
          accent: "#9AE6B4",
          neutral: "#0E0E12",
          error: "#FC8181",
          success: "#48BB78",
          warning: "#ECC94B",
          "base-100": "#26262E",
          "base-200": "#202027",
          "base-300": "#1A1A20",

          "--rounded-btn": "0.25rem",
        },
      },
    ],
  },
} satisfies Config;

export default config;
