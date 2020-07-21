import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'

import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import VueRouter from 'vue-router'
localVue.use(VueRouter)
const router = new VueRouter()

import ProductSearch from '@/views/ProductSearch.vue'

jest.mock('@/services/api')

describe('ProductSearch.vue', () => {
  it('renders search for Ascot', async () => {
    router.push({ name: 'search', params: { query: 'Ascot' } })

    const wrapper = mount(ProductSearch, {
      localVue,
      router
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})