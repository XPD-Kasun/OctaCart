import { create } from 'zustand'

export type OrderStatus = 'pending' | 'processing' | 'shipped' | 'delivered' | 'cancelled'

export interface Order {
  id: string
  customerId: string
  status: OrderStatus
  total: number
  createdAt: string
}

export interface OrdersState {
  orders: Order[]
  selectedOrderId: string | null
  filter: OrderStatus | 'all'

  setOrders: (orders: Order[]) => void
  selectOrder: (id: string | null) => void
  setFilter: (filter: OrderStatus | 'all') => void
}

export const useOrdersStore = create<OrdersState>((set) => ({
  orders: [],
  selectedOrderId: null,
  filter: 'all',

  setOrders: (orders) => set({ orders }),
  selectOrder: (id) => set({ selectedOrderId: id }),
  setFilter: (filter) => set({ filter }),
}))
