<template>
  <div id="comment-edit-form">
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

      <div id="form-errors">
        <b-alert show dismissible variant="light" v-for="error in errors" :key="error.message">
          <hr />
          <p>{{ error.message}}</p>
          <hr />
        </b-alert>
      </div>

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
import errhandler from "../../utility/errhandler.js";

export default {
  name: "comment-edit-form",

  props: {
    id: { type: Number, required: true },

    useruid: { type: String, required: true },
    newsuid: { type: String, required: true },
    edited: { type: String, required: true },
    body: { type: String, required: true },
    edit: { type: Boolean, required: true },
  },

  data() {
    return {
      show: true,
      errors: [],
      form: {
        body: this.body,
      },
    };
  },

  methods: {
    update(uid, id, data) {
      this.$http({
        url: `news/${uid}/comments/${id}`,
        data: data,
        method: "PATCH",
      })
        .then((response) => {
          this.$emit("update:body", response.data.body);
          this.$emit("update:edited", response.data.edited);
          this.$emit("update:edit", false);
        })
        .catch((error) => {
          const { message, code } = errhandler.handle(error);
          const title =
            "Comment editing error" + (code ? ` with code ${code}` : "");
          this.$bvToast.toast(message, {
            title: title,
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        });
    },

    submit(event) {
      event.preventDefault();
      this.update(this.newsuid, this.id, this.form);
    },

    reset(event) {
      event.preventDefault();

      this.form.body = this.body;
      this.show = false;
      this.$nextTick(() => (this.show = true));
    },
  },
};
</script>