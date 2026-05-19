/**
 * SupportStatus model class.
 */

import { Support } from '/static/js/models/Support.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class SupportStatus {
  /**
   * @param {number} id
   * @param {string} title
   * @param {Support[]} [supports]
   * @param {string|null} createdAt ISO8601
   * @param {string|null} updatedAt ISO8601
   */
  constructor(id, name, supports = [], createdAt = null, updatedAt = null) {
    this.id = id;
    this.name = name;
    this.title = name;
    this.supports = supports;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
  }

  /** @param {Object} json @returns {SupportStatus} */
  static fromJSON(json) {
    return new SupportStatus(
      json.id,
      json.name ?? json.title ?? '',
      (json.supports ?? []).map(Support.fromJSON),
      json.createdAt ?? null,
      json.updatedAt ?? null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
