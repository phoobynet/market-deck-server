import numeral from 'numeral'

export const formatPercent = (percent: number): string => {
  return numeral(percent).format('0.00%')
}
