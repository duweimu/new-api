import { describe, expect, mock, test } from 'bun:test'
import type { ReactNode } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'

mock.module('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

mock.module('sonner', () => ({
  toast: {
    success: () => {},
    error: () => {},
  },
}))

mock.module('@/lib/currency', () => ({
  formatLocalCurrencyAmount: (amount: number) => `¥${Number(amount).toFixed(2)}`,
}))

mock.module('@/lib/format', () => ({
  formatQuota: (value: number) => String(value),
}))

function wrap(children: ReactNode) {
  return <div>{children}</div>
}

mock.module('@/components/ui/alert', () => ({
  Alert: ({ children }: { children?: ReactNode }) => wrap(children),
  AlertDescription: ({ children }: { children?: ReactNode }) => wrap(children),
}))

mock.module('@/components/ui/button', () => ({
  Button: ({
    children,
    disabled,
  }: {
    children?: ReactNode
    disabled?: boolean
  }) => <button disabled={disabled}>{children}</button>,
}))

mock.module('@/components/ui/dialog', () => ({
  Dialog: ({
    children,
    open,
  }: {
    children?: ReactNode
    open?: boolean
  }) => (open ? <div>{children}</div> : null),
  DialogContent: ({ children }: { children?: ReactNode }) => wrap(children),
  DialogHeader: ({ children }: { children?: ReactNode }) => wrap(children),
  DialogTitle: ({ children }: { children?: ReactNode }) => wrap(children),
}))

mock.module('@/components/ui/select', () => ({
  Select: ({ children }: { children?: ReactNode }) => wrap(children),
  SelectContent: ({ children }: { children?: ReactNode }) => wrap(children),
  SelectGroup: ({ children }: { children?: ReactNode }) => wrap(children),
  SelectItem: ({ children }: { children?: ReactNode }) => wrap(children),
  SelectTrigger: ({ children }: { children?: ReactNode }) => wrap(children),
  SelectValue: ({ children }: { children?: ReactNode }) => wrap(children),
}))

mock.module('@/components/ui/separator', () => ({
  Separator: () => <hr />,
}))

mock.module('@/components/group-badge', () => ({
  GroupBadge: ({ group }: { group: string }) => <span>{group}</span>,
}))

import {
  formatEpayPayableAmount,
  SubscriptionPurchaseDialog,
} from './subscription-purchase-dialog'

const basePlanRecord = {
  plan: {
    id: 1,
    title: '测试套餐',
    price_amount: 12,
    currency: 'USD',
    duration_unit: 'month' as const,
    duration_value: 1,
    quota_reset_period: 'never' as const,
    enabled: true,
    sort_order: 0,
    allow_balance_pay: true,
    max_purchase_per_user: 0,
    total_amount: 10000000,
  },
}

describe('SubscriptionPurchaseDialog', () => {
  test('formats epay payable amounts with RMB symbol', () => {
    expect(formatEpayPayableAmount(81.42)).toBe('¥81.42')
  })

  test('开关打开时同时展示余额支付和第三方支付入口', () => {
    const html = renderToStaticMarkup(
      <SubscriptionPurchaseDialog
        open={true}
        onOpenChange={() => {}}
        plan={basePlanRecord as never}
        enableOnlineTopUp={true}
        epayMethods={[{ type: 'alipay', name: '支付宝' }]}
        userQuota={25}
      />
    )

    expect(html).toContain('Pay with Balance')
    expect(html).toContain('Required')
    expect(html).toContain('Available')
    expect(html).toContain('Pay')
  })

  test('开关关闭时只禁用余额支付，不影响第三方支付入口', () => {
    const html = renderToStaticMarkup(
      <SubscriptionPurchaseDialog
        open={true}
        onOpenChange={() => {}}
        plan={
          {
            plan: {
              ...basePlanRecord.plan,
              allow_balance_pay: false,
            },
          } as never
        }
        enableOnlineTopUp={true}
        epayMethods={[{ type: 'alipay', name: '支付宝' }]}
        userQuota={25}
      />
    )

    expect(html).toContain('This plan does not allow balance redemption')
    expect(html).toContain('Pay with Balance')
    expect(html).toContain('Pay')
  })
})
