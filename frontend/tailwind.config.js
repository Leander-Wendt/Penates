/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        'brut-black': '#000000',
        'brut-white': '#FFFFFF',
        'brut-cream': '#FDF6E3',
        'brut-yellow': '#FFD400',
        'brut-blue': '#1D4ED8',
        'brut-red': '#B3001B',
        'brut-green': '#007A33',
        'brut-pink': '#FF2D78',
      },
      borderWidth: {
        3: '3px',
      },
      boxShadow: {
        brut: '4px 4px 0 0 #000000',
        'brut-sm': '2px 2px 0 0 #000000',
        'brut-lg': '8px 8px 0 0 #000000',
        'brut-pressed': '1px 1px 0 0 #000000',
      },
      fontFamily: {
        sans: [
          'Inter',
          'ui-sans-serif',
          'system-ui',
          '-apple-system',
          'Segoe UI',
          'Roboto',
          'Helvetica Neue',
          'Arial',
          'sans-serif',
        ],
      },
      borderRadius: {
        none: '0px',
        sm: '2px',
        DEFAULT: '2px',
      },
    },
  },
  plugins: [],
}
