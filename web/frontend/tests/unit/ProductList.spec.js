import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import router from '@/router'
import ProductList from '@/components/ProductList.vue'

jest.mock('@/services/api')

// Load mock data
let mockJson = require('fs').readFileSync(__dirname + '/../../../../testing/mock-data/products.json')
let mockProducts = JSON.parse(mockJson)

describe('ProductList.vue', () => {
  it('renders product list', async () => {
    const wrapper = mount(ProductList, {
      propsData: {
        products: mockProducts
      },
      global: {
        plugins: [router]
      }
    })

    await flushPromises()
    expect(wrapper.html()).toMatch('Cravat')
  })
})
