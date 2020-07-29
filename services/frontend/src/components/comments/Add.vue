<template>
  <div id="add-comment-form">
    <b-form @submit="submit" @reset="reset" v-if="show">
      <b-form-group
        id="input-group-comment-body"
        label="Comment:"
        label-for="textarea-comment-body"
      >
        <b-textarea
          id="textarea-comment-body"
          v-model="form.body"
          type="textarea"
          placeholder="Write your comment here..."
          rows="3"
          max-rows="6"
        />
      </b-form-group>

      <b-button class="mr-1" type="reset" variant="white">
        <b-icon-x-square-fill />
      </b-button>

      <b-button class="ml-1" type="submit" variant="white">
        <b-icon-check2-square />
      </b-button>
    </b-form>
  </div>
</template>

<script>
export default {
  name: "add-comment-form",

  props: {
    newsuid: { type: String, required: true },
  },

  data() {
    return {
      errors: [],
      show: true,
      form: {
        body: "",
      },
    };
  },

  methods: {
    post(data) {
      this.$http({
        url: `news/${this.newsuid}/comments/`,
        data: data,
        method: "POST",
      })
        .then(() => {
          this.$emit("reload-comments");
        })
        .catch((error) => {
          this.$bvToast.toast(error, {
            title: "Error",
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        });
    },

    submit(event) {
      event.preventDefault();

      this.errors = [];
      this.post({ body: this.form.body });
    },

    reset(event) {
      event.preventDefault();

      this.form.body = "";
      this.show = false;
      this.$nextTick(() => (this.show = true));
    },
  },
};
</script>