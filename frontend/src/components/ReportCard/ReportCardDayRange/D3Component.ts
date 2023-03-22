import { BaseType, Selection } from 'd3-selection'

export abstract class D3Component<GElement extends BaseType, Datum, PElement extends BaseType, PDatum> {
  protected el: Selection<GElement, null, null, undefined>

  protected constructor (container: Selection<SVGElement, null, null, undefined>, tag: string) {
    this.el = container.append(tag)
  }

  set x (value: number) {
    if (isNaN(value)) {
      return
    }
    this.el.attr('x', value)
  }

  set y (value: number) {
    if (isNaN(value)) {
      return
    }
    this.el.attr('y', value)
  }

  set width (value: number) {
    if (isNaN(value)) {
      return
    }
    this.el.attr('width', value)
  }

  set height (value: number) {
    this.el.attr('height', value)
  }

  get element() {
    return this.el
  }
}
