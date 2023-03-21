/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        up: '#00BFA6',
        'up-dark': '#00887a',
        down: '#f06e6e',
        'down-dark' : '#c43c3c',
      },
      fontSize: {
        'xxs': '0.8rem',
      }
    },
    fontFamily: {
      sans: ['Rubik', 'sans-serif'],
    },
  },
  plugins: [require('daisyui')],
}
