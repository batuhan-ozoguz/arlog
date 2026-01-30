/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    './src/renderer/index.html',
    './src/renderer/src/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#0070f3',
          hover: '#0761d1',
          light: '#e6f2ff',
        },
        success: '#10b981',
        warning: '#f59e0b',
        error: '#ef4444',
        border: '#e5e5e5',
        'border-light': '#f0f0f0',
        'bg-primary': '#ffffff',
        'bg-secondary': '#fafafa',
        'bg-tertiary': '#f5f5f5',
        'text-primary': '#171717',
        'text-secondary': '#6b6b6b',
        'text-tertiary': '#9ca3af',
      },
      borderRadius: {
        'radius-sm': '4px',
        'radius-md': '6px',
        'radius-lg': '8px',
        'radius-xl': '12px',
      },
      boxShadow: {
        'shadow-sm': '0 1px 2px rgba(0,0,0,0.05)',
        'shadow-md': '0 4px 6px rgba(0,0,0,0.1)',
        'shadow-lg': '0 10px 15px rgba(0,0,0,0.1)',
        'shadow-xl': '0 20px 25px rgba(0,0,0,0.15)',
      }
    },
  },
  plugins: [],
}



