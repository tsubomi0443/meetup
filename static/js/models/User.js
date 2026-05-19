/**
 * User model class.
 */

import { Role } from '/static/js/models/Role.js';
import { Support } from '/static/js/models/Support.js';
import { Answer } from '/static/js/models/Answer.js';
import { Memo } from '/static/js/models/Memo.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class User {
  /**
   * @param {number}   id
   * @param {string}   name
   * @param {string}   email
   * @param {string}   memo
   * @param {string}   pass
   * @param {number}   roleId
   * @param {Role|null} [role]
   * @param {Support[]} [supports]
   * @param {Answer[]}  [answers]
   * @param {Memo[]}    [memos]
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   */
  constructor(
    id,
    name,
    email,
    memo,
    pass,
    roleId,
    role = null,
    supports = [],
    answers = [],
    memos = [],
    createdAt = null,
    updatedAt = null,
  ) {
    this.id = id;
    this.name = name;
    this.email = email;
    this.memo = memo;
    this.pass = pass;
    this.roleId = roleId;
    this.role = role;
    this.supports = supports;
    this.answers = answers;
    this.memos = memos;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
  }

  /** @param {Object} json @returns {User} */
  static fromJSON(json) {
    return new User(
      json.id,
      json.name ?? '',
      json.email ?? '',
      json.memo ?? '',
      '',
      json.roleId,
      json.role ? Role.fromJSON(json.role) : null,
      (json.supports ?? []).map(Support.fromJSON),
      (json.answers ?? []).map(Answer.fromJSON),
      (json.memos ?? []).map(Memo.fromJSON),
      json.createdAt ?? null,
      json.updatedAt ?? null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
