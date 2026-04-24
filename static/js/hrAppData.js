import { Question, User, Tag } from '/static/js/model.js';
import { SSE_KEY } from './sse.js';

document.addEventListener('alpine:init', () => {
    Alpine.data('hrAppData', () => ({
        isSidebarOpen: true,
        currentView: 'home',
        viewMode: 'card',
        isManager: true,
        showAddUserModal: false,
        showEditUserModal: false,
        showTagModal: false,
        showEditTagModal: false,
        timeNow: new Date(),
        unreadNotices: 0,
        activeQuestion: {},
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

        statusColorMap: {
            '未対応': 'bg-red-100 text-red-800',
            '対応中': 'bg-blue-100 text-blue-800',
            '完了': 'bg-gray-100 text-gray-600'
        },

        questions: [
            { id: 1002, sender: '鈴木 一郎', department: 'ここに部署名', title: '通勤手当の経路変更について', content: '引越しに伴い、通勤経路が変更になります。申請手順と必要な書類を教えてください。', support: { supportStatusId: "1", supportStatus: { id: "1", title: '未対応' } }, tags: [{ id: 50, title: '諸手当' }], daysLeft: 1, date: '2026-04-08 09:30', dueDate: '2026-04-10', relatedQuestions: [] },
            { id: 1003, sender: '田中 花子', department: 'ここに部署名', title: '育児休業の延長申請', content: '現在取得中の育休を半年間延長したいと考えています。手続きの流れを教えてください。', support: { supportStatusId: "2", supportStatus: { id: "2", title: '対応中' } }, tags: [{ id: 60, title: '休暇' }], daysLeft: 5, date: '2026-04-04 14:00', dueDate: '2026-04-14', relatedQuestions: [] },
            { id: 1004, sender: '佐藤 次郎', department: 'ここに部署名', title: '健康診断の受診日変更', content: '指定された健康診断の日程ですが、出張と重なってしまいました。', support: { supportStatusId: "3", supportStatus: { id: "3", title: '完了' } }, tags: [{ id: 70, title: '健康診断' }], daysLeft: 0, date: '2026-04-01 11:15', dueDate: '2026-04-09', relatedQuestions: [] },
            { id: 1005, sender: '高橋 三郎', department: 'ここに部署名', title: '慶弔休暇の適用範囲', content: '配偶者の祖父母が亡くなった場合、忌引休暇の対象になりますでしょうか？', support: { supportStatusId: "1", supportStatus: { id: "1", title: '未対応' } }, tags: [{ id: 80, title: '規程' }, { id: 70, title: '健康診断' }], daysLeft: 2, date: '2026-04-07 16:45', dueDate: '2026-04-11', relatedQuestions: [] }
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

        init() {
            this.refreshIcons();

            document.addEventListener('connect', (event) => {
                this.isConnect = event.detail;
            });
            document.addEventListener('disconnect', (event) => {
                this.isConnect = event.detail;
            });

            document.addEventListener(SSE_KEY.system.timeTick, (event) => {
                this.timeNow = event.detail
            });

            document.addEventListener(SSE_KEY.data.create.notice, (event) => {
                const notice = event.detail;
                const index = this.notices.findIndex(n => n.id === notice.id)
                if (index === -1) {
                    this.notices.unshift(notice);
                    this.refreshIcons();
                } else {
                    if (!_.isEqual(this.notices[index], notice)) {
                        this.notices.splice(index, 1, {
                            ...this.notices[index],
                            ...notice,
                        });
                        this.refreshIcons();
                    }
                }
                if (notice.typeId !== 1) this.unreadNotices += 1;
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
                } else {
                    if (!_.isEqual(this.questions[index], question)) {
                        this.questions.splice(index, 1, {
                            ...this.questions[index],
                            ...question,
                        });
                    }
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
                                    return localVal
                                }
                                // server採用
                                return serverVal
                            }
                        );
                        this.originalQuestion = _.cloneDeep(question);
                        this.refreshIcons();
                    }
                }
            });

            document.addEventListener(SSE_KEY.data.delete.question, (event) => {
                const id = event.detail;
                this.questions = this.questions.filter((q) => q.id !== id);
                if (this.activeQuestion?.id === id) {
                    this.activeQuestion = {};
                    this.originalQuestion = {};
                }
            });

            document.addEventListener(SSE_KEY.data.create.user, (event) => {
                const user = event.detail;
                const index = this.users.findIndex(q => q.id === user.id);
                if (index === -1) {
                    this.users.unshift(user);
                    this.refreshIcons();
                } else {
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
                const index = this.tags.findIndex(t => t.id === tag.id);
                if (index === -1) {
                    this.tags.unshift(tag);
                    this.availableTags = this.tags.map((t) => Tag.fromJSON(t));
                    this.refreshIcons();
                } else {
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

            this.getTags()
        },

        getIcon(name) {
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
            const daysLeft = due
                ? Math.max(0, Math.ceil((due.getTime() - Date.now()) / (1000 * 60 * 60 * 24)))
                : 0;

            return {
                id: question.id,
                sender: question.support?.user?.name ?? '不明',
                department: 'ここに部署名',
                title: question.title ?? '',
                content: question.content ?? '',
                support: question.support ?? null,
                tags: question.tags,
                daysLeft,
                date: this.formatDateTime(createdAt),
                dueDate: due ? this.formatDate(due) : '',
                relatedQuestions: (question.subQuestions ?? []).map((sq) =>
                    typeof sq === 'object' && sq !== null ? sq.id : sq),
            };
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
            if (!v.support) {
                v.support = { supportStatusId: '1', supportStatus: { id: '1', title: '未対応' }, user: { name: '' } };
            }
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
            console.log(index)
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
            if (this.activeQuestion.relatedQuestions.includes(qId)) {
                this.activeQuestion.relatedQuestions = this.activeQuestion.relatedQuestions.filter(id => id !== qId);
            } else {
                this.activeQuestion.relatedQuestions.push(qId);
            }
            this.refreshIcons();
        },

        async updateQuestion() {
            const question = this.questions.filter(q => q.id === this.activeQuestion.id)[0];
            if (!_.isEqual(question, this.activeQuestion)) {
                await fetch("/api/v1/question", {
                    method: "PUT",
                    headers: this.apiHeaders(),
                    body: JSON.stringify(Question.toModel(this.activeQuestion))
                });
            }
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

        refreshIcons() {
            this.$nextTick(() => {
                if (typeof lucide !== 'undefined') {
                    lucide.createIcons();
                }
            });
        }
    }));
});
