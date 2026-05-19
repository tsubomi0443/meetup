/**
 * Notice model class.
 */

import { NoticeType } from '/static/js/models/NoticeType.js';
import { Question } from '/static/js/models/Question.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Notice {
  /**
   * @param {number} id
   * @param {number} typeId
   * @param {number|null} questionId
   * @param {string|null} content
   * @param {string|null} displayDue ISO8601 or null
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   * @param {NoticeType|null} [noticeType]
   * @param {Question|null} [question]
   */
  constructor(id, typeId, questionId, content, displayDue, createdAt = null, updatedAt = null, noticeType = null, question = null) {
    this.id = id;
    this.typeId = typeId;
    this.questionId = questionId;
    this.content = content;
    this.displayDue = displayDue ? new Date(displayDue) : null;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
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
      json.createdAt ?? null,
      json.updatedAt ?? null,
      json.noticeType ? NoticeType.fromJSON(json.noticeType) : null,
      json.question ? Question.fromJSON(json.question) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
