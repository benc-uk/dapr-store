import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import router from '@/router'
import Account from '@/views/Account.vue'

jest.mock('@/services/api')
jest.mock('@/services/auth')

describe('Account.vue', () => {
  it('renders user profile', async () => {
    const wrapper = mount(Account, {
      global: {
        plugins: [router]
      }
    })

    await flushPromises()
    expect(wrapper.html()).toMatch('Mock User')
  })
})
