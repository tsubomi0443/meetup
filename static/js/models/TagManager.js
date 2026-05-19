/**
 * TagManager model class.
 */

import { Tag } from '/static/js/models/Tag.js';
import { Question } from '/static/js/models/Question.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class TagManager {
  /**
   * @param {number}       id
   * @param {number}       tagId
   * @param {number}       questionId
   * @param {string|null}  createdAt
   * @param {string|null}  updatedAt
   * @param {Tag|null}     [tag]
   * @param {Question|null} [question]
   */
  constructor(id, tagId, questionId, createdAt = null, updatedAt = null, tag = null, question = null) {
    this.id = id;
    this.tagId = tagId;
    this.questionId = questionId;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
    this.tag = tag;
    this.question = question;
  }

  /** @param {Object} json @returns {TagManager} */
  static fromJSON(json) {
    return new TagManager(
      json.id,
      json.tagId,
      json.questionId,
      json.createdAt ?? null,
      json.updatedAt ?? null,
      json.tag ? Tag.fromJSON(json.tag) : null,
      json.question ? Question.fromJSON(json.question) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
