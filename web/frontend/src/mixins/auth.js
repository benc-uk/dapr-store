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
        }
      })

      accessTokenRequest = {
        scopes: [ `api://${clientId}/store-api` ]
      }
    },

    authLogout: async function() {
      let userName = user.userName
      this.authUnsetUser()
      if (userName != demoUserName) {
        await msalApp.logout()
      }
    },

    authLogin: async function(doLogin = false) {
      if (!msalApp) {
        user = new User('', { name: 'Demo User' }, demoUserName)
        localStorage.setItem('user', user.userName)
        return
      }

      if (doLogin) {
        await msalApp.loginPopup({
          scopes: [ 'user.read' ],
          prompt: 'select_account'
        })
      }

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

      //return tokenResp.accessToken
      user = new User(tokenResp.accessToken, msalApp.getAccount(), msalApp.getAccount().userName || msalApp.getAccount().preferred_username)
      console.log(user)

      localStorage.setItem('user', user.userName)
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