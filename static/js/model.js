/**
 * model.js
 * Go の internal/infrastracture/model/ 配下の構造体に対応する
 * クライアントサイドクラス定義。
 * SSE 経由で受信した JSON オブジェクトを fromJSON() で復元できます。
 *
 * 使用例:
 *   const es = new EventSource('/sse');
 *   es.addEventListener('question', (e) => {
 *     const q = Question.fromJSON(JSON.parse(e.data));
 *     console.log(q.title);
 *   });
 */

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
    this.id       = id;
    this.roleName = roleName;
    this.users    = users;
  }

  /** @param {Object} json @returns {Role} */
  static fromJSON(json) {
    return new Role(
      json.ID,
      json.RoleName,
      (json.Users ?? []).map(User.fromJSON),
    );
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
    this.id       = id;
    this.title    = title;
    this.supports = supports;
  }

  /** @param {Object} json @returns {SupportStatus} */
  static fromJSON(json) {
    return new SupportStatus(
      json.ID,
      json.Title,
      (json.Supports ?? []).map(Support.fromJSON),
    );
  }
}

// ---------------------------------------------------------------------------
// Support
// ---------------------------------------------------------------------------
export class Support {
  /**
   * @param {number}        id
   * @param {number}        userId
   * @param {number}        supportStatusId
   * @param {User|null}     [user]
   * @param {SupportStatus|null} [supportStatus]
   * @param {Question[]}   [questions]
   */
  constructor(id, userId, supportStatusId, user = null, supportStatus = null, questions = []) {
    this.id              = id;
    this.userId          = userId;
    this.supportStatusId = supportStatusId;
    this.user            = user;
    this.supportStatus   = supportStatus;
    this.questions       = questions;
  }

  /** @param {Object} json @returns {Support} */
  static fromJSON(json) {
    return new Support(
      json.ID,
      json.UserID,
      json.SupportStatusID,
      json.User          ? User.fromJSON(json.User)                 : null,
      json.SupportStatus ? SupportStatus.fromJSON(json.SupportStatus) : null,
      (json.Questions ?? []).map(Question.fromJSON),
    );
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
    this.id       = id;
    this.name     = name;
    this.email    = email;
    this.roleId   = roleId;
    this.role     = role;
    this.supports = supports;
    this.answers  = answers;
    this.memos    = memos;
  }

  /** @param {Object} json @returns {User} */
  static fromJSON(json) {
    return new User(
      json.ID,
      json.Name,
      json.Email,
      json.RoleID,
      json.Role     ? Role.fromJSON(json.Role)               : null,
      (json.Supports ?? []).map(Support.fromJSON),
      (json.Answers  ?? []).map(Answer.fromJSON),
      (json.Memos    ?? []).map(Memo.fromJSON),
    );
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
    this.id           = id;
    this.categoryName = categoryName;
    this.tags         = tags;
  }

  /** @param {Object} json @returns {Category} */
  static fromJSON(json) {
    return new Category(
      json.ID,
      json.CategoryName,
      (json.Tags ?? []).map(Tag.fromJSON),
    );
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
    this.id         = id;
    this.title      = title;
    this.usage      = usage;
    this.categoryId = categoryId;
    this.category   = category;
    this.questions  = questions;
  }

  /** @param {Object} json @returns {Tag} */
  static fromJSON(json) {
    return new Tag(
      json.ID,
      json.Title,
      json.Usage,
      json.CategoryID,
      json.Category  ? Category.fromJSON(json.Category)         : null,
      (json.Questions ?? []).map(Question.fromJSON),
    );
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
    this.id      = id;
    this.title   = title;
    this.url     = url;
    this.answers = answers;
  }

  /** @param {Object} json @returns {Refer} */
  static fromJSON(json) {
    return new Refer(
      json.ID,
      json.Title,
      json.URL,
      (json.Answers ?? []).map(Answer.fromJSON),
    );
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
    this.id         = id;
    this.questionId = questionId;
    this.userId     = userId;
    this.content    = content;
    this.question   = question;
    this.user       = user;
  }

  /** @param {Object} json @returns {Memo} */
  static fromJSON(json) {
    return new Memo(
      json.ID,
      json.QuestionID,
      json.UserID,
      json.Content,
      json.Question ? Question.fromJSON(json.Question) : null,
      json.User     ? User.fromJSON(json.User)         : null,
    );
  }
}

// ---------------------------------------------------------------------------
// Answer
// ---------------------------------------------------------------------------
export class Answer {
  /**
   * @param {number}       id
   * @param {number}       userId
   * @param {number}       questionId
   * @param {string}       content
   * @param {string|null}  answeredAt  ISO8601 文字列 or null
   * @param {string}       createdAt   ISO8601 文字列
   * @param {User|null}    [user]
   * @param {Question|null} [question]
   * @param {Refer[]}      [refers]
   */
  constructor(id, userId, questionId, content, answeredAt, createdAt, user = null, question = null, refers = []) {
    this.id         = id;
    this.userId     = userId;
    this.questionId = questionId;
    this.content    = content;
    this.answeredAt = answeredAt ? new Date(answeredAt) : null;
    this.createdAt  = new Date(createdAt);
    this.user       = user;
    this.question   = question;
    this.refers     = refers;
  }

  /** @param {Object} json @returns {Answer} */
  static fromJSON(json) {
    return new Answer(
      json.ID,
      json.UserID,
      json.QuestionID,
      json.Content,
      json.AnsweredAt ?? null,
      json.CreatedAt,
      json.User     ? User.fromJSON(json.User)         : null,
      json.Question ? Question.fromJSON(json.Question) : null,
      (json.Refers ?? []).map(Refer.fromJSON),
    );
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
    this.id             = id;
    this.fromQuestionId = fromQuestionId;
    this.toQuestionId   = toQuestionId;
    this.escalatedAt    = new Date(escalatedAt);
    this.fromQuestion   = fromQuestion;
    this.toQuestion     = toQuestion;
  }

  /** @param {Object} json @returns {Escalation} */
  static fromJSON(json) {
    return new Escalation(
      json.ID,
      json.FromQuestionID,
      json.ToQuestionID,
      json.EscalatedAt,
      json.FromQuestion ? Question.fromJSON(json.FromQuestion) : null,
      json.ToQuestion   ? Question.fromJSON(json.ToQuestion)   : null,
    );
  }
}

// ---------------------------------------------------------------------------
// Question
// ---------------------------------------------------------------------------
export class Question {
  /**
   * @param {number}         id
   * @param {number|null}    messageId
   * @param {number|null}    originQuestionId
   * @param {number|null}    supportId
   * @param {string}         title
   * @param {string|null}    due              ISO8601 文字列 or null
   * @param {string}         createdAt        ISO8601 文字列
   * @param {Question|null}  [originQuestion]
   * @param {Question[]}     [subQuestions]
   * @param {Support|null}   [support]
   * @param {Answer[]}       [answers]
   * @param {Memo[]}         [memos]
   * @param {Tag[]}          [tags]
   * @param {Escalation[]}   [escalationsFrom]
   * @param {Escalation[]}   [escalationsTo]
   */
  constructor(
    id, messageId, originQuestionId, supportId, title, due, createdAt,
    originQuestion = null, subQuestions = [], support = null,
    answers = [], memos = [], tags = [],
    escalationsFrom = [], escalationsTo = [],
  ) {
    this.id               = id;
    this.messageId        = messageId;
    this.originQuestionId = originQuestionId;
    this.supportId        = supportId;
    this.title            = title;
    this.due              = due ? new Date(due) : null;
    this.createdAt        = new Date(createdAt);
    this.originQuestion   = originQuestion;
    this.subQuestions     = subQuestions;
    this.support          = support;
    this.answers          = answers;
    this.memos            = memos;
    this.tags             = tags;
    this.escalationsFrom  = escalationsFrom;
    this.escalationsTo    = escalationsTo;
  }

  /** @param {Object} json @returns {Question} */
  static fromJSON(json) {
    return new Question(
      json.ID,
      json.MessageID        ?? null,
      json.OriginQuestionID ?? null,
      json.SupportID        ?? null,
      json.Title,
      json.Due     ?? null,
      json.CreatedAt,
      json.OriginQuestion ? Question.fromJSON(json.OriginQuestion) : null,
      (json.SubQuestions    ?? []).map(Question.fromJSON),
      json.Support          ? Support.fromJSON(json.Support)       : null,
      (json.Answers         ?? []).map(Answer.fromJSON),
      (json.Memos           ?? []).map(Memo.fromJSON),
      (json.Tags            ?? []).map(Tag.fromJSON),
      (json.EscalationsFrom ?? []).map(Escalation.fromJSON),
      (json.EscalationsTo   ?? []).map(Escalation.fromJSON),
    );
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
    this.id       = id;
    this.answerId = answerId;
    this.referId  = referId;
    this.answer   = answer;
    this.refer    = refer;
  }

  /** @param {Object} json @returns {ReferManager} */
  static fromJSON(json) {
    return new ReferManager(
      json.ID,
      json.AnswerID,
      json.ReferID,
      json.Answer ? Answer.fromJSON(json.Answer) : null,
      json.Refer  ? Refer.fromJSON(json.Refer)   : null,
    );
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
    this.id         = id;
    this.tagId      = tagId;
    this.questionId = questionId;
    this.tag        = tag;
    this.question   = question;
  }

  /** @param {Object} json @returns {TagManager} */
  static fromJSON(json) {
    return new TagManager(
      json.ID,
      json.TagID,
      json.QuestionID,
      json.Tag      ? Tag.fromJSON(json.Tag)           : null,
      json.Question ? Question.fromJSON(json.Question) : null,
    );
  }
}
