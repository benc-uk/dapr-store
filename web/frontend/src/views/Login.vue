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
      Real user sign-in disabled, as AUTH_CLIENT_ID was not set. <br>Running in demo mode, using a dummy user account
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
import api from '../mixins/api'
import auth from '../mixins/auth'
import ErrorBox from '../components/ErrorBox'

export default {
  name: 'Login',

  components: {
    'error-box': ErrorBox
  },

  mixins: [ api, auth ],

  data() {
    return {
      error: null,
      inprogress: false,
      demoMode: false
    }
  },

  created() {
    // Demo mode is on when no AUTH_CLIENT_ID is provided
    this.demoMode = this.$config.AUTH_CLIENT_ID ? false : true
  },

  methods: {
    async register() {
      this.error = null
      this.inprogress = true

      try {
        await this.authLogin()

        let resp = await this.apiUserRegister({
          'username': this.user().userName,
          'displayName': this.user().name || 'Unknown Name',
          'profileImage': 'img/placeholder-profile.jpg'
        })
        if (resp.data && resp.data.registrationStatus == 'success') {
          console.log(`## Registered user ${this.user().userName}`)
        } else {
          throw new Error('Something went wrong while registering user')
        }
        this.$forceUpdate()
        this.$router.replace({ path: '/' })
      } catch (err) {
        console.error(err)

        this.authUnsetUser()
        let errMsg = this.apiDecodeError(err)
        this.error = JSON.stringify(errMsg).includes('already registered') ? 'You have already registered, please sign-in' : errMsg
      }
    },

    async login() {
      this.inprogress = true
      this.error = null

      try {
        await this.authLogin(true)

        if (this.user() && this.user().userName) {
          try {
            await this.apiUserCheckReg(this.user().userName)
          } catch (err) {
            this.authUnsetUser()
            throw new Error('Sorry, you aren\'t a registered user, please use the registration option below')
          }
          this.$forceUpdate()
          this.$router.replace({ path: '/' })
        }
      } catch (err) {
        this.error = this.apiDecodeError(err)
      }
    },
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