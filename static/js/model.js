/**
 * model.js
 * Go の _mac_infrastructure（Form / entity）に対応するクライアント側クラス。
 * サーバは camelCase の Form JSON を返す。SSE / API の JSON を fromJSON() で復元する。
 *
 * 使用例:
 *   const es = new EventSource('/sse');
 *   es.addEventListener('question', (e) => {
 *     const q = Question.fromJSON(JSON.parse(e.data));
 *     console.log(q.title);
 *   });
 */

export { formToApiJson } from '/static/js/models/helpers.js';
export { Role } from '/static/js/models/Role.js';
export { SupportStatus } from '/static/js/models/SupportStatus.js';
export { Support } from '/static/js/models/Support.js';
export { User } from '/static/js/models/User.js';
export { Category } from '/static/js/models/Category.js';
export { Tag } from '/static/js/models/Tag.js';
export { Refer } from '/static/js/models/Refer.js';
export { Memo } from '/static/js/models/Memo.js';
export { Answer } from '/static/js/models/Answer.js';
export { Escalation } from '/static/js/models/Escalation.js';
export { NoticeType } from '/static/js/models/NoticeType.js';
export { Notice } from '/static/js/models/Notice.js';
export { RelatedQuestion } from '/static/js/models/RelatedQuestion.js';
export { Sender } from '/static/js/models/Sender.js';
export { SenderTalk } from '/static/js/models/SenderTalk.js';
export { Question } from '/static/js/models/Question.js';
export { ReferManager } from '/static/js/models/ReferManager.js';
export { TagManager } from '/static/js/models/TagManager.js';
