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
        down: '#fb4545',
        'down-dark': '#c43c3c',
      },
      fontSize: {
        'xxs': '0.8rem',
        'xxxs': '0.6rem',
      },
    },
    fontFamily: {
      sans: ['Rubik', 'sans-serif'],
    },
  },
}
