/**
 * Escalation model class.
 */

import { Question } from '/static/js/models/Question.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Escalation {
  /**
   * @param {number}       id
   * @param {number}       fromQuestionId
   * @param {number}       toQuestionId
   * @param {string}       escalatedAt    ISO8601 文字列
   * @param {string|null}  createdAt
   * @param {string|null}  updatedAt
   * @param {Question|null} [fromQuestion]
   * @param {Question|null} [toQuestion]
   */
  constructor(id, fromQuestionId, toQuestionId, escalatedAt, createdAt = null, updatedAt = null, fromQuestion = null, toQuestion = null) {
    this.id = id;
    this.fromQuestionId = fromQuestionId;
    this.toQuestionId = toQuestionId;
    this.escalatedAt = new Date(escalatedAt);
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
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
      json.createdAt ?? null,
      json.updatedAt ?? null,
      json.fromQuestion ? Question.fromJSON(json.fromQuestion) : null,
      json.toQuestion ? Question.fromJSON(json.toQuestion) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
