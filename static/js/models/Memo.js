/**
 * Memo model class.
 */

import { Question } from '/static/js/models/Question.js';
import { User } from '/static/js/models/User.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Memo {
  /**
   * @param {number}       id
   * @param {number}       questionId
   * @param {number}       userId
   * @param {string}       content
   * @param {Question|null} [question]
   * @param {User|null}    [user]
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   */
  constructor(id, questionId, userId, content, question = null, user = null, createdAt = null, updatedAt = null) {
    this.id = id;
    this.questionId = questionId;
    this.userId = userId;
    this.content = content;
    this.question = question;
    this.user = user;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
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
      json.createdAt ?? null,
      json.updatedAt ?? null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
