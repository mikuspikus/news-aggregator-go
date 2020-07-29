<template>
  <div id="comment">
    <b-card class="text-left m-0 p-0" no-body>
      <b-card-header header-border-variant="dark" header-bg-variant="dark" class="m-0 p-0">
        <b-navbar type="dark">
          <b-navbar-brand :to="{name: 'User', params: {uuid: useruid }}">{{ user }}</b-navbar-brand>
          <template v-if="isOwner">
            <!-- Right aligned nav items -->
            <b-navbar-nav class="ml-auto">
              <b-button
                class="m-0 p-0"
                variant="dark"
                v-b-tooltip.hover.left="'Edit comment'"
                @click="edit = !edit"
              >
                <b-icon icon="pencil-square" />
              </b-button>

              <b-button
                class="m-0 p-0"
                variant="dark"
                v-b-tooltip.hover.right="'Delete comment'"
                @click="deleteComment"
              >
                <b-icon icon="x-square-fill" />
              </b-button>
            </b-navbar-nav>
          </template>
        </b-navbar>
      </b-card-header>

      <b-card-body class="text-left">
        <template v-if="edit">
          <comment-edit-form
            :id="id"
            :useruid="useruid"
            :newsuid="newsuid"
            :edited.sync="edited"
            :body.sync="body"
            :edit.sync="edit"
          />
        </template>

        <template v-else>
          <b-card-text>{{ body }}</b-card-text>
        </template>
      </b-card-body>
      <template v-slot:footer>
        <small class="text-muted p-0 m-0">{{ footer }}</small>
      </template>
    </b-card>
  </div>
</template>

<script>
import CommentEditForm from "../comments/Edit.vue";

export default {
  name: "comment-full",

  components: { CommentEditForm },

  data() {
    return {
      username: "",
      edit: false,
    };
  },

  props: {
    idx: { type: Number, required: true },

    id: { type: Number, required: true },
    useruid: { type: String, required: true },
    newsuid: { type: String, required: true },
    body: { type: String, required: true },
    created: { type: String, required: true },
    edited: { type: String, required: true },
  },

  created() {
    this.fetchUser(this.useruid);
  },

  methods: {
    delete(newsuid, id) {
      this.$http({
        url: `news/${newsuid}/comments/${id}`,
        method: "DELETE",
      })
        .then(() => {
          this.$emit("delete-comment-by-index", this.idx);
        })
        .catch((error) => {
          console.log(error);
          this.$bvToast.toast(error, {
            title: "Full comment deleting error",
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        });
    },

    deleteComment() {
      this.delete(this.newsuid, this.id);
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
          this.$bvToast.toast(error, {
            title: "Full comment fetching error",
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        });
    },
  },

  computed: {
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

    isOwner() {
      return this.$store.getters.uid === this.useruid;
    },
  },
};
</script>
