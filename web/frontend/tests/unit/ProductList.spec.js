import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import VueRouter from 'vue-router'
localVue.use(VueRouter)
const router = new VueRouter()

import ProductList from '@/components/ProductList.vue'

jest.mock('@/mixins/api')

// Load mock data
let mockJson = require('fs').readFileSync(__dirname+'/../../../../testing/mock-data/products.json')
let mockProducts = JSON.parse(mockJson)

describe('ProductList.vue', () => {
  it('renders product list', async () => {
    const wrapper = mount(ProductList, {
      localVue,
      router,
      propsData: {
        products: mockProducts
      }
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})