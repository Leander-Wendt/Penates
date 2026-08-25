export interface Organisation {
  id: string
  name: string
  description: string
}

export type OrganisationInput = Omit<Organisation, 'id'>
