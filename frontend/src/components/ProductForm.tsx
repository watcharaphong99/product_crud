import { useEffect, useState, type FormEvent } from 'react'
import type { Product, ProductInput } from '../types/product'

interface ProductFormProps {
  editingProduct: Product | null
  onSubmit: (input: ProductInput) => Promise<void>
  onCancel: () => void
  isSubmitting: boolean
}

const emptyForm: ProductInput = {
  name: '',
  description: '',
  price: 0,
  stock: 0,
}

export function ProductForm({ editingProduct, onSubmit, onCancel, isSubmitting }: ProductFormProps) {
  const [form, setForm] = useState<ProductInput>(emptyForm)

  useEffect(() => {
    if (editingProduct) {
      setForm({
        name: editingProduct.name,
        description: editingProduct.description,
        price: editingProduct.price,
        stock: editingProduct.stock,
      })
    } else {
      setForm(emptyForm)
    }
  }, [editingProduct])

  const handleChange = (field: keyof ProductInput, value: string) => {
    setForm((prev) => ({
      ...prev,
      [field]: field === 'name' || field === 'description' ? value : Number(value),
    }))
  }

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    await onSubmit(form)
    if (!editingProduct) {
      setForm(emptyForm)
    }
  }

  return (
    <section className="card form-card">
      <h2>{editingProduct ? 'Edit Product' : 'Add Product'}</h2>
      <form onSubmit={handleSubmit} className="product-form">
        <div className="form-group">
          <label htmlFor="name">Name</label>
          <input
            id="name"
            type="text"
            value={form.name}
            onChange={(e) => handleChange('name', e.target.value)}
            placeholder="Product name"
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="description">Description</label>
          <textarea
            id="description"
            value={form.description}
            onChange={(e) => handleChange('description', e.target.value)}
            placeholder="Product description"
            rows={3}
          />
        </div>

        <div className="form-row">
          <div className="form-group">
            <label htmlFor="price">Price</label>
            <input
              id="price"
              type="number"
              min="0"
              step="0.01"
              value={form.price}
              onChange={(e) => handleChange('price', e.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="stock">Stock</label>
            <input
              id="stock"
              type="number"
              min="0"
              step="1"
              value={form.stock}
              onChange={(e) => handleChange('stock', e.target.value)}
              required
            />
          </div>
        </div>

        <div className="form-actions">
          <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
            {isSubmitting ? 'Saving...' : editingProduct ? 'Update Product' : 'Add Product'}
          </button>
          {editingProduct && (
            <button type="button" className="btn btn-secondary" onClick={onCancel} disabled={isSubmitting}>
              Cancel
            </button>
          )}
        </div>
      </form>
    </section>
  )
}
