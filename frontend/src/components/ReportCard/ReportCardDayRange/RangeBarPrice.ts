import { D3Component } from '@/components/ReportCard/ReportCardDayRange/D3Component'
import { Selection } from 'd3-selection'
import { formatMoneyNoSymbol } from '@/libs/helpers/formatMoney'

export class RangeBarPrice extends D3Component<SVGTextElement, null, null, undefined> {
  constructor (container: Selection<SVGElement, null, null, undefined>) {
    super(container, 'text')
    this.el.attr('stroke', 'rgb(180, 198, 239)')
      .attr('fill', 'rgb(180, 198, 239)')
      .style('font-variant-numeric', 'tabular-nums')
      .style('font-size', '0.7rem')
  }

  set price (value: number) {
    this.el.text(formatMoneyNoSymbol(value))
  }

  get computedWidth (): number {
    return this.el.node()?.getComputedTextLength() ?? 0
  }
}
