/**
 * Role model class.
 */

import { User } from '/static/js/models/User.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Role {
  /**
   * @param {number} id
   * @param {string} roleName
   * @param {User[]} [users]
   * @param {string|null} createdAt ISO8601
   * @param {string|null} updatedAt ISO8601
   */
  constructor(id, name, users = [], createdAt = null, updatedAt = null) {
    this.id = id;
    this.name = name;
    this.roleName = name;
    this.users = users;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
  }

  /** @param {Object} json @returns {Role} */
  static fromJSON(json) {
    return new Role(
      json.id,
      json.name ?? json.roleName ?? '',
      (json.users ?? []).map(User.fromJSON),
      json.createdAt ?? null,
      json.updatedAt ?? null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
