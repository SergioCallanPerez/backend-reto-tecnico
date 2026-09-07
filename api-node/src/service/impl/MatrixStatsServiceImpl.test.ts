import { describe, expect, it } from 'vitest';
import { MatrixStatsServiceImpl } from './MatrixStatsServiceImpl';
import { AppError } from '../../model/AppError';

describe('MatrixStatsServiceImpl', () => {
  const service = new MatrixStatsServiceImpl();

  it('calcula max, min, promedio y suma sobre ambas matrices combinadas', () => {
    const stats = service.computeStats(
      [[1, 2], [3, 4]],
      [[5, 6], [7, 8]],
    );
    expect(stats.max).toBe(8);
    expect(stats.min).toBe(1);
    expect(stats.sum).toBe(36);
    expect(stats.average).toBe(4.5);
  });

  it('detecta q diagonal y r no diagonal por separado', () => {
    const stats = service.computeStats(
      [[1, 0], [0, 1]],
      [[2, 3], [0, 4]],
    );
    expect(stats.diagonalCheck).toEqual({ q: true, r: false, anyDiagonal: true });
  });

  it('detecta ambas diagonales', () => {
    const stats = service.computeStats(
      [[1, 0], [0, 1]],
      [[5, 0], [0, 7]],
    );
    expect(stats.diagonalCheck).toEqual({ q: true, r: true, anyDiagonal: true });
  });

  it('reporta anyDiagonal en false cuando ninguna lo es', () => {
    const stats = service.computeStats(
      [[1, 2], [3, 4]],
      [[5, 6], [7, 8]],
    );
    expect(stats.diagonalCheck).toEqual({ q: false, r: false, anyDiagonal: false });
  });

  it('una matriz no cuadrada nunca es diagonal, sin importar su contenido', () => {
    const stats = service.computeStats(
      [[1, 0], [0, 1], [0, 0]],
      [[9]],
    );
    expect(stats.diagonalCheck.q).toBe(false);
  });

  it('lanza AppError de validacion si no hay valores para calcular', () => {
    expect(() => service.computeStats([[]], [[]])).toThrowError(AppError);
  });
});
