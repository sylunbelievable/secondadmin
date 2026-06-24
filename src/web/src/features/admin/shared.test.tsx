/* @vitest-environment jsdom */

import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { DataTable, QueryState } from './shared'

describe('admin shared components', () => {
  it('keeps stale content visible while refetching', () => {
    render(<QueryState pending={false} refetching error={undefined} empty={false}>内容</QueryState>)

    expect(screen.getByText('正在刷新…')).toBeTruthy()
    expect(screen.getByText('内容')).toBeTruthy()
  })

  it('renders tables inside a horizontal overflow container', () => {
    const { container } = render(<DataTable headers={['名称']} rows={[{ key: 1, cells: ['Second Admin'] }]} />)

    expect(container.querySelector('.overflow-x-auto')).toBeTruthy()
  })
})
