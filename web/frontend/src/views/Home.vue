<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Home page
// ----------------------------------------------------------------------------
-->

<template>
  <b-container>
    <b-row class="m-2">
      <b-col>
        <b-card
          overlay
          img-src="photo/home.jpg"
          img-alt="Dapr Menswear"
          text-variant="white"
          title="Dapr Menswear"
        >
          <b-card-text>
            Dapper clothes with Dapr microservices
          </b-card-text>
        </b-card>
      </b-col>
    </b-row>

    <b-row class="m-2">
      <b-col>
        <b-link to="/catalog">
          <b-card
            overlay
            img-src="photo/catalog.jpg"
            img-alt="Catalog"
            text-variant="white"
            title="Catalog"
          />
        </b-link>
      </b-col>
      <b-col>
        <b-link to="/offers">
          <b-card
            overlay
            img-src="photo/sale.jpg"
            img-alt="Sale"
            text-variant="white"
            title="On Sale"
          />
        </b-link>
      </b-col>
      <b-col>
        <b-link to="/about">
          <b-card
            overlay
            img-src="photo/about.jpg"
            img-alt="About"
            text-variant="white"
            title="About Us"
          />
        </b-link>
      </b-col>
    </b-row>
  </b-container>
</template>

<style scoped>
.card {
  font-size: 1.5rem;
  text-shadow: 3px 3px 3px rgba(0,0,0,0.8);
}
.card-body {
  display: flex;
  flex-direction: column;
  padding: auto 0.7rem;
}
.card-text {
  font-variant: small-caps;
  color: rgba(255,255,255, 0.7) !important;
}
.card-title, .card-text {
  margin: auto;
  font-size: 4rem;
  text-align: center;
}
@media  (max-width: 1200px) {
  .card-title, .card-text {
    font-size: 3rem;
    /* color: yellow; */
  }
}
@media  (max-width: 992px) {
  .card-title, .card-text {
    font-size: 2rem;
    /* color: red; */
  }
}
@media  (max-width: 768px) {
  .card-title, .card-text {
    font-size: 1.3rem;
    /* color: blue; */
  }
}
@media  (max-width: 576px) {
  .card-title, .card-text {
    font-size: 1.0rem;
    /* color: green; */
  }
}
</style>