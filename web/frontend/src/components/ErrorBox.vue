<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Reusable error alert box component
// ----------------------------------------------------------------------------
-->

<template>
  <div v-if="error" class="alert alert-danger" show>
    <h4>There was a problem 😥</h4>
    <div class="errmsg">
      {{ error }}
    </div>
  </div>
</template>

<script>
export default {
  name: 'ErrorBox',

  props: {
    error: {
      validator: (prop) => typeof prop === 'string' || typeof prop === 'object' || prop === null,
      required: true
    }
  }
}
</script>

<style scoped>
.errmsg {
  color: #222;
  font-family: 'Consolas', 'Courier New', Courier, monospace;
}
</style>
