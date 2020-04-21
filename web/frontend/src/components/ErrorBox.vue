<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Reusable error alert box component
// ----------------------------------------------------------------------------
-->

<template>
  <b-alert v-if="error" show variant="warning">
    <h4>There was a problem 😥</h4>
    <div class="errmsg">
      {{ error }}
    </div>
  </b-alert>
</template>

<script>
export default {
  name: 'ErrorBox',

  props: {
    error: {
      validator: (prop) => typeof prop === 'string' || typeof prop === 'object' || prop === null,
      required: true
    }
  },
}
</script>

<style scoped>
.errmsg {
  color: #222;
  font-family: 'Consolas', 'Courier New', Courier, monospace;
}
</style>