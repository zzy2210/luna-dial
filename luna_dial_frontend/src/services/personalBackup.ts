import apiClient from './api';
import { ApiResponse, PersonalBackupImportResult } from '../types';

const FALLBACK_FILENAME_PREFIX = 'luna-dial-personal-backup';

const parseFilenameFromContentDisposition = (value: string | undefined): string | null => {
  if (!value) {
    return null;
  }

  // 常见格式：
  // - attachment; filename="xxx.tar.gz"
  // - attachment; filename=xxx.tar.gz
  // - attachment; filename*=UTF-8''xxx.tar.gz
  const parts = value.split(';').map(p => p.trim());

  for (const part of parts) {
    if (part.toLowerCase().startsWith('filename*=')) {
      const raw = part.slice('filename*='.length).trim();
      const normalized = raw.replace(/^UTF-8''/i, '');
      try {
        return decodeURIComponent(normalized.replace(/^"|"$/g, ''));
      } catch {
        return normalized.replace(/^"|"$/g, '');
      }
    }
  }

  for (const part of parts) {
    if (part.toLowerCase().startsWith('filename=')) {
      return part.slice('filename='.length).trim().replace(/^"|"$/g, '');
    }
  }

  return null;
};

class PersonalBackupService {
  async exportArchive(): Promise<{ blob: Blob; filename: string }> {
    const response = await apiClient.post<Blob>(
      '/api/v1/personal-backup/export',
      undefined,
      {
        responseType: 'blob',
        timeout: 60000,
      }
    );

    const contentDisposition = (response.headers?.['content-disposition'] as string | undefined) ?? undefined;
    const filenameFromHeader = parseFilenameFromContentDisposition(contentDisposition);
    const filename = filenameFromHeader || `${FALLBACK_FILENAME_PREFIX}.tar.gz`;

    return { blob: response.data, filename };
  }

  async importOverwrite(file: File, force: boolean): Promise<PersonalBackupImportResult> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await apiClient.post<ApiResponse<PersonalBackupImportResult>>(
      '/api/v1/personal-backup/import-overwrite',
      formData,
      {
        params: { force: force ? 'true' : 'false' },
        timeout: 60000,
      }
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '导入备份失败');
  }
}

export default new PersonalBackupService();

