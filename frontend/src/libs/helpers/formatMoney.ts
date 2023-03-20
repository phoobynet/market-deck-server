const numberFormat = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
})

export const formatMoney = (amount?: number): string => {
  return numberFormat.format(amount ?? 0)
}

export const formatMoneyNoSymbol = (amount?: number): string => {
  return numberFormat.format(amount ?? 0).replace('$', '')
}
