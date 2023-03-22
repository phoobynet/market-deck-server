import { D3Component } from '@/components/ReportCard/ReportCardDayRange/D3Component'
import { easeLinear } from 'd3-ease'
import { Selection } from 'd3-selection'
import { symbol as d3Symbol, symbolTriangle } from 'd3-shape'

const TRIANGLE_Y_POS = 31
const TRIANGLE_SIZE = 75

export class CurrentPriceIndicator extends D3Component<SVGPathElement, null, null, undefined> {
  constructor (container: Selection<SVGElement, null, null, undefined>) {
    super(container, 'path')

    const triangle = d3Symbol().type(symbolTriangle).size(TRIANGLE_SIZE)

    this.el
      .attr('d', triangle)
      .attr('stroke', 'white')
      .attr('fill', 'white')
  }

  set priceScaled (value: number) {
    this.el
      .transition()
      .duration(500)
      .ease(easeLinear)
      .attr('transform', `translate(${value}, ${TRIANGLE_Y_POS})`)
  }
}
