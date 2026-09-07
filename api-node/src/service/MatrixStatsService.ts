import { MatrixStatsResponse } from '../model/MatrixStats';

export interface MatrixStatsService {
  computeStats(q: number[][], r: number[][]): MatrixStatsResponse;
}
