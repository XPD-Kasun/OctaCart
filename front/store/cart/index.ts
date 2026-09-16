import { create } from 'zustand'

export interface CartItem {
  productId: string
  name: string
  quantity: number
  unitPrice: number
}

export interface CartState {
  items: CartItem[]

  addItem: (item: CartItem) => void
  removeItem: (productId: string) => void
  updateQty: (productId: string, quantity: number) => void
  clearCart: () => void
  total: () => number
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],

  addItem: (item) =>
    set((s) => {
      const existing = s.items.find((i) => i.productId === item.productId)
      if (existing) {
        return {
          items: s.items.map((i) =>
            i.productId === item.productId
              ? { ...i, quantity: i.quantity + item.quantity }
              : i,
          ),
        }
      }
      return { items: [...s.items, item] }
    }),

  removeItem: (productId) =>
    set((s) => ({ items: s.items.filter((i) => i.productId !== productId) })),

  updateQty: (productId, quantity) =>
    set((s) => ({
      items: s.items.map((i) =>
        i.productId === productId ? { ...i, quantity } : i,
      ),
    })),

  clearCart: () => set({ items: [] }),

  total: () =>
    get().items.reduce((sum, i) => sum + i.unitPrice * i.quantity, 0),
}))
