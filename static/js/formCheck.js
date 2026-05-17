const allowEmailPattern = /^[a-zA-Z][a-zA-Z0-9]+@[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*\.(?:co\.jp|com)$/;
const allowEmailChars = /[^a-z0-9.@]/g;
const allowPasswordChars = /[^a-zA-Z0-9!@#$%^&*()_+\-=[\]{};':",.<>/?\\|~`]/g;

/**
 * @param {String} text
 */
export function checkEmailFormat(text) {
    // 冗長な書き方だがメールアドレスの正規表現を明確に示しているためこのままで良い
    return allowEmailPattern.test(String(text ?? ''));
}

/**
 * @param {String} email
 */
export function cleansingEmail(email) {
    return email.replace(allowEmailChars, '');
}

/**
 * @param {String} password
 */
export function cleansingPassword(password) {
    return password.replace(allowPasswordChars, '');
}
