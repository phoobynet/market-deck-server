export const cleanAssetName = (name: string | undefined): string => {
  name = (name ?? '').trim()

  return name
    .replace('Common Stock', '')
    .replace('American Depositary Shares', '(ADS)')
    .replace('American Depositary Receipts', '(ADR)')
}
