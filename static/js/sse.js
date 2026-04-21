import { Question, User, Tag } from '/static/js/model.js';

document.addEventListener('DOMContentLoaded', () => {
    const es = new EventSource('/sse');

    es.onopen = () => {
        document.dispatchEvent(new CustomEvent('connect', {
            detail: true,
        }));
    };

    es.onerror = () => {
        document.dispatchEvent(new CustomEvent('disconnect', {
            detail: false,
        }));
    };

    es.addEventListener('notice', (event) => {
        try {
            document.dispatchEvent(new CustomEvent('sse-notice', {
                detail: JSON.parse(event.data),
            }));
        } catch (_) { }
    });

    es.addEventListener('question', (event) => {
        try {
            const question = Question.fromJSON(JSON.parse(event.data));
            document.dispatchEvent(new CustomEvent('sse-question', {
                detail: question,
            }));
        } catch (_) { }
    });

    es.addEventListener('user', (event) => {
        try {
            const user = User.fromJSON(JSON.parse(event.data));
            document.dispatchEvent(new CustomEvent('sse-user', {
                detail: user,
            }));
        } catch (fail) { console.log(fail) }
    });

    es.addEventListener('tag', (event) => {
        try {
            const tag = Tag.fromJSON(JSON.parse(event.data));
            document.dispatchEvent(new CustomEvent('sse-tag', {
                detail: tag,
            }));
        } catch (fail) { console.log(fail) }
    });
});
