import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import ProductOffers from '@/views/ProductOffers.vue'

jest.mock('@/mixins/api')

describe('ProductOffers.vue', () => {
  it('renders products on offer', async () => {
    const wrapper = mount(ProductOffers, {
      localVue
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})