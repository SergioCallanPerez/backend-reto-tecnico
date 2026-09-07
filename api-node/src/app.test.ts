import { describe, expect, it } from 'vitest';
import request from 'supertest';
import { createApp } from './app';

describe('api-node (integracion)', () => {
  const app = createApp();

  it('GET /health responde ok', async () => {
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body).toEqual({ status: 'ok' });
  });

  it('POST /api/v1/matrix/stats con datos validos responde 200 con el contrato esperado', async () => {
    const res = await request(app)
      .post('/api/v1/matrix/stats')
      .send({ q: [[1, 0], [0, 1]], r: [[2, 3], [0, 4]] });

    expect(res.status).toBe(200);
    expect(res.body).toEqual({
      max: 4,
      min: 0,
      average: 1.375,
      sum: 11,
      diagonalCheck: { q: true, r: false, anyDiagonal: true },
    });
  });

  it('POST /api/v1/matrix/stats sin "r" responde 400 VALIDATION_ERROR', async () => {
    const res = await request(app)
      .post('/api/v1/matrix/stats')
      .send({ q: [[1, 0], [0, 1]] });

    expect(res.status).toBe(400);
    expect(res.body.code).toBe('VALIDATION_ERROR');
  });

  it('POST /api/v1/matrix/stats con matriz no rectangular responde 400 VALIDATION_ERROR', async () => {
    const res = await request(app)
      .post('/api/v1/matrix/stats')
      .send({ q: [[1, 2], [3]], r: [[1]] });

    expect(res.status).toBe(400);
    expect(res.body.code).toBe('VALIDATION_ERROR');
  });
});
