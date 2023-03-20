import numeral from 'numeral'

export const formatLargeNumber = (number: number): string => {
  return numeral(number).format('0.00a')
}
