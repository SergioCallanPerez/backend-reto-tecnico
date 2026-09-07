export type ErrorCode = 'VALIDATION_ERROR' | 'UPSTREAM_ERROR' | 'INTERNAL_ERROR';

const HTTP_STATUS_BY_CODE: Record<ErrorCode, number> = {
  VALIDATION_ERROR: 400,
  UPSTREAM_ERROR: 502,
  INTERNAL_ERROR: 500,
};

export class AppError extends Error {
  constructor(
    public readonly code: ErrorCode,
    message: string,
    public readonly detail?: string,
  ) {
    super(message);
  }

  get httpStatus(): number {
    return HTTP_STATUS_BY_CODE[this.code];
  }

  toJSON(): { code: ErrorCode; message: string; detail?: string } {
    return { code: this.code, message: this.message, detail: this.detail };
  }

  static validation(message: string, detail?: string): AppError {
    return new AppError('VALIDATION_ERROR', message, detail);
  }

  static internal(message: string): AppError {
    return new AppError('INTERNAL_ERROR', message);
  }
}
