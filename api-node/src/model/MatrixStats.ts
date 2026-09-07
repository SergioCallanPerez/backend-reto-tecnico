export interface MatrixStatsRequest {
  q: number[][];
  r: number[][];
}

export interface DiagonalCheck {
  q: boolean;
  r: boolean;
  anyDiagonal: boolean;
}

export interface MatrixStatsResponse {
  max: number;
  min: number;
  average: number;
  sum: number;
  diagonalCheck: DiagonalCheck;
}
