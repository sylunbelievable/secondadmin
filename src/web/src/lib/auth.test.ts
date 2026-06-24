import { describe, expect, it } from 'vitest'
import { can } from './auth'

describe('can', () => {
  it('checks button permissions from the current menu response', () => {
    expect(can({ menus: [], permissions: ['system:user:create'] }, 'system:user:create')).toBe(true)
    expect(can({ menus: [], permissions: [] }, 'system:user:create')).toBe(false)
  })
})
