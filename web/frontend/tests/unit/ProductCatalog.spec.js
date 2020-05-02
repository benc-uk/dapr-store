import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import ProductCatalog from '@/views/ProductCatalog.vue'

jest.mock('@/mixins/api')

describe('ProductCatalog.vue', () => {
  it('renders products in catalog', async () => {
    const wrapper = mount(ProductCatalog, { localVue,
      propsData: { }, sync: true
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })

  it('renders error messages', async () => {
    const wrapper = mount(ProductCatalog, {
      localVue,
      data: () => { return { error:'Bad thing' } }
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})