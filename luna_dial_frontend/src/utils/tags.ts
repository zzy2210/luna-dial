// 从 JSON 字符串解析 tags
export const parseTagsFromJsonString = (tags?: string): string[] => {
  if (!tags) {
    return [];
  }

  try {
    const parsed: unknown = JSON.parse(tags);
    if (!Array.isArray(parsed)) {
      return [];
    }

    return parsed.filter((item): item is string => typeof item === 'string');
  } catch {
    return [];
  }
};

// 将 tags 序列化为 JSON 字符串
export const stringifyTagsToJson = (tags?: string[]): string => {
  return JSON.stringify(tags ?? []);
};
