/**
 * Answer model class.
 */

import { User } from '/static/js/models/User.js';
import { Refer } from '/static/js/models/Refer.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Answer {
  /**
   * @param {number}       id
   * @param {number|string} userId
   * @param {string}       content
   * @param {boolean}      isFinal
   * @param {string}       createdAt   ISO8601 文字列
   * @param {string|null}  updatedAt   ISO8601 文字列
   * @param {User|null}    [user]
   * @param {Refer[]}      [refers]
   */
  constructor(id, userId, content, isFinal, createdAt, updatedAt, user = null, refers = []) {
    this.id = id;
    this.userId = userId;
    this.content = content;
    this.isFinal = Boolean(isFinal);
    this.createdAt = new Date(createdAt);
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
    this.user = user;
    this.refers = refers;
  }

  /** @param {Object} json @returns {Answer} */
  static fromJSON(json) {
    return new Answer(
      json.id,
      json.userId,
      json.content,
      json.isFinal ?? false,
      json.createdAt ?? new Date().toISOString(),
      json.updatedAt ?? null,
      json.user ? User.fromJSON(json.user) : null,
      (json.refers ?? []).map(Refer.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
