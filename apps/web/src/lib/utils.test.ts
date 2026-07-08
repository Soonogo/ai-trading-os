import { describe, it, expect } from 'vitest';
import { cn } from './utils';

describe('cn helper', () => {
  it('merges class strings', () => {
    const result = cn('px-2', 'py-2 bg-black', 'px-4');
    expect(typeof result).toBe('string');
    expect(result.length).toBeGreaterThan(0);
  });
});
