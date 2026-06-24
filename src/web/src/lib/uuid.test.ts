import { describe, expect, it } from 'vitest'
import { createUUID, getOrCreateDeviceId } from './uuid'

describe('createUUID', () => {
  it('falls back when randomUUID is unavailable', () => {
    let value = 0
    const id = createUUID({
      getRandomValues(bytes) {
        const array = bytes as unknown as Uint8Array
        for (let index = 0; index < array.length; index++) array[index] = value++
        return bytes
      },
    })

    expect(id).toMatch(/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-8[0-9a-f]{3}-[0-9a-f]{12}$/)
  })
})

describe('getOrCreateDeviceId', () => {
  it('reuses the stored device id', () => {
    const storage = new Map<string, string>()

    expect(getOrCreateDeviceId({
      getItem: (key) => storage.get(key) ?? null,
      setItem: (key, value) => storage.set(key, value),
    })).toBe(storage.get('deviceId'))
  })
})
