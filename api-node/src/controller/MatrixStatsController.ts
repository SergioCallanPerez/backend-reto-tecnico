import { Request, Response, Router } from 'express';
import { MatrixStatsService } from '../service/MatrixStatsService';
import { AppError } from '../model/AppError';
import { MatrixStatsResponse } from '../model/MatrixStats';

const OUTPUT_PRECISION = 6;
const MAX_MATRIX_DIMENSION = 500;

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

// No confía en que solo Go la llame: valida tamaño igual que el handler de Go.
function validateMatrix(value: unknown, name: string): number[][] {
  if (!Array.isArray(value) || value.length === 0 || !Array.isArray(value[0])) {
    throw AppError.validation(`"${name}" must be a non-empty array of arrays of numbers`);
  }
  if (value.length > MAX_MATRIX_DIMENSION || value[0].length > MAX_MATRIX_DIMENSION) {
    throw AppError.validation(`"${name}" exceeds the maximum allowed size`, 'matrices are limited to 500x500');
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
