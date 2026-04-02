/** @type {import('tailwindcss').Config} */
export default {
  // HTMLとJSファイルをスキャン対象にする
  content: ["./index.html", "./src/**/*.{html,js,ts}"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  // daisyUIの設定（任意）
  daisyui: {
    themes: ["light", "dark", "cupcake"],
  },
}