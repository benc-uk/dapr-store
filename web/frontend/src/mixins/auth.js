// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Generic reusable MSAL authentication helper for Vue.js
// ----------------------------------------------------------------------------

import * as msal from 'msal'

let msalApp

const LOGIN_SCOPES = [ 'user.read', 'openid', 'profile', 'email' ]

export default {
  methods: {
    //
    // Call to set up MSAL when using authentication with AzureAD
    //
    authConfigure(clientId) {
      if (msalApp) { return }

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
    // Logout any stored user
    //
    authLogout: async function() {
      await msalApp.logout()
    },

    //
    // Call to get user, probably cached and stored locally by MSAL
    //
    authGetAccount: function() {
      if (!msalApp) { return null }

      return msalApp.getAccount()
    },

    //
    // Call through to acquireTokenSilent or acquireTokenPopup
    //
    authAcquireToken: async function(scopes = [ 'user.read' ]) {
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
    // Login a user with a popup
    //
    authLogin: async function() {
      await msalApp.loginPopup({
        scopes: LOGIN_SCOPES,
        prompt: 'select_account'
      })
    },

    //
    // Clear any stored/cached user
    //
    authClearLocalUser() {
      if (msalApp) {
        msalApp.cacheStorage.clear()
      }
    }
  }
}