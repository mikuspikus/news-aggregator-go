<template>
  <div id="single-news">
    <b-container fluid class="pt-5">
      <!-- News -->
      <b-row class="mb-1">
        <b-col>
          <template v-if="this.news.loading">
            <b-row>
              <b-col>
                <b-button variant="white" disabled>
                  <b-spinner small type="grow" variant="dark" label="Loading" />Loading news...
                </b-button>
              </b-col>
            </b-row>
          </template>

          <service-error-card
            v-else-if="!this.news.ok"
            :header="'News service error'"
            :lead="'Service is probably unavalible'"
          />

          <template v-else>
            <b-row>
              <b-col>
                <full-news
                  :uid="uid"
                  :useruid="this.news.user"
                  :title="this.news.title"
                  :url="this.news.url"
                  :created="this.news.created"
                  :edited="this.news.edited"
                />
              </b-col>
            </b-row>
          </template>
        </b-col>
      </b-row>

      <b-row class="mt-2 mb-2">
        <b-col>
          <!-- Comments list -->
          <template v-if="this.comment.loading">
            <b-row>
              <b-col>
                <b-button variant="white" disabled>
                  <b-spinner small type="grow" variant="dark" label="Loading" />Loading comments...
                </b-button>
              </b-col>
            </b-row>
          </template>

          <service-error-card
            v-else-if="!this.comment.ok"
            :header="'Comment service error'"
            :lead="'Service is probably unavalible'"
          />

          <template v-else>
            <template v-if="!this.comment.comments.length">
              <b-row class="text-left">
                <b-col>
                  <service-error-card
                    :header="'No comments yet'"
                    :lead="'Maybe you should add some?'"
                  />
                </b-col>
              </b-row>
            </template>

            <template v-else>
              <comment-full
                v-for="(comment, idx) in this.comment.comments"
                :key="comment.id"
                :idx="idx"
                :id="comment.id"
                :useruid="comment.user_uuid"
                :newsuid="comment.news_uuid"
                :body="comment.body"
                :created="comment.created"
                :edited="comment.edited"
                v-on:delete-comment-by-index="deleteCommentByIndex"
              />
            </template>
          </template>
        </b-col>
      </b-row>

      <!-- Comment form -->
      <b-row v-if="isLogged" class="mt-1 text-left">
        <b-col>
          <b-card>
            <add-comment-form
              :newsuid="this.uid"
              v-on:reload-comments="fetchComments(uid, page, size)"
            />
          </b-card>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import AddCommentForm from "../components/comments/Add.vue";
import ServiceErrorCard from "../components/utility/ServiceErrorCard.vue";
import FullNews from "../components/news/Full.vue";
import CommentFull from "../components/comments/Full.vue";

export default {
  name: "single-news-view",

  components: { AddCommentForm, ServiceErrorCard, FullNews, CommentFull },

  data() {
    return {
      uid: this.$route.params.uid,

      page: this.$route.query.page ? this.$route.query.page : 0,
      size: 25,

      news: {
        loading: true,
        ok: true,
        title: "",
        url: "",
        user: "",
        created: "",
        edited: "",
      },

      comment: {
        comments: [],
        pageCount: 1,
        loading: true,
        ok: true,
      },
    };
  },

  created() {
    this.fetchSingleNews(this.uid);
    this.fetchComments(this.uid, this.page, this.size);
  },

  methods: {
    deleteCommentByIndex(index) {
      this.comment.comments.splice(index, 1);
    },

    fetchSingleNews(newsuid) {
      this.news.loading = true;

      this.$http({
        url: `news/${newsuid}`,
        method: "GET",
      })
        .then((response) => {
          this.news.title = response.data.title;
          this.news.user = response.data.user;
          this.news.url = response.data.uri;
          this.news.created = response.data.created;
          this.news.edited = response.data.edited;
        })
        .catch((error) => {
          this.news.ok = false;

          this.$bvToast.toast(error, {
            title: "News fetching error",
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        })
        .finally(() => (this.news.loading = false));
    },

    fetchComments(newsuid, page, size) {
      this.comment.loading = true;

      this.$http({
        url: `news/${newsuid}/comments/`,
        params: { page: page, size: size },
        method: "GET",
      })
        .then((respone) => {
          console.log(respone.data.comments);
          this.comment.comments = respone.data.comments;
          this.comment.pageCount = respone.data.page_count;
        })
        .catch((error) => {
          this.comment.ok = false;
          this.$bvToast.toast(error, {
            title: "Comments fetching error",
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        })
        .finally(() => (this.comment.loading = false));
    },
  },

  computed: {
    isLogged() {
      return this.$store.getters.isLogged;
    },
  },
};
</script>