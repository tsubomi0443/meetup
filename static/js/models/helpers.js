/**
 * Shared JSON helpers for model classes.
 */
export function formToApiJson(form) {
  return JSON.parse(
    JSON.stringify(form, (key, value) => {
      if (key === 'uid' || key === 'talkroomId') return undefined;
      return value instanceof Date ? value.toISOString() : value;
    }),
  );
}
