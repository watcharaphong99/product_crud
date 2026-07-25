import type { ApiError, Product, ProductInput } from '../types/product'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api'

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let message = `Request failed with status ${response.status}`
    try {
      const data = (await response.json()) as ApiError
      if (data.error) {
        message = data.error
      }
    } catch {
      // Response body is not JSON
    }
    throw new Error(message)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

export async function getProducts(): Promise<Product[]> {
  const response = await fetch(`${API_BASE_URL}/products`)
  return handleResponse<Product[]>(response)
}

export async function getProduct(id: number): Promise<Product> {
  const response = await fetch(`${API_BASE_URL}/products/${id}`)
  return handleResponse<Product>(response)
}

export async function createProduct(input: ProductInput): Promise<Product> {
  const response = await fetch(`${API_BASE_URL}/products`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  return handleResponse<Product>(response)
}

export async function updateProduct(id: number, input: ProductInput): Promise<Product> {
  const response = await fetch(`${API_BASE_URL}/products/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  return handleResponse<Product>(response)
}

export async function deleteProduct(id: number): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/products/${id}`, {
    method: 'DELETE',
  })
  await handleResponse<{ message: string }>(response)
}
