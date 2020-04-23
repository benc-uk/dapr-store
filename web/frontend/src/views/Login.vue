<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Login and account registration
// ----------------------------------------------------------------------------
-->

<template>
  <b-container>
    <error-box :error="error" />
    <b-alert v-if="demoMode" show dismissible>
      Real user sign-in disabled, due to lack of client-id configuration. <br>Running in demo mode, using a dummy user account
    </b-alert>
    <b-overlay :show="inprogress && !error" rounded="sm">
      <b-row class="m-1">
        <!-- sign-in -->
        <b-col>
          <b-card v-if="!demoMode" header="🙍‍♂️ Existing User - Sign In">
            <b-card-body class="d-flex flex-column" style="height: 100%">
              <div class="flex-grow-1 text-center">
                If you have already registered on Dapr eShop, sign in using your Microsoft indentity.
              </div>
              <b-button size="lg" variant="dark" @click="login">
                <img src="../assets/img/ms-tiny-logo.png"> &nbsp; Sign-in using Microsoft Account
              </b-button>
            </b-card-body>
          </b-card>
          <b-card v-else header="🙍‍♂️ Demo User - Sign In">
            <b-card-body class="d-flex flex-column" style="height: 100%">
              <div class="flex-grow-1 text-center">
                <h3>DEMO MODE</h3>
                If you have registered the demo account, then you can sign-in with it.
              </div>
              <b-button size="lg" variant="dark" @click="login">
                <fa icon="user" /> &nbsp; Sign in with Demo User
              </b-button>
            </b-card-body>
          </b-card>
          <br>
        </b-col>

        <!-- registration -->
        <b-col>
          <b-card v-if="!demoMode" header="🚩 New users - Register">
            <b-card-body class="d-flex flex-column" style="height: 100%">
              <div class="flex-grow-1 text-center">
                If you do not have an Dapr eShop account yet, register using your Microsoft indentity.<br>You only need to do this once.
              </div>
              <b-button size="lg" variant="dark" @click="register">
                <img src="../assets/img/ms-tiny-logo.png"> &nbsp; Register using Microsoft Account
              </b-button>
            </b-card-body>
          </b-card>
          <b-card v-else header="🚩 Demo User - Register">
            <b-card-body class="d-flex flex-column" style="height: 100%">
              <div class="flex-grow-1 text-center">
                <h3>DEMO MODE</h3>
                Register the demo user account. You only need to do this once.
              </div>
              <b-button size="lg" variant="dark" @click="register">
                <fa icon="user" /> &nbsp; Register Demo User Account
              </b-button>
            </b-card-body>
          </b-card>
        </b-col>
      </b-row>
    </b-overlay>
  </b-container>
</template>

<script>
import { userProfile, msalApp, accessTokenRequest } from '../main'
import api from '../mixins/api'
import User from '../user'
import ErrorBox from '../components/ErrorBox'

export default {
  name: 'Login',

  components: {
    'error-box': ErrorBox
  },

  mixins: [ api ],

  data() {
    return {
      error: null,
      inprogress: false,
      demoMode: false
    }
  },

  created() {
    // Demo mode is on when no AUTH_CLIENT_ID is provided
    this.demoMode = process.env.VUE_APP_AUTH_CLIENT_ID ? false : true
  },

  methods: {
    async register() {
      this.error = null
      this.inprogress = true
      let authUser = await this.authenicateUser()
      let regUserRequest = {
        'username': authUser.userName,
        'displayName': authUser.account.name || 'Unknown Name',
        'profileImage': 'img/placeholder-profile.jpg'
      }

      try {
        // Moved in front of apiUserRegister call so we have a bearer token
        Object.assign(userProfile, authUser)
        localStorage.setItem('user', userProfile.userName)

        // Sidetracked by getting user's photo, resulted in Base64 encoding hell
        // let graphTokenResp = await msalApp.acquireTokenSilent({ scopes: [ 'user.read' ] })
        // let graphPhoto = await axios.get('https://graph.microsoft.com/beta/me/photo/$value', { headers: { Authorization: `Bearer ${graphTokenResp.accessToken}` } })
        await this.apiUserRegister(regUserRequest)
        this.$router.replace({ path: '/' })
      } catch (err) {
        Object.assign(userProfile, new User())
        localStorage.removeItem('user')
        let errMsg = this.apiDecodeError(err)
        console.log(JSON.stringify(errMsg))

        this.error = JSON.stringify(errMsg).includes('already registered') ? 'You have already registered, please sign-in' : errMsg
      }
    },

    async login() {
      this.inprogress = true
      this.error = null
      let authUser = await this.authenicateUser()

      if (authUser && authUser.userName) {
        try {
          try {
            await this.apiUserCheckReg(authUser.userName)
          } catch (err) {
            throw new Error('Sorry, you aren\'t a registered user, please use the registation option below')
          }

          Object.assign(userProfile, authUser)
          localStorage.setItem('user', userProfile.userName)

          this.$router.replace({ path: '/' })
        } catch (err) {
          Object.assign(userProfile, new User())
          localStorage.removeItem('user')
          this.error = this.apiDecodeError(err)
        }
      }
    },

    async authenicateUser() {
      this.error = null

      // In demo mode only one fake account is supported, it has no token
      if (this.demoMode) {
        let dummyUser = new User('', { name: 'Demo User' }, 'demo@example.net')
        dummyUser.dummy = true
        return dummyUser
      }

      let loginRequest = { scopes: [ 'user.read' ], prompt: 'select_account' }

      try {
        let tokenResp

        // 1. Login with popup
        await msalApp.loginPopup(loginRequest)
        console.log('### MSAL loginPopup was successful')
        try {
          // 2. Try to aquire token silently
          tokenResp = await msalApp.acquireTokenSilent(accessTokenRequest)
          console.log('### MSAL acquireTokenSilent was successful')
        } catch (tokenErr) {
          // 3. Silent process might have failed so try via popup
          tokenResp = await msalApp.acquireTokenPopup(accessTokenRequest)
          console.log('### MSAL acquireTokenPopup was successful')
        }

        // Just in case check, probably never triggers
        if (!tokenResp.accessToken) {
          throw new Error('Failed to aquire access token')
        }

        // Build user object to return
        let authUser = new User(tokenResp.accessToken, msalApp.getAccount(), msalApp.getAccount().userName || msalApp.getAccount().preferred_username)
        console.log(`### MSAL user ${authUser.userName} has authenticated`)
        return authUser
      } catch (err) {
        console.error(`### MSAL error! ${err.toString()}`)
        this.error = this.apiDecodeError(err)
      }
    }
  }
}
</script>

<style scoped>
.card {
  height: 350px;
}
.card-body {
  font-size: 20px;
}
.btn {
  height: 80px
}
</style>