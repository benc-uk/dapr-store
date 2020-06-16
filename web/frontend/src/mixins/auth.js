import * as msal from 'msal'

let msalApp = null
let accessTokenRequest = null
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

export default {
  methods: {
    authInitMsal(clientId) {
      msalApp = new msal.UserAgentApplication({
        auth: {
          clientId: clientId,
          redirectUri: window.location.origin
        },
        cache: {
          cacheLocation: 'localStorage'
        }
      })

      accessTokenRequest = {
        scopes: [ `api://${clientId}/store-api` ]
      }
    },

    authLogout: async function() {
      this.authUnsetUser()
      if (msalApp) {
        await msalApp.logout()
      }
    },

    authTryCachedUser: async function() {
      // Skip real login and return cached demo/dummy user
      if (!msalApp) {
        let storedUser = localStorage.getItem('user')
        if (storedUser) {
          user = JSON.parse(storedUser)
        }
        return
      }

      try {
        // Only try if there is a cached user
        if (msalApp.getAccount()) {
          let tokenResp = await msalApp.acquireTokenSilent(accessTokenRequest)
          console.log('### MSAL acquireTokenSilent from cache was successful')

          if (tokenResp){
            user = new User(tokenResp.accessToken, msalApp.getAccount(), msalApp.getAccount().userName || msalApp.getAccount().preferred_username)
          } else {
            this.authUnsetUser()
          }
        }
      } catch (err) {
        console.log('### authTryCachedUser failed, which is OK')
      }
    },

    authLogin: async function() {
      // Skip real login and set demo user, also save in storage
      if (!msalApp) {
        user = new User('', { name: 'Demo User' }, demoUserName)
        localStorage.setItem('user', JSON.stringify(user))
        return
      }

      // Login step
      await msalApp.loginPopup({
        scopes: [ 'user.read' ],
        prompt: 'select_account'
      })

      let tokenResp
      try {
        // 1. Try to acquire token silently
        tokenResp = await msalApp.acquireTokenSilent(accessTokenRequest)
        console.log('### MSAL acquireTokenSilent was successful')
      } catch (tokenErr) {
        // 2. Silent process might have failed so try via popup
        tokenResp = await msalApp.acquireTokenPopup(accessTokenRequest)
        console.log('### MSAL acquireTokenPopup was successful')
      }

      // Just in case check, probably never triggers
      if (!tokenResp.accessToken) {
        throw new Error('Failed to acquire access token')
      }

      user = new User(tokenResp.accessToken, msalApp.getAccount(), msalApp.getAccount().userName || msalApp.getAccount().preferred_username)
    },

    authUnsetUser() {
      user = null
      localStorage.removeItem('user')
      msalApp.cacheStorage.clear()
    },

    user() {
      return user
    },
  }
}