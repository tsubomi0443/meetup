/**
 * ReferManager model class.
 */

import { Answer } from '/static/js/models/Answer.js';
import { Refer } from '/static/js/models/Refer.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class ReferManager {
  /**
   * @param {number}      id
   * @param {number}      answerId
   * @param {number}      referId
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   * @param {Answer|null} [answer]
   * @param {Refer|null}  [refer]
   */
  constructor(id, answerId, referId, createdAt = null, updatedAt = null, answer = null, refer = null) {
    this.id = id;
    this.answerId = answerId;
    this.referId = referId;
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
    this.answer = answer;
    this.refer = refer;
  }

  /** @param {Object} json @returns {ReferManager} */
  static fromJSON(json) {
    return new ReferManager(
      json.id,
      json.answerId,
      json.referId,
      json.createdAt ?? null,
      json.updatedAt ?? null,
      json.answer ? Answer.fromJSON(json.answer) : null,
      json.refer ? Refer.fromJSON(json.refer) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
