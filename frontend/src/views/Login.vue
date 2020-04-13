<template>
  <div>
    <error-box :error="error" />

    <b-card header="🙍‍♂️ Existing User - Sign In">
      <b-card-body class="text-center">
        <b-button size="lg" variant="dark" @click="login">
          <img src="../assets/img/ms-tiny-logo.png"> &nbsp; Sign in with Microsoft Account
        </b-button>
      </b-card-body>
    </b-card>

    <br>

    <b-card header="🚩 New Users - Register">
      <b-card-body class="text-center">
        <b-button size="lg" variant="dark" @click="register">
          <img src="../assets/img/ms-tiny-logo.png"> &nbsp; Register with Microsoft Account
        </b-button>
      </b-card-body>
    </b-card>
  </div>
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
      error: null
    }
  },

  methods: {
    async register() {
      let authUser = await this.authenicateUser()
      let regUserRequest = {
        'username': authUser.userName,
        'displayName': authUser.account.name || 'New User'
      }

      try {
        // Sidetracked by getting user's photo, resulted in Base64 encoding hell
        // let graphTokenResp = await msalApp.acquireTokenSilent({ scopes: [ 'user.read' ] })
        // let graphPhoto = await axios.get('https://graph.microsoft.com/beta/me/photo/$value', { headers: { Authorization: `Bearer ${graphTokenResp.accessToken}` } })
        await this.apiUserRegister(regUserRequest)
        Object.assign(userProfile, authUser)
        localStorage.setItem('user', userProfile.userName)
        this.$router.replace({ path: '/' })
      } catch (err) {
        this.failed = true
        this.error = this.apiDecodeError(err)
      }
    },

    async login() {
      let authUser = await this.authenicateUser()

      if (authUser && authUser.userName) {
        try {
          let userCheck = await this.apiUserGet(authUser.userName)
          if (!userCheck && !userCheck.userName) {
            throw new Error('Please register first')
          }
          Object.assign(userProfile, authUser)
          localStorage.setItem('user', userProfile.userName)
          this.$router.replace({ path: '/' })
        } catch (err) {
          this.failed = true
          this.error = this.apiDecodeError(err)
        }
      }
    },

    async authenicateUser() {
      this.loginFailed = false

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
        // authUser.token = tokenResp.accessToken
        // authUser.account = msalApp.getAccount()
        // authUser.userName = authUser.account.userName || authUser.account.preferred_username
        let authUser = new User(tokenResp.accessToken, msalApp.getAccount(), msalApp.getAccount().userName || msalApp.getAccount().preferred_username)
        console.log(`### MSAL user ${authUser.userName} has authenticated`)
        return authUser
      } catch (err) {
        //if (!user.token) {
        console.error(`### MSAL error! ${err.toString()}`)
        this.failed = true
        this.error = this.apiDecodeError(err)
        //}
      }
    }
  }
}


</script>