import numeral from 'numeral'

export const formatPercent = (percent: number): string => {
  return numeral(percent).format('0.00%')
}

export const formatPercentAbs = (percent: number): string => {
  return numeral(Math.abs(percent)).format('0.00%')
}
