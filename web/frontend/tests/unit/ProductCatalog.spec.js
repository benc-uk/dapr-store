import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import ProductCatalog from '@/views/ProductCatalog.vue'

jest.mock('@/services/api')

describe('ProductCatalog.vue', () => {
  it('renders products in catalog', async () => {
    const wrapper = mount(ProductCatalog, { localVue,
      propsData: { }, sync: true
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})