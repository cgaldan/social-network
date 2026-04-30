/** Format API date for HTML date input (YYYY-MM-DD). */
export function toDateInputValue(value) {
  if (!value) return "";
  if (typeof value === "string") {
    return value.length >= 10 ? value.slice(0, 10) : "";
  }
  return "";
}
