import { Question, User, Tag, Notice } from '/static/js/model.js';

export const SSE = {
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
    create() {
        return Object.values(this.data.create).map((v) => v);
    },
    update() {
        return Object.values(this.data.update).map((v) => v);
    },
    get() {
        return Object.values(this.data.get).map((v) => v);
    },
    delete() {
        return Object.values(this.data.delete).map((v) => v);
    },
};

/** 既存インポート互換 */
export const SSE_KEY = SSE;

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

    es.addEventListener(SSE.system.timeTick, (event) => {
        document.dispatchEvent(new CustomEvent(SSE.system.timeTick, {
            detail: new Date(event.data),
        }));
    });

    const bindSSE = (eventSource, eventNames, fromJSON, normalizedEvent, deleteEvent = null) => {
        eventNames.forEach((eventName) => {
            eventSource.addEventListener(eventName, (event) => {
                try {
                    if (deleteEvent && eventName.startsWith('delete-')) {
                        const id = Number(event.data);
                        document.dispatchEvent(new CustomEvent(deleteEvent, {
                            detail: Number.isNaN(id) ? event.data : id,
                        }));
                        return;
                    }

                    const model = fromJSON(JSON.parse(event.data));
                    console.log(eventName, model);
                    let domEvent = normalizedEvent;
                    if (eventName.startsWith('get-') || eventName.startsWith('update-')) {
                        if (eventName.endsWith('-notice')) domEvent = SSE.data.update.notice;
                        else if (eventName.endsWith('-question')) domEvent = SSE.data.update.question;
                        else if (eventName.endsWith('-user')) domEvent = SSE.data.update.user;
                        else if (eventName.endsWith('-tag')) domEvent = SSE.data.update.tag;
                    } else if (eventName.startsWith('create-')) {
                        if (eventName.endsWith('-notice')) domEvent = SSE.data.create.notice;
                        else if (eventName.endsWith('-question')) domEvent = SSE.data.create.question;
                        else if (eventName.endsWith('-user')) domEvent = SSE.data.create.user;
                        else if (eventName.endsWith('-tag')) domEvent = SSE.data.create.tag;
                    }
                    document.dispatchEvent(new CustomEvent(domEvent, {
                        detail: model,
                    }));
                } catch (fail) {
                    console.log(fail);
                }
            });
        });
    };

    bindSSE(es, SSE.notice(true), Notice.fromJSON, SSE.data.create.notice, SSE.data.delete.notice);
    bindSSE(es, SSE.question(true), Question.fromJSON, SSE.data.create.question, SSE.data.delete.question);
    bindSSE(es, SSE.user(true), User.fromJSON, SSE.data.create.user, SSE.data.delete.user);
    bindSSE(es, SSE.tag(true), Tag.fromJSON, SSE.data.create.tag, SSE.data.delete.tag);
});
