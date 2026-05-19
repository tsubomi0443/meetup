import {createAvatar} from '../vendor/dicebear-core.min.js';
import {create, meta, schema} from '../vendor/dicebear-notionists.min.js';

export const AVATAR_SIZE = {
    ZERO: 0,
    SMALL: 25,
    MIDDLE: 50,
    LARGE: 100,
    XLARGE: 250,
}

/**
 * textを元にAvatarの画像を生成します。
 * @param {String} text
 * @param {Number} size
 * @returns 
 */
export function createAvatarBy(text = "", size = 500) {
    const avatar = createAvatar({"create": create, "meta": meta, "schema": schema}, {
      seed: text,
      size: size
    })
    const svg = avatar.toString();
    return svg
  }