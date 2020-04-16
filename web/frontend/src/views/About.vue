
<template>
  <div class="about">
    <img src="../assets/img/logo.svg" class="logo">
    <h1>Dapr eShop</h1>
    A demonstration app built using microservices and the Dapr framework<br>
    <b>Ben Coleman, 2020 [MIT License]</b><br>
    <br>
    v{{ version }}
    <br><br>
    &bull; Learn more about Dapr - <a href="https://dapr.io/" target="_blank">dapr.io</a><br>
    &bull; GitHub project - <a href="https://github.com/benc-uk/dapr-store" target="_blank">github.com/benc-uk/dapr-store</a>
  </div>
</template>

<script>
export default {
  data() {
    return {
      version: require('../../package.json').version,
    }
  }
}
</script>


<style scoped>
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
