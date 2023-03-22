import numeral from 'numeral'

export const formatMoney = (amount?: number): string => {
  return numeral(amount).format('$0,0.00')
}

export const formatMoneyNoSymbol = (amount?: number): string => {
  return numeral(amount).format('0,0.00')
}
