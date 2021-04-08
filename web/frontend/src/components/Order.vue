<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Shows details of an order
// ----------------------------------------------------------------------------
-->

<template>
  <b-card class="order">
    <b-link :to="`/order/` + order.id">
      <h2>{{ order.title }}</h2>
    </b-link>
    <h2>
      Status: <span class="text-capitalize order-status" :class="['order-' + order.status]">{{ order.status }}</span>
    </h2>
    <ul>
      <li>Order Id: {{ order.id }}</li>
      <li>Amount: £{{ order.amount }}</li>
      <li>Items:</li>
      <ul>
        <li v-for="(line, index) in order.lineItems" :key="index">{{ line.count }} x {{ line.product.name }} &mdash; £{{ line.product.cost }}</li>
      </ul>
    </ul>
  </b-card>
</template>

<script>
export default {
  name: 'Order',

  props: {
    order: {
      type: Object,
      required: true
    }
  }
}
</script>

<style scoped>
.order {
  font-size: 140%;
  margin-bottom: 1rem;
}
.order-received {
  color: rgb(115, 8, 119);
  background-color: rgb(227, 189, 236);
}
.order-received {
  color: rgb(129, 66, 8);
  background-color: rgb(231, 202, 138);
}
.order-processing {
  color: rgb(23, 38, 173);
  background-color: rgb(194, 214, 240);
}
.order-complete {
  color: rgb(10, 107, 10);
  background-color: rgb(119, 223, 150);
}
.order-status {
  padding: 6px 20px;
  margin: 6px;
  display: inline-block;
  border-radius: 5px;
}
</style>
