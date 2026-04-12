export const generateAlias = (tableName: string) => {
  const normalized = tableName.replace(/[^a-zA-Z0-9]/g, "");

  if (!normalized) return "";
  if (normalized.length === 1) return normalized.toLowerCase();
  if (normalized.length === 2) return normalized.toLowerCase();

  const first = normalized[0];
  const middle = normalized[Math.floor(normalized.length / 2)];
  const last = normalized[normalized.length - 1];

  return `${first}${middle}${last}`.toLowerCase();
};
