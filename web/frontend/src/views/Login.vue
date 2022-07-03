<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Login and account registration
// ----------------------------------------------------------------------------
-->

<template>
  <div>
    <error-box :error="error" />
    <div v-if="demoMode" class="alert alert-info" show dismissible>
      Real user sign-in disabled, as AUTH_CLIENT_ID is not set. <br />Running in demo mode, using a dummy user account
    </div>
    <div class="row">
      <!-- Login -->
      <div class="col">
        <div v-if="!demoMode" class="card">
          <div class="card-header"><i class="fa-solid fa-right-to-bracket"></i> Existing User - Sign In</div>
          <div class="card-body text-center">
            <div class="text-center">If you have already registered on Dapr eShop, sign in using your Microsoft identity</div>
            <div class="btn btn-dark btn-lg mt-3" @click="login"><img src="../assets/img/ms-tiny-logo.png" /> &nbsp; Sign-in with Microsoft Account</div>
          </div>
        </div>
        <div v-else class="card">
          <div class="card-header"><i class="fa-solid fa-right-to-bracket"></i> Demo User - Sign In</div>
          <div class="card-body text-center">
            <h3>DEMO MODE</h3>
            <div class="text-center">Sign in with the demo user account</div>
            <div class="btn btn-dark btn-lg mt-3" @click="login"><i class="fa-solid fa-user-pen"></i> &nbsp; Sign in with Demo User</div>
          </div>
        </div>
      </div>

      <!-- Register -->
      <div class="col">
        <div class="card">
          <div v-if="!demoMode" class="card">
            <div class="card-header"><i class="fa-solid fa-address-card"></i> New users - Register</div>
            <div class="card-body text-center">
              <div class="text-center">Register a new account using your Microsoft identity, either MSA or Work & School.</div>
              <div class="btn btn-dark btn-lg mt-3" @click="register">
                <img src="../assets/img/ms-tiny-logo.png" /> &nbsp; Register With Microsoft Account
              </div>
            </div>
          </div>
          <div v-else class="card">
            <div class="card-header"><i class="fa-solid fa-address-card"></i> Demo User - Register</div>
            <div class="card-body text-center">
              <h3>DEMO MODE</h3>
              <div class="text-center">Register demo user account. You only need to do this once.</div>
              <div class="btn btn-dark btn-lg mt-3" @click="register"><i class="fa-solid fa-user-pen"></i> &nbsp; Register Demo User Account</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
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
    this.demoMode = false;//auth.clientId() ? false : true
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
  line-height: 70px;
}
</style>
