export interface Product {
  id: number
  name: string
  description: string
  price: number
  stock: number
}

export interface ProductInput {
  name: string
  description: string
  price: number
  stock: number
}

export interface ApiError {
  error: string
}
