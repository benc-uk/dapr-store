import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import router from '@/router'
import ProductOffers from '@/views/ProductOffers.vue'

jest.mock('@/services/api')

describe('ProductOffers.vue', () => {
  it('renders products on offer', async () => {
    const wrapper = mount(ProductOffers, {
      global: {
        plugins: [router]
      }
    })

    await flushPromises()
    expect(wrapper.html()).toMatch('Black')
    expect(wrapper.html()).not.toMatch('Cravat')
  })
})
