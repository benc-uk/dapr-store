import * as msal from 'msal'

// MSAL object used for signing in users with MS identity platform
let msalApp

export default {
  //
  // Configure with clientId or empty string/null to set in "demo" mode
  //
  async configure(clientId) {
    // Can only call configure once
    if (msalApp) { return }

    // If no clientId then create a mock MSAL UserAgentApplication
    // Allows us to run without Azure AD for demos & local dev
    if (!clientId) {
      console.log('### Azure AD sign-in: disabled. Will run in demo mode with dummy demo@example.net account')

      const dummyUser = {
        accountIdentifier: 'e11d4d0c-1c70-430d-a644-aed03a60e059',
        homeAccountIdentifier: '',
        userName: 'demo@example.net',
        name: 'Demo User',
        idToken: null,
        idTokenClaims: null,
        sid: '',
        environment: ''
      }

      msalApp = {
        clientId: null,

        loginPopup() {
          localStorage.setItem('dummyAccount', JSON.stringify(dummyUser))
          return new Promise((resolve) => resolve())
        },
        logout() {
          localStorage.removeItem('dummyAccount')
          window.location.href = '/'
          return new Promise((resolve) => resolve())
        },
        acquireTokenSilent() {
          return new Promise((resolve) => resolve())
        },
        cacheStorage: {
          clear() {
            localStorage.removeItem('dummyAccount')
          }
        },
        getAccount() {
          return JSON.parse(localStorage.getItem('dummyAccount'))
        }
      }
      return
    }

    console.log(`### Azure AD sign-in: enabled. Using clientId: ${clientId}`)
    msalApp = new msal.UserAgentApplication({
      auth: {
        clientId: clientId
      },
      cache: {
        cacheLocation: 'localStorage'
      }
    })
  },

  //
  // Return the configured client id
  //
  clientId() {
    if (!msalApp) { return null }

    return msalApp.clientId
  },

  //
  // Login a user with a popup
  //
  async login() {
    if (!msalApp) { return }

    const LOGIN_SCOPES = [ 'user.read', 'openid', 'profile', 'email' ]
    await msalApp.loginPopup({
      scopes: LOGIN_SCOPES,
      prompt: 'select_account'
    })
  },

  //
  // Logout any stored user
  //
  async logout() {
    if (!msalApp) { return }

    await msalApp.logout()
  },

  //
  // Call to get user, probably cached and stored locally by MSAL
  //
  user() {
    if (!msalApp) { return null }

    return msalApp.getAccount()
  },

  //
  // Call through to acquireTokenSilent or acquireTokenPopup
  //
  async acquireToken(scopes = [ 'user.read' ]) {
    if (!msalApp) { return null }

    // Set scopes for token request
    const accessTokenRequest = {
      scopes
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
      throw new Error('### accessToken not found in response, that\'s bad')
    }

    return tokenResp.accessToken
  },

  //
  // Clear any stored/cached user
  //
  clearLocal() {
    if (msalApp) {
      msalApp.cacheStorage.clear()
    }
  }
}