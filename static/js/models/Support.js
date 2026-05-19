/**
 * Support model class.
 */

import { User } from '/static/js/models/User.js';
import { SupportStatus } from '/static/js/models/SupportStatus.js';
import { Question } from '/static/js/models/Question.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Support {
  /**
   * @param {number}        id
   * @param {number}        userId
   * @param {string}        supportStatusId
   * @param {User|null}     [user]
   * @param {SupportStatus|null} [supportStatus]
   * @param {Question|null} [question]
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   */
  constructor(
    id,
    userId,
    supportStatusId,
    user = null,
    supportStatus = null,
    question = null,
    createdAt = null,
    updatedAt = null,
  ) {
    this.id = id;
    this.userId = userId;
    this.supportStatusId = supportStatusId;
    this.user = user;
    this.supportStatus = supportStatus;
    this.question = question;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
  }

  /** @param {Object} json @returns {Support} */
  static fromJSON(json) {
    return new Support(
      json.id,
      json.userId,
      json.supportStatusId,
      json.user ? User.fromJSON(json.user) : null,
      json.supportStatus ? SupportStatus.fromJSON(json.supportStatus) : null,
      json.question ? Question.fromJSON(json.question) : null,
      json.createdAt ?? null,
      json.updatedAt ?? null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
