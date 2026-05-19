/**
 * SenderTalk model class.
 */

import { Sender } from '/static/js/models/Sender.js';
import { formToApiJson } from '/static/js/models/helpers.js';

export class SenderTalk {
  /**
   * @param {number} id
   * @param {string} content
   * @param {string|number} senderId
   * @param {string|number} questionId
   * @param {string|null} createdAt
   * @param {string|null} updatedAt
   * @param {Sender|null} [sender]
   */
  constructor(id, content, senderId, questionId, createdAt = null, updatedAt = null, sender = null) {
    this.id = id;
    this.content = content;
    this.senderId = senderId != null ? String(senderId) : '';
    this.questionId = questionId != null ? String(questionId) : '';
    this.createdAt = createdAt ? new Date(createdAt) : null;
    this.updatedAt = updatedAt ? new Date(updatedAt) : null;
    this.sender = sender;
  }

  /** @param {Object} json @returns {SenderTalk} */
  static fromJSON(json) {
    return new SenderTalk(
      json.id,
      json.content ?? '',
      json.senderId ?? '',
      json.questionId ?? '',
      json.createdAt ?? null,
      json.updatedAt ?? null,
      json.sender ? Sender.fromJSON(json.sender) : null,
    );
  }

  /** @returns {Object} */
  static toModel(form) {
    return formToApiJson(form);
  }
}
