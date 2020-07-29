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

      <template v-if="isLoading">
        <b-button variant="white" disabled>
          <b-spinner small type="grow" variant="dark" label="Loading" />Loading...
        </b-button>
      </template>

      <template v-else-if="!ok">
        <b-row class="text-left">
          <b-col>
            <service-error-card />
          </b-col>
        </b-row>
      </template>

      <template v-else>
        <template v-if="!news.length || !pageCount">
          <b-row class="text-left">
            <b-col>
              <service-error-card :header="'No news yet'" :lead="'Maybe you should add some?'" />
            </b-col>
          </b-row>
        </template>
        <template v-else>
          <template v-for="snews in news">
            <b-row :key="snews.uid">
              <b-col>
                <news-short :uid="snews.uid" :title="snews.title" :created="snews.created" />
              </b-col>
            </b-row>
          </template>

          <b-row class="mt-2">
            <b-col>
              <b-pagination-nav
                align="fill"
                :link-gen="linkGenerator"
                :number-of-pages="pageCount"
              />
            </b-col>
          </b-row>
        </template>
      </template>
    </b-container>
  </div>
</template>

<script>
// @ is an alias to /src
import AddNewsForm from "@/components/news/Add.vue";
import ServiceErrorCard from "../components/utility/ServiceErrorCard.vue";
import NewsShort from "../components/news/Short.vue";

export default {
  name: "news-index",

  components: { ServiceErrorCard, NewsShort, AddNewsForm },

  data() {
    return {
      pageSize: 25,
      isLoading: true,
      ok: true,
      news: [],
      pageCount: 1,
      page: this.$route.query.page ? this.$route.query.page : 0,
    };
  },

  computed: {
    isLogged() {
      return this.$store.getters.isLogged;
    },
  },

  created() {
    this.fetch(this.page);
  },

  methods: {
    linkGenerator(pageNumber) {
      return pageNumber === 1 ? "?" : `?page=${pageNumber}`;
    },

    fetch(page) {
      this.$http({
        url: "news/",
        params: { page: page, size: this.pageSize },
        method: "GET",
      })
        .then((response) => {
          this.news = response.data.news;
          this.pageCount = response.data.page_count;
        })
        .catch((error) => {
          this.ok = false;

          this.$bvToast.toast(error, {
            title: "Error",
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        })
        .finally(() => {
          this.isLoading = false;
        });
    },
  },
};
</script>

<style>
.pagination > li > a {
  background-color: white;
  color: #343a40;
}

.pagination > li > a:focus,
.pagination > li > a:hover,
.pagination > li > span:focus,
.pagination > li > span:hover {
  color: #343a40;
  background-color: #eee;
  border-color: #ddd;
}

.pagination > .active > a {
  color: white;
  background-color: #343a40 !important;
  border: solid 1px #343a40 !important;
}

.pagination > .active > a:hover {
  background-color: #343a40 !important;
  border: solid 1px #343a40;
}
</style>