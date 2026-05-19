import { checkEmailFormat, cleansingEmail, cleansingPassword } from './formCheck.js';

document.addEventListener('alpine:init', () => {
  Alpine.data('hrAppLogin', () => ({
    login: {},

    init() {
      if (typeof lucide !== 'undefined') {
        lucide.createIcons();
      }
    },

    checkEmailFmt(text) {
      return checkEmailFormat(text);
    },

    checkEmail(text) {
      return checkEmailFormat(text);
    },

    passwordCleansing(text) {
      return cleansingPassword(text);
    },

    emailCleansing(text) {
      return cleansingEmail(text);
    },

    async doLogin() {
      const res = await fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: this.login.email,
          pass: this.login.pass,
        }),
      });
      if (!res.ok) {
        window.notice.show({
          message: 'メールアドレスかパスワードが間違っています。',
          type: 'error',
          duration: 1000,
        });
        return;
      }
      const body = await res.json();
      const redirect = body.redirect;
      location.href = redirect;
    },
  }));
});
