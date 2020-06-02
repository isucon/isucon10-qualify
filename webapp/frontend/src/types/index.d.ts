export interface Estate {
  id: string
  name: string
  thumbnail: string
  address: string
  description: string
  doorHeight: number
  doorWidth: number
  features: string
  latitude: number
  longitude: number
  rent: number
}

export interface Chair {
  id: string
  name: string
  thumbnail: string
  description: string
  height: number
  width: number
  depth: number
  features: string
  price: number
  color: string
  kind: string
}

export interface Coordinate {
  latitude: number
  longitude: number
}

export interface Range {
  id: number
  min: number
  max: number
}

export interface RangeList {
  prefix: string
  suffix: string
  ranges: Range[]
}

export interface EstateRangeMap {
  doorWidth: RangeList
  doorHeight: RangeList
  rent: RangeList
}

export interface EstateSearchCondition {
  doorWidthRangeId: string
  doorHeightRangeId: string
  rentRangeId: string
  features: string
  page: number
  perPage: number
}

export interface EstateSearchResponse {
  estates: Estate[]
  count: number
}

export interface ChairRangeMap {
  price: RangeList
  height: RangeList
  width: RangeList
  depth: RangeList
}

export interface ChairSearchCondition {
  priceRangeId: string
  heightRangeId: string
  widthRangeId: string
  depthRangeId: string
  color: string
  kind: string
  features: string
  page: number
  perPage: number
}

export interface ChairSearchResponse {
  chairs: Chair[]
  count: number
}
