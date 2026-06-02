/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React from 'react';
import {
  Banner,
  Button,
  Card,
  Divider,
  Modal,
  Select,
  Tooltip,
  Typography,
} from '@douyinfe/semi-ui';
import { IconCreditCard } from '@douyinfe/semi-icons';
import { CalendarClock, Crown, Package } from 'lucide-react';
import { SiStripe } from 'react-icons/si';
import { renderQuota } from '../../../helpers';
import {
  formatSubscriptionDuration,
  formatSubscriptionResetPeriod,
} from '../../../helpers/subscriptionFormat';

const { Text } = Typography;
const EPAY_CURRENCY_SYMBOL = '\u00A5';

function formatLocalPayAmount(amount) {
  const numericAmount = Number(amount);
  if (!Number.isFinite(numericAmount)) {
    return '-';
  }
  return `${EPAY_CURRENCY_SYMBOL}${numericAmount.toFixed(2)}`;
}

const SubscriptionPurchaseModal = ({
  t,
  visible,
  onCancel,
  selectedPlan,
  paying,
  selectedEpayMethod,
  setSelectedEpayMethod,
  epayMethods = [],
  enableOnlineTopUp = false,
  enableStripeTopUp = false,
  enableCreemTopUp = false,
  epayAmount = null,
  epayAmountLoading = false,
  epayAmountError = '',
  purchaseLimitInfo = null,
  onPayStripe,
  onPayCreem,
  onPayEpay,
}) => {
  const plan = selectedPlan?.plan;
  const totalAmount = Number(plan?.total_amount || 0);
  const planPrice = plan ? Number(plan.price_amount || 0).toFixed(2) : '0.00';
  const hasStripe = enableStripeTopUp && !!plan?.stripe_price_id;
  const hasCreem = enableCreemTopUp && !!plan?.creem_product_id;
  const hasEpay = enableOnlineTopUp && epayMethods.length > 0;
  const hasAnyPayment = hasStripe || hasCreem || hasEpay;
  const purchaseLimit = Number(purchaseLimitInfo?.limit || 0);
  const purchaseCount = Number(purchaseLimitInfo?.count || 0);
  const purchaseLimitReached =
    purchaseLimit > 0 && purchaseCount >= purchaseLimit;
  const payableAmount = formatLocalPayAmount(epayAmount);

  return (
    <Modal
      title={
        <div className='flex items-center'>
          <Crown className='mr-2' size={18} />
          {t('\u8d2d\u4e70\u8ba2\u9605\u5957\u9910')}
        </div>
      }
      visible={visible}
      onCancel={onCancel}
      footer={null}
      size='small'
      centered
    >
      {plan ? (
        <div className='space-y-4 pb-10'>
          <Card className='!rounded-xl !border-0 bg-slate-50 dark:bg-slate-800'>
            <div className='space-y-3'>
              <div className='flex items-center justify-between'>
                <Text strong className='text-slate-700 dark:text-slate-200'>
                  {t('\u5957\u9910\u540d\u79f0')}:
                </Text>
                <Typography.Text
                  ellipsis={{ rows: 1, showTooltip: true }}
                  className='text-slate-900 dark:text-slate-100'
                  style={{ maxWidth: 200 }}
                >
                  {plan.title}
                </Typography.Text>
              </div>

              <div className='flex items-center justify-between'>
                <Text strong className='text-slate-700 dark:text-slate-200'>
                  {t('\u6709\u6548\u671f')}:
                </Text>
                <div className='flex items-center'>
                  <CalendarClock size={14} className='mr-1 text-slate-500' />
                  <Text className='text-slate-900 dark:text-slate-100'>
                    {formatSubscriptionDuration(plan, t)}
                  </Text>
                </div>
              </div>

              {formatSubscriptionResetPeriod(plan, t) !==
                t('\u4e0d\u91cd\u7f6e') && (
                <div className='flex items-center justify-between'>
                  <Text strong className='text-slate-700 dark:text-slate-200'>
                    {t('\u91cd\u7f6e\u5468\u671f')}:
                  </Text>
                  <Text className='text-slate-900 dark:text-slate-100'>
                    {formatSubscriptionResetPeriod(plan, t)}
                  </Text>
                </div>
              )}

              <div className='flex items-center justify-between'>
                <Text strong className='text-slate-700 dark:text-slate-200'>
                  {t('\u603b\u989d\u5ea6')}:
                </Text>
                <div className='flex items-center'>
                  <Package size={14} className='mr-1 text-slate-500' />
                  {totalAmount > 0 ? (
                    <Tooltip
                      content={`${t('\u539f\u751f\u989d\u5ea6')}: ${totalAmount}`}
                    >
                      <Text className='text-slate-900 dark:text-slate-100'>
                        {renderQuota(totalAmount)}
                      </Text>
                    </Tooltip>
                  ) : (
                    <Text className='text-slate-900 dark:text-slate-100'>
                      {t('\u4e0d\u9650')}
                    </Text>
                  )}
                </div>
              </div>

              {plan?.upgrade_group ? (
                <div className='flex items-center justify-between'>
                  <Text strong className='text-slate-700 dark:text-slate-200'>
                    {t('\u5347\u7ea7\u5206\u7ec4')}:
                  </Text>
                  <Text className='text-slate-900 dark:text-slate-100'>
                    {plan.upgrade_group}
                  </Text>
                </div>
              ) : null}

              <Divider margin={8} />

              <div className='flex items-center justify-between'>
                <Text strong className='text-slate-700 dark:text-slate-200'>
                  {t('\u5957\u9910\u6807\u4ef7 (USD)')}:
                </Text>
                <Text strong className='text-xl text-purple-600'>
                  ${planPrice}
                </Text>
              </div>

              {hasEpay && (
                <div className='flex items-center justify-between'>
                  <Text strong className='text-slate-700 dark:text-slate-200'>
                    {t('\u5f85\u652f\u4ed8\u91d1\u989d')}:
                  </Text>
                  <Text
                    strong
                    className='text-base text-slate-900 dark:text-slate-100'
                  >
                    {epayAmountLoading ? '...' : payableAmount}
                  </Text>
                </div>
              )}
            </div>
          </Card>

          {purchaseLimitReached && (
            <Banner
              type='warning'
              description={`${t('\u5df2\u8fbe\u5230\u8d2d\u4e70\u4e0a\u9650')} (${purchaseCount}/${purchaseLimit})`}
              className='!rounded-xl'
              closeIcon={null}
            />
          )}

          {hasAnyPayment ? (
            <div className='space-y-3'>
              <Text size='small' type='tertiary'>
                {t('\u9009\u62e9\u652f\u4ed8\u65b9\u5f0f')}:
              </Text>

              {(hasStripe || hasCreem) && (
                <div className='flex gap-2'>
                  {hasStripe && (
                    <Button
                      theme='light'
                      className='flex-1'
                      icon={<SiStripe size={14} color='#635BFF' />}
                      onClick={onPayStripe}
                      loading={paying}
                      disabled={purchaseLimitReached}
                    >
                      Stripe
                    </Button>
                  )}
                  {hasCreem && (
                    <Button
                      theme='light'
                      className='flex-1'
                      icon={<IconCreditCard />}
                      onClick={onPayCreem}
                      loading={paying}
                      disabled={purchaseLimitReached}
                    >
                      Creem
                    </Button>
                  )}
                </div>
              )}

              {hasEpay && (
                <div className='space-y-2'>
                  {epayAmountError && (
                    <Banner
                      type='danger'
                      description={epayAmountError}
                      className='!rounded-xl'
                      closeIcon={null}
                    />
                  )}
                  <div className='flex gap-2'>
                    <Select
                      value={selectedEpayMethod}
                      onChange={setSelectedEpayMethod}
                      style={{ flex: 1 }}
                      size='default'
                      placeholder={t('\u9009\u62e9\u652f\u4ed8\u65b9\u5f0f')}
                      optionList={epayMethods.map((m) => ({
                        value: m.type,
                        label: m.name || m.type,
                      }))}
                      disabled={purchaseLimitReached}
                    />
                    <Button
                      theme='solid'
                      type='primary'
                      onClick={onPayEpay}
                      loading={paying}
                      disabled={
                        !selectedEpayMethod ||
                        purchaseLimitReached ||
                        epayAmountLoading ||
                        !!epayAmountError ||
                        epayAmount === null
                      }
                    >
                      {t('\u652f\u4ed8')}
                    </Button>
                  </div>
                </div>
              )}
            </div>
          ) : (
            <Banner
              type='info'
              description={t(
                '\u7ba1\u7406\u5458\u672a\u5f00\u542f\u5728\u7ebf\u652f\u4ed8\u529f\u80fd\uff0c\u8bf7\u8054\u7cfb\u7ba1\u7406\u5458\u914d\u7f6e\u3002',
              )}
              className='!rounded-xl'
              closeIcon={null}
            />
          )}
        </div>
      ) : null}
    </Modal>
  );
};

export default SubscriptionPurchaseModal;
