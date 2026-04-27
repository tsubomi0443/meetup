import { Question, User, Tag, Notice } from '/static/js/model.js';
import { SSE_KEY } from './sse.js';

/** API / モックで support はあるが user が無い場合がある。テンプレは activeQuestion.support.user.name を前提にする */
function ensureSupportForView(support) {
    const defaults = {
        supportStatusId: '1',
        supportStatus: { id: '1', title: '未対応' },
        user: { name: '' }
    };
    if (!support) return defaults;
    const o = { ...support };
    o.user = support.user && typeof support.user === 'object'
        ? { name: String(support.user.name ?? '') }
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
        answer: null,
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

/** Question API/SSE から answer を view 用に正規化 */
function normalizeAnswerForViewModel(question) {
    const a = question.answer;
    if (!a) return null;
    return {
        id: a.id,
        userId: a.userId,
        content: a.content ?? '',
        createdAt: a.createdAt ?? null,
        user: a.user && typeof a.user === 'object'
            ? { id: a.user.id, name: String(a.user.name ?? '') }
            : { id: a.userId, name: '' },
    };
}

document.addEventListener('alpine:init', () => {
    Alpine.data('hrAppData', () => ({
        timeNow: new Date(),
        isSidebarOpen: true,
        currentView: 'home',
        viewMode: 'card',
        isManager: true,
        showAddUserModal: false,
        showEditUserModal: false,
        showTagModal: false,
        showEditTagModal: false,
        unreadNotices: 0,
        activeQuestion: emptyActiveQuestion(),
        originalQuestion: {},
        activeUser: {},
        activeTag: {},
        isConnect: false,

        selectedTags: [],
        availableTags: [],
        searchQuery: '',
        statusFilter: 'all',
        sortOrder: 'date_desc',

        relatedSearchQuery: '',

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

        statusColorMap: {
            '未対応': 'bg-red-100 text-red-800',
            '対応中': 'bg-blue-100 text-blue-800',
            '完了': 'bg-gray-100 text-gray-600'
        },

        questions: [
            { id: 1002, sender: '鈴木 一郎', department: 'ここに部署名', title: '通勤手当の経路変更について', content: '引越しに伴い、通勤経路が変更になります。申請手順と必要な書類を教えてください。', support: { supportStatusId: "1", supportStatus: { id: "1", title: '未対応' } }, tags: [{ id: 50, title: '諸手当' }], daysLeft: 1, date: '2026-04-08 09:30', dueDate: '2026-04-10', relatedQuestions: [], memos: [], answer: null },
            { id: 1003, sender: '田中 花子', department: 'ここに部署名', title: '育児休業の延長申請', content: '現在取得中の育休を半年間延長したいと考えています。手続きの流れを教えてください。', support: { supportStatusId: "2", supportStatus: { id: "2", title: '対応中' } }, tags: [{ id: 60, title: '休暇' }], daysLeft: 5, date: '2026-04-04 14:00', dueDate: '2026-04-14', relatedQuestions: [], memos: [], answer: null },
            { id: 1004, sender: '佐藤 次郎', department: 'ここに部署名', title: '健康診断の受診日変更', content: '指定された健康診断の日程ですが、出張と重なってしまいました。', support: { supportStatusId: "3", supportStatus: { id: "3", title: '完了' } }, tags: [{ id: 70, title: '健康診断' }], daysLeft: 0, date: '2026-04-01 11:15', dueDate: '2026-04-09', relatedQuestions: [], memos: [], answer: null },
            { id: 1005, sender: '高橋 三郎', department: 'ここに部署名', title: '慶弔休暇の適用範囲', content: '配偶者の祖父母が亡くなった場合、忌引休暇の対象になりますでしょうか？', support: { supportStatusId: "1", supportStatus: { id: "1", title: '未対応' } }, tags: [{ id: 80, title: '規程' }, { id: 70, title: '健康診断' }], daysLeft: 2, date: '2026-04-07 16:45', dueDate: '2026-04-11', relatedQuestions: [], memos: [], answer: null }
        ],

        notices: [],

        users: [],
        tags: [],

        get filteredQuestions() {
            let result = this.questions.filter(q => {
                const matchesTag = this.selectedTags.length === 0 || q.tags.some(tag => this.selectedTags.includes(tag.title));
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

        apiHeaders(withJSON = true) {
            const h = {};
            if (withJSON) h['Content-Type'] = 'application/json';
            return h;
        },

        initData() {
            this.getLoginUser();
            this.getQuestions();
            this.getUsers();
            this.getTags();
        },

        init() {
            this.refreshIcons();

            document.addEventListener('connect', (event) => {
                this.isConnect = event.detail;
                this.initData()
            });
            document.addEventListener('disconnect', (event) => {
                this.isConnect = event.detail;
            });

            document.addEventListener(SSE_KEY.system.timeTick, (event) => {
                this.timeNow = event.detail
            });

            document.addEventListener(SSE_KEY.data.create.notice, (event) => {
                const notice = event.detail;
                const index = this.notices.findIndex((n) => n.id === notice.id);
                if (index === -1) {
                    this.notices.unshift(notice);
                    if (notice.typeId !== 1) this.unreadNotices += 1;
                    this.refreshIcons();
                }
            });

            document.addEventListener(SSE_KEY.data.update.notice, (event) => {
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

            document.addEventListener(SSE_KEY.data.delete.notice, (event) => {
                const id = event.detail;
                this.notices = this.notices.filter((n) => n.id !== id);
            });

            document.addEventListener(SSE_KEY.data.create.question, (event) => {
                const question = this.toQuestionViewModel(event.detail);
                const index = this.questions.findIndex(q => q.id === question.id);
                if (index === -1) {
                    this.questions.unshift(question);
                    this.refreshIcons();
                }
            });

            document.addEventListener(SSE_KEY.data.update.question, (event) => {
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
                                        if (key === "support") return serverVal; // support.supportStatusの更新が必ずLocalWinになってしまうため例外
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

            document.addEventListener(SSE_KEY.data.delete.question, (event) => {
                const id = event.detail;
                if (this.activeQuestion?.id === id) {
					this.questions = this.questions.filter((q) => q.id !== id);
                    this.activeQuestion = emptyActiveQuestion();
                    this.originalQuestion = {};
                }
            });

            document.addEventListener(SSE_KEY.data.create.user, (event) => {
                const user = event.detail;
                const index = this.users.findIndex((u) => u.id === user.id);
                if (index === -1) {
                    this.users.unshift(user);
                    this.refreshIcons();
                }
            });

            document.addEventListener(SSE_KEY.data.update.user, (event) => {
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

            document.addEventListener(SSE_KEY.data.delete.user, (event) => {
                const id = event.detail;
                this.users = this.users.filter((u) => u.id !== id);
            });

            document.addEventListener(SSE_KEY.data.create.tag, (event) => {
                const tag = event.detail;
                const index = this.tags.findIndex((t) => t.id === tag.id);
                if (index === -1) {
                    this.tags.unshift(tag);
                    this.availableTags = this.tags.map((t) => Tag.fromJSON(t));
                    this.refreshIcons();
                }
            });

            document.addEventListener(SSE_KEY.data.update.tag, (event) => {
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

            document.addEventListener(SSE_KEY.data.delete.tag, (event) => {
                const id = event.detail;
                const newTags = this.tags.filter((t) => t.id !== id);
                this.tags = newTags;
                this.availableTags = newTags.map((t) => Tag.fromJSON(t));
            });
        },

        getIcon(name) {
        },

        // --- mock5: 表示用ヘルパー（HTMLのAlpine式からロジックを分離） ---

        hasSelectedTags() {
            return this.selectedTags.length > 0;
        },

        questionStatusTitle(q) {
            if (!q?.support?.supportStatus?.title) return '';
            const id = Number(q?.support?.supportStatusId);
            const title = this.statusTitleMap[id];

            this.activeQuestion.support.supportStatus = {id: id, title: title}
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
            if (u?.role && u.role.roleName) return String(u.role.roleName);
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

        chatAnswerContent() {
            return String(this.activeQuestion?.answer?.content ?? '');
        },

        isSelfMemo(memo) {
            if (this.activeUser?.id == null || memo?.userId == null) return false;
            return Number(memo.userId) === Number(this.activeUser.id);
        },

        isSelfAnswer() {
            const a = this.activeQuestion?.answer;
            if (this.activeUser?.id == null || a?.userId == null) return false;
            return Number(a.userId) === Number(this.activeUser.id);
        },

        chatAvatarSeed(name) {
            return String(name ?? 'User');
        },

        chatTimelineItems() {
            const memos = (this.activeQuestion?.memos ?? []).map((m, i) => ({
                kind: 'memo',
                userId: m.userId,
                userName: m.user?.name ?? '',
                content: this.chatMemoContent(m),
                createdAt: m.createdAt ? new Date(m.createdAt) : null,
                _origIdx: i,
                memo: m,
            }));
            const a = this.activeQuestion?.answer;
            const items = [...memos];
            if (a) {
                items.push({
                    kind: 'answer',
                    userId: a.userId,
                    userName: a.user?.name ?? '',
                    content: this.chatAnswerContent(),
                    createdAt: a.createdAt ? new Date(a.createdAt) : null,
                    _origIdx: memos.length,
                    answer: a,
                });
            }
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
            if (item?.kind === 'answer') return this.isSelfAnswer();
            return false;
        },

        activeQuestionTagButtonLabel() {
            const tags = this.activeQuestion?.tags;
            if (!Array.isArray(tags) || tags.length === 0) return 'タグを選択';
            return tags.map((t) => t.title).join(', ');
        },

        isActiveQuestionTagSelected(tag) {
            const tags = this.activeQuestion?.tags;
            if (!Array.isArray(tags) || !tag) return false;
            return tags.filter((t) => t.title === tag.title).length === 1;
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
            if (!n?.question?.due) return '';
            return `残り時間: ${this.calcRemainingTimeAndString(n.question)}`;
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
            const title = tag?.title != null ? String(tag.title) : '';
            if (!title) return 0;
            return this.questions.filter((q) => {
                const tags = q?.tags;
                if (!Array.isArray(tags)) return false;
                return tags.some((t) => t && t.title === title);
            }).length;
        },

        tagUsageCountLabel(tag) {
            return `${this.countQuestionsUsingTagTitle(tag)}件`;
        },

        tagCategoryDisplayName(tag) {
            if (tag?.category && tag.category.categoryName) {
                return String(tag.category.categoryName);
            }
            return '';
        },

        currentTimeFormatted(dt = new Date()) {
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
                sender: question.support?.user?.name ?? '不明',
                department: 'ここに部署名',
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
                answer: normalizeAnswerForViewModel(question),
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
            location.href = "/login";
        },

        setView(view) {
            this.currentView = view;
            this.refreshIcons();
        },

        openDetail(q) {
            const v = _.cloneDeep(q);
            v.support = ensureSupportForView(v.support);
            this.originalQuestion = _.cloneDeep(v);
            this.activeQuestion = v;
            this.relatedSearchQuery = '';
            this.setView('detail');
        },

        openDetailByID(id = 0) {
            if (id === 0) {
                window.alert(`invalid support id is ${id}.`)
            }
            const index = this.questions.findIndex(q => q.id == id)
            if (index !== -1) {
                this.openDetail(this.questions[index])
            }
        },

        toggleTag(tag) {
            if (this.selectedTags.includes(tag.title)) {
                this.selectedTags = this.selectedTags.filter(t => t !== tag.title);
            } else {
                this.selectedTags.push(tag.title);
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
            this.refreshIcons();
        },

        async updateQuestion() {
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

                const payload = {
                    ...aq,
                    createdAt: createdAtISO,
                    due: dueISO,
                    relatedQuestions: relatedQuestionsPayload,
                };
                await fetch("/api/v1/question", {
                    method: "PUT",
                    headers: this.apiHeaders(),
                    body: JSON.stringify(Question.toModel(payload))
                });
            }
        },

        async getQuestions() {
            const res = await fetch('/api/v1/question', { headers: this.apiHeaders(false) });
            const json = await res.json();
            const list = Array.isArray(json) ? json : [json];
            console.log(list);
            this.questions = list.map((q) => this.toQuestionViewModel(q));
        },

        async getLoginUser() {
            const res = await fetch('/api/v1/user/t', { headers: this.apiHeaders(false) });
            const json = await res.json();
            this.activeUser = User.fromJSON(json);
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

        async updateUser() {
            const data = JSON.stringify(User.toModel(this.activeUser));
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

        async registerTag(name, cateogryId, categoryText) {
            const data = Tag.toModel({ title: name, categoryId: cateogryId, usage: 0, category: { categoryName: categoryText } });
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

        async updateTag() {
            const data = JSON.stringify(Tag.toModel(this.activeTag));
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
            const [hours, minutes] = this.calcRemainingTime(new Date(question.due), this.timeNow)
            return `${hours}時間${minutes}分`;
        },

        calcRemainingTime(baseDate, d = new Date()) {
            const diff = baseDate - d;
            const hours = Math.floor(diff / (1000 * 60 * 60));
            const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

            return [hours, minutes];
        },

        async registerUser(email, name, pass, role) {
            const data = User.toModel({ email: email, name: name, password: pass, roleId: role });
            await fetch("/api/v1/user", {
                method: "POST",
                body: JSON.stringify(data),
                headers: this.apiHeaders(),
            }).catch({});
        },

        calcRemainingTimeAndString(question) {
            const [hours, minutes] = this.calcRemainingTime(new Date(question.due), this.timeNow)
            return `${hours}時間${minutes}分`;
        },


        calcRemainingTime(baseDate, d = new Date()) {
            const diff = baseDate - d;
            const hours = Math.floor(diff / (1000 * 60 * 60));
            const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

            return [hours, minutes];
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
