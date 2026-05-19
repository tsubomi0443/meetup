/**
 * Sender model class.
 */

import { SenderTalk } from '/static/js/models/SenderTalk.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class Sender {
  /**
   * @param {number} id
   * @param {string} name
   * @param {string} departmentName
   * @param {SenderTalk[]} [senderTalks]
   */
  constructor(id, name, departmentName, senderTalks = []) {
    this.id = id;
    this.name = name;
    this.departmentName = departmentName;
    this.senderTalks = senderTalks;
  }

  /** @param {Object} json @returns {Sender} */
  static fromJSON(json) {
    return new Sender(
      json.id,
      json.name ?? '',
      json.departmentName ?? '',
      (json.senderTalks ?? []).map(SenderTalk.fromJSON),
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
