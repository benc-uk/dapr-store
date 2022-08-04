<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Shows details of an order
// ----------------------------------------------------------------------------
-->

<template>
  <div class="card order p-4">
    <div class="d-flex justify-content-between mb-2">
      <h2>{{ order.title }}</h2>
      <router-link v-if="!hideDetailsButton" class="btn btn-primary" :to="`/order/` + order.id"><i class="fa-solid fa-list-check"></i> &nbsp; ORDER DETAILS</router-link>
    </div>
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
  </div>
</template>

<script>
export default {
  name: 'Order',

  props: {
    hideDetailsButton: {
      type: Boolean,
      default: true
    },
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
