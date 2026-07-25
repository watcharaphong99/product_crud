import { useCallback, useEffect, useState } from 'react'
import { ProductForm } from './components/ProductForm'
import { ProductList } from './components/ProductList'
import {
  createProduct,
  deleteProduct,
  getProducts,
  updateProduct,
} from './services/productService'
import type { Product, ProductInput } from './types/product'
import './App.css'

function App() {
  const [products, setProducts] = useState<Product[]>([])
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isDeletingId, setIsDeletingId] = useState<number | null>(null)

  const fetchProducts = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getProducts()
      setProducts(data)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to load products'
      setError(message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void fetchProducts()
  }, [fetchProducts])

  const handleSubmit = async (input: ProductInput) => {
    setIsSubmitting(true)
    setError(null)
    try {
      if (editingProduct) {
        await updateProduct(editingProduct.id, input)
        setEditingProduct(null)
      } else {
        await createProduct(input)
      }
      await fetchProducts()
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to save product'
      setError(message)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleEdit = (product: Product) => {
    setEditingProduct(product)
    setError(null)
  }

  const handleCancelEdit = () => {
    setEditingProduct(null)
    setError(null)
  }

  const handleDelete = async (product: Product) => {
    const confirmed = window.confirm(`Are you sure you want to delete "${product.name}"?`)
    if (!confirmed) {
      return
    }

    setIsDeletingId(product.id)
    setError(null)
    try {
      await deleteProduct(product.id)
      if (editingProduct?.id === product.id) {
        setEditingProduct(null)
      }
      await fetchProducts()
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to delete product'
      setError(message)
    } finally {
      setIsDeletingId(null)
    }
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>Product Management</h1>
        <p>Full-stack CRUD with React + Go Fiber</p>
      </header>

      <main className="app-main">
        {error && (
          <div className="alert alert-error" role="alert">
            {error}
          </div>
        )}

        <ProductForm
          editingProduct={editingProduct}
          onSubmit={handleSubmit}
          onCancel={handleCancelEdit}
          isSubmitting={isSubmitting}
        />

        {loading ? (
          <div className="loading">Loading products...</div>
        ) : (
          <ProductList
            products={products}
            onEdit={handleEdit}
            onDelete={handleDelete}
            isDeletingId={isDeletingId}
          />
        )}
      </main>
    </div>
  )
}

export default App
