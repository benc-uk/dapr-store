import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import VueRouter from 'vue-router'
localVue.use(VueRouter)
const router = new VueRouter()

import ProductSingle from '@/views/ProductSingle.vue'

jest.mock('@/mixins/api')

describe('ProductSingle.vue', () => {
  it('renders product details', async () => {
    router.push({ name: 'single-product', params: { id: '1' } })

    const wrapper = mount(ProductSingle, {
      localVue,
      router
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})