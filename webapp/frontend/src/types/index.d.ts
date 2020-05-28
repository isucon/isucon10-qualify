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

export interface Coordinate {
  latitude: number
  longitude: number
}
