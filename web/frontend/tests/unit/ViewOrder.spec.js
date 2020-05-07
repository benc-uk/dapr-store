import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import VueRouter from 'vue-router'
localVue.use(VueRouter)
const router = new VueRouter()

import ViewOrder from '@/views/ViewOrder.vue'

jest.mock('@/mixins/api')

const orderId = 'order123'

describe('ProductSingle.vue', () => {
  it('renders product details', async () => {
    router.push({ name: 'view-order', params: { id: orderId } })

    const wrapper = mount(ViewOrder, {
      localVue,
      router
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})