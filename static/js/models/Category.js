/**
 * Category model class.
 */

import { Tag } from '/static/js/models/Tag.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Category {
  /**
   * @param {number} id
   * @param {string} categoryName
   * @param {Tag[]}  [tags]
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   */
  constructor(id, name, tags = [], createdAt = null, updatedAt = null) {
    this.id = id;
    this.name = name;
    this.categoryName = name;
    this.tags = tags;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
  }

  /** @param {Object} json @returns {Category} */
  static fromJSON(json) {
    return new Category(
      json.id,
      json.name ?? json.categoryName ?? '',
      (json.tags ?? []).map(Tag.fromJSON),
      json.createdAt ?? null,
      json.updatedAt ?? null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
