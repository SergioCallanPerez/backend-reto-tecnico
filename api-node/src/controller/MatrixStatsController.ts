import { Request, Response, Router } from 'express';
import { MatrixStatsService } from '../service/MatrixStatsService';
import { AppError } from '../model/AppError';
import { MatrixStatsResponse } from '../model/MatrixStats';

const OUTPUT_PRECISION = 6;

export class MatrixStatsController {
  constructor(private readonly statsService: MatrixStatsService) {

  }

  registerRoutes(router: Router): void {
    router.post('/matrix/stats', this.postMatrixStats);
  }

  private postMatrixStats = (req: Request, res: Response): void => {
    const body: unknown = req.body ?? {};
    const { q, r } = body as { q?: unknown; r?: unknown };
    const stats = this.statsService.computeStats(validateMatrix(q, 'q'), validateMatrix(r, 'r'));
    res.json(roundStats(stats));
  };
}

function validateMatrix(value: unknown, name: string): number[][] {
  if (!Array.isArray(value) || value.length === 0 || !Array.isArray(value[0])) {
    throw AppError.validation(`"${name}" must be a non-empty array of arrays of numbers`);
  }
  const width = value[0].length;
  for (const row of value) {
    if (!Array.isArray(row) || row.length !== width || row.some((v: unknown) => typeof v !== 'number' || !Number.isFinite(v))) {
      throw AppError.validation(`"${name}" must be a rectangular matrix of finite numbers`);
    }
  }
  return value as number[][];
}

function roundStats(stats: MatrixStatsResponse): MatrixStatsResponse {
  const factor = 10 ** OUTPUT_PRECISION;
  const round = (value: number): number => Math.round(value * factor) / factor;
  return {
    ...stats,
    max: round(stats.max),
    min: round(stats.min),
    average: round(stats.average),
    sum: round(stats.sum),
  };
}
