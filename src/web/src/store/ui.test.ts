import { beforeEach, describe, expect, it } from 'vitest'
import { useUIStore } from './ui'

describe('ui tabs', () => {
  beforeEach(() => useUIStore.getState().resetTabs())

  it('deduplicates tabs and caps the cache list', () => {
    for (let index = 1; index <= 12; index++) {
      useUIStore.getState().openTab({ path: `/page-${index}`, title: `Page ${index}` })
    }

    const tabs = useUIStore.getState().tabs
    expect(tabs).toHaveLength(10)
    expect(tabs.at(-1)?.path).toBe('/page-12')
    useUIStore.getState().openTab({ path: '/page-12', title: 'Duplicate' })
    expect(useUIStore.getState().tabs).toHaveLength(10)
  })
})
