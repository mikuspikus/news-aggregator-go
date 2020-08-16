<template>
  <div id="news-page">
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
            <b-pagination-nav align="fill" :link-gen="linkGenerator" :number-of-pages="pageCount" />
          </b-col>
        </b-row>
      </template>
    </template>
  </div>
</template>

<script>
// @ is an alias to /src
import ServiceErrorCard from "../utility/ServiceErrorCard.vue";
import NewsShort from "../news/Short.vue";
import errhandler from "../../utility/errhandler.js";

export default {
  name: "news-index",

  components: { ServiceErrorCard, NewsShort },

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

          const { message, code } = errhandler.handle(error);
          const title =
            "News fetching error" + (code ? ` with code ${code}` : "");

          this.$bvToast.toast(message, {
            title: title,
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