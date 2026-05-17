import { Question, User, Tag, Notice } from '/static/js/model.js';
import { SSE } from './sse.js';
import { createAvatarBy, AVATAR_SIZE } from './createAvatar.js';
import { cleansingPassword, cleansingEmail, checkEmailFormat } from './formCheck.js';

/** API / モックで support はあるが user が無い場合がある。テンプレは activeQuestion.support.user.name を前提にする */
function ensureSupportForView(support) {
    const defaults = {
        supportStatusId: '1',
        supportStatus: { id: '1', title: '未対応' },
        user: { name: '' }
    };
    if (!support) return defaults;
    const o = { ...support };
    o.supportStatus.name = String(o.supportStatus.name ?? o.supportStatus.title ?? '');
    o.supportStatus.title = String(o.supportStatus.title ?? o.supportStatus.name ?? '');
    o.user = support.user && typeof support.user === 'object'
        ? { id: support.user.id, name: String(support.user.name ?? '') }
        : { name: '' };
    if (o.supportStatusId == null) o.supportStatusId = defaults.supportStatusId;
    if (!o.supportStatus) o.supportStatus = { ...defaults.supportStatus };
    return o;
}

/** 詳細ビュー (mock5_view_detail) が x-show 中でも式評価するため、空状態でも support 階層を持つ */
function emptyActiveQuestion() {
    return {
        support: ensureSupportForView(null),
        memos: [],
        answers: [],
        senderTalks: [],
    };
}

/** Question API/SSE から memos を view 用に正規化 */
function normalizeMemosForViewModel(question) {
    const raw = question.memos ?? [];
    return raw.map((m) => ({
        id: m.id,
        questionId: m.questionId,
        userId: m.userId,
        content: m.content ?? '',
        createdAt: m.createdAt ?? null,
        user: m.user && typeof m.user === 'object'
            ? { id: m.user.id, name: String(m.user.name ?? '') }
            : { id: m.userId, name: '' },
    }));
}

/** PUT 用ペイロードから、未対応 (supportStatusId === 1) のとき support / supportId を除去する */
function stripSupportForPutIfUnassigned(payload) {
    if (!payload || typeof payload !== 'object') return;
    const sid = payload.support?.supportStatusId;
    if (Number(sid) !== 1) return;
    delete payload.supportId;
    delete payload.support;
}

/** PUT 用ペイロードで、対応中 (supportStatusId === 2) かつ担当者未設定なら loginUser を担当として埋める */
function ensureSupportAssigneeForInProgressPut(payload, loginUser) {
    if (!payload || typeof payload !== 'object') return;
    const sid = payload.support?.supportStatusId;
    if (Number(sid) !== 2) return;
    const uid = loginUser?.id;
    if (uid == null || uid === '') return;
    const userName = String(loginUser?.name ?? '');
    const existing = payload.support;
    const existingUid = existing?.userId != null && existing.userId !== '' ? Number(existing.userId) : 0;
    if (existingUid > 0) return;
    payload.support = {
        id: 0,
        userId: String(uid),
        supportStatusId: '2',
        user: { id: uid, name: userName },
        supportStatus: { id: 2, title: '対応中' },
    };
    delete payload.supportId;
}

/** Question API/SSE から answers を view 用に正規化 */
function normalizeAnswersForViewModel(question) {
    const raw = Array.isArray(question.answers) ? question.answers : [];
    return raw.map((a) => ({
        isFinal: Boolean(a.isFinal),
        id: a.id,
        userId: a.userId,
        content: a.content ?? '',
        createdAt: a.createdAt ?? null,
        user: a.user && typeof a.user === 'object'
            ? { id: a.user.id, name: String(a.user.name ?? '') }
            : { id: a.userId, name: '' },
    }));
}

function firstSenderTalk(question) {
    const talks = Array.isArray(question?.senderTalks) ? question.senderTalks : [];
    return talks.length > 0 ? talks[0] : null;
}

function senderNameFromQuestion(question) {
    return firstSenderTalk(question)?.sender?.name
        ?? question?.sender
        ?? question?.support?.user?.name
        ?? '不明';
}

function senderDepartmentFromQuestion(question) {
    return firstSenderTalk(question)?.sender?.departmentName
        ?? question?.department
        ?? '';
}

/** 同一タブ内の画面復元用（sessionStorage）。Alpine.store はリロードで消えるため併用しない */
const HR_NAV_VIEW_KEY = 'meetup.hr.currentView';
const HR_NAV_DETAIL_Q_KEY = 'meetup.hr.detailQuestionId';
const HR_NAV_VIEWS = ['home', 'detail', 'notice', 'users', 'tags', 'settings'];

function clearHrNavSession() {
    try {
        sessionStorage.removeItem(HR_NAV_VIEW_KEY);
        sessionStorage.removeItem(HR_NAV_DETAIL_Q_KEY);
    } catch (_) {
        /* ignore */
    }
}

document.addEventListener('alpine:init', () => {
    Alpine.data('hrAppData', () => ({
        timeNow: new Date(),
        isSidebarOpen: true,
        scrollBoxFuncName: {
            question: 'question',
            user: 'user',
            tag: 'tag',
        },
        currentView: 'home',
        viewMode: 'card',
        showAddUserModal: false,
        showEditUserModal: false,
        showTagModal: false,
        showEditTagModal: false,
        unreadNotices: 0,
        activeQuestion: emptyActiveQuestion(),
        originalQuestion: {},
        activeUser: {},
        loginUser: {},
        activeTag: {},
        isConnect: false,

        /** 初回SSE接続〜初回 initData 成功まで true（リロード時も再度 true から始まる） */
        isInitialLoadingVisible: true,
        hasInitialLoadCompleted: false,
        initialLoadTotal: 4,
        initialLoadCompleted: 0,
        initialLoadLabel: '',

        selectedTags: [],
        availableTags: [],
        searchQuery: '',
        statusFilter: 'all',
        sortOrder: 'date_desc',

        relatedSearchQuery: '',

        /** タグ管理画面（mock5_view_tag）専用。質問一覧の searchQuery とは独立 */
        tagSearchQuery: '',
        tagCategoryFilter: 'all',

        /** ユーザー管理画面（mock5_view_users）専用 */
        userSearchQuery: '',
        userRoleFilter: 'all',

        /** 詳細画面の回答・メモ入力下書き（mock5_view_detail） */
        detailComposerDraft: '',

        titleMap: {
            home: '質問一覧',
            detail: '質問詳細',
            notice: '通知センター',
            users: 'ユーザー管理',
            tags: 'タグ管理',
            settings: 'プロフィール設定'
        },

        statusTitleMap: {
            1: '未対応',
            2: '対応中',
            3: '完了'
        },

        categoryTitleMap: {
            1: '総務',
            2: '人事',
            3: 'その他',
        },

        statusColorMap: {
            '未対応': 'bg-red-100 text-red-800',
            '対応中': 'bg-blue-100 text-blue-800',
            '完了': 'bg-gray-100 text-gray-600'
        },

        questions: [],

        notices: [],

        users: [],
        tags: [],

        get filteredQuestions() {
            let result = this.questions.filter(q => {
                const matchesTag = this.selectedTags.length === 0 || q.tags.some(tag => this.selectedTags.includes(this.tagDisplayName(tag)));
                const matchesSearch = q.title.includes(this.searchQuery) || q.content.includes(this.searchQuery) || q.sender.includes(this.searchQuery);
                const st = q.support?.supportStatus?.title ?? '';
                let matchesStatus = true;
                if (this.statusFilter === 'unanswered') matchesStatus = st === '未対応';
                else if (this.statusFilter === 'progress') matchesStatus = st === '対応中';
                else if (this.statusFilter === 'open') matchesStatus = (st === '未対応' || st === '対応中');
                else if (this.statusFilter === 'completed') matchesStatus = st === '完了';
                return matchesTag && matchesSearch && matchesStatus;
            });

            result.sort((a, b) => {
                if (this.sortOrder === 'date_desc') {
                    return new Date(b.date) - new Date(a.date);
                } else if (this.sortOrder === 'date_asc') {
                    return new Date(a.date) - new Date(b.date);
                } else if (this.sortOrder === 'dueDate_asc') {
                    return new Date(a.dueDate) - new Date(b.dueDate);
                } else if (this.sortOrder === 'dueDate_desc') {
                    return new Date(b.dueDate) - new Date(a.dueDate);
                }
                return 0;
            });

            return result;
        },

        get availableRelatedQuestions() {
            if (!this.activeQuestion || !this.activeQuestion.id) return [];
            let qs = this.questions.filter(q => q.id !== this.activeQuestion.id);
            if (this.relatedSearchQuery) {
                const query = this.relatedSearchQuery.toLowerCase();
                qs = qs.filter(q => q.title.toLowerCase().includes(query) || q.id.toString().includes(query));
            }
            return qs;
        },

        get filteredTags() {
            const q = String(this.tagSearchQuery ?? '').trim().toLowerCase();
            const cat = this.tagCategoryFilter;
            return this.tags.filter((tag) => {
                if (cat && cat !== 'all') {
                    const cid =
                        tag.categoryId != null && tag.categoryId !== ''
                            ? Number(tag.categoryId)
                            : Number(tag.category?.id ?? NaN);
                    if (Number(cat) !== cid) return false;
                }
                if (!q) return true;
                const name = this.tagDisplayName(tag).toLowerCase();
                const cname = this.tagCategoryDisplayName(tag).toLowerCase();
                return name.includes(q) || cname.includes(q);
            });
        },

        get filteredUsers() {
            const q = String(this.userSearchQuery ?? '').trim().toLowerCase();
            const role = this.userRoleFilter;
            return this.users.filter((u) => {
                if (role && role !== 'all') {
                    if (Number(u.roleId) !== Number(role)) return false;
                }
                if (!q) return true;
                const name = this.userDisplayName(u).toLowerCase();
                const email = this.userDisplayEmail(u).toLowerCase();
                return name.includes(q) || email.includes(q);
            });
        },

        apiHeaders(withJSON = true) {
            const h = {};
            if (withJSON) h['Content-Type'] = 'application/json';
            return h;
        },

        /**
         * @param {{ blocking?: boolean }} [options]
         * blocking: 初回のみ。全画面ローディング＋逐次進捗。再接続時は false で並列取得のみ。
         */
        async initData(options = {}) {
            const blocking = Boolean(options.blocking);
            if (blocking) {
                this.isInitialLoadingVisible = true;
                this.initialLoadCompleted = 0;
                this.initialLoadTotal = 4;
                this.initialLoadLabel = 'サーバーに接続しています…';
            }
            const steps = [
                { label: 'ログインユーザーを読み込んでいます', fn: () => this.getLoginUser() },
                { label: '質問一覧を読み込んでいます', fn: () => this.getQuestions() },
                { label: 'ユーザーを読み込んでいます', fn: () => this.getUsers() },
                { label: 'タグを読み込んでいます', fn: () => this.getTags() },
            ];
            try {
                if (blocking) {
                    for (const step of steps) {
                        this.initialLoadLabel = step.label;
                        await step.fn();
                        this.initialLoadCompleted += 1;
                    }
                } else {
                    await Promise.all(steps.map((s) => s.fn()));
                }
                this.applyPersistedNavigation();
                if (blocking) {
                    this.hasInitialLoadCompleted = true;
                    this.isInitialLoadingVisible = false;
                }
            } catch (e) {
                console.error('initData', e);
                if (blocking) {
                    this.isInitialLoadingVisible = false;
                    if (typeof window !== 'undefined' && window.notice?.show) {
                        window.notice.show({
                            message: 'データの読み込みに失敗しました。接続が回復すると自動で再試行されます。',
                            type: 'error',
                            duration: 6000,
                        });
                    }
                }
            }
            this.refreshIcons();
        },

        init() {

            document.addEventListener('connect', (event) => {
                this.isConnect = event.detail;
                void this.initData({ blocking: !this.hasInitialLoadCompleted });
                this.refreshIcons();
            });
            document.addEventListener('disconnect', (event) => {
                this.isConnect = event.detail;
            });

            document.addEventListener(SSE.system.timeTick, (event) => {
                this.timeNow = event.detail
            });

            document.addEventListener(SSE.data.create.notice, (event) => {
                const notice = event.detail;
                const index = this.notices.findIndex((n) => n.id === notice.id);
                if (index === -1) {
                    this.notices.unshift(notice);
                    if (notice.typeId !== 1) this.unreadNotices += 1;
                    this.refreshIcons();
                }
            });

            document.addEventListener(SSE.data.update.notice, (event) => {
                const notice = event.detail;
                const index = this.notices.findIndex((n) => n.id === notice.id);
                if (index !== -1) {
                    if (!_.isEqual(this.notices[index], notice)) {
                        this.notices.splice(index, 1, {
                            ...this.notices[index],
                            ...notice,
                        });
                        this.refreshIcons();
                    }
                }
            });

            document.addEventListener(SSE.data.delete.notice, (event) => {
                const id = event.detail;
                this.notices = this.notices.filter((n) => n.id !== id);
            });

            document.addEventListener(SSE.data.create.question, (event) => {
                const question = this.toQuestionViewModel(event.detail);
                const index = this.questions.findIndex(q => q.id === question.id);
                if (index === -1) {
                    this.questions.unshift(question);
                    this.refreshIcons();
                }
            });

            document.addEventListener(SSE.data.update.question, (event) => {
                const question = this.toQuestionViewModel(event.detail);
                const index = this.questions.findIndex(q => q.id === question.id);

                if (index !== -1) {
                    if (!_.isEqual(this.questions[index], question)) {
                        this.questions.splice(index, 1, {
                            ...this.questions[index],
                            ...question,
                        });
                    }
					if (this.activeQuestion?.id === question.id) {
						if (!_.isEqual(this.activeQuestion, question)) {
							this.activeQuestion = _.mergeWith(
								_.cloneDeep(question),
								this.activeQuestion,
								(serverVal, localVal, key) => {
									const originalVal = _.get(this.originalQuestion, key)
									// ローカル変更されてたら優先
									if (!_.isEqual(localVal, originalVal)) {
										switch (key) {
											// support.supportStatusの更新が必ずLocalWinになってしまうため例外
											case 'support':
											case 'memos':
												return serverVal;
										}
										return localVal;
									}
									// server採用
									return serverVal
								}
							);
							this.originalQuestion = _.cloneDeep(question);
							this.refreshIcons();
						}
					}
                }
            });

            document.addEventListener(SSE.data.delete.question, (event) => {
                const id = event.detail;
                if (this.activeQuestion?.id === id) {
					this.questions = this.questions.filter((q) => q.id !== id);
                    this.activeQuestion = emptyActiveQuestion();
                    this.originalQuestion = {};
                }
            });

            document.addEventListener(SSE.data.create.user, (event) => {
                const user = event.detail;
                const index = this.users.findIndex((u) => u.id === user.id);
                if (index === -1) {
                    this.users.unshift(user);
                    this.refreshIcons();
                }
            });

            document.addEventListener(SSE.data.update.user, (event) => {
                const user = event.detail;
                const index = this.users.findIndex((u) => u.id === user.id);
                if (index !== -1) {
                    if (!_.isEqual(this.users[index], user)) {
                        this.users.splice(index, 1, {
                            ...this.users[index],
                            ...user,
                        });
                        this.refreshIcons();
                    }
                }
            });

            document.addEventListener(SSE.data.delete.user, (event) => {
                const id = event.detail;
                this.users = this.users.filter((u) => u.id !== id);
            });

            document.addEventListener(SSE.data.create.tag, (event) => {
                const tag = event.detail;
                const id = tag.id;
                const index =
                    id != null && id !== 0
                        ? this.tags.findIndex((t) => t.id === id)
                        : -1;
                if (index === -1) {
                    this.tags.unshift(tag);
                    this.availableTags = this.tags.map((t) => Tag.fromJSON(t));
                    this.refreshIcons();
                }
            });

            document.addEventListener(SSE.data.update.tag, (event) => {
                const tag = event.detail;
                const index = this.tags.findIndex((t) => t.id === tag.id);
                if (index !== -1) {
                    if (!_.isEqual(this.tags[index], tag)) {
                        this.tags.splice(index, 1, {
                            ...this.tags[index],
                            ...tag,
                        });
                        this.availableTags = this.tags.map((t) => Tag.fromJSON(t));
                        this.refreshIcons();
                    }
                }
            });

            document.addEventListener(SSE.data.delete.tag, (event) => {
                const id = event.detail;
                const newTags = this.tags.filter((t) => t.id !== id);
                this.tags = newTags;
                this.availableTags = newTags.map((t) => Tag.fromJSON(t));
            });

            this.$nextTick(() => this.refreshIcons());
        },

        getIcon(name) {
        },

        // --- mock5: 表示用ヘルパー（HTMLのAlpine式からロジックを分離） ---

        hasSelectedTags() {
            return this.selectedTags.length > 0;
        },

        questionStatusTitle(q) {
            if (!q?.support?.supportStatus?.title && !q?.support?.supportStatus?.name) return '';
            const id = Number(q?.support?.supportStatusId);
            const title = this.statusTitleMap[id];

            this.activeQuestion.support.supportStatus = {id: id, name: title, title: title}
            return title;
        },

        questionStatusBadgeClass(q) {
            const t = this.questionStatusTitle(q);
            if (!t) return '';
            return this.statusColorMap[t] || '';
        },

        questionTagsList(q) {
            return Array.isArray(q?.tags) ? q.tags : [];
        },

        isQuestionStatusCompleted(q) {
            return this.questionStatusTitle(q) === '完了';
        },

        isQuestionDueUrgent(q) {
            return (
                (q?.daysLeft ?? 99) <= 1 &&
                !this.isQuestionStatusCompleted(q)
            );
        },

        questionDueLabel(q) {
            return this.isQuestionStatusCompleted(q)
                ? '完了'
                : `あと${q?.daysLeft ?? 0}日`;
        },

        userInitialChar(u) {
            return String(u?.name ?? '').charAt(0);
        },

        userDisplayName(u) {
            return u?.name != null && u.name !== '' ? String(u.name) : '';
        },

        userDisplayEmail(u) {
            return u?.email != null && u.email !== '' ? String(u.email) : '';
        },

        userRoleDisplayName(u) {
            if (u?.role?.name) return String(u.role.name);
            return '';
        },

        activeQuestionHeaderLine() {
            const s = this.activeQuestion?.sender != null
                ? String(this.activeQuestion.sender)
                : '';
            const d = this.activeQuestion?.department != null
                ? String(this.activeQuestion.department)
                : '';
            return `${s} • ${d}`;
        },

        senderInitialChar() {
            const s = this.activeQuestion?.sender != null
                ? String(this.activeQuestion.sender)
                : '';
            return s.charAt(0);
        },

        supportUserInitialChar() {
            const s = this.activeQuestion?.support?.user?.name != null
                ? String(this.activeQuestion.support.user.name)
                : '';
            return s.charAt(0);
        },

        chatQuestionContent() {
            return String(this.activeQuestion?.content ?? '');
        },

        chatMemoContent(memo) {
            return String(memo?.content ?? '');
        },

        chatAnswerContent(answer) {
            return String(answer?.content ?? '');
        },

        isSelfMemo(memo) {
            if (this.loginUser?.id == null || memo?.userId == null) return false;
            return Number(memo.userId) === Number(this.loginUser.id);
        },

        isSelfAnswer(answer) {
            if (this.loginUser?.id == null || answer?.userId == null) return false;
            return Number(answer.userId) === Number(this.loginUser.id);
        },

        makeAvatar(text) {
            return createAvatarBy(text, AVATAR_SIZE.ZERO);
        },

        chatBubbleColorByKind(kind) {
            switch (kind) {
                case 'memo':
                    return 'chat-bubble-accent';
                case 'answer':
                    return '';
                case 'talk':
                    return 'chat-bubble-secondary';
                default:
                    return '';
            }
        },

        hasActiveQuestionAnswers() {
            return Array.isArray(this.activeQuestion?.answers) && this.activeQuestion.answers.length > 0;
        },

        tagDisplayName(tag) {
            return String(tag?.name ?? tag?.title ?? '');
        },

        chatAvatarSeed(name) {
            return encodeURIComponent(String(name ?? 'User'));
        },

        /**
         * @param {string} text
         * @returns {string}
         */
        chatAvatarURI(text) {
            return `https://api.dicebear.com/7.x/notionists/svg?seed=${this.chatAvatarSeed(text)}&backgroundColor=eff6ff`;
        },

        chatTimelineItems() {
            const memos = (this.activeQuestion?.memos ?? []).sort((m1, m2) => Number(m1.id) - Number(m2.id)).map((memo, i) => ({
                kind: 'memo',
                userId: memo.userId,
                userName: memo.user?.name ?? '',
                content: this.chatMemoContent(memo),
                createdAt: memo.createdAt ? new Date(memo.createdAt) : null,
                createdAtStr: memo.createdAt ? this.timeFormattMMDDHHMM(memo.createdAt) : '',
                _origIdx: i,
                memo: memo,
            }));

            const answers = (this.activeQuestion?.answers ?? []).map((answer, i) => ({
                kind: 'answer',
                userId: answer.userId,
                userName: answer.user?.name ?? '',
                content: this.chatAnswerContent(answer),
                createdAt: answer.createdAt ? new Date(answer.createdAt) : null,
                createdAtStr: answer.createdAt ? this.timeFormattMMDDHHMM(answer.createdAt) : '',
                _origIdx: memos.length + i,
                answer: answer,
            }));

            const senderTalks = (this.activeQuestion?.senderTalks ?? []).map((talk, i) => {
                const sender = talk.sender;
                return ({
                    kind: 'talk',
                    userId: sender.id,
                    userName: sender.name,
                    content: talk.content,
                    createdAt: talk.createdAt ? new Date(talk.createdAt) : null,
                    createdAtStr: talk.createdAt ? this.timeFormattMMDDHHMM(talk.createdAt) : '',
                    _origIdx: memos.length + answers.length + i,
                    talk: talk,
                });
            });

            const items = [...memos, ...answers, ...senderTalks];
            items.sort((x, y) => {
                const timeKey = (it) => {
                    if (!it.createdAt) return Number.POSITIVE_INFINITY;
                    const t = it.createdAt.getTime();
                    return Number.isFinite(t) ? t : Number.POSITIVE_INFINITY;
                };
                const xt = timeKey(x);
                const yt = timeKey(y);
                if (xt !== yt) return xt - yt;
                return x._origIdx - y._origIdx;
            });
            return items;
        },

        isSelfChatItem(item) {
            if (item?.kind === 'memo') return this.isSelfMemo(item.memo);
            if (item?.kind === 'answer') return this.isSelfAnswer(item.answer);
            return false;
        },

        activeQuestionTagButtonLabel() {
            const tags = this.activeQuestion?.tags;
            if (!Array.isArray(tags) || tags.length === 0) return 'タグを選択';
            return tags.map((t) => this.tagDisplayName(t)).filter(Boolean).join(', ');
        },

        isActiveQuestionTagSelected(tag) {
            const tags = this.activeQuestion?.tags;
            if (!Array.isArray(tags) || !tag) return false;
            return tags.filter((t) => t.id === tag.id || this.tagDisplayName(t) === this.tagDisplayName(tag)).length === 1;
        },

        activeQuestionDateOrPlaceholder() {
            return this.activeQuestion?.date
                ? String(this.activeQuestion.date)
                : '---';
        },

        relatedQuestionsIsEmpty() {
            const rq = this.activeQuestion?.relatedQuestions;
            return !Array.isArray(rq) || rq.length === 0;
        },

        relatedQuestionPickerButtonLabel() {
            if (this.relatedQuestionsIsEmpty()) return '質問を選択';
            const n = this.activeQuestion.relatedQuestions.length;
            return `${n}件選択中`;
        },

        /** Form の relatedQuestions 行またはレガシーな数値IDから関連先質問IDを取り出す */
        relatedQuestionTargetId(rowOrId) {
            if (rowOrId == null) return 0;
            if (typeof rowOrId === 'object' && rowOrId.relatedQuestionId != null && rowOrId.relatedQuestionId !== '') {
                const n = Number(rowOrId.relatedQuestionId);
                return Number.isNaN(n) ? 0 : n;
            }
            const n = Number(rowOrId);
            return Number.isNaN(n) ? 0 : n;
        },

        relatedQuestionRowKey(row, idx) {
            const rid = row && typeof row === 'object' ? String(row.relatedQuestionId ?? '') : '';
            const idPart = row && typeof row === 'object' && row.id != null ? String(row.id) : '';
            return `${idPart}-${rid}-${idx}`;
        },

        isRelatedQuestionChecked(rId) {
            const rq = this.activeQuestion?.relatedQuestions;
            if (!Array.isArray(rq)) return false;
            const id = Number(rId);
            return rq.some((row) => this.relatedQuestionTargetId(row) === id);
        },

        hasActiveRelatedQuestionIds() {
            const rq = this.activeQuestion?.relatedQuestions;
            return Array.isArray(rq) && rq.length > 0;
        },

        relatedQuestionLineTitle(rowOrId) {
            const id = this.relatedQuestionTargetId(rowOrId);
            let title = '';
            if (rowOrId && typeof rowOrId === 'object' && rowOrId.relatedQuestion && rowOrId.relatedQuestion.title != null) {
                title = String(rowOrId.relatedQuestion.title);
            }
            const q = this.questions.find((x) => x.id === id);
            const t = title || (q?.title != null ? String(q.title) : '');
            return `#${id} ${t}`;
        },

        findQuestionById(id) {
            return this.questions.find((q) => q.id === id);
        },

        openNoticeIfQuestion(n) {
            if (n?.question) {
                this.openDetailByID(n.questionId);
            }
        },

        hasNoticeQuestion(n) {
            return Boolean(n?.question);
        },

        noticeItemLayoutClass(n) {
            return n?.question && n.question.support ? 'flex-col' : '';
        },

        noticeIconWrapperClass(n) {
            return [2, 3].includes(n?.typeId)
                ? 'bg-red-50 text-red-500'
                : 'bg-slate-50 text-slate-500';
        },

        noticeIconName(n) {
            return [2, 3].includes(n?.typeId) ? 'alert-circle' : 'info';
        },

        noticeSupportStatusTitle(n) {
            return (
                n?.question?.support?.supportStatus?.title
                    ? String(n.question.support.supportStatus.title)
                    : ''
            );
        },

        noticeSupportStatusBadgeClass(n) {
            const t = this.noticeSupportStatusTitle(n);
            if (!t) return '';
            return this.statusColorMap[t] || '';
        },

        noticeType3Heading(n) {
            if (!n?.question) return '';
            if (n.question.title) {
                return `回答期限接近中: ${n.question.title}`;
            }
            return `回答期限接近中: 質問ID#${n.question.id}`;
        },

        noticeDueTimeLabel(n) {
            if (!n?.question) return '';
            const due = n.question.due;
            if (due == null) return '';
            const dueDate = due instanceof Date ? due : new Date(due);
            if (Number.isNaN(dueDate.getTime())) return '';
            const diffMs = dueDate - this.timeNow;
            if (diffMs >= 0) {
                return `残り時間: ${this.calcRemainingTimeAndString(n.question)}`;
            }
            return this.calcRemainingTimeAndString(n.question);
        },

        noticeQuestionTagsList(n) {
            if (!n?.question) return [];
            return Array.isArray(n.question.tags) ? n.question.tags : [];
        },

        noticeQuestionContent(n) {
            if (!n?.question) return '';
            return n.question.content != null
                ? String(n.question.content)
                : '';
        },

        countQuestionsUsingTagTitle(tag) {
            const title = tag?.name != null ? String(tag.name) : String(tag?.title ?? '');
            if (!title) return 0;
            return this.questions.filter((q) => {
                const tags = q?.tags;
                if (!Array.isArray(tags)) return false;
                return tags.some((t) => t && (t.name === title || t.title === title));
            }).length;
        },

        tagUsageCountLabel(tag) {
            return `${this.countQuestionsUsingTagTitle(tag)}件`;
        },

        tagCategoryDisplayName(tag) {
            if (tag?.category && (tag.category.name || tag.category.categoryName)) {
                return String(tag.category.name ?? tag.category.categoryName);
            }
            return '';
        },

        /**
         * 
         * @param {Date|String} dt 
         * @param {String[]} requireUnit 
         * @returns {String} MM/DD hh:mm
         */
        timeFormattMMDDHHMM(dt) {
            const formatter = new Intl.DateTimeFormat('ja-JP', {
                month: '2-digit',  // "01月" （2桁固定）
                day: '2-digit',    // "02日" （2桁固定）
                hour: '2-digit',   // "15時"
                minute: '2-digit', // "04分"
                hour12: false      // 24時間表記
            });
            return formatter.format(dt instanceof Date ? dt : new Date(dt));
        },

        currentTimeFormattedHHMMSS(dt = new Date()) {
            return new Intl.DateTimeFormat('ja-JP', {
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit'
            }).format(dt);
        },

        toQuestionViewModel(question) {
            const createdAt = question.createdAt instanceof Date
                ? question.createdAt
                : new Date(question.createdAt);
            const due = question.due instanceof Date
                ? question.due
                : (question.due ? new Date(question.due) : null);
            const createdAtValid = !Number.isNaN(createdAt.getTime());
            const dueValid = due != null && !Number.isNaN(due.getTime());
            const dueForCalc = dueValid ? due : null;
            const daysLeft = dueForCalc
                ? Math.max(0, Math.ceil((dueForCalc.getTime() - Date.now()) / (1000 * 60 * 60 * 24)))
                : 0;
            return {
                id: question.id,
                userId: firstSenderTalk(question)?.sender?.id ?? question.support?.user?.id ?? null,
                sender: senderNameFromQuestion(question),
                department: senderDepartmentFromQuestion(question),
                title: question.title ?? '',
                content: question.content ?? '',
                support: ensureSupportForView(question.support),
                tags: question.tags,
                daysLeft,
                date: createdAtValid ? this.formatDateTime(createdAt) : '',
                dueDate: dueValid ? this.formatDate(due) : '',
                createdAt: createdAtValid ? createdAt : null,
                due: dueValid ? due : null,
                relatedQuestions: this.normalizeRelatedQuestionsForViewModel(question),
                memos: normalizeMemosForViewModel(question),
                senderTalks: question.senderTalks ?? [],
                answers: normalizeAnswersForViewModel(question),
            };
        },

        normalizeRelatedQuestionsForViewModel(question) {
            const raw = question.relatedQuestions;
            if (Array.isArray(raw) && raw.length > 0) {
                return raw.map((row) => {
                    if (row && typeof row === 'object' && (row.relatedQuestionId != null || row.relatedQuestion)) {
                        return {
                            id: row.id ?? 0,
                            questionId: row.questionId != null ? String(row.questionId) : String(question.id ?? ''),
                            relatedQuestionId: row.relatedQuestionId != null ? String(row.relatedQuestionId) : '',
                            question: row.question ?? null,
                            relatedQuestion: row.relatedQuestion ?? null,
                        };
                    }
                    const n = typeof row === 'number' ? row : Number(row);
                    return {
                        id: 0,
                        questionId: String(question.id ?? ''),
                        relatedQuestionId: String(Number.isNaN(n) ? 0 : n),
                        question: null,
                        relatedQuestion: null,
                    };
                });
            }
            return (question.subQuestions ?? []).map((sq) => ({
                id: 0,
                questionId: String(question.id ?? ''),
                relatedQuestionId: String(
                    typeof sq === 'object' && sq !== null && sq.id != null ? sq.id : sq,
                ),
                question: null,
                relatedQuestion: null,
            }));
        },

        formatDateTime(date) {
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1).padStart(2, '0');
            const day = String(date.getDate()).padStart(2, '0');
            const hour = String(date.getHours()).padStart(2, '0');
            const minute = String(date.getMinutes()).padStart(2, '0');
            return `${year}-${month}-${day} ${hour}:${minute}`;
        },

        formatDate(date) {
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1).padStart(2, '0');
            const day = String(date.getDate()).padStart(2, '0');
            return `${year}-${month}-${day}`;
        },

        logout() {
            clearHrNavSession();
            location.href = "/login";
        },

        persistNavToSession() {
            try {
                const v = this.currentView;
                if (!HR_NAV_VIEWS.includes(v)) return;
                sessionStorage.setItem(HR_NAV_VIEW_KEY, v);
                if (
                    v === 'detail'
                    && this.activeQuestion?.id != null
                    && String(this.activeQuestion.id) !== ''
                    && Number(this.activeQuestion.id) > 0
                ) {
                    sessionStorage.setItem(HR_NAV_DETAIL_Q_KEY, String(this.activeQuestion.id));
                } else {
                    sessionStorage.removeItem(HR_NAV_DETAIL_Q_KEY);
                }
            } catch (_) {
                /* ignore */
            }
        },

        applyPersistedNavigation() {
            let view;
            let qidRaw;
            try {
                view = sessionStorage.getItem(HR_NAV_VIEW_KEY);
                qidRaw = sessionStorage.getItem(HR_NAV_DETAIL_Q_KEY);
            } catch (_) {
                return;
            }
            if (!view || !HR_NAV_VIEWS.includes(view)) {
                if (view) clearHrNavSession();
                return;
            }
            if (view === 'detail') {
                const qid = qidRaw != null && qidRaw !== '' ? Number(qidRaw) : NaN;
                if (!Number.isFinite(qid) || qid <= 0) {
                    clearHrNavSession();
                    this.currentView = 'home';
                    return;
                }
                const index = this.questions.findIndex((q) => q.id == qid);
                if (index === -1) {
                    clearHrNavSession();
                    this.currentView = 'home';
                    return;
                }
                this.openDetail(this.questions[index]);
                return;
            }
            this.setView(view);
        },

        setView(view) {
            this.currentView = view;
            this.persistNavToSession();
            this.refreshIcons();
        },

        openDetail(q) {
            const v = _.cloneDeep(q);
            v.support = ensureSupportForView(v.support);
            this.originalQuestion = _.cloneDeep(v);
            this.activeQuestion = v;
            this.relatedSearchQuery = '';
            this.detailComposerDraft = '';
            this.setView('detail');
        },

        openDetailByID(id = 0) {
            if (id === 0) {
                window.notice.show({ message: `invalid support id is ${id}.`, type: 'error' });
            }
            const index = this.questions.findIndex(q => q.id == id)
            if (index !== -1) {
                this.openDetail(this.questions[index])
            }
        },

        toggleTag(tag) {
            const name = this.tagDisplayName(tag);
            if (!name) return;
            if (this.selectedTags.includes(name)) {
                this.selectedTags = this.selectedTags.filter(t => t !== name);
            } else {
                this.selectedTags.push(name);
            }
            this.refreshIcons();
        },

        toggleQuestionTag(tag) {
            if (!this.activeQuestion.tags) {
                this.activeQuestion.tags = [];
            }
            if (this.activeQuestion.tags.filter(t => t.id === tag.id).length === 1) {
                this.activeQuestion.tags = this.activeQuestion.tags.filter(t => t.id !== tag.id);
            } else {
                this.activeQuestion.tags.push(tag);
            }
            this.updateQuestion()
            this.refreshIcons();
        },

        toggleRelatedQuestion(qId) {
            if (!this.activeQuestion.relatedQuestions) {
                this.activeQuestion.relatedQuestions = [];
            }
            const q = this.activeQuestion;
            const target = Number(qId);
            const idx = q.relatedQuestions.findIndex(
                (row) => this.relatedQuestionTargetId(row) === target
            );
            if (idx !== -1) {
                q.relatedQuestions.splice(idx, 1);
            } else {
                q.relatedQuestions.push({
                    id: 0,
                    questionId: String(q.id),
                    relatedQuestionId: String(target),
                    question: null,
                    relatedQuestion: null,
                });
            }
            void this.updateQuestion();
            this.refreshIcons();
        },

        /**
         * メモを activeQuestion.memos に追加し PUT する。サーバは memos を全件差し替えのため既存分は aq.memos ごと送る。
         * Go MemoForm は questionId / userId が文字列想定のため updateQuestion 内で正規化する。
         */
        async submitMemo() {
            const text = String(this.detailComposerDraft ?? '').trim();
            if (!text) return;

            const uid = this.loginUser?.id;
            if (uid == null || uid === '') {
                window.notice.show({ message: 'ログイン情報がありません', type: 'warning' });
                return;
            }

            const qid = this.activeQuestion?.id;
            if (!qid) return;

            if (!Array.isArray(this.activeQuestion.memos)) {
                this.activeQuestion.memos = [];
            }

            const newMemo = {
                id: 0,
                questionId: String(qid),
                userId: String(uid),
                content: text,
                user: {
                    id: uid,
                    name: String(this.loginUser?.name ?? ''),
                },
            };
            this.activeQuestion.memos.push(newMemo);

            const ok = await this.updateQuestion();
            if (ok) {
                this.detailComposerDraft = '';
            } else {
                this.activeQuestion.memos.pop();
                window.notice.show({ message: 'メモの送信に失敗しました', type: 'error' });
            }
            this.refreshIcons();
        },

        /** 詳細画面から回答を送信し、PUT で永続化する */
        async submitAnswer() {
            const text = String(this.detailComposerDraft ?? '').trim();
            if (!text) return;

            const uid = this.loginUser?.id;
            if (uid == null || uid === '') {
                window.notice.show({ message: 'ログイン情報がありません', type: 'warning' });
                return;
            }

            const qid = this.activeQuestion?.id;
            if (!qid) return;

            if (!Array.isArray(this.activeQuestion.memos)) {
                this.activeQuestion.memos = [];
            }
            if (!Array.isArray(this.activeQuestion.answers)) {
                this.activeQuestion.answers = [];
            }

            const newAnswer = {
                id: 0,
                questionId: String(qid),
                userId: String(uid),
                content: text,
                user: {
                    id: uid,
                    name: String(this.loginUser?.name ?? ''),
                },
                createdAt: new Date(),
            };
            this.activeQuestion.answers.push(newAnswer);

            const ok = await this.updateQuestion();
            if (ok) {
                this.detailComposerDraft = '';
            } else {
                this.activeQuestion.answers.pop();
                window.notice.show({ message: '回答の送信に失敗しました', type: 'error' });
            }
            this.refreshIcons();
        },

        /** @returns {Promise<boolean>} PUT が成功したか、または送信不要で問題ない場合 true */
        async updateQuestion() {
            // 未対応 (supportStatusId === 1) のとき stripSupportForPutIfUnassigned で payload から support を落とす（サーバ側で detach も実施）。
            // TODO; Answer送信後、ステータスを完了へと更新する？一撃で帰ってこなかった場合の処理を考えないと終わらなさそう。
            const question = this.questions.filter(q => q.id === this.activeQuestion.id)[0];

            if (!_.isEqual(question, this.activeQuestion)) {
                const aq = this.activeQuestion;
                let dueISO = null;
                if (aq.dueDate) {
                    const sameAsOriginal = aq.due instanceof Date && this.formatDate(aq.due) === aq.dueDate;
                    if (sameAsOriginal) {
                        dueISO = aq.due.toISOString();
                    } else {
                        const parts = aq.dueDate.split('-').map(Number);
                        if (parts.length === 3 && parts.every((n) => !Number.isNaN(n))) {
                            const [y, m, d] = parts;
                            dueISO = new Date(y, m - 1, d).toISOString();
                        }
                    }
                }
                const createdAtISO = aq.createdAt instanceof Date && !Number.isNaN(aq.createdAt.getTime())
                    ? aq.createdAt.toISOString()
                    : null;

                /** PUT は RelatedQuestionForm のみ（ネストした relatedQuestion が Question 実体だと originQuestionId が数値になり Go が拒否する） */
                const relatedQuestionsPayload = Array.isArray(aq.relatedQuestions)
                    ? aq.relatedQuestions.map((row) => ({
                        id: typeof row?.id === 'number' ? row.id : Number(row?.id) || 0,
                        questionId:
                            row?.questionId != null && row.questionId !== ''
                                ? String(row.questionId)
                                : String(aq.id ?? ''),
                        relatedQuestionId:
                            row?.relatedQuestionId != null && row.relatedQuestionId !== ''
                                ? String(row.relatedQuestionId)
                                : String(this.relatedQuestionTargetId(row)),
                    }))
                    : [];

                /** PUT: MemoForm は questionId / userId が文字列（Go の json との整合） */
                const memosPayload = Array.isArray(aq.memos)
                    ? aq.memos.map((m) => ({
                        id: m?.id != null ? Number(m.id) || 0 : 0,
                        questionId:
                            m?.questionId != null && m.questionId !== ''
                                ? String(m.questionId)
                                : String(aq.id ?? ''),
                        userId:
                            m?.userId != null && m.userId !== ''
                                ? String(m.userId)
                                : '',
                        content:
                            m?.content != null ? String(m.content) : '',
                    }))
                    : [];

                const answersPayload = Array.isArray(aq.answers)
                    ? aq.answers.map((a) => ({
                        id: a?.id != null ? Number(a.id) || 0 : 0,
                        userId:
                            a?.userId != null && a.userId !== ''
                                ? String(a.userId)
                                : '',
                        content:
                            a?.content != null ? String(a.content) : '',
                        isFinal: Boolean(a?.isFinal),
                        createdAt: a?.createdAt ?? null,
                        updatedAt: a?.updatedAt ?? null,
                    }))
                    : [];

                const payload = {
                    ...aq,
                    createdAt: createdAtISO,
                    due: dueISO,
                    relatedQuestions: relatedQuestionsPayload,
                    memos: memosPayload,
                    answers: answersPayload,
                };
                stripSupportForPutIfUnassigned(payload);
                ensureSupportAssigneeForInProgressPut(payload, this.loginUser);
                const res = await fetch("/api/v1/question", {
                    method: "PUT",
                    headers: this.apiHeaders(),
                    body: JSON.stringify(Question.toModel(payload))
                });
                if (!res.ok) {
                    const errText = await res.text().catch(() => '');
                    console.error('updateQuestion failed:', res.status, errText);
                    return false;
                }
                return true;
            }
            return true;
        },

        async getQuestions() {
            const res = await fetch('/api/v1/question', { headers: this.apiHeaders(false) });
            const json = await res.json();
            const list = Array.isArray(json) ? json : [json];
            this.questions = list.map((q) => this.toQuestionViewModel(q));
        },

        async getLoginUser() {
            const res = await fetch('/api/v1/user/t', { headers: this.apiHeaders(false) });
            const json = await res.json();
            this.loginUser = User.fromJSON(json);
        },

        async getUsers() {
            const res = await fetch('/api/v1/user', { headers: this.apiHeaders(false) });
            const json = await res.json();
            this.users = Array.isArray(json) ? json.map((u) => User.fromJSON(u)) : [User.fromJSON(json)];
        },

        async deleteUser(id) {
            await fetch(`/api/v1/user/${id}`, {
                method: "DELETE",
                headers: this.apiHeaders(false)
            }).then(res => {
                if (res.ok) {
                    this.users = [...this.users.filter(u => u.id !== id)];
                }
            });
        },

        openUser(u) {
            this.activeUser = u;
        },

        async updateUser(role, roleName) {
            this.activeUser.role = { id: Number(role), name: roleName };
            const data = JSON.stringify(User.toModel(this.activeUser));
            await fetch("/api/v1/user", {
                method: "PUT",
                headers: this.apiHeaders(),
                body: data
            });
        },

        async updateLoginUser() {
            const data = JSON.stringify(User.toModel(this.loginUser));
            await fetch("/api/v1/user", {
                method: "PUT",
                headers: this.apiHeaders(),
                body: data
            });
        },

        async getTags() {
            const res = await fetch("/api/v1/tag", { headers: this.apiHeaders(false) });
            const json = await res.json();
            const list = Array.isArray(json) ? json : [json];
            this.tags = list.map((t) => Tag.fromJSON(t));
            this.availableTags = list.map((t) => Tag.fromJSON(t));
        },

        async registerTag(name, categoryId, categoryText) {
            const data = Tag.toModel({ name: name, categoryId: categoryId, usage: 0, category: { id: Number(categoryId), name: categoryText } });
            const res = await fetch('/api/v1/tag', {
                method: "POST",
                headers: this.apiHeaders(),
                body: JSON.stringify(data)
            });
            await res.json();
        },

        async deleteTag(id) {
            await fetch(`/api/v1/tag/${id}`, {
                method: "DELETE",
                headers: this.apiHeaders(false)
            }).then(res => {
                if (res.ok) {
                    const newTags = [...this.tags.filter(t => t.id !== id)]
                    this.tags = newTags
                    this.availableTags = newTags.map((t) => Tag.fromJSON(t));
                }
            });
        },

        openEditTag(t) {
            this.activeTag = {
                ...t,
                categoryId: t.categoryId != null ? String(t.categoryId) : '0',
            };
        },

        async updateTag(categoryId, categoryName) {
            this.activeTag.category = { id: Number(categoryId), name: categoryName };
            const tag = Tag.toModel(this.activeTag);
            const data = JSON.stringify(tag);
            await fetch('/api/v1/tag', {
                method: 'PUT',
                headers: this.apiHeaders(),
                body: data
            });
        },

        toggleSidebar() {
            this.isSidebarOpen = !this.isSidebarOpen;
        },

        calcRemainingTimeAndString(question) {
            const baseDate = new Date(question.due);
            const d = this.timeNow;
            const diff = baseDate - d;

            if (diff >= 0) {
                const [hours, minutes] = this.calcRemainingTime(baseDate, d);
                return `${hours}時間${minutes}分`;
            }

            const overdueMs = Math.abs(diff);
            const overdueHours = Math.floor(overdueMs / (1000 * 60 * 60));

            if (overdueHours < 24) {
                return `${overdueHours}時間超過中`;
            }
            const overdueDays = Math.floor(overdueHours / 24);
            return `${overdueDays}日超過中`;
        },

        calcRemainingTime(baseDate, d = new Date()) {
            const diff = baseDate - d;
            const hours = Math.floor(diff / (1000 * 60 * 60));
            const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

            return [hours, minutes];
        },

        checkEmailFmt(text) {
            return checkEmailFormat(text);
        },

        passwordCleansing(text) {
            return cleansingPassword(text);
        },

        emailCleansing(text) {
            return cleansingEmail(text);
        },

        async registerUser(email, name, pass, role, roleName) {
            const data = User.toModel({ email: email, name: name, pass: pass, roleId: role, role: { id: Number(role), name: roleName } });
            await fetch("/api/v1/user", {
                method: "POST",
                body: JSON.stringify(data),
                headers: this.apiHeaders(),
            }).catch({});
        },

        /**
         * @param {string} funcName
         */
        scrollBottom(funcName) {
            this.$nextTick(() => {
                const el = this.$refs[`scrollBox-${funcName}`];
                if (el) {
                    el.scrollTop = el.scrollHeight;
                }
            });
        },

        generatePassword(length) {
            let passLength = typeof length === 'string' ? parseInt(length, 10) : Number(length);
            if (!Number.isFinite(passLength) || passLength < 10) {
                passLength = 10;
            } else if (passLength > 64) {
                passLength = 64;
            }
            const ALLOWED_CHARS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-=[]{};\':"./<>?\\|~`';
            return Array.from(new Array(passLength), () => _.sample(ALLOWED_CHARS)).join('');
        },

        refreshIconsByElement(el, query) {
            this.$nextTick(() => {
                if (typeof lucide === 'undefined' || !el) return;
                const sel = query === '' || query === undefined ? '[data-lucide]' : `${query} [data-lucide]`;
                lucide.createIcons({ nodes: el.querySelectorAll(sel) });
            });
        },

        refreshIcons() {
            this.$nextTick(() => {
                if (typeof lucide !== 'undefined') {
                    lucide.createIcons();
                }
            });
        }
    }));
});
