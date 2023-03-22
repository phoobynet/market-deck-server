import { Selection } from 'd3-selection'
import { D3Component } from '@/components/ReportCard/ReportCardDayRange/D3Component'

export class RangeBar extends D3Component<SVGRectElement, null, null, undefined> {
  private _height = 20

  constructor (container: Selection<SVGElement, null, null, undefined>) {
    super(container, 'rect')

    this.el.attr('height', this._height)
  }

  set signSymbol (value: string) {
    this.el
      .attr('stroke-width', 1)
      .attr('stroke', value === '+'
        ? '#a8f9a8'
        : value === '-'
          ? '#f96a6a'
          : 'white')
  }
}

