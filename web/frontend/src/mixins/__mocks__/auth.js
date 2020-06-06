let user = { }

// User is a simple user state store
class User {
  constructor(token, account, userName) {
    this.token = token
    this.account = account
    this.userName = userName
  }
}
const demoUserName = 'demo@example.net'
user = new User('', { name: 'Demo User' }, demoUserName)

export default {
  methods: {
    authLogout: async function() {
      this.authUnsetUser()
    },

    authUnsetUser() {
      user = null
      localStorage.removeItem('user')
    },

    user() {
      return user
    },

    authStoredUsername() {
      return localStorage.getItem('user')
    }
  }
}