/**
 * Refer model class.
 */

import { Answer } from '/static/js/models/Answer.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Refer {
  /**
   * @param {number}   id
   * @param {string}   title
   * @param {string}   url
   * @param {Answer[]} [answers]
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   */
  constructor(id, title, url, answers = [], createdAt = null, updatedAt = null) {
    this.id = id;
    this.title = title;
    this.url = url;
    this.answers = answers;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
  }

  /** @param {Object} json @returns {Refer} */
  static fromJSON(json) {
    return new Refer(
      json.id,
      json.title,
      json.url,
      (json.answers ?? []).map(Answer.fromJSON),
      json.createdAt ?? null,
      json.updatedAt ?? null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
