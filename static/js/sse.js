import { Question, User, Tag, Notice } from '/static/js/model.js';

export const SSE_TIME_TICKER = "time-ticker";
export const SSE_ADD_NOTICE = "sse-notice";
export const SSE_ADD_USER = "sse-user";
export const SSE_ADD_QUESTION = "sse-question";
export const SSE_ADD_TAG = "sse-tag";

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

    es.addEventListener('time-tick', (event) => {
        document.dispatchEvent(new CustomEvent(SSE_TIME_TICKER, {
            detail: new Date(event.data),
        }));
    });

    es.addEventListener('notice', (event) => {
        try {
            const json = JSON.parse(event.data);
            const notice = Notice.fromJSON(json);
            document.dispatchEvent(new CustomEvent(SSE_ADD_NOTICE, {
                detail: notice,
            }));
        } catch (_) { }
    });

    es.addEventListener('question', (event) => {
        try {
            const question = Question.fromJSON(JSON.parse(event.data));
            document.dispatchEvent(new CustomEvent(SSE_ADD_QUESTION, {
                detail: question,
            }));
        } catch (_) { }
    });

    es.addEventListener('user', (event) => {
        try {
            const user = User.fromJSON(JSON.parse(event.data));
            document.dispatchEvent(new CustomEvent(SSE_ADD_USER, {
                detail: user,
            }));
        } catch (fail) { console.log(fail) }
    });

    es.addEventListener('tag', (event) => {
        try {
            const tag = Tag.fromJSON(JSON.parse(event.data));
            document.dispatchEvent(new CustomEvent(SSE_ADD_TAG, {
                detail: tag,
            }));
        } catch (fail) { console.log(fail) }
    });
});
