/**
 * RelatedQuestion model class.
 */

import { Question } from '/static/js/models/Question.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class RelatedQuestion {
  /**
   * @param {number} id
   * @param {string} questionId
   * @param {string} relatedQuestionId
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   * @param {Question|null} [question]
   * @param {Question|null} [relatedQuestion]
   */
  constructor(id, questionId, relatedQuestionId, createdAt = null, updatedAt = null, question = null, relatedQuestion = null) {
    this.id = id;
    this.questionId = questionId;
    this.relatedQuestionId = relatedQuestionId;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
    this.question = question;
    this.relatedQuestion = relatedQuestion;
  }

  /** @param {Object} json @returns {RelatedQuestion} */
  static fromJSON(json) {
    return new RelatedQuestion(
      json.id ?? 0,
      json.questionId != null ? String(json.questionId) : '',
      json.relatedQuestionId != null ? String(json.relatedQuestionId) : '',
      json.createdAt ?? null,
      json.updatedAt ?? null,
      json.question ? Question.fromJSON(json.question) : null,
      json.relatedQuestion ? Question.fromJSON(json.relatedQuestion) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
