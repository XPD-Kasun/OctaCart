import { create } from 'zustand'

export interface Customer {
  id: string
  name: string
  email: string
  totalOrders: number
  totalSpent: number
}

export interface CustomersState {
  customers: Customer[]
  selectedCustomerId: string | null
  searchQuery: string

  setCustomers: (customers: Customer[]) => void
  selectCustomer: (id: string | null) => void
  setSearchQuery: (q: string) => void
}

export const useCustomersStore = create<CustomersState>((set) => ({
  customers: [],
  selectedCustomerId: null,
  searchQuery: '',

  setCustomers: (customers) => set({ customers }),
  selectCustomer: (id) => set({ selectedCustomerId: id }),
  setSearchQuery: (q) => set({ searchQuery: q }),
}))
