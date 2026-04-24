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

// ---------------------------------------------------------------------------
// JSON helpers (camelCase + ISO8601)
// ---------------------------------------------------------------------------

/**
 * API 送信用のプレーンオブジェクト（Date は ISO8601）。
 * @param {object} form
 * @returns {Object}
 */
export function formToApiJson(form) {
  return JSON.parse(
    JSON.stringify(form, (_key, value) =>
      value instanceof Date ? value.toString() : value,
    ),
  );
}

// ---------------------------------------------------------------------------
// Role
// ---------------------------------------------------------------------------
export class Role {
  /**
   * @param {number} id
   * @param {string} roleName
   * @param {User[]} [users]
   */
  constructor(id, roleName, users = []) {
    this.id = id;
    this.roleName = roleName;
    this.users = users;
  }

  /** @param {Object} json @returns {Role} */
  static fromJSON(json) {
    return new Role(
      json.id,
      json.roleName,
      (json.users ?? []).map(User.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// SupportStatus
// ---------------------------------------------------------------------------
export class SupportStatus {
  /**
   * @param {number} id
   * @param {string} title
   * @param {Support[]} [supports]
   */
  constructor(id, title, supports = []) {
    this.id = id;
    this.title = title;
    this.supports = supports;
  }

  /** @param {Object} json @returns {SupportStatus} */
  static fromJSON(json) {
    return new SupportStatus(
      json.id,
      json.title,
      (json.supports ?? []).map(Support.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// Support
// ---------------------------------------------------------------------------
export class Support {
  /**
   * @param {number}        id
   * @param {number}        userId
   * @param {string}        supportStatusId
   * @param {User|null}     [user]
   * @param {SupportStatus|null} [supportStatus]
   * @param {Question|null} [question]
   */
  constructor(id, userId, supportStatusId, user = null, supportStatus = null, question = null) {
    this.id = id;
    this.userId = userId;
    this.supportStatusId = supportStatusId;
    this.user = user;
    this.supportStatus = supportStatus;
    this.question = question;
  }

  /** @param {Object} json @returns {Support} */
  static fromJSON(json) {
    return new Support(
      json.id,
      json.userId,
      json.supportStatusId,
      json.user ? User.fromJSON(json.user) : null,
      json.supportStatus ? SupportStatus.fromJSON(json.supportStatus) : null,
      json.question ? Question.fromJSON(json.question) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// User
// ---------------------------------------------------------------------------
export class User {
  /**
   * @param {number}   id
   * @param {string}   name
   * @param {string}   email
   * @param {number}   roleId
   * @param {Role|null} [role]
   * @param {Support[]} [supports]
   * @param {Answer[]}  [answers]
   * @param {Memo[]}    [memos]
   */
  constructor(id, name, email, roleId, role = null, supports = [], answers = [], memos = []) {
    this.id = id;
    this.name = name;
    this.email = email;
    this.roleId = roleId;
    this.role = role;
    this.supports = supports;
    this.answers = answers;
    this.memos = memos;
  }

  /** @param {Object} json @returns {User} */
  static fromJSON(json) {
    // const rid = json.roleId;
    // const roleIdNum = 
    //   rid === null || rid === undefined || rid === ''
    //     ? rid
    //     : typeof rid === 'string'
    //       ? parseInt(rid, 10)
    //       : rid;
    return new User(
      json.id,
      json.name,
      json.email,
      json.roleId,
      json.role ? Role.fromJSON(json.role) : null,
      (json.supports ?? []).map(Support.fromJSON),
      (json.answers ?? []).map(Answer.fromJSON),
      (json.memos ?? []).map(Memo.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// Category
// ---------------------------------------------------------------------------
export class Category {
  /**
   * @param {number} id
   * @param {string} categoryName
   * @param {Tag[]}  [tags]
   */
  constructor(id, categoryName, tags = []) {
    this.id = id;
    this.categoryName = categoryName;
    this.tags = tags;
  }

  /** @param {Object} json @returns {Category} */
  static fromJSON(json) {
    return new Category(
      json.id,
      json.categoryName,
      (json.tags ?? []).map(Tag.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// Tag
// ---------------------------------------------------------------------------
export class Tag {
  /**
   * @param {number}      id
   * @param {string}      title
   * @param {number}      usage
   * @param {number}      categoryId
   * @param {Category|null} [category]
   * @param {Question[]}  [questions]
   */
  constructor(id, title, usage, categoryId, category = null, questions = []) {
    this.id = id;
    this.title = title;
    this.usage = usage;
    this.categoryId = categoryId;
    this.category = category;
    this.questions = questions;
  }

  /** @param {Object} json @returns {Tag} */
  static fromJSON(json) {
    // const cid = json.categoryId;
    // const categoryIdVal =
    //   cid === null || cid === undefined || cid === ''
    //     ? cid
    //     : typeof cid === 'string'
    //       ? parseInt(cid, 10)
    //       : cid;
    return new Tag(
      json.id,
      json.title,
      json.usage,
      json.categoryId,
      json.category ? Category.fromJSON(json.category) : null,
      (json.questions ?? []).map(Question.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// Refer
// ---------------------------------------------------------------------------
export class Refer {
  /**
   * @param {number}   id
   * @param {string}   title
   * @param {string}   url
   * @param {Answer[]} [answers]
   */
  constructor(id, title, url, answers = []) {
    this.id = id;
    this.title = title;
    this.url = url;
    this.answers = answers;
  }

  /** @param {Object} json @returns {Refer} */
  static fromJSON(json) {
    return new Refer(
      json.id,
      json.title,
      json.url,
      (json.answers ?? []).map(Answer.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// Memo
// ---------------------------------------------------------------------------
export class Memo {
  /**
   * @param {number}       id
   * @param {number}       questionId
   * @param {number}       userId
   * @param {string}       content
   * @param {Question|null} [question]
   * @param {User|null}    [user]
   */
  constructor(id, questionId, userId, content, question = null, user = null) {
    this.id = id;
    this.questionId = questionId;
    this.userId = userId;
    this.content = content;
    this.question = question;
    this.user = user;
  }

  /** @param {Object} json @returns {Memo} */
  static fromJSON(json) {
    return new Memo(
      json.id,
      json.questionId,
      json.userId,
      json.content,
      json.question ? Question.fromJSON(json.question) : null,
      json.user ? User.fromJSON(json.user) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// Answer
// ---------------------------------------------------------------------------
export class Answer {
  /**
   * @param {number}       id
   * @param {number|string} userId
   * @param {string}       content
   * @param {string|null}  answeredAt  ISO8601 文字列 or null
   * @param {string}       createdAt   ISO8601 文字列
   * @param {User|null}    [user]
   * @param {Refer[]}      [refers]
   */
  constructor(id, userId, content, answeredAt, createdAt, user = null, refers = []) {
    this.id = id;
    this.userId = userId;
    this.content = content;
    this.answeredAt = answeredAt ? new Date(answeredAt) : null;
    this.createdAt = new Date(createdAt);
    this.user = user;
    this.refers = refers;
  }

  /** @param {Object} json @returns {Answer} */
  static fromJSON(json) {
    return new Answer(
      json.id,
      json.userId,
      json.content,
      json.answeredAt ?? null,
      json.createdAt,
      json.user ? User.fromJSON(json.user) : null,
      (json.refers ?? []).map(Refer.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// Escalation
// ---------------------------------------------------------------------------
export class Escalation {
  /**
   * @param {number}       id
   * @param {number}       fromQuestionId
   * @param {number}       toQuestionId
   * @param {string}       escalatedAt    ISO8601 文字列
   * @param {Question|null} [fromQuestion]
   * @param {Question|null} [toQuestion]
   */
  constructor(id, fromQuestionId, toQuestionId, escalatedAt, fromQuestion = null, toQuestion = null) {
    this.id = id;
    this.fromQuestionId = fromQuestionId;
    this.toQuestionId = toQuestionId;
    this.escalatedAt = new Date(escalatedAt);
    this.fromQuestion = fromQuestion;
    this.toQuestion = toQuestion;
  }

  /** @param {Object} json @returns {Escalation} */
  static fromJSON(json) {
    return new Escalation(
      json.id,
      json.fromQuestionId,
      json.toQuestionId,
      json.escalatedAt ?? new Date().toISOString(),
      json.fromQuestion ? Question.fromJSON(json.fromQuestion) : null,
      json.toQuestion ? Question.fromJSON(json.toQuestion) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// NoticeType
// ---------------------------------------------------------------------------
export class NoticeType {
  /**
   * @param {number} id
   * @param {string} name
   * @param {Notice[]} [notices]
   */
  constructor(id, name, notices = []) {
    this.id = id;
    this.name = name;
    this.notices = notices;
  }

  /** @param {Object} json @returns {NoticeType} */
  static fromJSON(json) {
    return new NoticeType(
      json.id,
      json.name,
      (json.notices ?? []).map(Notice.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// Notice
// ---------------------------------------------------------------------------
export class Notice {
  /**
   * @param {number} id
   * @param {number} typeId
   * @param {number|null} questionId
   * @param {string|null} content
   * @param {string|null} displayDue ISO8601 or null
   * @param {NoticeType|null} [noticeType]
   * @param {Question|null} [question]
   */
  constructor(id, typeId, questionId, content, displayDue, noticeType = null, question = null) {
    this.id = id;
    this.typeId = typeId;
    this.questionId = questionId;
    this.content = content;
    this.displayDue = displayDue ? new Date(displayDue) : null;
    this.noticeType = noticeType;
    this.question = question;
  }

  /** @param {Object} json @returns {Notice} */
  static fromJSON(json) {
    return new Notice(
      json.id,
      json.typeId,
      json.questionId ?? null,
      json.content ?? null,
      json.displayDue ?? null,
      json.noticeType ? NoticeType.fromJSON(json.noticeType) : null,
      json.question ? Question.fromJSON(json.question) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// Question
// ---------------------------------------------------------------------------
export class Question {
  /**
   * @param {number}         id
   * @param {number|null}    originQuestionId
   * @param {number|null}    answerId
   * @param {number|null}    supportId
   * @param {string}         title
   * @param {string}         content
   * @param {boolean}        deleted
   * @param {string|null}    due              ISO8601 文字列 or null
   * @param {string}         createdAt        ISO8601 文字列
   * @param {Question|null}  [originQuestion]
   * @param {Question[]}     [subQuestions]
   * @param {Support|null}   [support]
   * @param {Answer|null}    [answer]
   * @param {Memo[]}         [memos]
   * @param {Tag[]}          [tags]
   * @param {Notice[]}       [notices]
   * @param {Escalation[]}   [escalationsFrom]
   * @param {Escalation[]}   [escalationsTo]
   */
  constructor(
    id, originQuestionId, answerId, supportId, title, content, deleted, due, createdAt,
    originQuestion = null, subQuestions = [], support = null,
    answer = null, memos = [], tags = [], notices = [],
    escalationsFrom = [], escalationsTo = [],
  ) {
    this.id = id;
    this.originQuestionId = originQuestionId;
    this.answerId = answerId;
    this.supportId = supportId;
    this.title = title;
    this.content = content;
    this.deleted = deleted;
    this.due = due ? new Date(due) : null;
    this.createdAt = new Date(createdAt);
    this.originQuestion = originQuestion;
    this.subQuestions = subQuestions;
    this.support = support;
    this.answer = answer;
    this.memos = memos;
    this.tags = tags;
    this.notices = notices;
    this.escalationsFrom = escalationsFrom;
    this.escalationsTo = escalationsTo;
  }

  /** @param {Object} json @returns {Question} */
  static fromJSON(json) {
    let oid = json.originQuestionId ?? null;
    if (oid !== null && typeof oid === 'string' && oid !== '') {
      oid = parseInt(oid, 10);
    }
    return new Question(
      json.id,
      oid,
      json.answerId ?? null,
      json.supportId ?? null,
      json.title,
      json.content,
      json.deleted ?? false,
      json.due ?? null,
      json.createdAt,
      json.originQuestion ? Question.fromJSON(json.originQuestion) : null,
      (json.subQuestions ?? []).map(Question.fromJSON),
      json.support ? Support.fromJSON(json.support) : null,
      json.answer ? Answer.fromJSON(json.answer) : null,
      (json.memos ?? []).map(Memo.fromJSON),
      (json.tags ?? []).map(Tag.fromJSON),
      (json.notices ?? []).map(Notice.fromJSON),
      (json.escalationsFrom ?? []).map(Escalation.fromJSON),
      (json.escalationsTo ?? []).map(Escalation.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// ReferManager  (中間テーブル refer_managers)
// ---------------------------------------------------------------------------
export class ReferManager {
  /**
   * @param {number}      id
   * @param {number}      answerId
   * @param {number}      referId
   * @param {Answer|null} [answer]
   * @param {Refer|null}  [refer]
   */
  constructor(id, answerId, referId, answer = null, refer = null) {
    this.id = id;
    this.answerId = answerId;
    this.referId = referId;
    this.answer = answer;
    this.refer = refer;
  }

  /** @param {Object} json @returns {ReferManager} */
  static fromJSON(json) {
    return new ReferManager(
      json.id,
      json.answerId,
      json.referId,
      json.answer ? Answer.fromJSON(json.answer) : null,
      json.refer ? Refer.fromJSON(json.refer) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}

// ---------------------------------------------------------------------------
// TagManager  (中間テーブル tag_managers)
// ---------------------------------------------------------------------------
export class TagManager {
  /**
   * @param {number}       id
   * @param {number}       tagId
   * @param {number}       questionId
   * @param {Tag|null}     [tag]
   * @param {Question|null} [question]
   */
  constructor(id, tagId, questionId, tag = null, question = null) {
    this.id = id;
    this.tagId = tagId;
    this.questionId = questionId;
    this.tag = tag;
    this.question = question;
  }

  /** @param {Object} json @returns {TagManager} */
  static fromJSON(json) {
    return new TagManager(
      json.id,
      json.tagId,
      json.questionId,
      json.tag ? Tag.fromJSON(json.tag) : null,
      json.question ? Question.fromJSON(json.question) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
