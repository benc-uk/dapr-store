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
      Real user sign-in disabled, as AUTH_CLIENT_ID is not set. <br />Running in demo mode, using a dummy user account
    </b-alert>
    <b-overlay :show="inprogress && !error" rounded="sm">
      <b-row class="m-1">
        <!-- sign-in -->
        <b-col>
          <b-card v-if="!demoMode" header="🙍‍♂️ Existing User - Sign In">
            <b-card-body class="d-flex flex-column" style="height: 100%">
              <div class="flex-grow-1 text-center">If you have already registered on Dapr eShop, sign in using your Microsoft indentity.</div>
              <b-button size="lg" variant="dark" @click="login">
                <img src="../assets/img/ms-tiny-logo.png" /> &nbsp; Sign-in using Microsoft Account
              </b-button>
            </b-card-body>
          </b-card>
          <b-card v-else header="🙍‍♂️ Demo User - Sign In">
            <b-card-body class="d-flex flex-column" style="height: 100%">
              <div class="flex-grow-1 text-center">
                <h3>DEMO MODE</h3>
                If you have registered the demo account, then you can sign-in with it.
              </div>
              <b-button size="lg" variant="dark" @click="login"> <fa icon="user" /> &nbsp; Sign in with Demo User </b-button>
            </b-card-body>
          </b-card>
          <br />
        </b-col>

        <!-- registration -->
        <b-col>
          <b-card v-if="!demoMode" header="🚩 New users - Register">
            <b-card-body class="d-flex flex-column" style="height: 100%">
              <div class="flex-grow-1 text-center">
                If you do not have an Dapr eShop account yet, register using your Microsoft indentity.<br />You only need to do this once.
              </div>
              <b-button size="lg" variant="dark" @click="register">
                <img src="../assets/img/ms-tiny-logo.png" /> &nbsp; Register using Microsoft Account
              </b-button>
            </b-card-body>
          </b-card>
          <b-card v-else header="🚩 Demo User - Register">
            <b-card-body class="d-flex flex-column" style="height: 100%">
              <div class="flex-grow-1 text-center">
                <h3>DEMO MODE</h3>
                Register the demo user account. You only need to do this once.
              </div>
              <b-button size="lg" variant="dark" @click="register"> <fa icon="user" /> &nbsp; Register Demo User Account </b-button>
            </b-card-body>
          </b-card>
        </b-col>
      </b-row>
    </b-overlay>
  </b-container>
</template>

<script>
import api from '../services/api'
import auth from '../services/auth'
import ErrorBox from '../components/ErrorBox'

export default {
  name: 'Login',

  components: {
    'error-box': ErrorBox
  },

  data() {
    return {
      error: null,
      inprogress: false,
      demoMode: false
    }
  },

  created() {
    this.demoMode = auth.clientId() ? false : true
  },

  methods: {
    async register() {
      this.error = null
      this.inprogress = true

      try {
        await auth.login()
        const user = auth.user()

        let resp = await api.userRegister({
          username: user.username,
          displayName: user.name || 'Unknown Name',
          profileImage: 'img/placeholder-profile.jpg'
        })

        if (resp && resp.registrationStatus == 'success') {
          console.log(`## Registered user ${user.username}`)
        } else {
          throw new Error('Something went wrong while registering user')
        }

        // !Important! Let parent component (App) know login has finished
        this.$emit('loginComplete')
        this.$router.replace({ path: '/' })
      } catch (err) {
        auth.clearLocal()
        this.error = err.message.includes('already registered') ? 'You have already registered, please sign-in' : err
      }
    },

    async login() {
      this.inprogress = true
      this.error = null

      try {
        await auth.login()
        const user = auth.user()
        if (user && user.username) {
          try {
            await api.userCheckReg(user.username)
          } catch (err) {
            auth.clearLocal()
            throw new Error("Sorry, you aren't a registered user, please use the registration option below")
          }
          this.$emit('loginComplete')
          this.$router.replace({ path: '/' })
        }
      } catch (err) {
        this.error = err
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
  height: 80px;
}
</style>
