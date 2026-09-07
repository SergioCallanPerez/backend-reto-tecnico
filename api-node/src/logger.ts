type LogLevel = 'debug' | 'info' | 'warn' | 'error';

const LEVEL_RANK: Record<LogLevel, number> = { debug: 0, info: 1, warn: 2, error: 3 };
const configuredLevel = (process.env.LOG_LEVEL as LogLevel) ?? 'info';

// JSON a stdout, sin librerías externas - Cloud Run lo captura como entrada estructurada.
function log(level: LogLevel, message: string, meta: Record<string, unknown> = {}): void {
  if (LEVEL_RANK[level] < LEVEL_RANK[configuredLevel]) {
    return;
  }
  console.log(JSON.stringify({ level, message, time: new Date().toISOString(), ...meta }));
}

export const logger = {
  debug: (message: string, meta?: Record<string, unknown>) => log('debug', message, meta),
  info: (message: string, meta?: Record<string, unknown>) => log('info', message, meta),
  warn: (message: string, meta?: Record<string, unknown>) => log('warn', message, meta),
  error: (message: string, meta?: Record<string, unknown>) => log('error', message, meta),
};
