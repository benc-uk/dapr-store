<!--
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - Home page
// ----------------------------------------------------------------------------
-->

<template>
  <div>
    <div class="row m-2">
      <div class="col"><!-- empty padding --></div>
      <div class="col-9">
        <div class="card">
          <img src="photo/home.jpg" class="card-img" alt="Dapr Menswear" />
          <div class="card-img-overlay">
            <div class="card-text">Dapper clothes with Dapr microservices</div>
          </div>
        </div>
      </div>
      <div class="col"><!-- empty padding --></div>
    </div>

    <div class="row m-2">
      <div class="col">
        <div class="card">
          <router-link to="/catalog">
            <img src="photo/catalog.jpg" class="card-img" alt="Catalog" />
            <div class="card-img-overlay">
              <div class="card-text">Catalog</div>
            </div>
          </router-link>
        </div>
      </div>
      <div class="col">
        <div class="card">
          <router-link to="/offers">
            <img src="photo/sale.jpg" class="card-img" alt="Sale" />
            <div class="card-img-overlay">
              <div class="card-text">Sale</div>
            </div>
          </router-link>
        </div>
      </div>
      <div class="col">
        <div class="card">
          <router-link to="/about">
            <img src="photo/about.jpg" class="card-img" alt="About Us" />
            <div class="card-img-overlay">
              <div class="card-text">About Us</div>
            </div>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  font-size: 1.5rem;
  text-shadow: 3px 3px 3px rgba(0, 0, 0, 0.8);
}
.card-body {
  display: flex;
  flex-direction: column;
  padding: auto 0.7rem;
}
.card-text {
  font-variant: small-caps;
  color: rgba(255, 255, 255, 0.7) !important;
}
.card-title,
.card-text {
  margin: auto;
  font-size: 4rem;
  text-align: center;
}
@media (max-width: 1200px) {
  .card-title,
  .card-text {
    font-size: 3rem;
    /* color: yellow; */
  }
}
@media (max-width: 992px) {
  .card-title,
  .card-text {
    font-size: 2rem;
    /* color: red; */
  }
}
@media (max-width: 768px) {
  .card-title,
  .card-text {
    font-size: 1.3rem;
    /* color: blue; */
  }
}
@media (max-width: 576px) {
  .card-title,
  .card-text {
    font-size: 1rem;
    /* color: green; */
  }
}
</style>
