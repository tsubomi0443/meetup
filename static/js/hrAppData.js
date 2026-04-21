import { Question, User, Tag } from '/static/js/model.js';

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
        unreadNotices: 1,
        activeQuestion: {},
        activeUser: {},
        activeTag: {},
        // TODO; InitでTokenからTokenを取得し、通信のタイミングで投げるように変更する。
        // authorizationHeader: { "Authorization": "Bearer " + token },
        isConnect: false,

        selectedTags: [],
        availableTags: ['諸手当', '休暇', '規程', '健康診断'],
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
            { id: 1002, sender: '鈴木 一郎', department: 'ここに部署名', title: '通勤手当の経路変更について', content: '引越しに伴い、通勤経路が変更になります。申請手順と必要な書類を教えてください。', status: '未対応', tags: ['諸手当'], daysLeft: 1, date: '2026-04-08 09:30', dueDate: '2026-04-10', relatedQuestions: [] },
            { id: 1003, sender: '田中 花子', department: 'ここに部署名', title: '育児休業の延長申請', content: '現在取得中の育休を半年間延長したいと考えています。手続きの流れを教えてください。', status: '対応中', tags: ['休暇'], daysLeft: 5, date: '2026-04-04 14:00', dueDate: '2026-04-14', relatedQuestions: [] },
            { id: 1004, sender: '佐藤 次郎', department: 'ここに部署名', title: '健康診断の受診日変更', content: '指定された健康診断の日程ですが、出張と重なってしまいました。', status: '完了', tags: ['健康診断'], daysLeft: 0, date: '2026-04-01 11:15', dueDate: '2026-04-09', relatedQuestions: [] },
            { id: 1005, sender: '高橋 三郎', department: 'ここに部署名', title: '慶弔休暇の適用範囲', content: '配偶者の祖父母が亡くなった場合、忌引休暇の対象になりますでしょうか？', status: '未対応', tags: ['規程'], daysLeft: 2, date: '2026-04-07 16:45', dueDate: '2026-04-11', relatedQuestions: [] }
        ],

        notices: [
            { id: 1, type: 'alert', title: '期日接近: 質問 #1002', time: '10分前', content: '「通勤手当の経路変更」の期日が明日です。早急な対応をお願いします。' },
            { id: 2, type: 'info', title: 'システムアップデート完了', time: '1時間前', content: 'FAQ検索エンジンの精度を向上させるアップデートを適用しました。' }
        ],

        users: [],
        tags: [],

        get filteredQuestions() {
            let result = this.questions.filter(q => {
                const matchesTag = this.selectedTags.length === 0 || q.tags.some(tag => this.selectedTags.includes(tag));
                const matchesSearch = q.title.includes(this.searchQuery) || q.content.includes(this.searchQuery) || q.sender.includes(this.searchQuery);
                let matchesStatus = true;
                if (this.statusFilter === 'unanswered') matchesStatus = q.status === '未対応';
                else if (this.statusFilter === 'progress') matchesStatus = q.status === '対応中';
                else if (this.statusFilter === 'open') matchesStatus = (q.status === '未対応' || q.status === '対応中');
                else if (this.statusFilter === 'completed') matchesStatus = q.status === '完了';
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
            const t = localStorage.getItem('access_token');
            if (t) h['Authorization'] = 'Bearer ' + t;
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

            document.addEventListener('sse-notice', (event) => {
                const notice = event.detail;
                this.notices.unshift(notice);
                if (notice.type === 'alert') this.unreadNotices += 1;
                this.refreshIcons();
            });

            document.addEventListener('sse-question', (event) => {
                const question = this.toQuestionViewModel(event.detail);
                const index = this.questions.findIndex(q => q.id === question.id);
                if (index === -1) {
                    this.questions.unshift(question);
                    this.refreshIcons();
                } else {
                    if (!_.isEqual(this.questions[index], question)) {
                        this.questions.splice(index, 1, {
                            ...this.questions[index],
                            ...question,
                        });
                        this.refreshIcons();
                    }
                }

                if (question.status === '完了') {
                    const before = this.notices.length;
                    this.notices = this.notices.filter(
                        n => !(n.type === 'alert' && n.title.includes(`#${question.id}`))
                    );
                    this.unreadNotices = Math.max(0, this.unreadNotices - (before - this.notices.length));
                }

                if (this.activeQuestion?.id === question.id) {
                    this.activeQuestion = {
                        ...this.activeQuestion,
                        ...question,
                    };
                }

                this.refreshIcons();
            });

            document.addEventListener('sse-user', (event) => {
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

            document.addEventListener('sse-tag', (event) => {
                const tag = event.detail;
                const index = this.tags.findIndex(q => q.id === tag.id);
                if (index === -1) {
                    this.tags.unshift(tag);
                    this.availableTags = this.tags.map((t) => t.title);
                    this.refreshIcons();
                } else {
                    if (!_.isEqual(this.tags[index], tag)) {
                        this.tags.splice(index, 1, {
                            ...this.tags[index],
                            ...tag,
                        });
                        this.availableTags = this.tags.map((t) => t.title);
                        this.refreshIcons();
                    }
                }
            });
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
                status: question.support?.supportStatus?.title ?? '未対応',
                tags: (question.tags ?? []).map(tag => tag.title),
                daysLeft,
                date: this.formatDateTime(createdAt),
                dueDate: due ? this.formatDate(due) : '',
                relatedQuestions: (question.subQuestions ?? []).map(q => q.id),
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
            localStorage.removeItem("access_token");
            location.href = "/login";
        },

        setView(view) {
            this.currentView = view;
            this.refreshIcons();
        },

        openDetail(q) {
            this.activeQuestion = q;
            this.relatedSearchQuery = '';
            this.setView('detail');
            return
        },

        toggleTag(tag) {
            if (this.selectedTags.includes(tag)) {
                this.selectedTags = this.selectedTags.filter(t => t !== tag);
            } else {
                this.selectedTags.push(tag);
            }
            this.refreshIcons();
        },

        toggleQuestionTag(tag) {
            if (!this.activeQuestion.tags) {
                this.activeQuestion.tags = [];
            }
            if (this.activeQuestion.tags.includes(tag)) {
                this.activeQuestion.tags = this.activeQuestion.tags.filter(t => t !== tag);
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
            console.log(question, this.activeQuestion, _.isEqual(question, this.activeQuestion))
            if (!_.isEqual(question, this.activeQuestion)) {
                console.log(JSON.stringify(Question.toModel(this.activeQuestion)))
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
                    this.users = [...this.users.filter(u => u.id !== id)]
                }
            })
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
            })
        },

        async getTags() {
            const res = await fetch("/api/v1/tag", { headers: this.apiHeaders(false) });
            const json = await res.json();
            const list = Array.isArray(json) ? json : [json];
            this.availableTags = list.map((t) => Tag.fromJSON(t));
        },

        async registerTag(name, cateogryId) {
            const data = Tag.toModel({title: name, categoryId: cateogryId, usage: 0});
            const res = await fetch('/api/v1/tag', {
                method: "POST",
                headers: this.apiHeaders(),
                body: JSON.stringify(data)
            });
            const json = await res.json();
            console.log(json);
        },

        async deleteTag(id) {
            await fetch(`/api/v1/tag/${id}`, {
                method: "DELETE",
                headers: this.apiHeaders(false)
            }).then(res => {
                if (res.ok) {
                    const newTags = [...this.tags.filter(t => t.id !== id)]
                    this.tags = newTags
                    this.availableTags = newTags.map((t) => t.title);
                }
            })
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
            })
        },

        toggleSidebar() {
            this.isSidebarOpen = !this.isSidebarOpen;
        },

        async registerUser(email, name, pass, role) {
            const data = User.toModel({email: email, name: name, password: pass, roleId: role});
            await fetch("/api/v1/user", {
                method: "POST",
                body: JSON.stringify(data),
                headers: this.apiHeaders(),
            }).catch({})
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
