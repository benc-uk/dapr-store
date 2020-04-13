<template>
  <b-container>
    <b-row class="m-3">
      <b-col cols="12">
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

    <b-row class="m-3">
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
            img-src="https://i.picsum.photos/id/5/600/600.jpg"
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
.card-title {
  font-size: 1.8rem;
}
@media  (max-width: 800px) {
.card-title {
  font-size: 1.1rem;
}
}
</style>