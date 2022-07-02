import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import router from '@/router'
import ProductCatalog from '@/views/ProductCatalog.vue'

jest.mock('@/services/api')

describe('ProductCatalog.vue', () => {
  it('renders products in catalog', async () => {
    const wrapper = mount(ProductCatalog, {
      propsData: {},
      sync: true,
      global: {
        plugins: [router]
      }
    })

    await flushPromises()
    expect(wrapper.html()).toMatch('Tie')
    expect(wrapper.html()).toMatch('Cravat')
  })
})
