<template>
  <div id="full-news">
    <b-card no-body class="text-left m-0 p-0">
      <b-card-header header-border-variant="dark" header-bg-variant="light" class="m-0 p-0">
        <b-navbar type="light">
          <b-navbar-nav variant="light">
            <b-navbar-brand class="ml-0 pl-0">{{ this.title }}</b-navbar-brand>
            <b-nav-text class="mr-3">by</b-nav-text>
            <b-navbar-brand
              :to="{ name: 'User', params: { uuid: this.useruid }}"
              variant="light"
            >{{ user }}</b-navbar-brand>
          </b-navbar-nav>

          <template v-if="isOwner">
            <!-- Right aligned nav items -->

            <b-navbar-nav class="ml-auto">
              <b-button
                class="m-0 p-0"
                variant="light"
                border-variant="light"
                v-b-tooltip.hover.left="'Edit news'"
                @click="edit = !edit"
              >
                <b-icon icon="pencil-square"></b-icon>
              </b-button>

              <b-button
                class="m-0 p-0"
                variant="light"
                border-variant="light"
                v-b-tooltip.hover.right="'Delete news'"
                @click="deleteNews"
              >
                <b-icon icon="x-square-fill" aria-hidden="true"></b-icon>
              </b-button>
            </b-navbar-nav>
          </template>
        </b-navbar>
      </b-card-header>

      <b-card-body class="text-left">
        <template v-if="edit">
          <news-edit-form
            :uid="uid"
            :title.sync="title"
            :uri.sync="url"
            :edited.sync="edited"
            :edit.sync="edit"
          />
        </template>

        <template v-else>
          <b-card-text>
            <b-link :href="this.url">{{ this.url }}</b-link>
          </b-card-text>
        </template>
      </b-card-body>
      <template v-slot:footer>
        <small class="text-muted">{{ footer }}</small>
      </template>
    </b-card>
  </div>
</template>

<script>
import errhandler from '../../utility/errhandler.js'
import NewsEditForm from "../news/Edit.vue";

export default {
  name: "news-full",

  components: { NewsEditForm },

  props: {
    uid: { type: String, required: true },
    useruid: { type: String, required: true },
    title: { type: String, required: true },
    url: { type: String, required: true },
    created: { type: String, required: true },
    edited: { type: String, required: true },
  },

  data() {
    return {
      edit: false,
      username: "",
    };
  },

  created() {
    this.fetchUser(this.useruid);
  },

  computed: {
    isOwner() {
      return this.$store.getters.uid === this.useruid;
    },
    user() {
      return this.username ? this.username : this.useruid;
    },

    footer() {
      const created = this.$moment(this.created).format(
        "dddd, MMMM Do YYYY, HH:mm:ss"
      );
      const edited = this.$moment(this.edited).format(
        "dddd, MMMM Do YYYY, HH:mm:ss"
      );

      if (edited !== created) {
        return `Created at ${created} | Edited at ${edited}`;
      } else {
        return `Created at ${created}`;
      }
    },
  },

  methods: {
    deleteNews() {
      this.delete(this.uid);
    },

    delete(uid) {
      this.$http({
        url: `news/${uid}`,
        method: "DELETE",
      })
        .then(() => {
          this.$router.push({ name: "Home" });
        })
        .catch((error) => {
          const { message, code } = errhandler.handle(error);
          const title = "News deleting error" + (code ? ` with code ${code}` : "");
          this.$bvToast.toast(message, {
            title: title,
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        });
    },

    fetchUser(useruid) {
      this.$http({
        url: `user/${useruid}`,
        method: "GET",
      })
        .then((response) => {
          this.username = response.data.username;
        })
        .catch((error) => {
          const { message, code } = errhandler.handle(error);
          const title = "User fetching error" + (code ? ` with code ${code}` : "");
          this.$bvToast.toast(message, {
            title: title,
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        });
    },
  },
};
</script>