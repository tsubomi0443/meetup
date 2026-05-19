/**
 * NoticeType model class.
 */

import { Notice } from '/static/js/models/Notice.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class NoticeType {
  /**
   * @param {number} id
   * @param {string} name
   * @param {Notice[]} [notices]
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   */
  constructor(id, name, notices = [], createdAt = null, updatedAt = null) {
    this.id = id;
    this.name = name;
    this.notices = notices;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
  }

  /** @param {Object} json @returns {NoticeType} */
  static fromJSON(json) {
    return new NoticeType(
      json.id,
      json.name,
      (json.notices ?? []).map(Notice.fromJSON),
      json.createdAt ?? null,
      json.updatedAt ?? null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
