export interface Item {
  inventoryNumber: string
  name: string
  description: string
  imagePath: string | null
  category: string
  location: string
  amount: number
  electricalAppliance: boolean
  note: string
  manufacturer?: string
  serialNumber?: string
  lastTechnicalInspectionDate?: string
  lastElectricalInspectionDate?: string
  dateOfPurchase?: string
  price?: number
  resolutionNumber?: string
}

export type ItemInput = Omit<Item, 'inventoryNumber' | 'imagePath'>

export interface ItemSearchParams {
  search?: string
  category?: string
}
