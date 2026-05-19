/**
 * Tag model class.
 */

import { Category } from '/static/js/models/Category.js';
import { Question } from '/static/js/models/Question.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Tag {
  /**
   * @param {number}      id
   * @param {string}      title
   * @param {number}      usage
   * @param {number}      categoryId
   * @param {Category|null} [category]
   * @param {Question[]}  [questions]
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   */
  constructor(id, name, usage, categoryId, category = null, questions = [], createdAt = null, updatedAt = null) {
    this.id = id;
    this.name = name;
    this.title = name;
    this.usage = usage;
    this.categoryId = categoryId;
    this.category = category;
    this.questions = questions;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
  }

  /** @param {Object} json @returns {Tag} */
  static fromJSON(json) {
    return new Tag(
      json.id,
      json.name ?? json.title ?? '',
      json.usage,
      json.categoryId,
      json.category ? Category.fromJSON(json.category) : null,
      (json.questions ?? []).map(Question.fromJSON),
      json.createdAt ?? null,
      json.updatedAt ?? null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
