/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ"],
  theme: {
    fontFamily: {
      sans: ["Poppins", "sans-serif"],
    },
    extend: {
      colors: {
        primary: "#c6c3c3",
        secondary: "#ffffff",
      },
    },
  },
  plugins: [require("daisyui")],
};
