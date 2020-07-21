import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'

import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import VueRouter from 'vue-router'
localVue.use(VueRouter)
const router = new VueRouter()

import Account from '@/views/Account.vue'

jest.mock('@/services/api')
jest.mock('@/services/auth')

describe('Account.vue', () => {
  it('renders user profile', async () => {
    const wrapper = mount(Account, {
      localVue,
      router
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})