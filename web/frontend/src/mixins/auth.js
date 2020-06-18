// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Generic reusable MSAL authentication mixin for Vue.js
// ----------------------------------------------------------------------------

import * as msal from 'msal'

let msalApp
let accessTokenRequest
let accessToken

const LOGIN_SCOPES = [ 'user.read' ]

export default {
  // Store the accessToken here for use by components (which is unlikely)
  data: function () {
    return {
      // For weird reasons we need to do this, I think it's a Vuejs bug with data in mixins
      accessToken: accessToken
    }
  },

  methods: {
    //
    // Call to set up MSAL when using real authentication with AzureAD
    // To run in demo/dummy mode simply never call this
    //
    authInitMsal(clientId, apiScopes = []) {
      msalApp = new msal.UserAgentApplication({
        auth: {
          clientId: clientId,
          redirectUri: window.location.origin
        },
        cache: {
          cacheLocation: 'localStorage'
        }
      })

      // Build accessTokenRequest from provided list of API scopes
      // These scopes need to be set up on registered app in AzureAD
      // https://docs.microsoft.com/en-us/azure/active-directory/develop/scenario-protected-web-api-app-registration
      accessTokenRequest = {
        scopes: [ ]
      }
      for (let scope in apiScopes) {
        accessTokenRequest.scopes.push(`api://${clientId}/${scope}`)
      }
    },

    //
    // Logout any stored user
    //
    authLogout: async function() {
      this.authUnsetUser()
      if (msalApp) {
        await msalApp.logout()
      }
    },

    //
    // Call when starting app (main.js) to restore any session from cache
    //
    authRestoreUser: async function() {
      try {
        // Only try if there is a cached user
        if (msalApp.getAccount()) {
          let tokenResp = await msalApp.acquireTokenSilent(accessTokenRequest)
          console.log('### MSAL acquireTokenSilent from cache was successful')

          if (tokenResp) {
            accessToken = tokenResp.accessToken
            //user = new User(tokenResp.accessToken, msalApp.getAccount(), msalApp.getAccount().userName || msalApp.getAccount().preferred_username)
          } else {
            this.authUnsetUser()
          }
        }
      } catch (err) {
        console.log('### authTryCachedUser failed, which is OK')
      }
    },

    //
    // Login a user with a pop
    //
    authLogin: async function() {
      // Skip real login and set demo user, also save in storage
      if (!msalApp) {
        // Invented but static user info
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

        // Store dummy user
        localStorage.setItem('dummyAccount', JSON.stringify(dummyUser))
        return
      }

      // Login step
      await msalApp.loginPopup({
        scopes: LOGIN_SCOPES,
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

      accessToken = tokenResp.accessToken
    },

    //
    // Clear any stored/cached user
    //
    authUnsetUser() {
      accessToken = null
      localStorage.removeItem('dummyAccount')
      if (msalApp) { msalApp.cacheStorage.clear() }
    },

    //
    // Get the logged in user, returns null if no user logged in
    //
    user() {
      if (msalApp) {
        return msalApp.getAccount()
      } else {
        let storedUser = localStorage.getItem('dummyAccount')
        if (storedUser) {
          return JSON.parse(storedUser)
        }
        return null
      }
    },

    // Aliases for user()
    authUser() { return this.user() },
    authGetUser() { return this.user() }
  }
}