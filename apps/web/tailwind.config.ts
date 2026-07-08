import type { Config } from 'tailwindcss';
import aitosPreset from '@aitos/design-system/tailwind-preset';

const config: Config = {
  darkMode: 'class',
  content: [
    './src/**/*.{js,ts,jsx,tsx,mdx}',
    '../../packages/design-system/src/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {},
  },
  presets: [aitosPreset],
  plugins: [],
};

export default config;
