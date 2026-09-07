import { MatrixStatsService } from '../MatrixStatsService';
import { MatrixStatsResponse } from '../../model/MatrixStats';
import { AppError } from '../../model/AppError';

const DIAGONAL_EPSILON = 1e-9;

export class MatrixStatsServiceImpl implements MatrixStatsService {
  computeStats(q: number[][], r: number[][]): MatrixStatsResponse {
    const values = [...q.flat(), ...r.flat()];
    if (values.length === 0) {
      throw AppError.validation('no values found to compute statistics from');
    }

    const sum = values.reduce((acc, value) => acc + value, 0);
    const diagonalQ = isDiagonal(q);
    const diagonalR = isDiagonal(r);

    return {
      max: Math.max(...values),
      min: Math.min(...values),
      average: sum / values.length,
      sum,
      diagonalCheck: {
        q: diagonalQ,
        r: diagonalR,
        anyDiagonal: diagonalQ || diagonalR,
      },
    };
  }
}

// Solo definida para matrices cuadradas
function isDiagonal(matrix: number[][]): boolean {
  const rows = matrix.length;
  const cols = matrix[0]?.length ?? 0;
  if (rows === 0 || rows !== cols) {
    return false;
  }
  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      if (i !== j && Math.abs(matrix[i][j]) > DIAGONAL_EPSILON) {
        return false;
      }
    }
  }
  return true;
}
