<template>
  <div id="news">
    <b-container class="mt-2" fluid>
      <b-row v-if="isLogged" class="mb-2 text-left">
        <b-col>
          <b-card>
            <add-news-form />
          </b-card>
        </b-col>
      </b-row>

      <news-page />
    </b-container>
  </div>
</template>

<script>
// @ is an alias to /src
import AddNewsForm from "@/components/news/Add.vue";
import NewsPage from "../components/news/NewsPage.vue";

export default {
  name: "index-view",

  components: { AddNewsForm, NewsPage },

  computed: {
    isLogged() {
      return this.$store.getters.isLogged;
    },
  },
};
</script>