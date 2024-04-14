/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ"],
  theme: {
    fontFamily: {
      sans: ["Poppins", "sans-serif"],
    },
    extend: {},
  },
  daisyui: {
    themes: ["cupcake"],
  },
  plugins: [require("daisyui")],
};
