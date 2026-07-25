import type { Product } from '../types/product'

interface ProductListProps {
  products: Product[]
  onEdit: (product: Product) => void
  onDelete: (product: Product) => void
  isDeletingId: number | null
}

export function ProductList({ products, onEdit, onDelete, isDeletingId }: ProductListProps) {
  if (products.length === 0) {
    return (
      <section className="card">
        <p className="empty-message">No products found. Add your first product above.</p>
      </section>
    )
  }

  return (
    <section className="card table-card">
      <h2>Product List</h2>
      <div className="table-wrapper">
        <table className="product-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Price</th>
              <th>Stock</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {products.map((product) => (
              <tr key={product.id}>
                <td data-label="Name">{product.name}</td>
                <td data-label="Description">{product.description}</td>
                <td data-label="Price">${product.price.toFixed(2)}</td>
                <td data-label="Stock">{product.stock}</td>
                <td data-label="Actions" className="actions-cell">
                  <button
                    type="button"
                    className="btn btn-edit"
                    onClick={() => onEdit(product)}
                    disabled={isDeletingId === product.id}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    className="btn btn-delete"
                    onClick={() => onDelete(product)}
                    disabled={isDeletingId === product.id}
                  >
                    {isDeletingId === product.id ? 'Deleting...' : 'Delete'}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  )
}
