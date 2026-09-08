import express, { Express, NextFunction, Request, Response, Router } from 'express';
import { MatrixStatsController } from './controller/MatrixStatsController';
import { MatrixStatsServiceImpl } from './service/impl/MatrixStatsServiceImpl';
import { AppError } from './model/AppError';
import { logger } from './logger';

export function createApp(): Express {
  const app = express();
  app.use(express.json());

  app.use((req: Request, res: Response, next: NextFunction) => {
    const start = Date.now();
    res.on('finish', () => {
      logger.info('http_request', {
        method: req.method,
        path: req.path,
        status: res.statusCode,
        durationMs: Date.now() - start,
      });
    });
    next();
  });

  app.get('/health', (_req: Request, res: Response) => {
    res.json({ status: 'ok' });
  });

  const statsController = new MatrixStatsController(new MatrixStatsServiceImpl());
  const v1 = Router();
  statsController.registerRoutes(v1);
  app.use('/api/v1', v1);

  app.use((err: unknown, _req: Request, res: Response, _next: NextFunction) => {
    const appError = err instanceof AppError ? err : AppError.internal(err instanceof Error ? err.message : 'error desconocido');
    const logFn = appError.code === 'VALIDATION_ERROR' ? logger.warn : logger.error;
    logFn(appError.code === 'VALIDATION_ERROR' ? 'request_rejected' : 'request_failed', appError.toJSON());
    res.status(appError.httpStatus).json(appError.toJSON());
  });

  return app;
}
