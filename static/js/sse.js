import { Question, User, Tag, Notice } from '/static/js/model.js';

export const SSE_KEY = {
    system: {
        timeTick: 'time-tick',
    },
    data: {
        create: {
            user: 'create-user',
            question: 'create-question',
            tag: 'create-tag',
            notice: 'create-notice',
        },
        update: {
            user: 'update-user',
            question: 'update-question',
            tag: 'update-tag',
            notice: 'update-notice',
        },
        delete: {
            user: 'delete-user',
            question: 'delete-question',
            tag: 'delete-tag',
            notice: 'delete-notice',
        },
        get: {
            user: 'get-user',
            question: 'get-question',
            tag: 'get-tag',
            notice: 'get-notice',
        },
    },
    compose(...keys) {
        return keys.filter(Boolean);
    },
    user(includeDelete = false) {
        return this.compose(
            this.data.create.user,
            this.data.update.user,
            this.data.get.user,
            includeDelete ? this.data.delete.user : '',
        );
    },
    question(includeDelete = false) {
        return this.compose(
            this.data.create.question,
            this.data.update.question,
            this.data.get.question,
            includeDelete ? this.data.delete.question : '',
        );
    },
    tag(includeDelete = false) {
        return this.compose(
            this.data.create.tag,
            this.data.update.tag,
            this.data.get.tag,
            includeDelete ? this.data.delete.tag : '',
        );
    },
    notice(includeDelete = false) {
        return this.compose(
            this.data.create.notice,
            this.data.update.notice,
            this.data.get.notice,
            includeDelete ? this.data.delete.notice : '',
        );
    },
};

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

    es.addEventListener(SSE_KEY.system.timeTick, (event) => {
        document.dispatchEvent(new CustomEvent(SSE_KEY.system.timeTick, {
            detail: new Date(event.data),
        }));
    });

    const bindSSE = (eventNames, fromJSON, normalizedEvent, deleteEvent = null) => {
        eventNames.forEach((eventName) => {
            es.addEventListener(eventName, (event) => {
                try {
                    if (deleteEvent && eventName.startsWith('delete-')) {
                        const id = Number(event.data);
                        document.dispatchEvent(new CustomEvent(deleteEvent, {
                            detail: Number.isNaN(id) ? event.data : id,
                        }));
                        return;
                    }

                    const model = fromJSON(JSON.parse(event.data));
                    document.dispatchEvent(new CustomEvent(normalizedEvent, {
                        detail: model,
                    }));
                } catch (fail) {
                    console.log(fail);
                }
            });
        });
    };

    bindSSE(SSE_KEY.notice(true), Notice.fromJSON, SSE_KEY.data.create.notice, SSE_KEY.data.delete.notice);
    bindSSE(SSE_KEY.question(true), Question.fromJSON, SSE_KEY.data.create.question, SSE_KEY.data.delete.question);
    bindSSE(SSE_KEY.user(true), User.fromJSON, SSE_KEY.data.create.user, SSE_KEY.data.delete.user);
    bindSSE(SSE_KEY.tag(true), Tag.fromJSON, SSE_KEY.data.create.tag, SSE_KEY.data.delete.tag);
});
