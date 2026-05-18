/**
 * Mock5 画面上部用 Notice（Alpine.store + window.notice API）
 * DaisyUI alert + Lucide。duration はミリ秒。省略時 4000（DEFAULT_WAIT_MSEC）。null または 0 で自動消去なし。
 * dismissible: false で閉じるボタンを出さない（window.notice.dismiss / clear は有効）。
 * widthClass: 各 Notice 行に付与する Tailwind の幅クラス（例: w-full, w-96, w-full max-w-md）。省略時は w-full。
 */

const NOTICE_STORE_NAME = 'appNotice';

/** @type {Map<number, ReturnType<typeof setTimeout>>} */
const _noticeTimers = new Map();

/** @type {Array<() => void>} */
const _pendingNoticeOps = [];

function clearNoticeTimer(id) {
    const t = _noticeTimers.get(id);
    if (t != null) {
        clearTimeout(t);
        _noticeTimers.delete(id);
    }
}

function scheduleNoticeTimer(id, durationMs, onExpire) {
    clearNoticeTimer(id);
    if (durationMs == null || Number(durationMs) <= 0) return;
    const t = setTimeout(() => {
        _noticeTimers.delete(id);
        onExpire();
    }, Number(durationMs));
    _noticeTimers.set(id, t);
}

function refreshNoticeIcons() {
    queueMicrotask(() => {
        if (typeof lucide !== 'undefined') {
            lucide.createIcons();
        }
    });
}

function normalizeNoticeType(type) {
    const t = String(type ?? 'info').toLowerCase();
    if (t === 'alert') return 'warning';
    if (t === 'info' || t === 'success' || t === 'warning' || t === 'error') return t;
    return 'info';
}

function defaultIconForType(type) {
    switch (type) {
        case 'success':
            return 'circle-check';
        case 'warning':
            return 'alert-triangle';
        case 'error':
            return 'circle-x';
        default:
            return 'info';
    }
}

function alertClassForType(type) {
    switch (type) {
        case 'success':
            return 'alert alert-success shadow-md';
        case 'warning':
            return 'alert alert-warning shadow-md';
        case 'error':
            return 'alert alert-error shadow-md';
        default:
            return 'alert alert-info shadow-md';
    }
}

/** Tailwind の幅系クラス（w-*, min-w-*, max-w-* など）を1行にまとめる。空なら w-full */
function normalizeWidthClass(raw) {
    const s = String(raw ?? '').trim().replace(/\s+/g, ' ');
    if (!s) return 'w-full';
    return s.slice(0, 240);
}

function getNoticeStore() {
    if (typeof Alpine === 'undefined' || typeof Alpine.store !== 'function') return null;
    try {
        return Alpine.store(NOTICE_STORE_NAME);
    } catch {
        return null;
    }
}

function flushPendingNoticeOps() {
    while (_pendingNoticeOps.length) {
        const fn = _pendingNoticeOps.shift();
        try {
            fn();
        } catch (e) {
            console.error('[notice] pending op failed', e);
        }
    }
}

function enqueueOrRun(fn) {
    const s = getNoticeStore();
    if (s) {
        fn();
        return;
    }
    _pendingNoticeOps.push(fn);
}

document.addEventListener('alpine:init', () => {
    const DEFAULT_WAIT_MSEC = 4000;

    Alpine.store(NOTICE_STORE_NAME, {
        items: [],
        _seq: 0,

        /**
         * @param {{ message?: string, text?: string, type?: string, duration?: number|null, icon?: string, dismissible?: boolean, widthClass?: string, width?: string }} opts
         * @returns {number|null}
         */
        show(opts = {}) {
            const message = String(opts.message ?? opts.text ?? '').trim();
            if (!message) return null;
            const id = ++this._seq;
            const type = normalizeNoticeType(opts.type);
            const icon = opts.icon ? String(opts.icon) : defaultIconForType(type);
            const rawDur = opts.duration;
            const duration =
                rawDur === undefined ? DEFAULT_WAIT_MSEC : rawDur === null ? 0 : Number(rawDur);
            const dismissible = opts.dismissible !== false;
            const widthClass = normalizeWidthClass(opts.widthClass ?? opts.width);

            this.items.push({ id, message, type, icon, duration, dismissible, widthClass });
            scheduleNoticeTimer(id, duration, () => {
                this.dismiss(id);
            });
            refreshNoticeIcons();
            return id;
        },

        /**
         * @param {number} id
         * @param {{ message?: string, type?: string, duration?: number|null, icon?: string, dismissible?: boolean, widthClass?: string, width?: string }} patch
         */
        update(id, patch = {}) {
            const idx = this.items.findIndex((n) => n.id === id);
            if (idx === -1) return;
            const cur = this.items[idx];
            if (patch.message != null) cur.message = String(patch.message);
            if (patch.type != null) {
                cur.type = normalizeNoticeType(patch.type);
                if (patch.icon === undefined) {
                    cur.icon = defaultIconForType(cur.type);
                }
            }
            if (patch.icon != null) cur.icon = String(patch.icon);
            if (patch.dismissible !== undefined) {
                cur.dismissible = Boolean(patch.dismissible);
            }
            if (patch.widthClass !== undefined || patch.width !== undefined) {
                cur.widthClass = normalizeWidthClass(patch.widthClass ?? patch.width);
            }

            if (patch.duration !== undefined) {
                clearNoticeTimer(id);
                const d =
                    patch.duration === null ? 0 : Number(patch.duration);
                cur.duration = d;
                if (d > 0) {
                    scheduleNoticeTimer(id, d, () => {
                        this.dismiss(id);
                    });
                }
            }
            refreshNoticeIcons();
        },

        /** @param {number} id */
        dismiss(id) {
            clearNoticeTimer(id);
            this.items = this.items.filter((n) => n.id !== id);
            refreshNoticeIcons();
        },

        clear() {
            for (const n of this.items) {
                clearNoticeTimer(n.id);
            }
            this.items = [];
            refreshNoticeIcons();
        },

        /** @param {{ type: string }} n */
        alertClass(n) {
            return alertClassForType(normalizeNoticeType(n.type));
        },
    });

    flushPendingNoticeOps();
});

window.notice = {
    /**
     * @param {{ message?: string, text?: string, type?: string, duration?: number|null, icon?: string, dismissible?: boolean, widthClass?: string, width?: string }} opts
     * @returns {number|null}
     */
    show(opts) {
        let result = null;
        enqueueOrRun(() => {
            const s = getNoticeStore();
            if (s) result = s.show(opts);
        });
        return result;
    },

    /**
     * @param {number} id
     * @param {{ message?: string, type?: string, duration?: number|null, icon?: string, dismissible?: boolean, widthClass?: string, width?: string }} patch
     */
    update(id, patch) {
        enqueueOrRun(() => {
            const s = getNoticeStore();
            if (s) s.update(id, patch);
        });
    },

    /** @param {number} id */
    dismiss(id) {
        enqueueOrRun(() => {
            const s = getNoticeStore();
            if (s) s.dismiss(id);
        });
    },

    clear() {
        enqueueOrRun(() => {
            const s = getNoticeStore();
            if (s) s.clear();
        });
    },
};
