/** @type {import('tailwindcss').Config} */
export default {
  content: ["./web/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/forms")({ strategy: "base" })],
};
