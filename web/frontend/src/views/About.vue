<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - About page
// ----------------------------------------------------------------------------
-->

<template>
  <div class="about">
    <img src="../assets/img/logo.svg" class="logo">
    <h1>Dapr eShop</h1>
    A demonstration app built using microservices and the Dapr framework<br>
    <b>Ben Coleman, 2020 [MIT License]</b><br>
    <br>
    Version: {{ version }}<br>
    Build Info: <pre>{{ buildInfo }}</pre>
    <br>
    &bull; Learn more about Dapr - <a href="https://dapr.io/" target="_blank">dapr.io</a><br>
    &bull; GitHub project - <a href="https://github.com/benc-uk/dapr-store" target="_blank">github.com/benc-uk/dapr-store</a>
  </div>
</template>

<script>
export default {
  data() {
    return {
      version: require('../../package.json').version,
      buildInfo: (process.env.VUE_APP_BUILD_INFO || 'N/A VUE_APP_BUILD_INFO is unset').replace(/\s\/\s/g, '\n')
    }
  }
}
</script>


<style scoped>
pre {
  margin-left:2rem;
  background-color: rgba(0, 0, 0, 0.1);
  padding: 1rem;
  width: 50%;
  font-size: 70%
}
h1 {
  font-weight: bold;
  font-size: 250%;
}
.about {
  font-size: 150%;
  background-color: #e2e2e2;
  padding: 2rem;
  border-radius: 0.6rem;
}
.logo {
  float: right;
  animation-timing-function: ease;
  animation-name: fadein;
  animation-duration: 1s;
  height: 12rem;
}
@keyframes fadein {
  from { opacity: 0; transform: scaleY(.1);}
  to   { opacity: 1; transform: scaleY(1);}
}
</style>
