/**
 * Question model class.
 */

import { Support } from '/static/js/models/Support.js';
import { Answer } from '/static/js/models/Answer.js';
import { Memo } from '/static/js/models/Memo.js';
import { Tag } from '/static/js/models/Tag.js';
import { Notice } from '/static/js/models/Notice.js';
import { Escalation } from '/static/js/models/Escalation.js';
import { RelatedQuestion } from '/static/js/models/RelatedQuestion.js';
import { SenderTalk } from '/static/js/models/SenderTalk.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Question {
  /**
   * @param {number}         id
   * @param {number|null}    originQuestionId
   * @param {number|null}    supportId
   * @param {string}         title
   * @param {string}         content
   * @param {string|null}    due              ISO8601 文字列 or null
   * @param {string}         createdAt        ISO8601 文字列
   * @param {string|null}    updatedAt
   * @param {Question|null}  [originQuestion]
   * @param {Question[]}     [subQuestions]
   * @param {Support|null}   [support]
   * @param {Answer[]}       [answers]
   * @param {Memo[]}         [memos]
   * @param {Tag[]}          [tags]
   * @param {Notice[]}       [notices]
   * @param {Escalation[]}   [escalationsFrom]
   * @param {Escalation[]}   [escalationsTo]
   * @param {RelatedQuestion[]} [relatedQuestions]
   * @param {SenderTalk[]} [senderTalks]
   */
  constructor(
    id,
    originQuestionId,
    supportId,
    title,
    content,
    due,
    createdAt,
    updatedAt,
    originQuestion = null,
    subQuestions = [],
    support = null,
    answers = [],
    memos = [],
    tags = [],
    notices = [],
    escalationsFrom = [],
    escalationsTo = [],
    relatedQuestions = [],
    senderTalks = [],
  ) {
    this.id = id;
    this.originQuestionId = originQuestionId;
    this.supportId = supportId;
    this.title = title;
    this.content = content;
    this.due = due ? new Date(due) : null;
    this.createdAt = new Date(createdAt);
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
    this.originQuestion = originQuestion;
    this.subQuestions = subQuestions;
    this.support = support;
    this.answers = answers;
    this.memos = memos;
    this.tags = tags;
    this.notices = notices;
    this.escalationsFrom = escalationsFrom;
    this.escalationsTo = escalationsTo;
    this.relatedQuestions = relatedQuestions;
    this.senderTalks = senderTalks;
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
      json.supportId ?? null,
      json.title,
      json.content,
      json.due ?? null,
      json.createdAt,
      json.updatedAt ?? null,
      json.originQuestion ? Question.fromJSON(json.originQuestion) : null,
      (json.subQuestions ?? []).map(Question.fromJSON),
      json.support ? Support.fromJSON(json.support) : null,
      (json.answers ?? []).map(Answer.fromJSON),
      (json.memos ?? []).map(Memo.fromJSON),
      (json.tags ?? []).map(Tag.fromJSON),
      (json.notices ?? []).map(Notice.fromJSON),
      (json.escalationsFrom ?? []).map(Escalation.fromJSON),
      (json.escalationsTo ?? []).map(Escalation.fromJSON),
      (json.relatedQuestions ?? []).map(RelatedQuestion.fromJSON),
      (json.senderTalks ?? []).map(SenderTalk.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
