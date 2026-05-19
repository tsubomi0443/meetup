/** @type {import('tailwindcss').Config} */
export default {
  // HTMLとJSファイルをスキャン対象にする
  content: [
    "../static/js/**/*.js",
    "../templates/**/*.html",
  ],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  // daisyUIの設定（任意）
  daisyui: {
    themes: ["light", "dark", "cupcake"],
  },
}