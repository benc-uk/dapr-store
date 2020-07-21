import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'

import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import VueRouter from 'vue-router'
localVue.use(VueRouter)
const router = new VueRouter()

import Cart from '@/views/Cart.vue'

jest.mock('@/services/api')
jest.mock('@/services/auth')


describe('Cart.vue', () => {
  it('renders cart contents', async () => {
    const wrapper = mount(Cart, {
      localVue,
      router
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})