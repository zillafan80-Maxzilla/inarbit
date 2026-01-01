/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,jsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#9D8B7E',
        secondary: '#5F5550',
        background: '#F5F1E8',
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
